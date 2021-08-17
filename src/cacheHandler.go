package main

import (
	"fmt"
	"log"
	"time"
)

type CacheList map[string]cache

type cacheInterface interface {
	clear()
	path(id string) string
	has(id string) bool
	get(id string, getNewData func(id string) interface{}) interface{}
	add(id string, data interface{})
	delete(id string)
}

type Cache struct {
	list     CacheList
	basePath string
	id       int
}

type cache struct {
	Expire int64
	Cache  interface{}
}

var userCache = Cache{
	list:     make(CacheList),
	basePath: userLocation,
	id:       0,
}
var sessionsCache = Cache{
	list:     make(CacheList),
	basePath: sessionLocation,
	id:       1,
}
var articleCache = Cache{
	list:     make(CacheList),
	basePath: articleLocation,
	id:       2,
}
var allowedLinkCache = Cache{
	list:     make(CacheList),
	basePath: allowedLinkLocation,
	id:       3,
}

var cacheList = map[int]cacheInterface{
	userCache.id:        userCache,
	sessionsCache.id:    sessionsCache,
	articleCache.id:     articleCache,
	allowedLinkCache.id: allowedLinkCache,
}

const ServerCommunicationEndpoint = "SERVER-COMM"
const ServerPollRate = 2 // 2 seconds to prevent spams

const expireTime = 60 * 2 // 2 minute cache // maybe put this in db and have it be configurable by admin

var autoClearCache chan struct{}

const autoClearPoint = ServerCommunicationEndpoint + "/autoClear"

var clearingCache bool

// cache for communicating between servers
var actionCache = Cache{
	list:     make(CacheList),
	basePath: ServerCommunicationEndpoint + "/cacheClear",
	id:       10,
}

// auto clear caches to save on ram for rarely visited pages that got cached
func initCacheClearing() {
	// completely clear all cache every 2 hours
	ticker := time.NewTicker(2 * time.Hour)
	autoClearCache = make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				wipeCache()
			case <-autoClearCache:
				ticker.Stop()
				return
			}
		}
	}()
	setAutoClear(true)
}

func stopClearingCache() {
	close(autoClearCache)
	clearingCache = false
}

func setAutoClear(val bool) {
	clearingCache = val
	err := setEntry(dataBase, autoClearPoint, val)
	if err != nil {
		log.Println("failed to set autoClear in db to", val)
	}
}

func wipeCache() {
	for _, c := range cacheList {
		c.clear()
	}
}

func (c Cache) clear() {
	c.list = make(CacheList)
}

func (c Cache) path(id string) string {
	return fmt.Sprintf(c.basePath+"/%s", id)
}

func (c Cache) has(id string) bool {
	_, exists := c.list[id]
	return exists
}

func (c Cache) get(id string, getNewData func(id string) interface{}) interface{} {
	if data, ok := c.list[id]; ok {
		if data.Expire-time.Now().Unix() < 0 {
			delete(c.list, id)
		} else {
			return data.Cache
		}
	}
	data := getNewData(id)
	if data != nil {
		c.add(id, data)
	}
	return data
}

func (c Cache) add(id string, data interface{}) {
	c.list[id] = cache{
		Expire: time.Now().Unix() + expireTime,
		Cache:  data,
	}
}

func (c Cache) delete(id string) {
	deleteID := fmt.Sprintf("%d-%s", time.Now().Unix(), id)
	deleteDATA := map[string]interface{}{
		"database": c.id,
		"item":     id,
	}
	actionCache.add(deleteID, deleteDATA)
	err := setEntry(dataBase, actionCache.path(deleteID), deleteDATA)
	if err != nil && debug {
		//addLog(3, updater.UID, "cache delete failed", map[string]interface{}{"id": id, "cacheList": c.id})
		log.Println("failed to send cache clear req for", c.id, "on", id, err)
	}
	delete(c.list, id)
}

var serverPolled int64 = 0

func getServerComms() {
	timeNow := time.Now().Unix()
	if (serverPolled+ServerPollRate)-timeNow < 0 {
		serverPolled = timeNow
		// get list of cache clears
		// get poll var

	}
}

/*todo function that gets all cache events
err := setEntry(dataBase, dataCache.path(id), nil)
if err != nil && debug {
	log.Println("failed to clear cache for ",dataCache.id,err)
}
*/
