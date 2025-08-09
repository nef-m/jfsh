// Package mpv provides functions for playing jellyfin items in mpv
package mpv

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/hacel/jfsh/jellyfin"
)

func secondsToTicks(seconds float64) int64 {
	return int64(seconds * 10_000_000)
}

func ticksToSeconds(ticks int64) float64 {
	return float64(ticks) / 10_000_000
}

type request struct {
	Command any `json:"command"`
	ID      int `json:"request_id,omitempty"`
}

type response struct {
	Error      string `json:"error"`
	ID         int    `json:"request_id,omitempty"`
	Event      string `json:"event,omitempty"`
	Name       string `json:"name,omitempty"`
	Reason     string `json:"reason,omitempty"`
	Data       any    `json:"data"`
	PlaylistID int    `json:"playlist_entry_id,omitempty"`
}

type mpv struct {
	conn    net.Conn
	scanner *bufio.Scanner
	cmd     *exec.Cmd
	socket  string
}

func createMpv() (*mpv, error) {
	socket := filepath.Join(os.TempDir(), fmt.Sprintf("jfsh-mpv-socket-%d", time.Now().UnixNano()))
	cmd := exec.Command("mpv", "--idle", "--input-ipc-server="+socket)
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to create mpv: %w", err)
	}

	// Wait for socket to be created
	var conn net.Conn
	var err error
	for range 300 {
		conn, err = net.Dial("unix", socket)
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		cmd.Process.Kill()
		return nil, fmt.Errorf("failed to connect to mpv socket: %w", err)
	}
	return &mpv{
		conn:    conn,
		scanner: bufio.NewScanner(conn),
		cmd:     cmd,
		socket:  socket,
	}, nil
}

func (c *mpv) close() error {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return fmt.Errorf("failed to close mpv socket connection: %w", err)
		}
	}
	if c.cmd != nil && c.cmd.Process != nil {
		if err := c.cmd.Process.Signal(os.Interrupt); err != nil {
			return fmt.Errorf("failed to interrupt mpv process: %w", err)
		}
		if err := c.cmd.Wait(); err != nil {
			return fmt.Errorf("failed to wait for mpv process: %w", err)
		}
	}
	if err := os.Remove(c.socket); err != nil {
		return fmt.Errorf("failed to remove mpv socket: %w", err)
	}
	return nil
}

func (c *mpv) send(command []any) error {
	req := request{Command: command}
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(append(data, '\n'))
	return err
}

func (c *mpv) setProperty(name string, value any) error {
	return c.send([]any{"set_property", name, value})
}

func (c *mpv) observeProperty(name string) error {
	return c.send([]any{"observe_property", 1, name})
}

func (c *mpv) prependFile(url, title string) error {
	cmd := []any{"loadfile", url, "insert-at", 0, map[string]any{
		"force-media-title": title,
	}}
	return c.send(cmd)
}

func (c *mpv) appendFile(url, title string) error {
	cmd := []any{"loadfile", url, "append", 0, map[string]any{
		"force-media-title": title,
	}}
	return c.send(cmd)
}

func (c *mpv) playFile(url, title string, start float64) error {
	cmd := []any{"loadfile", url, "replace", 0, map[string]any{
		"force-media-title": title,
		"start":             strconv.FormatFloat(start, 'f', 6, 64),
	}}
	return c.send(cmd)
}

func Play(client *jellyfin.Client, items []jellyfin.Item, index int) {
	mpv, err := createMpv()
	if err != nil {
		panic(fmt.Sprintf("failed to create mpv client: %v", err))
	}
	defer mpv.close()

	// makes mpv report file posisiton
	if err := mpv.observeProperty("time-pos"); err != nil {
		panic(fmt.Sprintf("failed to observe time-pos: %v", err))
	}

	// keeps track of the playlist index of items as they get loaded into mpv
	playlistIDs := make([]int, 0, len(items))

	// load file specified by index
	url := jellyfin.GetStreamingURL(client.Host, items[index])
	start := ticksToSeconds(jellyfin.GetResumePosition(items[index]))
	title := jellyfin.GetMediaTitle(items[index])
	if err := mpv.playFile(url, title, start); err != nil {
		panic(fmt.Sprintf("failed to load file: %v", err))
	}
	playlistIDs = append(playlistIDs, index)

	// append to playlist the files after the index
	for i := index + 1; i < len(items); i++ {
		url := jellyfin.GetStreamingURL(client.Host, items[i])
		title := jellyfin.GetMediaTitle(items[i])
		if err := mpv.appendFile(url, title); err != nil {
			panic(fmt.Sprintf("failed to append file: %v", err))
		}
		playlistIDs = append(playlistIDs, i)
	}

	// prepend to playlist the files before the index
	for i := index - 1; i >= 0; i-- {
		url := jellyfin.GetStreamingURL(client.Host, items[i])
		title := jellyfin.GetMediaTitle(items[i])
		if err := mpv.prependFile(url, title); err != nil {
			panic(fmt.Sprintf("failed to load file: %v", err))
		}
		playlistIDs = append(playlistIDs, i)
	}

	pos := float64(0)
	lastProgressUpdate := time.Now()
	item := items[index]
	for mpv.scanner.Scan() {
		line := mpv.scanner.Text()
		if line == "" {
			continue
		}
		var response response
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			slog.Error("failed to unmarshal response", "line", line, "err", err)
			continue
		}

		switch response.Event {
		case "property-change":
			switch response.Name {
			case "time-pos":
				if time.Since(lastProgressUpdate) < 3*time.Second {
					// debounce
					continue
				}
				data, ok := response.Data.(float64)
				if !ok {
					slog.Error("failed to parse time-pos data as float64", "line", line, "data", response.Data)
					continue
				}
				pos = data
				if err := client.ReportPlaybackProgress(item, secondsToTicks(pos)); err != nil {
					slog.Error("failed to report playback progress", "err", err)
					continue
				}
				slog.Info("reported progress", "item", item.GetName(), "pos", pos)
				lastProgressUpdate = time.Now()
			}

		case "start-file":
			// figure out what item is being played
			id := response.PlaylistID - 1
			if id >= len(playlistIDs) {
				slog.Error("start-file event for unknown playlist id", "id", response.PlaylistID)
				// user probably loaded something manually
				return
			}
			item = items[playlistIDs[response.PlaylistID-1]]
			slog.Info("received", "event", response.Event, "playlist_id", response.PlaylistID, "index", playlistIDs[response.PlaylistID-1], "item", item.GetName())
			// report playback start
			if err := client.ReportPlaybackStart(item, secondsToTicks(pos)); err != nil {
				slog.Error("failed to report playback progress", "err", err)
				continue
			}
			slog.Info("reported progress", "item", item.GetName(), "pos", pos)

		case "seek":
			slog.Info("received", "event", response.Name, "item", item.GetName())
			lastProgressUpdate = time.Time{}

		case "end-file", "shutdown":
			slog.Info("received", "event", response.Event, "item", item.GetName())
			if err := client.ReportPlaybackStopped(item, secondsToTicks(pos)); err != nil {
				slog.Error("failed to report playback stopped", "err", err)
			}
		}
	}
}
