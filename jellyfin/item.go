package jellyfin

import (
	"fmt"
	"strings"

	"github.com/sj14/jellyfin-go/api"
)

// Item is a type alias just because it looks nicer
type Item = api.BaseItemDto

func getItemRuntime(ticks int64) string {
	minutes := ticks / 600_000_000
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}
	hours := minutes / 60
	minutes -= hours * 60
	return fmt.Sprintf("%dh%dm", hours, minutes)
}

// Helpers to try to contain BaseItemDto implementation from leaking out of the package without casting it to something else
// NOTE: just cast it to something else?

func GetResumePosition(item Item) (secs int64) {
	if item.UserData.IsSet() {
		secs = *item.UserData.Get().PlaybackPositionTicks / 10000000
	}
	return
}

func GetStreamingURL(host string, item Item) string {
	url := fmt.Sprintf("%s/videos/%s/stream?static=true", host, *item.Id)
	return fmt.Sprintf("edl://%%%d%%%s", len(url), url)
}

func GetMediaTitle(item Item) string {
	title := item.GetPath()
	switch item.GetType() {
	case api.BASEITEMKIND_MOVIE:
		title = fmt.Sprintf("%s (%d)", item.GetName(), item.GetProductionYear())
	case api.BASEITEMKIND_EPISODE:
		title = fmt.Sprintf("%s - S%d:E%d - %s (%d)", item.GetSeriesName(), item.GetParentIndexNumber(), item.GetIndexNumber(), item.GetName(), item.GetProductionYear())
	}
	return title
}

func GetItemTitle(i Item) string {
	str := &strings.Builder{}
	switch *i.Type {
	case api.BASEITEMKIND_MOVIE:
		fmt.Fprintf(str, "%s (%d)", i.GetName(), i.GetProductionYear())
		if i.UserData.IsSet() && i.UserData.Get().PlayedPercentage.IsSet() {
			fmt.Fprintf(str, " [%.f%%]", *i.GetUserData().PlayedPercentage.Get())
		}
	case api.BASEITEMKIND_EPISODE:
		fmt.Fprintf(str, "%s S%.2dE%.2d (%d)", i.GetSeriesName(), i.GetParentIndexNumber(), i.GetIndexNumber(), i.GetProductionYear())
		if i.UserData.IsSet() && i.UserData.Get().PlayedPercentage.IsSet() {
			fmt.Fprintf(str, " [%.f%%]", *i.GetUserData().PlayedPercentage.Get())
		}
	case api.BASEITEMKIND_SERIES:
		fmt.Fprintf(str, "%s (%d)", i.GetName(), i.GetProductionYear())
	case api.BASEITEMKIND_VIDEO:
		fmt.Fprintf(str, "%s (%d)", i.GetName(), i.GetProductionYear())
	}
	return str.String()
}

func GetItemDescription(i Item) string {
	str := &strings.Builder{}
	switch *i.Type {
	case api.BASEITEMKIND_MOVIE:
		fmt.Fprintf(str, "Movie  | Rating: %.1f | Runtime: %s", i.GetCommunityRating(), getItemRuntime(i.GetRunTimeTicks()))
	case api.BASEITEMKIND_SERIES:
		fmt.Fprintf(str, "Series | Rating: %.1f", i.GetCommunityRating())
	case api.BASEITEMKIND_EPISODE:
		fmt.Fprintf(str, "%s", i.GetName())
	case api.BASEITEMKIND_VIDEO:
		fmt.Fprintf(str, "Video  | Rating: %.1f | Runtime: %s", i.GetCommunityRating(), getItemRuntime(i.GetRunTimeTicks()))
	}
	return str.String()
}

func IsSeries(i Item) bool {
	return *i.Type == api.BASEITEMKIND_SERIES
}

func Watched(i Item) bool {
	if i.UserData.IsSet() {
		if i.UserData.Get().GetPlayed() {
			return true
		}
	}
	return false
}
