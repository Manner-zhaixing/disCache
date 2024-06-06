package lru

import "container/list"

type Cache struct {
	maxBytes int                           // 缓存最大容量
	nBytes   int                           // 已经占用的容量
	ll       *list.List                    // 双向链表
	cache    map[string]*list.Element      // 缓存的键值对
	onEvict  func(key string, value value) // 某条记录被移除时的回调函数
}

// 双向链表节点的数据类型
type entry struct {
	key   string
	value value
}

type value interface {
	Len() int
}

// New 初始化LRU缓存
func New(maxBytes int, onEvict func(key string, value value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
		onEvict:  onEvict,
	}
}

// Get 查找操作
func (c *Cache) Get(key string) (value value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return
}

// RemoveOldest 删除操作
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes = c.nBytes - len(kv.key) - kv.value.Len()
		if c.onEvict != nil {
			c.onEvict(kv.key, kv.value)
		}
	}
}

// Add 添加/修改操作
func (c *Cache) Add(key string, value value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes = c.nBytes - kv.value.Len() + value.Len()
		kv.value = value
		return
	} else {
		ele := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.nBytes = c.nBytes + value.Len()
	}
	for c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

// Len 获取缓存数量
func (c *Cache) Len() int {
	return c.ll.Len()
}
