package fun

import (
	"context"
	"encoding/json"
	"errors"
	"html"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type sourceMeta struct {
	name     string
	category string
	hint     string
}

func (m sourceMeta) Name() string {
	return m.name
}

func (m sourceMeta) Category() string {
	return m.category
}

func (m sourceMeta) Hint() string {
	return m.hint
}

func (m sourceMeta) buildResult(title string, url string) RandomJumpResult {
	return RandomJumpResult{
		Title:     normalizeText(title),
		URL:       strings.TrimSpace(url),
		Category:  m.category,
		Hint:      m.hint,
		Source:    m.name,
		FetchedAt: time.Now().Format(time.RFC3339),
	}
}

type htmlHotRankSource struct {
	sourceMeta
	pageURL     string
	itemPattern *regexp.Regexp
	urlMapper   func(string) string
	itemFilter  func(string, string) bool
}

func (s htmlHotRankSource) Fetch(ctx context.Context, client *http.Client) ([]RandomJumpResult, error) {
	body, err := fetchBody(ctx, client, s.pageURL)
	if err != nil {
		return nil, err
	}

	matches := s.itemPattern.FindAllStringSubmatch(string(body), -1)
	if len(matches) == 0 {
		return nil, errors.New("热榜页面没有抓到内容")
	}

	results := make([]RandomJumpResult, 0, 12)
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		url := strings.TrimSpace(match[1])
		if s.urlMapper != nil {
			url = s.urlMapper(url)
		}

		title := normalizeText(match[2])
		if s.itemFilter != nil && !s.itemFilter(url, title) {
			continue
		}

		result := s.buildResult(title, url)
		if result.Title != "" && result.URL != "" {
			results = append(results, result)
		}
	}

	if len(results) == 0 {
		return nil, errors.New("热榜页面没有抓到可用内容")
	}

	return results, nil
}

type hackerNewsSource struct {
	sourceMeta
	topStoriesURL string
	itemURLPrefix string
}

func (s hackerNewsSource) Fetch(ctx context.Context, client *http.Client) ([]RandomJumpResult, error) {
	body, err := fetchBody(ctx, client, s.topStoriesURL)
	if err != nil {
		return nil, err
	}

	var ids []int64
	if err = json.Unmarshal(body, &ids); err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, errors.New("hn 源没有返回内容")
	}

	limit := len(ids)
	if limit > 18 {
		limit = 18
	}

	order, err := randomPermutation(limit, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	resultsCh := make(chan RandomJumpResult, limit)
	var wg sync.WaitGroup

	for _, pickedIndex := range order {
		itemID := ids[pickedIndex]
		wg.Add(1)
		go func() {
			defer wg.Done()

			itemBody, fetchErr := fetchBody(ctx, client, s.itemURLPrefix+int64ToString(itemID)+".json")
			if fetchErr != nil {
				return
			}

			var item hackerNewsItem
			if unmarshalErr := json.Unmarshal(itemBody, &item); unmarshalErr != nil {
				return
			}

			url := strings.TrimSpace(item.URL)
			if url == "" {
				url = "https://news.ycombinator.com/item?id=" + int64ToString(item.ID)
			}

			result := s.buildResult(item.Title, url)
			if result.Title == "" || result.URL == "" {
				return
			}

			select {
			case resultsCh <- result:
			case <-ctx.Done():
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	results := make([]RandomJumpResult, 0, 8)
	for result := range resultsCh {
		results = append(results, result)
		if len(results) >= 8 {
			cancel()
			break
		}
	}

	if len(results) == 0 {
		return nil, errors.New("hn 源没有抓到可用链接")
	}

	return results, nil
}

type bilibiliPopularSource struct {
	sourceMeta
	apiURL string
}

func (s bilibiliPopularSource) Fetch(ctx context.Context, client *http.Client) ([]RandomJumpResult, error) {
	body, err := fetchBody(ctx, client, s.apiURL)
	if err != nil {
		return nil, err
	}

	var payload bilibiliPopularResponse
	if err = json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	if payload.Code != 0 {
		return nil, errors.New("视频热榜返回了异常状态")
	}

	results := make([]RandomJumpResult, 0, len(payload.Data.List))
	for _, item := range payload.Data.List {
		url := strings.TrimSpace(item.RedirectURL)
		if url == "" && item.Bvid != "" {
			url = "https://www.bilibili.com/video/" + item.Bvid
		}

		result := s.buildResult(item.Title, url)
		if result.Title == "" || result.URL == "" {
			continue
		}
		results = append(results, result)
	}

	if len(results) == 0 {
		return nil, errors.New("视频源没有抓到可用内容")
	}

	return results, nil
}

func newHotRankCrawlerSources() []crawlerSource {
	return []crawlerSource{
		newQidianRankSource(),
		newBilibiliPopularSource(),
		newCnblogsTopDiggsSource(),
		newHackerNewsTopSource(),
	}
}

func newQidianRankSource() crawlerSource {
	return htmlHotRankSource{
		sourceMeta: sourceMeta{
			name:     "起点排行榜",
			category: "novel",
			hint:     "这次直接从小说热榜页抓书名和链接。",
		},
		pageURL:     "https://m.qidian.com/rank/",
		itemPattern: regexp.MustCompile(`(?s)<a[^>]+href="(//m\.qidian\.com/book/\d+/)"[^>]*>.*?<h2[^>]*>([^<]+)</h2>`),
		urlMapper: func(rawURL string) string {
			if strings.HasPrefix(rawURL, "//") {
				return "https:" + rawURL
			}
			return rawURL
		},
	}
}

func newBilibiliPopularSource() crawlerSource {
	return bilibiliPopularSource{
		sourceMeta: sourceMeta{
			name:     "B 站热门",
			category: "video",
			hint:     "视频链接来自 B 站热门榜接口。",
		},
		apiURL: "https://api.bilibili.com/x/web-interface/popular?pn=1&ps=20",
	}
}

func newCnblogsTopDiggsSource() crawlerSource {
	return htmlHotRankSource{
		sourceMeta: sourceMeta{
			name:     "博客园推荐排行",
			category: "blog",
			hint:     "博客内容直接来自博客园推荐排行页。",
		},
		pageURL:     "https://www.cnblogs.com/aggsite/topdiggs",
		itemPattern: regexp.MustCompile(`href="(https://www\.cnblogs\.com/[^"]+)"[^>]*target="_blank">([^<]+)</a>`),
		itemFilter: func(url string, title string) bool {
			return strings.Contains(url, "/p/") && title != "»"
		},
	}
}

func newHackerNewsTopSource() crawlerSource {
	return hackerNewsSource{
		sourceMeta: sourceMeta{
			name:     "Hacker News",
			category: "news",
			hint:     "新闻链接来自 Hacker News 热榜。",
		},
		topStoriesURL: "https://hacker-news.firebaseio.com/v0/topstories.json",
		itemURLPrefix: "https://hacker-news.firebaseio.com/v0/item/",
	}
}

func fetchBody(ctx context.Context, client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json, application/xml, text/xml, */*")
	req.Header.Set("User-Agent", "luangao-chaos-bot/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.New("远程源返回状态异常")
	}

	return io.ReadAll(resp.Body)
}

func normalizeText(value string) string {
	normalized := strings.TrimSpace(html.UnescapeString(value))
	normalized = strings.ReplaceAll(normalized, "\n", " ")
	normalized = strings.Join(strings.Fields(normalized), " ")
	return normalized
}

func int64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

type hackerNewsItem struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type bilibiliPopularResponse struct {
	Code int `json:"code"`
	Data struct {
		List []bilibiliPopularItem `json:"list"`
	} `json:"data"`
}

type bilibiliPopularItem struct {
	Title       string `json:"title"`
	Bvid        string `json:"bvid"`
	RedirectURL string `json:"redirect_url"`
}
