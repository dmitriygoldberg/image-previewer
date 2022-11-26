package cache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	newCacheItem := &cacheItem{key, value}

	if queueItem, ok := cache.items[newCacheItem.key]; ok {
		queueItem.Value = newCacheItem
		cache.queue.MoveToFront(queueItem)

		return ok
	}

	if cache.queue.Len() == cache.capacity {
		toRemove := cache.queue.Back()
		delete(cache.items, toRemove.Value.(*cacheItem).key)
		cache.queue.Remove(toRemove)
	}

	cache.queue.PushFront(newCacheItem)
	cache.items[newCacheItem.key] = cache.queue.Front()

	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	if queueItem, ok := cache.items[key]; ok {
		cache.queue.MoveToFront(queueItem)

		return queueItem.Value.(*cacheItem).value, ok
	}

	return nil, false
}

func (cache *lruCache) Clear() {
	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
