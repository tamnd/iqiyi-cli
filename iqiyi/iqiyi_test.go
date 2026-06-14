package iqiyi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tamnd/iqiyi-cli/iqiyi"
)

func mockResponse(list []map[string]any) []byte {
	resp := map[string]any{
		"code": "A00000",
		"data": map[string]any{
			"list":     list,
			"has_next": 0,
		},
	}
	b, _ := json.Marshal(resp)
	return b
}

func TestList(t *testing.T) {
	items := []map[string]any{
		{
			"tvId":        int64(12345),
			"albumId":     int64(12345),
			"channelId":   1,
			"name":        "Test Movie",
			"playUrl":     "https://www.iqiyi.com/v_test.html",
			"description": "A test movie",
			"imageUrl":    "https://example.com/img.jpg",
			"duration":    "120min",
			"videoCount":  1,
			"latestOrder": "",
			"categories":  []map[string]any{{"name": "Action"}, {"name": "Drama"}},
			"period":      "2024",
			"score":       8.5,
			"focus":       "Test Focus",
			"people": map[string]any{
				"main_charactor": []map[string]any{
					{"id": 1, "name": "Actor One"},
				},
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") == "" {
			t.Error("request carried no User-Agent")
		}
		if r.URL.Path != "/search/recommend/list" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(mockResponse(items))
	}))
	defer srv.Close()

	cfg := iqiyi.DefaultConfig()
	cfg.BaseURL = srv.URL
	cfg.Rate = 0

	c := iqiyi.NewClient(cfg)
	shows, err := c.List(context.Background(), 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(shows) != 1 {
		t.Fatalf("got %d shows, want 1", len(shows))
	}
	s := shows[0]
	if s.Name != "Test Movie" {
		t.Errorf("Name = %q, want %q", s.Name, "Test Movie")
	}
	if s.Rank != 1 {
		t.Errorf("Rank = %d, want 1", s.Rank)
	}
	if s.Channel != "movie" {
		t.Errorf("Channel = %q, want %q", s.Channel, "movie")
	}
	if s.Categories != "Action/Drama" {
		t.Errorf("Categories = %q, want %q", s.Categories, "Action/Drama")
	}
	if s.Score != 8.5 {
		t.Errorf("Score = %v, want 8.5", s.Score)
	}
	if s.URL != "https://www.iqiyi.com/v_test.html" {
		t.Errorf("URL = %q", s.URL)
	}
}

func TestChannels(t *testing.T) {
	channels := iqiyi.Channels()
	if len(channels) == 0 {
		t.Fatal("Channels() returned empty list")
	}
	for _, ch := range channels {
		if ch.ID == 0 {
			t.Errorf("channel %q has zero ID", ch.Name)
		}
		if ch.Name == "" {
			t.Errorf("channel ID %d has empty name", ch.ID)
		}
	}
}

func TestChannelID(t *testing.T) {
	id, ok := iqiyi.ChannelID("movie")
	if !ok {
		t.Fatal("ChannelID(movie) not found")
	}
	if id != 1 {
		t.Errorf("ChannelID(movie) = %d, want 1", id)
	}

	_, ok = iqiyi.ChannelID("unknown_channel")
	if ok {
		t.Error("ChannelID(unknown_channel) should return false")
	}
}
