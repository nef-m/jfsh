// Package mpv provides functions for playing jellyfin items in mpv
package mpv

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/hacel/jfsh/jellyfin"
)

type request struct {
	Command any `json:"command"`
	ID      int `json:"request_id,omitempty"`
}

type response struct {
	Error  string `json:"error"`
	ID     int    `json:"request_id,omitempty"`
	Event  string `json:"event,omitempty"`
	Name   string `json:"name,omitempty"`
	Reason string `json:"reason,omitempty"`
	Data   any    `json:"data"`
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

func (c *mpv) loadFile(url string, pos int64) error {
	cmd := []any{"loadfile", url, "replace"}
	if pos > 0 {
		cmd = append(cmd, "0", "start="+strconv.FormatInt(pos, 10))
	}
	return c.send(cmd)
}

func Play(client *jellyfin.Client, item jellyfin.Item) {
	mpv, err := createMpv()
	if err != nil {
		panic(fmt.Sprintf("failed to create mpv client: %v", err))
	}
	defer mpv.close()

	mpv.setProperty("force-media-title", jellyfin.GetMediaTitle(item))

	if err := mpv.observeProperty("time-pos"); err != nil {
		panic(fmt.Sprintf("failed to observe time-pos: %v", err))
	}

	url := jellyfin.GetStreamingURL(client.Host, item)
	pos := jellyfin.GetResumePosition(item)
	if err := mpv.loadFile(url, pos); err != nil {
		panic(fmt.Sprintf("failed to load file: %v", err))
	}

	var progress int64
	for mpv.scanner.Scan() {
		line := mpv.scanner.Text()
		if line == "" {
			continue
		}
		var response response
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			continue
		}
		switch response.Event {
		case "property-change":
			if response.Name == "time-pos" && response.Data != nil {
				if pos, ok := response.Data.(float64); ok {
					progress = int64(pos)
					client.ReportPlaybackProgress(item, progress)
				}
			}
		case "end-file":
			client.ReportPlaybackStopped(item, progress)
			return
		case "shutdown":
			client.ReportPlaybackStopped(item, progress)
			return
		}
	}
}
