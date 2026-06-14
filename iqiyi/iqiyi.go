// Package iqiyi is the library behind the iqiyi command line:
// the HTTP client, request shaping, and the typed data models for iQIYI.
//
// The Client here is the spine every command shares. It sets a real
// User-Agent, paces requests so a busy session stays polite, and retries the
// transient failures (429 and 5xx) that any public site throws under load.
package iqiyi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// DefaultUserAgent identifies the client to iQIYI.
const DefaultUserAgent = "iqiyi-cli/dev (+https://github.com/tamnd/iqiyi-cli)"

// channelIDs maps channel name to channel_id.
var channelIDs = map[string]int{
	"movie":       1,
	"tv_drama":    2,
	"animation":   3,
	"kids":        4,
	"documentary": 6,
	"variety":     15,
}

// channelNames maps channel_id to channel name.
var channelNames = map[int]string{
	1:  "movie",
	2:  "tv_drama",
	3:  "animation",
	4:  "kids",
	6:  "documentary",
	15: "variety",
}

func channelName(id int) string {
	if n, ok := channelNames[id]; ok {
		return n
	}
	return fmt.Sprintf("channel_%d", id)
}

// ChannelID returns the numeric channel_id for a named channel.
// It returns 0 and false if the name is unknown.
func ChannelID(name string) (int, bool) {
	id, ok := channelIDs[name]
	return id, ok
}

// Channels returns the list of known channels in a stable order.
func Channels() []ChannelInfo {
	return []ChannelInfo{
		{ID: 1, Name: "movie"},
		{ID: 2, Name: "tv_drama"},
		{ID: 3, Name: "animation"},
		{ID: 4, Name: "kids"},
		{ID: 6, Name: "documentary"},
		{ID: 15, Name: "variety"},
	}
}

// Config holds the parameters for constructing a Client.
type Config struct {
	BaseURL   string
	UserAgent string
	Rate      time.Duration
	Retries   int
	Timeout   time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		BaseURL:   "https://pcw-api.iqiyi.com",
		UserAgent: DefaultUserAgent,
		Rate:      200 * time.Millisecond,
		Retries:   3,
		Timeout:   30 * time.Second,
	}
}

// Client talks to the iQIYI pcw-api over HTTP.
type Client struct {
	http      *http.Client
	userAgent string
	rate      time.Duration
	retries   int
	baseURL   string

	last time.Time
}

// NewClient returns a Client configured from cfg.
func NewClient(cfg Config) *Client {
	return &Client{
		http:      &http.Client{Timeout: cfg.Timeout},
		userAgent: cfg.UserAgent,
		rate:      cfg.Rate,
		retries:   cfg.Retries,
		baseURL:   cfg.BaseURL,
	}
}

// List fetches the recommend list for channelID and returns up to limit shows.
func (c *Client) List(ctx context.Context, channelID int, limit int) ([]Show, error) {
	if limit <= 0 {
		limit = 20
	}
	u, err := url.Parse(c.baseURL + "/search/recommend/list")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("channel_id", fmt.Sprintf("%d", channelID))
	q.Set("data_type", "1")
	q.Set("mode", "24")
	q.Set("page_id", "")
	q.Set("ret_num", fmt.Sprintf("%d", limit))
	q.Set("session", "")
	u.RawQuery = q.Encode()

	body, err := c.get(ctx, u.String())
	if err != nil {
		return nil, err
	}

	var resp wireResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if resp.Code != "A00000" {
		return nil, fmt.Errorf("API error code %s", resp.Code)
	}

	shows := make([]Show, 0, len(resp.Data.List))
	for i, w := range resp.Data.List {
		shows = append(shows, wireToShow(w, i+1))
	}
	return shows, nil
}

// get fetches url with pacing and retry.
func (c *Client) get(ctx context.Context, rawURL string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= c.retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff(attempt)):
			}
		}
		body, retry, err := c.do(ctx, rawURL)
		if err == nil {
			return body, nil
		}
		lastErr = err
		if !retry {
			return nil, err
		}
	}
	return nil, fmt.Errorf("get %s: %w", rawURL, lastErr)
}

func (c *Client) do(ctx context.Context, rawURL string) (body []byte, retry bool, err error) {
	c.pace()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, true, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
		return nil, true, fmt.Errorf("http %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("http %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, true, err
	}
	return b, false, nil
}

// pace blocks until at least Rate has passed since the previous request.
func (c *Client) pace() {
	if c.rate <= 0 {
		return
	}
	if wait := c.rate - time.Since(c.last); wait > 0 {
		time.Sleep(wait)
	}
	c.last = time.Now()
}

func backoff(attempt int) time.Duration {
	d := time.Duration(attempt) * 500 * time.Millisecond
	if d > 5*time.Second {
		d = 5 * time.Second
	}
	return d
}
