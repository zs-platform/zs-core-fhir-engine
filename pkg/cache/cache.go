package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir"
)

// Cache defines the interface for caching operations
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Invalidate(ctx context.Context, pattern string) error
	Health(ctx context.Context) error
}

// ResourceCache provides caching for FHIR resources
type ResourceCache struct {
	cache Cache
	ttl   time.Duration
}

// NewResourceCache creates a new resource cache
func NewResourceCache(cache Cache, defaultTTL time.Duration) *ResourceCache {
	return &ResourceCache{
		cache: cache,
		ttl:   defaultTTL,
	}
}

// GetResource retrieves a cached resource
func (rc *ResourceCache) GetResource(ctx context.Context, resourceType, id string) (*fhir.Resource, error) {
	key := rc.resourceKey(resourceType, id)
	
	data, err := rc.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	
	var resource fhir.Resource
	if err := json.Unmarshal(data, &resource); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached resource: %w", err)
	}
	
	log.Debugf("Cache hit for %s/%s", resourceType, id)
	return &resource, nil
}

// SetResource caches a resource
func (rc *ResourceCache) SetResource(ctx context.Context, resourceType, id string, resource *fhir.Resource) error {
	key := rc.resourceKey(resourceType, id)
	
	data, err := json.Marshal(resource)
	if err != nil {
		return fmt.Errorf("failed to marshal resource for cache: %w", err)
	}
	
	if err := rc.cache.Set(ctx, key, data, rc.ttl); err != nil {
		return err
	}
	
	log.Debugf("Cached %s/%s", resourceType, id)
	return nil
}

// DeleteResource removes a resource from cache
func (rc *ResourceCache) DeleteResource(ctx context.Context, resourceType, id string) error {
	key := rc.resourceKey(resourceType, id)
	return rc.cache.Delete(ctx, key)
}

// InvalidateResourceType invalidates all cached resources of a type
func (rc *ResourceCache) InvalidateResourceType(ctx context.Context, resourceType string) error {
	pattern := fmt.Sprintf("resource:%s:*", resourceType)
	return rc.cache.Invalidate(ctx, pattern)
}

// SearchCache provides caching for search results
type SearchCache struct {
	cache Cache
	ttl   time.Duration
}

// NewSearchCache creates a new search cache
func NewSearchCache(cache Cache, defaultTTL time.Duration) *SearchCache {
	return &SearchCache{
		cache: cache,
		ttl:   defaultTTL,
	}
}

// GetSearchResults retrieves cached search results
func (sc *SearchCache) GetSearchResults(ctx context.Context, resourceType string, params map[string][]string) (*fhir.Bundle, error) {
	key := sc.searchKey(resourceType, params)
	
	data, err := sc.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	
	var bundle fhir.Bundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached search results: %w", err)
	}
	
	log.Debugf("Search cache hit for %s", resourceType)
	return &bundle, nil
}

// SetSearchResults caches search results
func (sc *SearchCache) SetSearchResults(ctx context.Context, resourceType string, params map[string][]string, bundle *fhir.Bundle) error {
	key := sc.searchKey(resourceType, params)
	
	data, err := json.Marshal(bundle)
	if err != nil {
		return fmt.Errorf("failed to marshal search results for cache: %w", err)
	}
	
	if err := sc.cache.Set(ctx, key, data, sc.ttl); err != nil {
		return err
	}
	
	log.Debugf("Cached search results for %s", resourceType)
	return nil
}

// InvalidateSearch invalidates cached search results for a resource type
func (sc *SearchCache) InvalidateSearch(ctx context.Context, resourceType string) error {
	pattern := fmt.Sprintf("search:%s:*", resourceType)
	return sc.cache.Invalidate(ctx, pattern)
}

// TokenCache provides caching for authentication tokens
type TokenCache struct {
	cache Cache
	ttl   time.Duration
}

// NewTokenCache creates a new token cache
func NewTokenCache(cache Cache, defaultTTL time.Duration) *TokenCache {
	return &TokenCache{
		cache: cache,
		ttl:   defaultTTL,
	}
}

// GetToken retrieves a cached token
func (tc *TokenCache) GetToken(ctx context.Context, tokenID string) (map[string]interface{}, error) {
	key := tc.tokenKey(tokenID)
	
	data, err := tc.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	
	var token map[string]interface{}
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached token: %w", err)
	}
	
	return token, nil
}

// SetToken caches a token
func (tc *TokenCache) SetToken(ctx context.Context, tokenID string, token map[string]interface{}, ttl time.Duration) error {
	key := tc.tokenKey(tokenID)
	
	if ttl == 0 {
		ttl = tc.ttl
	}
	
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token for cache: %w", err)
	}
	
	return tc.cache.Set(ctx, key, data, ttl)
}

// DeleteToken removes a token from cache
func (tc *TokenCache) DeleteToken(ctx context.Context, tokenID string) error {
	key := tc.tokenKey(tokenID)
	return tc.cache.Delete(ctx, key)
}

// Key generation helpers

func (rc *ResourceCache) resourceKey(resourceType, id string) string {
	return fmt.Sprintf("resource:%s:%s", resourceType, id)
}

func (sc *SearchCache) searchKey(resourceType string, params map[string][]string) string {
	// Create a deterministic key from params
	key := fmt.Sprintf("search:%s:", resourceType)
	
	// Simple key generation - in production would use sorted params with hash
	for k, v := range params {
		key += fmt.Sprintf("%s=%s&", k, v)
	}
	
	return key
}

func (tc *TokenCache) tokenKey(tokenID string) string {
	return fmt.Sprintf("token:%s", tokenID)
}

// InMemoryCache implements Cache with in-memory storage (for development/testing)
type InMemoryCache struct {
	data map[string]cacheEntry
}

type cacheEntry struct {
	value     []byte
	expiresAt time.Time
}

// NewInMemoryCache creates a new in-memory cache
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]cacheEntry),
	}
}

// Get implements Cache
func (c *InMemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	entry, exists := c.data[key]
	if !exists {
		return nil, fmt.Errorf("cache miss: %s", key)
	}
	
	if time.Now().After(entry.expiresAt) {
		delete(c.data, key)
		return nil, fmt.Errorf("cache expired: %s", key)
	}
	
	return entry.value, nil
}

// Set implements Cache
func (c *InMemoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	c.data[key] = cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	return nil
}

// Delete implements Cache
func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	delete(c.data, key)
	return nil
}

// Invalidate implements Cache
func (c *InMemoryCache) Invalidate(ctx context.Context, pattern string) error {
	for key := range c.data {
		if matchPattern(key, pattern) {
			delete(c.data, key)
		}
	}
	return nil
}

// Health implements Cache
func (c *InMemoryCache) Health(ctx context.Context) error {
	// In-memory cache is always healthy
	return nil
}

// matchPattern checks if a key matches a pattern (simple wildcard matching)
func matchPattern(key, pattern string) bool {
	// Simple implementation - convert pattern to prefix matching
	if len(pattern) > 1 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(key) >= len(prefix) && key[:len(prefix)] == prefix
	}
	return key == pattern
}
