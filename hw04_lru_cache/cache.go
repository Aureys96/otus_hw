package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	if k, ok := l.items[key]; ok {
		updatedCacheItem := &cacheItem{key, value}
		k.Value = updatedCacheItem
		l.queue.MoveToFront(l.items[key])
		return ok
	}
	newCacheItem := newCacheItem(key, value)
	pushed := l.queue.PushFront(newCacheItem)
	l.items[key] = pushed
	if l.queue.Len() > l.capacity {
		back := l.queue.Back()
		l.queue.Remove(back)
		delete(l.items, back.Value.(*cacheItem).key)
	}
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	if k, ok := l.items[key]; ok {
		l.queue.MoveToFront(l.items[key])
		return k.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem)
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

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func newCacheItem(key Key, value interface{}) *cacheItem {
	return &cacheItem{
		key:   key,
		value: value,
	}
}
