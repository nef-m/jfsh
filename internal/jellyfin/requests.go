package jellyfin

import (
	"context"

	"github.com/sj14/jellyfin-go/api"
)

func (c *Client) GetResume() ([]Item, error) {
	res, _, err := c.api.ItemsAPI.GetResumeItems(context.Background()).UserId(c.UserID).Execute()
	if err != nil {
		return nil, err
	}
	return res.Items, nil
}

func (c *Client) GetNextUp() ([]Item, error) {
	res, _, err := c.api.TvShowsAPI.GetNextUp(context.Background()).
		EnableTotalRecordCount(false).
		DisableFirstEpisode(false).
		EnableResumable(false).
		EnableRewatching(false).
		Execute()
	if err != nil {
		return nil, err
	}
	return res.Items, nil
}

func (c *Client) GetRecentlyAdded() ([]Item, error) {
	res, _, err := c.api.ItemsAPI.GetItems(context.Background()).
		Recursive(true).
		IncludeItemTypes([]api.BaseItemKind{api.BASEITEMKIND_MOVIE, api.BASEITEMKIND_SERIES}).
		Limit(100).
		SortBy([]api.ItemSortBy{api.ITEMSORTBY_DATE_CREATED}).
		SortOrder([]api.SortOrder{api.SORTORDER_DESCENDING}).
		Execute()
	if err != nil {
		return nil, err
	}
	return res.Items, nil
}

func (c *Client) GetEpisodes(item Item) ([]Item, error) {
	seriesID := item.GetSeriesId()
	if item.GetType() == api.BASEITEMKIND_SERIES {
		seriesID = item.GetId()
	}
	res, _, err := c.api.TvShowsAPI.GetEpisodes(context.Background(), seriesID).
		Execute()
	if err != nil {
		return nil, err
	}
	return res.Items, nil
}

func (c *Client) Search(query string) ([]Item, error) {
	res, _, err := c.api.ItemsAPI.GetItems(context.Background()).
		SearchTerm(query).
		Recursive(true).
		IncludeItemTypes([]api.BaseItemKind{api.BASEITEMKIND_MOVIE, api.BASEITEMKIND_SERIES}).
		Limit(100).
		Execute()
	if err != nil {
		return nil, err
	}
	return res.Items, nil
}

func (c *Client) ReportPlaybackStart(item Item, ticks int64) error {
	_, err := c.api.PlaystateAPI.ReportPlaybackStart(context.Background()).PlaybackStartInfo(api.PlaybackStartInfo{
		ItemId:        item.Id,
		PositionTicks: *api.NewNullableInt64(&ticks),
	}).Execute()
	return err
}

func (c *Client) ReportPlaybackStopped(item Item, ticks int64) error {
	_, err := c.api.PlaystateAPI.ReportPlaybackStopped(context.Background()).PlaybackStopInfo(api.PlaybackStopInfo{
		ItemId:        item.Id,
		PositionTicks: *api.NewNullableInt64(&ticks),
	}).Execute()
	return err
}

func (c *Client) ReportPlaybackProgress(item Item, ticks int64) error {
	_, err := c.api.PlaystateAPI.ReportPlaybackProgress(context.Background()).PlaybackProgressInfo(api.PlaybackProgressInfo{
		ItemId:        item.Id,
		PositionTicks: *api.NewNullableInt64(&ticks),
	}).Execute()
	return err
}
