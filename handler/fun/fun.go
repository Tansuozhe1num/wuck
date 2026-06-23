package fun

import (
	"context"
	cryptorand "crypto/rand"
	"errors"
	"io"
	"math/big"
	"net/http"
	"sort"
	"sync"
	"time"

	"luangao/store"
)

type RandomJumpResult struct {
	Title        string `json:"title"`
	URL          string `json:"url"`
	Category     string `json:"category"`
	Hint         string `json:"hint"`
	Source       string `json:"source"`
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
	FetchedAt    string `json:"fetchedAt"`
}

type RandomJumpFinder interface {
	Pick(ctx context.Context) (*RandomJumpResult, error)
	PickN(ctx context.Context, n int, userID string) ([]RandomJumpResult, error)
}

type crawlerSource interface {
	Name() string
	Category() string
	Hint() string
	Fetch(ctx context.Context, client *http.Client) ([]RandomJumpResult, error)
}

type RandomJumpHandler struct {
	client    *http.Client
	sources   []crawlerSource
	cacheTTL  time.Duration
	userStore *store.UserStore

	mu          sync.RWMutex
	cachedItems []RandomJumpResult
	lastCrawled time.Time
	refreshing  bool
}

func NewRandomJumpHandler(userStore *store.UserStore) *RandomJumpHandler {
	return NewRandomJumpHandlerWithSources(newDefaultHTTPClient(), newHotRankCrawlerSources(), 8*time.Minute, userStore)
}

func NewRandomJumpHandlerWithSources(client *http.Client, sources []crawlerSource, cacheTTL time.Duration, userStore *store.UserStore) *RandomJumpHandler {
	if client == nil {
		client = newDefaultHTTPClient()
	}
	if cacheTTL <= 0 {
		cacheTTL = 8 * time.Minute
	}

	return &RandomJumpHandler{
		client:    client,
		sources:   sources,
		cacheTTL:  cacheTTL,
		userStore: userStore,
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

func (h *RandomJumpHandler) PickN(ctx context.Context, n int, userID string) ([]RandomJumpResult, error) {
	if n <= 0 {
		return nil, errors.New("选取数量必须大于0")
	}

	items, err := h.loadItems(ctx)
	if err != nil {
		return nil, err
	}

	// ── Personalization ─────────────────────────────
	clickedURLs := make(map[string]bool)
	categoryPrefs := make(map[string]float64)
	if userID != "" && h.userStore != nil {
		clickedURLs = h.userStore.GetClickedURLs(userID)
		categoryPrefs = h.userStore.GetCategoryPreference(userID)
	}

	// Filter out already-clicked URLs
	if len(clickedURLs) > 0 {
		filtered := items[:0]
		for _, item := range items {
			if !clickedURLs[item.URL] {
				filtered = append(filtered, item)
			}
		}
		if len(filtered) >= n {
			items = filtered
		}
		// If not enough items left, keep original pool
	}

	if n >= len(items) {
		results := make([]RandomJumpResult, len(items))
		copy(results, items)
		return results, nil
	}

	// ── Weighted selection ──────────────────────────
	if len(categoryPrefs) > 0 && len(items) > n*2 {
		return h.weightedPickN(items, n, categoryPrefs)
	}

	// ── Pure random (fallback) ──────────────────────
	return h.randomPickN(items, n)
}

func (h *RandomJumpHandler) weightedPickN(items []RandomJumpResult, n int, categoryPrefs map[string]float64) ([]RandomJumpResult, error) {
	bySource := groupBySource(items)
	sourceOrder := shuffledSourceNames(bySource)

	// Score each item by user's category preference + small random jitter
	type scored struct {
		item  RandomJumpResult
		score float64
	}

	scoreItem := func(item RandomJumpResult) float64 {
		base := categoryPrefs[item.Category]
		if base <= 0 {
			base = 0.1
		}
		jitter, _ := randomFloat()
		return base + jitter*0.2
	}

	results := make([]RandomJumpResult, 0, n)
	used := make(map[string]struct{}, n)

	// First pass: pick the highest-scored item from each source
	for _, src := range sourceOrder {
		if len(results) >= n {
			break
		}
		pool := bySource[src]
		bestIdx := 0
		bestScore := -1.0
		for i, item := range pool {
			if _, ok := used[item.URL]; ok {
				continue
			}
			s := scoreItem(item)
			if s > bestScore {
				bestScore = s
				bestIdx = i
			}
		}
		if bestScore < 0 {
			continue
		}
		item := pool[bestIdx]
		used[item.URL] = struct{}{}
		results = append(results, item)
	}

	// Fallback: if not enough sources, fill remaining from all items
	if len(results) < n {
		scores := make([]scored, 0, len(items))
		for _, item := range items {
			if _, ok := used[item.URL]; ok {
				continue
			}
			scores = append(scores, scored{item: item, score: scoreItem(item)})
		}
		sort.Slice(scores, func(i, j int) bool {
			return scores[i].score > scores[j].score
		})
		for _, s := range scores {
			if len(results) >= n {
				break
			}
			if _, ok := used[s.item.URL]; ok {
				continue
			}
			used[s.item.URL] = struct{}{}
			results = append(results, s.item)
		}
	}

	return results, nil
}

func (h *RandomJumpHandler) randomPickN(items []RandomJumpResult, n int) ([]RandomJumpResult, error) {
	bySource := groupBySource(items)
	sourceOrder := shuffledSourceNames(bySource)

	results := make([]RandomJumpResult, 0, n)
	used := make(map[string]struct{}, n)

	// First pass: pick one item from each source
	for _, src := range sourceOrder {
		if len(results) >= n {
			break
		}
		pool := bySource[src]
		idx, err := randomIndex(len(pool), cryptorand.Reader)
		if err != nil {
			continue
		}
		item := pool[idx]
		used[item.URL] = struct{}{}
		results = append(results, item)
	}

	// Fallback: if not enough distinct sources, fill from all items
	if len(results) < n {
		order, err := randomPermutation(len(items), cryptorand.Reader)
		if err != nil {
			return results, nil
		}
		for _, idx := range order {
			if len(results) >= n {
				break
			}
			item := items[idx]
			if _, exists := used[item.URL]; exists {
				continue
			}
			used[item.URL] = struct{}{}
			results = append(results, item)
		}
	}

	return results, nil
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

const maxItemsPerSource = 8

func (h *RandomJumpHandler) crawlSources(ctx context.Context) ([]RandomJumpResult, error) {
	if len(h.sources) == 0 {
		return nil, errors.New("没有可用的爬虫源")
	}

	type crawlResult struct {
		items []RandomJumpResult
		err   error
	}

	resultsCh := make(chan crawlResult, len(h.sources))

	// Fire all sources concurrently — each capped by our 4s HTTP client timeout.
	// Total wall-clock time ≈ slowest source (≤4s). This only runs on cache
	// miss/refresh (every 8 min); normal requests hit in μs.
	for _, source := range h.sources {
		go func() {
			items, fetchErr := source.Fetch(ctx, h.client)
			resultsCh <- crawlResult{items: items, err: fetchErr}
		}()
	}

	results := make([]RandomJumpResult, 0, 48)
	var firstErr error

	for range h.sources {
		result := <-resultsCh
		if result.err != nil {
			if firstErr == nil {
				firstErr = result.err
			}
			continue
		}
		// Cap each source to prevent large sources (Bilibili 20, Douyin 47)
		// from numerically dominating smaller sources (Qidian ~10, Weibo ~12).
		capped := capSourceItems(result.items, maxItemsPerSource)
		results = append(results, capped...)
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

func capSourceItems(items []RandomJumpResult, max int) []RandomJumpResult {
	if len(items) <= max {
		return items
	}
	order, err := randomPermutation(len(items), cryptorand.Reader)
	if err != nil {
		// On error, just truncate (better than returning nothing)
		return items[:max]
	}
	result := make([]RandomJumpResult, max)
	for i, idx := range order[:max] {
		result[i] = items[idx]
	}
	return result
}

func groupBySource(items []RandomJumpResult) map[string][]RandomJumpResult {
	bySource := make(map[string][]RandomJumpResult)
	for _, item := range items {
		bySource[item.Source] = append(bySource[item.Source], item)
	}
	return bySource
}

func shuffledSourceNames(bySource map[string][]RandomJumpResult) []string {
	names := make([]string, 0, len(bySource))
	for name := range bySource {
		names = append(names, name)
	}
	order, err := randomPermutation(len(names), cryptorand.Reader)
	if err != nil {
		return names
	}
	shuffled := make([]string, len(names))
	for i, idx := range order {
		shuffled[i] = names[idx]
	}
	return shuffled
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

func randomFloat() (float64, error) {
	b := make([]byte, 4)
	if _, err := cryptorand.Read(b); err != nil {
		return 0, err
	}
	// Convert 4 random bytes to a float in [0, 1)
	v := uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
	return float64(v) / float64(1<<32), nil
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

