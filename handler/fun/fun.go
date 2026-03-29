package fun

import (
	"context"
	cryptorand "crypto/rand"
	"errors"
	"io"
	"math/big"
	"net/http"
	"sync"
	"time"
)

type RandomJumpResult struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Category  string `json:"category"`
	Hint      string `json:"hint"`
	Source    string `json:"source"`
	FetchedAt string `json:"fetchedAt"`
}

type RandomJumpFinder interface {
	Pick(ctx context.Context) (*RandomJumpResult, error)
}

type crawlerSource interface {
	Name() string
	Category() string
	Hint() string
	Fetch(ctx context.Context, client *http.Client) ([]RandomJumpResult, error)
}

type RandomJumpHandler struct {
	client   *http.Client
	sources  []crawlerSource
	cacheTTL time.Duration

	mu          sync.RWMutex
	cachedItems []RandomJumpResult
	lastCrawled time.Time
	refreshing  bool
}

func NewRandomJumpHandler() *RandomJumpHandler {
	return NewRandomJumpHandlerWithSources(newDefaultHTTPClient(), newHotRankCrawlerSources(), 8*time.Minute)
}

func NewRandomJumpHandlerWithSources(client *http.Client, sources []crawlerSource, cacheTTL time.Duration) *RandomJumpHandler {
	if client == nil {
		client = newDefaultHTTPClient()
	}
	if cacheTTL <= 0 {
		cacheTTL = 8 * time.Minute
	}

	return &RandomJumpHandler{
		client:   client,
		sources:  sources,
		cacheTTL: cacheTTL,
	}
}

func (h *RandomJumpHandler) Pick(ctx context.Context) (*RandomJumpResult, error) {
	items, err := h.loadItems(ctx)
	if err != nil {
		return nil, err
	}

	index, err := randomIndex(len(items), cryptorand.Reader)
	if err != nil {
		return nil, err
	}

	result := items[index]
	return &result, nil
}

func (h *RandomJumpHandler) loadItems(ctx context.Context) ([]RandomJumpResult, error) {
	cachedItems, isFresh := h.getCachedItems()
	if isFresh {
		return cachedItems, nil
	}

	if len(cachedItems) > 0 {
		h.refreshCacheAsync()
		return cachedItems, nil
	}

	items, err := h.crawlSources(ctx)
	if err == nil && len(items) > 0 {
		h.setCachedItems(items)
		return items, nil
	}

	if len(cachedItems) > 0 {
		return cachedItems, nil
	}

	if err != nil {
		return nil, err
	}

	return nil, errors.New("随机爬虫暂时没有抓到内容")
}

func (h *RandomJumpHandler) crawlSources(ctx context.Context) ([]RandomJumpResult, error) {
	if len(h.sources) == 0 {
		return nil, errors.New("没有可用的爬虫源")
	}

	order, err := randomPermutation(len(h.sources), cryptorand.Reader)
	if err != nil {
		return nil, err
	}

	targetSources, err := randomSourceBatchSize(len(order), cryptorand.Reader)
	if err != nil {
		return nil, err
	}

	type crawlResult struct {
		items []RandomJumpResult
		err   error
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	resultsCh := make(chan crawlResult, len(order))

	for _, sourceIndex := range order {
		source := h.sources[sourceIndex]
		go func() {
			items, fetchErr := source.Fetch(ctx, h.client)
			resultsCh <- crawlResult{
				items: items,
				err:   fetchErr,
			}
		}()
	}

	results := make([]RandomJumpResult, 0, 24)
	successCount := 0
	var firstErr error

	for range order {
		result := <-resultsCh
		if result.err != nil {
			if firstErr == nil {
				firstErr = result.err
			}
			continue
		}

		results = append(results, result.items...)
		successCount++

		if successCount >= targetSources && hasEnoughUniqueResults(results, 6) {
			cancel()
			break
		}
	}

	results = uniqueResults(results)
	if len(results) > 0 {
		return results, nil
	}

	if firstErr != nil {
		return nil, firstErr
	}

	return nil, errors.New("随机爬虫暂时没有抓到内容")
}

func (h *RandomJumpHandler) getCachedItems() ([]RandomJumpResult, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.cachedItems) == 0 {
		return nil, false
	}

	items := append([]RandomJumpResult(nil), h.cachedItems...)
	isFresh := time.Since(h.lastCrawled) < h.cacheTTL
	return items, isFresh
}

func (h *RandomJumpHandler) setCachedItems(items []RandomJumpResult) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.cachedItems = append([]RandomJumpResult(nil), items...)
	h.lastCrawled = time.Now()
}

func newDefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 4 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        32,
			MaxIdleConnsPerHost: 8,
			IdleConnTimeout:     90 * time.Second,
			ForceAttemptHTTP2:   true,
		},
	}
}

func (h *RandomJumpHandler) refreshCacheAsync() {
	if !h.tryStartRefresh() {
		return
	}

	go func() {
		defer h.finishRefresh()

		refreshCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		items, err := h.crawlSources(refreshCtx)
		if err != nil || len(items) == 0 {
			return
		}

		h.setCachedItems(items)
	}()
}

func (h *RandomJumpHandler) tryStartRefresh() bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.refreshing {
		return false
	}

	h.refreshing = true
	return true
}

func (h *RandomJumpHandler) finishRefresh() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.refreshing = false
}

func randomSourceBatchSize(total int, reader io.Reader) (int, error) {
	if total <= 0 {
		return 0, errors.New("没有可选的爬虫源")
	}
	if total <= 2 {
		return total, nil
	}

	upper := total
	if upper > 4 {
		upper = 4
	}

	value, err := randomIndex(upper-1, reader)
	if err != nil {
		return 0, err
	}

	return value + 2, nil
}

func randomPermutation(total int, reader io.Reader) ([]int, error) {
	order := make([]int, total)
	for index := range order {
		order[index] = index
	}

	for index := total - 1; index > 0; index-- {
		pick, err := randomIndex(index+1, reader)
		if err != nil {
			return nil, err
		}
		order[index], order[pick] = order[pick], order[index]
	}

	return order, nil
}

func randomIndex(limit int, reader io.Reader) (int, error) {
	if limit <= 0 {
		return 0, errors.New("随机范围不能为空")
	}
	if reader == nil {
		reader = cryptorand.Reader
	}

	value, err := cryptorand.Int(reader, big.NewInt(int64(limit)))
	if err != nil {
		return 0, err
	}

	return int(value.Int64()), nil
}

func uniqueResults(items []RandomJumpResult) []RandomJumpResult {
	if len(items) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(items))
	results := make([]RandomJumpResult, 0, len(items))

	for _, item := range items {
		if item.URL == "" || item.Title == "" {
			continue
		}
		if _, exists := seen[item.URL]; exists {
			continue
		}
		seen[item.URL] = struct{}{}
		results = append(results, item)
	}

	return results
}

func hasEnoughUniqueResults(items []RandomJumpResult, target int) bool {
	if target <= 0 {
		return true
	}

	seen := make(map[string]struct{}, target)
	count := 0

	for _, item := range items {
		if item.URL == "" || item.Title == "" {
			continue
		}
		if _, exists := seen[item.URL]; exists {
			continue
		}
		seen[item.URL] = struct{}{}
		count++
		if count >= target {
			return true
		}
	}

	return false
}
