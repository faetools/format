package notion

import (
	"context"
	"sync"
)

// CacheClient is a client that caches results.
type CacheClient struct {
	cl    *Client
	pages *sync.Map
	// dbs    *sync.Map
	// blocks *sync.Map

	muxes *sync.Map
}

func (c *CacheClient) mutexForID(id Id) *sync.Mutex {
	mu, _ := c.muxes.LoadOrStore(id, &sync.Mutex{})
	return mu.(*sync.Mutex)
}

// GetPage returns a given page
func (c *CacheClient) GetPage(ctx context.Context, id Id) (*Page, error) {
	// lock for each call requesting the same ID
	mu := c.mutexForID(id)

	mu.Lock()
	defer mu.Unlock()

	if p, ok := c.pages.Load(id); ok {
		return p.(*Page), nil
	}

	p, err := c.cl.GetNotionPage(ctx, id)
	if err != nil {
		return nil, err
	}

	c.pages.Store(id, p)

	return p, nil
}
