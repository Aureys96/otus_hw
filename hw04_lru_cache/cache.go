package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	if k, ok := l.items[key]; ok {
		l.queue.MoveToFront(k)
		l.items[key].Value = value
		return ok
	} else {
		pushed := l.queue.PushFront(value)
		l.items[key] = pushed
		if l.queue.Len() > l.capacity {
			back := l.queue.Back()
			l.queue.Remove(back)
			delete(l.items, key)
		}
		return ok
	}
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	if k, ok := l.items[key]; ok {
		l.queue.MoveToFront(k)
		return k.Value, true
	} else {
		return nil, false
	}
}

func (l *lruCache) Clear() {
	//TODO implement me
	panic("implement me")
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
