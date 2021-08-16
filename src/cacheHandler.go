package main

import "time"

type cache struct {
	Expire int64
	Cache  interface{}
}

var userCache = map[string]cache{}
var sessionsCache = map[string]cache{}
var articleCache = map[string]cache{}
var allowedLinkCache = map[string]cache{}

var autoClearCache chan struct{}
var clearingCache bool

const expireTime = 60 * 2 // 2 minute cache

func initCacheClearing() {
	ticker := time.NewTicker(2 * time.Hour)
	autoClearCache = make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				clearCache()
			case <-autoClearCache:
				ticker.Stop()
				return
			}
		}
	}()
	clearingCache = true
}

/*
func stopClearingCache() {
	close(autoClearCache)
	clearingCache = false
}
*/

func clearCache() {
	userCache = make(map[string]cache)
	sessionsCache = make(map[string]cache)
	articleCache = make(map[string]cache)
}

func getCache(dataCache map[string]cache, id string, getNewData func(id string) interface{}) interface{} {
	if data, ok := dataCache[id]; ok {
		if data.Expire-time.Now().Unix() < 0 {
			delete(dataCache, id)
		} else {
			return data.Cache
		}
	}
	data := getNewData(id)
	if data != nil {
		addCache(dataCache, id, data)
	}
	return data
}

func addCache(dataCache map[string]cache, id string, data interface{}) {
	dataCache[id] = cache{
		Expire: time.Now().Unix() + expireTime,
		Cache:  data,
	}
}

func deleteCache(dataCache map[string]cache, id string) {
	delete(dataCache, id)
}
