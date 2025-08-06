package jellyfin

import (
	"context"
	"time"

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
	res, _, err := c.api.TvShowsAPI.GetNextUp(context.Background()).Execute()
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

func (c *Client) GetEpisodes(seriesID string) ([]Item, error) {
	res, _, err := c.api.ItemsAPI.GetItems(context.Background()).
		Recursive(true).
		ParentId(seriesID).
		IncludeItemTypes([]api.BaseItemKind{api.BASEITEMKIND_EPISODE}).
		SortBy([]api.ItemSortBy{api.ITEMSORTBY_PARENT_INDEX_NUMBER, api.ITEMSORTBY_INDEX_NUMBER}).
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

func (c *Client) ReportPlaybackStopped(item Item, pos int64) {
	posTicks := pos * 10000000
	if _, err := c.api.PlaystateAPI.ReportPlaybackStopped(context.Background()).PlaybackStopInfo(api.PlaybackStopInfo{
		ItemId:        item.Id,
		PositionTicks: *api.NewNullableInt64(&posTicks),
	}).Execute(); err != nil {
		panic(err)
	}
}

func (c *Client) ReportPlaybackProgress(item Item, pos int64) {
	if time.Since(c.lastProgressReport) < time.Second*3 { // debounce
		return
	}
	posTicks := pos * 10000000
	if _, err := c.api.PlaystateAPI.ReportPlaybackProgress(context.Background()).PlaybackProgressInfo(api.PlaybackProgressInfo{
		ItemId:        item.Id,
		PositionTicks: *api.NewNullableInt64(&posTicks),
	}).Execute(); err != nil {
		panic(err)
	}
}
