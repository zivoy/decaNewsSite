package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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
var articleListCache = Cache{
	list:     make(CacheList),
	basePath: articleLocation,
	id:       4,
}

var cacheList = map[int]cacheInterface{
	userCache.id:        &userCache,
	sessionsCache.id:    &sessionsCache,
	articleCache.id:     &articleCache,
	allowedLinkCache.id: &allowedLinkCache,
}

const ServerCommunicationEndpoint = "SERVER-COMM"
const ServerPollRate = 5 * time.Second // 2 seconds to prevent spams
const expireTime = 60 * 2              // 2 minute cache // maybe put this in db and have it be configurable by admin

const autoClearPoint = ServerCommunicationEndpoint + "/autoClear"
const actionHash = ServerCommunicationEndpoint + "/actionHash"

var clearingCache = false
var autoClearCache chan struct{}

// cache for communicating between servers
var actionCache = Cache{
	list:     make(CacheList),
	basePath: ServerCommunicationEndpoint + "/cacheAction",
	id:       -1,
}

const (
	clearValue = iota
	clearList
)

type cacheAction struct {
	CacheListId int    `json:"database"`
	ItemId      string `json:"item"`
	ActionType  int    `json:"type"`
	id          string
}

func (c *cacheAction) createCacheId() string {
	c.id = fmt.Sprintf("%d-%s", time.Now().Unix(), c.ItemId)
	return c.id
}

func (c cacheAction) execute() {
	targetCache := cacheList[c.CacheListId].(*Cache)
	switch c.ActionType {
	case clearValue:
		delete(targetCache.list, c.ItemId)
	case clearList:
		targetCache.clear()
	}
}

func actionCacheHash() string {
	var protoHash uint32 = 0
	for _, v := range actionCache.list {
		protoHash = protoHash + hashTo32(v.Cache.(cacheAction).id)
	}
	return strconv.Itoa(int(protoHash))
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
}

func stopClearingCache() {
	close(autoClearCache)
	clearingCache = false
}

func setAutoClear(val bool) {
	if clearingCache == val {
		return
	} else if clearingCache && !val {
		stopClearingCache()
	} else {
		initCacheClearing()
	}
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

func updateActionCache() {
	hashHere := actionCacheHash()
	if pathExists(dataBase, actionHash) {
		farHash, err := readValue(dataBase, actionHash)
		if err != nil {
			log.Println("failed to get cache action hash")
			return
		}
		if farHash.(string) == hashHere {
			return // no changes
		}
	}

	actions, err := readEntry(dataBase, actionCache.basePath)
	if err != nil {
		log.Println("failed to get cache action")
		return
	}
	// remove actions present locally
	for _, v := range actionCache.list {
		id := v.Cache.(cacheAction).id
		if _, ok := actions[v.Cache.(cacheAction).id]; ok {
			delete(actions, id)
		}
	}

	// add to local cache
	for id, data := range actions {
		v := data.(map[string]interface{})
		item := cacheAction{
			CacheListId: int(v["database"].(float64)),
			ItemId:      v["item"].(string),
			ActionType:  int(v["type"].(float64)),
			id:          id,
		}

		// get expire time
		rawTime := strings.Split(id, "-")[0]
		createdTime, err := strconv.ParseInt(rawTime, 10, 64)
		if err != nil {
			log.Println("Error decoding action time")
			continue
		}

		actionCache.list[id] = cache{
			Expire: createdTime + expireTime,
			Cache:  item,
		}

		// act on new items
		item.execute()
	}

	err = setEntry(dataBase, actionHash, hashHere)
	if err != nil {
		log.Println("failed to set cache action hash")
		return
	}
}

func (c *Cache) clear() {
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
	deleteDATA := cacheAction{
		CacheListId: c.id,
		ItemId:      id,
		ActionType:  clearValue,
	}
	sendAction(deleteDATA)
	delete(c.list, id)
}

// server heartbeat
func startServerComms() {
	ServerComms = true
	HearRateAlive = true
	go func() {
		for now := range time.Tick(ServerPollRate) {
			HearRateAlive = true
			// service action cache
			updateActionCache()
			for id, data := range actionCache.list {
				if data.Expire-now.Unix() < 0 {
					delete(actionCache.list, id)
					err := setEntry(dataBase, actionCache.path(id), nil)
					if err != nil {
						log.Println("failed to clear cache for ", id, err)
					}
				}
			}

			// maybe dont get this every time?
			value, err := readValue(dataBase, autoClearPoint)
			if err != nil {
				log.Println("failed to get autoClearVal")
			} else {
				setAutoClear(value.(bool))
			}
			if !ServerComms {
				break
			}
		}
	}()
}

func sendAction(action cacheAction) {
	actionCache.add(action.createCacheId(), action)
	err := setEntry(dataBase, actionCache.path(action.id), action)
	if err != nil && debug {
		//addLog(3, updater.UID, "cache delete failed", map[string]interface{}{"id": id, "cacheList": c.id})
		log.Println("failed to send cache clear req for", action.CacheListId, "on", action.ItemId, err)
	}
}
