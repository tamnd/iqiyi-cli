package iqiyi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Video is the unified record for a movie, series, or other iQIYI content.
type Video struct {
	QipuID       string   `json:"qipu_id"`
	Title        string   `json:"title"`
	Type         string   `json:"type"`
	ChannelID    int      `json:"channel_id"`
	Categories   []string `json:"categories"`
	Countries    []string `json:"countries"`
	Language     string   `json:"language"`
	Year         int      `json:"year"`
	Director     string   `json:"director"`
	Actors       []string `json:"actors"`
	Description  string   `json:"description"`
	EpisodeCount int      `json:"episode_count"`
	Score        float64  `json:"score"`
	ScoreCount   int64    `json:"score_count"`
	PlayCount    int64    `json:"play_count"`
	LikeCount    int64    `json:"like_count"`
	IsVIP        bool     `json:"is_vip"`
	Is4K         bool     `json:"is_4k"`
	IsDolby      bool     `json:"is_dolby"`
	OnlineDate   string   `json:"online_date"`
	CoverURL     string   `json:"cover_url"`
	URL          string   `json:"url"`
}

// RankItem is one row in a ranking list.
type RankItem struct {
	Rank       int     `json:"rank"`
	QipuID     string  `json:"qipu_id"`
	Title      string  `json:"title"`
	Type       string  `json:"type"`
	Score      float64 `json:"score"`
	PlayCount  int64   `json:"play_count"`
	OnlineDate string  `json:"online_date"`
	Director   string  `json:"director"`
	Actors     string  `json:"actors"`
	IsVIP      bool    `json:"is_vip"`
	CoverURL   string  `json:"cover_url"`
	URL        string  `json:"url"`
}

// Episode is one episode record within a series.
type Episode struct {
	TvID          string `json:"tv_id"`
	AlbumID       string `json:"album_id"`
	EpisodeNumber int    `json:"episode_number"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	DurationSecs  int    `json:"duration_seconds"`
	PlayCount     int64  `json:"play_count"`
	PublishDate   string `json:"publish_date"`
	IsVIP         bool   `json:"is_vip"`
	CoverURL      string `json:"cover_url"`
	URL           string `json:"url"`
}

// ─── wire types for mesh.if.iqiyi.com hot list ───────────────────────────────

type wireHotResp struct {
	Code string      `json:"code"`
	Data wireHotData `json:"data"`
}

type wireHotData struct {
	Rows []wireHotRow `json:"rows"`
}

type wireHotRow struct {
	Cells []wireHotCell `json:"cells"`
}

type wireHotCell struct {
	Data wireHotItem `json:"data"`
}

type wireHotItem struct {
	QipuID      string `json:"qipu_id"`
	Name        string `json:"name"`
	PlayCount   string `json:"play_count"`
	PlayTip     string `json:"play_count_tip"`
	Score       string `json:"score"`
	TvPicURL    string `json:"tv_pic_url"`
	ChannelName string `json:"channel_name"`
	ChannelID   int    `json:"channel_id"`
	Description string `json:"description"`
	Directors   string `json:"directors"`
	Stars       string `json:"stars"`
	OnlineTime  string `json:"online_time"`
	IsVip       int    `json:"is_vip"`
}

// ─── wire types for pcw-api.iqiyi.com ranking ────────────────────────────────

type wireRankResp struct {
	Code json.RawMessage `json:"code"`
	Data wireRankData    `json:"data"`
}

type wireRankData struct {
	List []wireRankItem `json:"list"`
}

type wireRankItem struct {
	QipuID      string  `json:"qipuId"`
	Name        string  `json:"name"`
	Score       float64 `json:"score"`
	PlayCount   int64   `json:"playCount"`
	TvPicURL    string  `json:"tvPicUrl"`
	ChannelName string  `json:"channelName"`
	OnlineTime  string  `json:"onlineTime"`
	Director    string  `json:"director"`
	MainActors  string  `json:"mainActors"`
	IsVip       int     `json:"isVip"`
	AlbumID     string  `json:"albumId"`
	Country     string  `json:"country"`
	Language    string  `json:"language"`
	Year        int     `json:"year"`
}

// ─── wire types for search.video.iqiyi.com ────────────────────────────────────

type wireSearchResp struct {
	Status int            `json:"status"`
	Data   wireSearchData `json:"data"`
}

type wireSearchData struct {
	Docinfos []wireDocInfo `json:"docinfos"`
	Total    int           `json:"total"`
}

type wireDocInfo struct {
	AlbumInfo wireSearchAlbum `json:"albumInfo"`
}

type wireSearchAlbum struct {
	QipuID      string  `json:"qipu_id"`
	Name        string  `json:"name"`
	ChannelName string  `json:"channel_name"`
	ChannelID   int     `json:"channel_id"`
	PlayTip     string  `json:"play_count_tip"`
	Score       float64 `json:"score"`
	IsVip       int     `json:"is_vip"`
	Director    string  `json:"director"`
	MainActor   string  `json:"main_actor"`
	Description string  `json:"description"`
	TvPicURL    string  `json:"tv_pic_url"`
	OnlineTime  string  `json:"online_time"`
}

// ─── wire types for pcw-api.iqiyi.com album baseinfo ─────────────────────────

type wireAlbumResp struct {
	Code json.RawMessage `json:"code"`
	Data wireAlbumData   `json:"data"`
}

type wireAlbumData struct {
	AlbumInfo wireAlbumInfo `json:"albumInfo"`
}

type wireAlbumInfo struct {
	QipuID       string  `json:"qipuId"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	ChannelName  string  `json:"channelName"`
	ChannelID    int     `json:"channelId"`
	Categories   string  `json:"categories"`
	Year         int     `json:"year"`
	Country      string  `json:"country"`
	Language     string  `json:"language"`
	TvPicURL     string  `json:"tvPicUrl"`
	OnlineTime   string  `json:"onlineTime"`
	EpisodeCount int     `json:"episodeCount"`
	Score        float64 `json:"score"`
	ScoreCount   int64   `json:"scoreCount"`
	PlayCount    int64   `json:"playCount"`
	LikeCount    int64   `json:"likeCount"`
	IsVip        int     `json:"isVip"`
	Is4K         int     `json:"is4K"`
	IsDolby      int     `json:"isDolby"`
	Director     string  `json:"director"`
	MainActors   string  `json:"mainActors"`
}

// ─── wire types for pcw-api.iqiyi.com episode list ───────────────────────────

type wireEpListResp struct {
	Code json.RawMessage `json:"code"`
	Data wireEpListData  `json:"data"`
}

type wireEpListData struct {
	Alists []wireEpisode `json:"alists"`
}

type wireEpisode struct {
	TvID        string `json:"tvId"`
	Order       int    `json:"order"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	PlayCount   int64  `json:"playCount"`
	PublishDate string `json:"publishDate"`
	IsVip       int    `json:"isVip"`
	TvPicURL    string `json:"tvPicUrl"`
}

// ─── helpers ─────────────────────────────────────────────────────────────────

// isSuccessCode accepts both integer 0 and string "A00000" as success.
func isSuccessCode(raw json.RawMessage) bool {
	var n int
	if json.Unmarshal(raw, &n) == nil {
		return n == 0
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s == "A00000"
	}
	return false
}

// ParsePlayCount parses abbreviated Chinese play-count strings like "5.6亿", "1234万".
func ParsePlayCount(s string) int64 {
	return parsePlayCount(s)
}

// parsePlayCount is the internal implementation.
func parsePlayCount(s string) int64 {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "亿") {
		n, _ := strconv.ParseFloat(strings.TrimSuffix(s, "亿"), 64)
		return int64(n * 100_000_000)
	}
	if strings.HasSuffix(s, "万") {
		n, _ := strconv.ParseFloat(strings.TrimSuffix(s, "万"), 64)
		return int64(n * 10_000)
	}
	n, _ := strconv.ParseInt(strings.ReplaceAll(s, ",", ""), 10, 64)
	return n
}

func splitComma(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func videoURL(albumID string) string {
	if albumID == "" {
		return ""
	}
	return fmt.Sprintf("https://www.iqiyi.com/a_%s.html", albumID)
}

// ─── Hot ─────────────────────────────────────────────────────────────────────

// Hot fetches the current trending content from the iQIYI hot list.
func (c *Client) Hot(ctx context.Context) ([]Video, error) {
	base := c.hotURL
	if base == "" {
		base = "https://mesh.if.iqiyi.com"
	}
	rawURL := base + "/tvg/row?key=iqiyi_hot_list"
	body, err := c.getWithHeaders(ctx, rawURL)
	if err != nil {
		return nil, err
	}
	var resp wireHotResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode hot response: %w", err)
	}
	if resp.Code != "A00000" {
		return nil, fmt.Errorf("API error code %s", resp.Code)
	}
	var videos []Video
	for _, row := range resp.Data.Rows {
		for _, cell := range row.Cells {
			it := cell.Data
			if it.QipuID == "" {
				continue
			}
			score, _ := strconv.ParseFloat(it.Score, 64)
			pc := parsePlayCount(it.PlayCount)
			if pc == 0 {
				pc = parsePlayCount(it.PlayTip)
			}
			videos = append(videos, Video{
				QipuID:      it.QipuID,
				Title:       it.Name,
				Type:        it.ChannelName,
				ChannelID:   it.ChannelID,
				Director:    it.Directors,
				Actors:      splitComma(it.Stars),
				Description: it.Description,
				Score:       score,
				PlayCount:   pc,
				IsVIP:       it.IsVip != 0,
				OnlineDate:  it.OnlineTime,
				CoverURL:    it.TvPicURL,
				URL:         videoURL(it.QipuID),
			})
		}
	}
	return videos, nil
}

// ─── Rank ─────────────────────────────────────────────────────────────────────

// Rank fetches the ranking list filtered by period, dimension, and content type.
// period: day|week|month; dimension: hot|playcount|score; contentType: tvSeries|movie|variety|cartoon|documentary
func (c *Client) Rank(ctx context.Context, period, dimension, contentType string) ([]RankItem, error) {
	base := c.baseURL
	if base == "" {
		base = "https://pcw-api.iqiyi.com"
	}
	rawURL := fmt.Sprintf(
		"%s/ranking/ranklist.json?period=%s&dimension=%s&type=%s",
		base, period, dimension, contentType,
	)
	body, err := c.getWithHeaders(ctx, rawURL)
	if err != nil {
		return nil, err
	}
	var resp wireRankResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode rank response: %w", err)
	}
	if !isSuccessCode(resp.Code) {
		return nil, fmt.Errorf("API error: rank request failed")
	}
	items := make([]RankItem, 0, len(resp.Data.List))
	for i, w := range resp.Data.List {
		id := w.QipuID
		if id == "" {
			id = w.AlbumID
		}
		items = append(items, RankItem{
			Rank:       i + 1,
			QipuID:     id,
			Title:      w.Name,
			Type:       w.ChannelName,
			Score:      w.Score,
			PlayCount:  w.PlayCount,
			OnlineDate: w.OnlineTime,
			Director:   w.Director,
			Actors:     w.MainActors,
			IsVIP:      w.IsVip != 0,
			CoverURL:   w.TvPicURL,
			URL:        videoURL(id),
		})
	}
	return items, nil
}

// ─── Search ───────────────────────────────────────────────────────────────────

// Search queries iQIYI for the given term and returns up to limit results.
func (c *Client) Search(ctx context.Context, query string, limit int) ([]Video, error) {
	if limit <= 0 {
		limit = 20
	}
	base := c.searchURL
	if base == "" {
		base = "https://search.video.iqiyi.com"
	}
	rawURL := fmt.Sprintf(
		"%s/o?key=%s&pageNum=1&pageSize=%d&video=1&album=1&fwds=1",
		base, query, limit,
	)
	body, err := c.getWithHeaders(ctx, rawURL)
	if err != nil {
		return nil, err
	}
	var resp wireSearchResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode search response: %w", err)
	}
	if resp.Status != 1 {
		return nil, nil
	}
	videos := make([]Video, 0, len(resp.Data.Docinfos))
	for _, doc := range resp.Data.Docinfos {
		a := doc.AlbumInfo
		if a.QipuID == "" {
			continue
		}
		pc := parsePlayCount(a.PlayTip)
		videos = append(videos, Video{
			QipuID:    a.QipuID,
			Title:     a.Name,
			Type:      a.ChannelName,
			ChannelID: a.ChannelID,
			Director:  a.Director,
			Actors:    splitComma(a.MainActor),
			Description: a.Description,
			Score:       a.Score,
			PlayCount:   pc,
			IsVIP:       a.IsVip != 0,
			OnlineDate:  a.OnlineTime,
			CoverURL:    a.TvPicURL,
			URL:         videoURL(a.QipuID),
		})
	}
	return videos, nil
}

// ─── Show ─────────────────────────────────────────────────────────────────────

// Show fetches the album/series detail for the given album ID (qipu_id).
func (c *Client) Show(ctx context.Context, albumID string) (*Video, error) {
	base := c.baseURL
	if base == "" {
		base = "https://pcw-api.iqiyi.com"
	}
	rawURL := fmt.Sprintf(
		"%s/albums/album/baseinfo?album_id=%s&recommend_albums_number=5",
		base, albumID,
	)
	body, err := c.getWithHeaders(ctx, rawURL)
	if err != nil {
		return nil, err
	}
	var resp wireAlbumResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("decode show response: %w", err)
	}
	if !isSuccessCode(resp.Code) {
		return nil, fmt.Errorf("API error: show request failed for album %s", albumID)
	}
	info := resp.Data.AlbumInfo
	id := info.QipuID
	if id == "" {
		id = albumID
	}
	cats := splitComma(info.Categories)
	actors := splitComma(info.MainActors)
	var countries []string
	if info.Country != "" {
		countries = []string{info.Country}
	}
	v := &Video{
		QipuID:       id,
		Title:        info.Name,
		Type:         info.ChannelName,
		ChannelID:    info.ChannelID,
		Categories:   cats,
		Countries:    countries,
		Language:     info.Language,
		Year:         info.Year,
		Director:     info.Director,
		Actors:       actors,
		Description:  info.Description,
		EpisodeCount: info.EpisodeCount,
		Score:        info.Score,
		ScoreCount:   info.ScoreCount,
		PlayCount:    info.PlayCount,
		LikeCount:    info.LikeCount,
		IsVIP:        info.IsVip != 0,
		Is4K:         info.Is4K != 0,
		IsDolby:      info.IsDolby != 0,
		OnlineDate:   info.OnlineTime,
		CoverURL:     info.TvPicURL,
		URL:          videoURL(id),
	}
	return v, nil
}

// ─── Episodes ─────────────────────────────────────────────────────────────────

// Episodes fetches the episode list for the given album ID.
// It paginates automatically until all episodes are collected or limit is reached.
func (c *Client) Episodes(ctx context.Context, albumID string, limit int) ([]Episode, error) {
	if limit <= 0 {
		limit = 200
	}
	base := c.baseURL
	if base == "" {
		base = "https://pcw-api.iqiyi.com"
	}
	var all []Episode
	page := 1
	size := 50
	for {
		rawURL := fmt.Sprintf(
			"%s/albums/album/avlistinfo?album_id=%s&page=%d&size=%d",
			base, albumID, page, size,
		)
		body, err := c.getWithHeaders(ctx, rawURL)
		if err != nil {
			return nil, err
		}
		var resp wireEpListResp
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, fmt.Errorf("decode episodes response: %w", err)
		}
		if !isSuccessCode(resp.Code) {
			break
		}
		if len(resp.Data.Alists) == 0 {
			break
		}
		for _, w := range resp.Data.Alists {
			ep := Episode{
				TvID:          w.TvID,
				AlbumID:       albumID,
				EpisodeNumber: w.Order,
				Title:         w.Name,
				Description:   w.Description,
				DurationSecs:  w.Duration,
				PlayCount:     w.PlayCount,
				PublishDate:   w.PublishDate,
				IsVIP:         w.IsVip != 0,
				CoverURL:      w.TvPicURL,
				URL:           fmt.Sprintf("https://www.iqiyi.com/v_%s.html", w.TvID),
			}
			all = append(all, ep)
			if len(all) >= limit {
				return all, nil
			}
		}
		page++
	}
	return all, nil
}

// ─── getWithHeaders ───────────────────────────────────────────────────────────

// getWithHeaders fetches rawURL adding iQIYI-required headers.
func (c *Client) getWithHeaders(ctx context.Context, rawURL string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= c.retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff(attempt)):
			}
		}
		body, retry, err := c.doWithHeaders(ctx, rawURL)
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

func (c *Client) doWithHeaders(ctx context.Context, rawURL string) (body []byte, retry bool, err error) {
	c.pace()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, false, err
	}
	ua := c.userAgent
	if ua == "" {
		ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	}
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Referer", "https://www.iqiyi.com/")
	req.Header.Set("Accept", "application/json")

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
