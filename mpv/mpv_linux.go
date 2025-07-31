package mpv

/*
#cgo pkg-config: mpv

#include <stdlib.h>
#include <stdio.h>
#include <mpv/client.h>
*/
import "C"

import (
	"strconv"
	"unsafe"

	"github.com/hacel/jfsh/jellyfin"
)

func setProperty(handle *C.mpv_handle, name string, format C.mpv_format, data []byte) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cdata := C.CString(string(data))
	defer C.free(unsafe.Pointer(cdata))
	status := C.mpv_set_property(handle, cname, format, unsafe.Pointer(&cdata))
	if status < 0 {
		// NOTE: possibly don't have to panic?
		panic("err in mpv_set_property")
	}
}

func observeProperty(handle *C.mpv_handle, name string, format C.mpv_format) {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	status := C.mpv_observe_property(handle, 0, n, format)
	if status < 0 {
		// NOTE: possibly don't have to panic?
		panic("err in mpv_observe_property")
	}
}

func loadFile(handle *C.mpv_handle, host string, item jellyfin.Item) {
	cmd := []string{"loadfile", getStreamingURL(host, item), "replace", "0", "start=" + strconv.Itoa(int(getResumePosition(item)))}
	ccmd := make([]*C.char, len(cmd)+1)
	for i := range cmd {
		ccmd[i] = C.CString(cmd[i])
		defer C.free(unsafe.Pointer(ccmd[i]))
	}
	status := C.mpv_command(handle, (**C.char)(&ccmd[0]))
	if status < 0 {
		// NOTE: possibly don't have to panic?
		panic("err in mpv_loadfile")
	}
}

func Play(client *jellyfin.Client, item jellyfin.Item) {
	handle := C.mpv_create()
	defer C.mpv_terminate_destroy(handle)

	setProperty(handle, "config", C.MPV_FORMAT_FLAG, []byte("1"))
	setProperty(handle, "osc", C.MPV_FORMAT_FLAG, []byte("1"))
	setProperty(handle, "input-default-bindings", C.MPV_FORMAT_FLAG, []byte("1"))
	setProperty(handle, "input-vo-keyboard", C.MPV_FORMAT_FLAG, []byte("1"))
	setProperty(handle, "force-media-title", C.MPV_FORMAT_STRING, []byte(getMediaTitle(item)))

	observeProperty(handle, "time-pos", C.MPV_FORMAT_INT64)

	status := C.mpv_initialize(handle)
	if status < 0 {
		// NOTE: possibly don't have to panic?
		panic("err in mpv_initialize")
	}

	loadFile(handle, client.Host, item)

	// TODO: should this communicate back to the main thread through a channel or something?
	var progress int64
	for {
		e := C.mpv_wait_event(handle, 1)
		switch e.event_id {
		case C.MPV_EVENT_SHUTDOWN, C.MPV_EVENT_END_FILE:
			client.ReportPlaybackStopped(item, progress)
			return
		// TODO: report progress immediately on seek?
		case C.MPV_EVENT_PROPERTY_CHANGE:
			data := (*C.mpv_event_property)(e.data)
			switch C.GoString(data.name) {
			case "time-pos":
				pos := (*int64)(data.data)
				if pos == nil {
					continue
				}
				progress = *pos
				client.ReportPlaybackProgress(item, progress)
			}
		}
	}
}
