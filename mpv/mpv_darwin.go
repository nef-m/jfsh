package mpv

import (
	"os/exec"

	"github.com/hacel/jfsh/jellyfin"
)

func Play(client *jellyfin.Client, item jellyfin.Item) {
	c := exec.Command("mpv", getStreamingURL(client.Host, item))
	c.Run()
	// TODO: finish the rest of this function
}
