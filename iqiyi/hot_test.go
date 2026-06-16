package iqiyi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tamnd/iqiyi-cli/iqiyi"
)

// marshalJSON encodes v to JSON, panicking on error (test helper only).
func marshalJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

func cfgForServer(srv *httptest.Server) iqiyi.Config {
	cfg := iqiyi.DefaultConfig()
	cfg.BaseURL = srv.URL
	cfg.HotURL = srv.URL
	cfg.SearchURL = srv.URL
	cfg.Rate = 0
	return cfg
}

// ─── Hot ─────────────────────────────────────────────────────────────────────

func TestHot(t *testing.T) {
	payload := map[string]any{
		"code": "A00000",
		"data": map[string]any{
			"rows": []map[string]any{
				{
					"cells": []map[string]any{
						{
							"data": map[string]any{
								"qipu_id":      "200247736",
								"name":         "花千骨",
								"play_count":   "5.6亿",
								"score":        "8.7",
								"tv_pic_url":   "https://example.com/img.jpg",
								"channel_name": "电视剧",
								"channel_id":   2,
								"is_vip":       0,
								"directors":    "刘国辉",
								"stars":        "赵丽颖,蒋欣",
							},
						},
					},
				},
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(marshalJSON(payload))
	}))
	defer srv.Close()

	c := iqiyi.NewClient(cfgForServer(srv))
	videos, err := c.Hot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(videos) != 1 {
		t.Fatalf("got %d videos, want 1", len(videos))
	}
	v := videos[0]
	if v.QipuID != "200247736" {
		t.Errorf("QipuID = %q, want %q", v.QipuID, "200247736")
	}
	if v.Title != "花千骨" {
		t.Errorf("Title = %q", v.Title)
	}
	if v.PlayCount != 560_000_000 {
		t.Errorf("PlayCount = %d, want 560000000", v.PlayCount)
	}
	if v.Score != 8.7 {
		t.Errorf("Score = %v, want 8.7", v.Score)
	}
	if v.IsVIP {
		t.Error("IsVIP should be false")
	}
	if v.URL != "https://www.iqiyi.com/a_200247736.html" {
		t.Errorf("URL = %q", v.URL)
	}
	if len(v.Actors) != 2 {
		t.Errorf("Actors = %v, want 2 entries", v.Actors)
	}
}

// ─── Rank ─────────────────────────────────────────────────────────────────────

func TestRank(t *testing.T) {
	payload := map[string]any{
		"code": 0,
		"data": map[string]any{
			"list": []map[string]any{
				{
					"qipuId":      "100001",
					"name":        "庆余年 第二季",
					"score":       9.1,
					"playCount":   int64(800_000_000),
					"tvPicUrl":    "https://example.com/img2.jpg",
					"channelName": "电视剧",
					"onlineTime":  "2024-05-16",
					"director":    "孙皓",
					"mainActors":  "张若昀,李沁",
					"isVip":       0,
					"albumId":     "100001",
				},
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ranking/ranklist.json" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(marshalJSON(payload))
	}))
	defer srv.Close()

	c := iqiyi.NewClient(cfgForServer(srv))
	items, err := c.Rank(context.Background(), "day", "hot", "tvSeries")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("got %d items, want 1", len(items))
	}
	it := items[0]
	if it.Rank != 1 {
		t.Errorf("Rank = %d, want 1", it.Rank)
	}
	if it.QipuID != "100001" {
		t.Errorf("QipuID = %q", it.QipuID)
	}
	if it.Title != "庆余年 第二季" {
		t.Errorf("Title = %q", it.Title)
	}
	if it.Score != 9.1 {
		t.Errorf("Score = %v, want 9.1", it.Score)
	}
	if it.PlayCount != 800_000_000 {
		t.Errorf("PlayCount = %d", it.PlayCount)
	}
	if it.URL != "https://www.iqiyi.com/a_100001.html" {
		t.Errorf("URL = %q", it.URL)
	}
}

// ─── Search ───────────────────────────────────────────────────────────────────

func TestSearch(t *testing.T) {
	payload := map[string]any{
		"status": 1,
		"data": map[string]any{
			"docinfos": []map[string]any{
				{
					"albumInfo": map[string]any{
						"qipu_id":        "200001",
						"name":           "庆余年",
						"channel_name":   "电视剧",
						"channel_id":     2,
						"play_count_tip": "3.8亿",
						"score":          9.1,
						"is_vip":         0,
						"director":       "孙皓",
						"main_actor":     "张若昀",
						"description":    "权谋剧",
						"tv_pic_url":     "https://example.com/qyn.jpg",
						"online_time":    "2019-11-26",
					},
				},
			},
			"total": 1,
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(marshalJSON(payload))
	}))
	defer srv.Close()

	c := iqiyi.NewClient(cfgForServer(srv))
	videos, err := c.Search(context.Background(), "庆余年", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(videos) != 1 {
		t.Fatalf("got %d videos, want 1", len(videos))
	}
	v := videos[0]
	if v.QipuID != "200001" {
		t.Errorf("QipuID = %q", v.QipuID)
	}
	if v.Title != "庆余年" {
		t.Errorf("Title = %q", v.Title)
	}
	if v.PlayCount != 380_000_000 {
		t.Errorf("PlayCount = %d, want 380000000", v.PlayCount)
	}
}

// ─── Show ─────────────────────────────────────────────────────────────────────

func TestShow(t *testing.T) {
	payload := map[string]any{
		"code": 0,
		"data": map[string]any{
			"albumInfo": map[string]any{
				"qipuId":       "200247736",
				"name":         "花千骨",
				"description":  "古装仙侠剧",
				"channelName":  "电视剧",
				"channelId":    2,
				"categories":   "仙侠,古装,爱情",
				"year":         2015,
				"country":      "中国大陆",
				"language":     "普通话",
				"tvPicUrl":     "https://example.com/hqg.jpg",
				"onlineTime":   "2015-06-26",
				"episodeCount": 58,
				"score":        8.7,
				"scoreCount":   int64(450_000),
				"playCount":    int64(560_000_000),
				"likeCount":    int64(230_000),
				"isVip":        0,
				"is4K":         0,
				"isDolby":      0,
				"director":     "刘国辉",
				"mainActors":   "赵丽颖,蒋欣,陈学冬",
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/albums/album/baseinfo" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(marshalJSON(payload))
	}))
	defer srv.Close()

	c := iqiyi.NewClient(cfgForServer(srv))
	video, err := c.Show(context.Background(), "200247736")
	if err != nil {
		t.Fatal(err)
	}
	if video == nil {
		t.Fatal("Show returned nil")
	}
	if video.Title != "花千骨" {
		t.Errorf("Title = %q", video.Title)
	}
	if video.EpisodeCount != 58 {
		t.Errorf("EpisodeCount = %d, want 58", video.EpisodeCount)
	}
	if len(video.Actors) != 3 {
		t.Errorf("Actors = %v, want 3", video.Actors)
	}
	if len(video.Categories) != 3 {
		t.Errorf("Categories = %v, want 3", video.Categories)
	}
	if video.Score != 8.7 {
		t.Errorf("Score = %v, want 8.7", video.Score)
	}
	if video.PlayCount != 560_000_000 {
		t.Errorf("PlayCount = %d", video.PlayCount)
	}
	if video.URL != "https://www.iqiyi.com/a_200247736.html" {
		t.Errorf("URL = %q", video.URL)
	}
}

// ─── Episodes ─────────────────────────────────────────────────────────────────

func TestEpisodes(t *testing.T) {
	payload := map[string]any{
		"code": 0,
		"data": map[string]any{
			"alists": []map[string]any{
				{
					"tvId":        "9876543210",
					"order":       1,
					"name":        "第1集",
					"description": "序章",
					"duration":    2520,
					"playCount":   int64(12_000_000),
					"publishDate": "2015-06-26",
					"isVip":       0,
					"tvPicUrl":    "https://example.com/ep1.jpg",
				},
				{
					"tvId":        "9876543211",
					"order":       2,
					"name":        "第2集",
					"description": "冒险开始",
					"duration":    2640,
					"playCount":   int64(10_000_000),
					"publishDate": "2015-06-27",
					"isVip":       0,
					"tvPicUrl":    "https://example.com/ep2.jpg",
				},
			},
		},
	}

	pageCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/albums/album/avlistinfo" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		pageCount++
		// Return empty on second page to stop pagination.
		if r.URL.Query().Get("page") == "2" {
			empty := map[string]any{
				"code": 0,
				"data": map[string]any{"alists": []map[string]any{}},
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(marshalJSON(empty))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(marshalJSON(payload))
	}))
	defer srv.Close()

	c := iqiyi.NewClient(cfgForServer(srv))
	episodes, err := c.Episodes(context.Background(), "200247736", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(episodes) != 2 {
		t.Fatalf("got %d episodes, want 2", len(episodes))
	}
	ep := episodes[0]
	if ep.TvID != "9876543210" {
		t.Errorf("TvID = %q", ep.TvID)
	}
	if ep.EpisodeNumber != 1 {
		t.Errorf("EpisodeNumber = %d, want 1", ep.EpisodeNumber)
	}
	if ep.DurationSecs != 2520 {
		t.Errorf("DurationSecs = %d, want 2520", ep.DurationSecs)
	}
	if ep.AlbumID != "200247736" {
		t.Errorf("AlbumID = %q", ep.AlbumID)
	}
	if ep.URL != "https://www.iqiyi.com/v_9876543210.html" {
		t.Errorf("URL = %q", ep.URL)
	}
}

// ─── parsePlayCount ───────────────────────────────────────────────────────────

func TestParsePlayCount(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"5.6亿", 560_000_000},
		{"1.2亿", 120_000_000},
		{"1234万", 12_340_000},
		{"8000", 8000},
		{"", 0},
	}
	for _, tc := range tests {
		got := iqiyi.ParsePlayCount(tc.input)
		if got != tc.want {
			t.Errorf("ParsePlayCount(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}
}
