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
)

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

func (c *mpv) close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			slog.Error("failed to close mpv socket connection", "err", err)
		}
	}
	if c.cmd != nil && c.cmd.Process != nil {
		if err := c.cmd.Process.Signal(os.Interrupt); err != nil {
			slog.Error("failed to interrupt mpv process", "err", err)
		}
		if err := c.cmd.Wait(); err != nil {
			slog.Error("failed to wait for mpv process", "err", err)
		}
	}
	if err := os.Remove(c.socket); err != nil {
		slog.Error("failed to remove mpv socket", "err", err)
	}
}

func (c *mpv) send(command []any) error {
	req := request{Command: command}
	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal mpv command: %w", err)
	}
	_, err = c.conn.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("failed to write mpv command to socket: %w", err)
	}
	return nil
}

func (c *mpv) observeProperty(name string) error {
	return c.send([]any{"observe_property", 1, name})
}

func (c *mpv) seekTo(pos float64) error {
	return c.send([]any{"seek", pos, "absolute"})
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
