package pokecache

import (
	"sync"
	"time"
	

)




type Cache struct{
	entrys	map[string]cacheEntry
	mu		sync.Mutex

}


type cacheEntry struct{
	createdAt 	time.Time
	val			[]byte
}



func (c *Cache)Add(key string,val []byte){
	
	c.mu.Lock()
	entry := cacheEntry{
		time.Now(),
		val,
	}
	c.entrys[key] = entry
	c.mu.Unlock()
}

func (c *Cache)Get(key string)([]byte,bool){
	c.mu.Lock()
	defer c.mu.Unlock()
	if entry, ok := c.entrys[key];ok{
		return entry.val,true
	}else {
		return nil,false
	}
}

func NewCache(interval time.Duration)*Cache{
	entry := make(map[string]cacheEntry)
	cache := Cache{
		entry,
		sync.Mutex{},
	}
	go cache.reapLoop(interval)
	
	return &cache
}

func (c *Cache)reapLoop(interval time.Duration){
	ticker := time.NewTicker(interval)
	
	defer ticker.Stop()
	for t := range ticker.C {
		c.mu.Lock()
		var keysToDelete []string
		for k,v := range c.entrys{
			
			creationInterval := t.Sub(v.createdAt)
			if interval < creationInterval{
				keysToDelete = append(keysToDelete, k)
			}
		}
		for _, k := range keysToDelete {
			delete(c.entrys,k)
		}
		c.mu.Unlock()
	}
	

}