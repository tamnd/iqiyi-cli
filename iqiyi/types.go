package iqiyi

import "strings"

// Show is the record emitted for each iQIYI title.
type Show struct {
	Rank       int     `json:"rank"`
	Name       string  `json:"name"`
	Channel    string  `json:"channel"`
	Categories string  `json:"categories"`
	Duration   string  `json:"duration"`
	Period     string  `json:"period"`
	Score      float64 `json:"score"`
	Focus      string  `json:"focus"`
	VideoCount int     `json:"video_count"`
	URL        string  `json:"url"`
}

// ChannelInfo describes a single channel.
type ChannelInfo struct {
	ID   int
	Name string
}

// ─── wire types from the pcw-api ─────────────────────────────────────────────

type wireResponse struct {
	Code string   `json:"code"`
	Data wireData `json:"data"`
}

type wireData struct {
	List    []wireShow `json:"list"`
	HasNext int        `json:"has_next"`
}

type wireShow struct {
	TvID        int64          `json:"tvId"`
	AlbumID     int64          `json:"albumId"`
	ChannelID   int            `json:"channelId"`
	Name        string         `json:"name"`
	PlayUrl     string         `json:"playUrl"`
	Description string         `json:"description"`
	ImageUrl    string         `json:"imageUrl"`
	Duration    string         `json:"duration"`
	VideoCount  int            `json:"videoCount"`
	LatestOrder string         `json:"latestOrder"`
	Categories  []wireCategory `json:"categories"`
	Period      string         `json:"period"`
	Score       float64        `json:"score"`
	Focus       string         `json:"focus"`
	People      wirePeople     `json:"people"`
}

type wireCategory struct {
	Name string `json:"name"`
}

type wirePeople struct {
	MainCharactor []wirePerson `json:"main_charactor"`
}

type wirePerson struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ─── conversion ──────────────────────────────────────────────────────────────

func wireToShow(w wireShow, rank int) Show {
	cats := make([]string, 0, len(w.Categories))
	for _, c := range w.Categories {
		if c.Name != "" {
			cats = append(cats, c.Name)
		}
	}
	ch := channelName(w.ChannelID)
	return Show{
		Rank:       rank,
		Name:       w.Name,
		Channel:    ch,
		Categories: strings.Join(cats, "/"),
		Duration:   w.Duration,
		Period:     w.Period,
		Score:      w.Score,
		Focus:      w.Focus,
		VideoCount: w.VideoCount,
		URL:        w.PlayUrl,
	}
}
