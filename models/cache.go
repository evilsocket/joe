package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/evilsocket/islazy/log"
	"sync"
	"time"
)

type CachePolicyType int

const (
	None CachePolicyType = iota
	ByKey
	ByTime
)

type Cached struct {
	At   time.Time
	Data interface{}
}

type CachePolicy struct {
	Type CachePolicyType `yaml:"type" json:"type"`
	Keys []string        `yaml:"keys" json:"keys"`
	TTL  int             `yaml:"ttl" json:"ttl"`

	cache sync.Map
}

func (c *CachePolicy) KeyFor(params map[string]interface{}) string {
	// results are cached regardless of the parameters for a given amount of time
	if c.Type == ByTime {
		return "time"
	}

	// results are cached by a key made of combining the value of the given parameters
	hash := sha256.New()
	for _, key := range c.Keys {
		what := ""
		if v, found := params[key]; found {
			what = fmt.Sprintf("%s:%s", key, v)
		} else {
			what = fmt.Sprintf("%s:<nil>", key)
		}
		hash.Write([]byte(what))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func (c *CachePolicy) Get(params map[string]interface{}) *Cached {
	// cache disabled
	if c.Type == None {
		return nil
	}
	key := c.KeyFor(params)
	if obj, found := c.cache.Load(key); !found {
		log.Debug("cache[%s] miss", key)
		// miss
		return nil
	} else if cached := obj.(*Cached); time.Since(cached.At).Seconds() > float64(c.TTL) {
		log.Debug("cache[%s] expired", key)
		// expired
		c.cache.Delete(key)
		return nil
	} else {
		log.Debug("cache[%s] hit", key)
		// found
		return cached
	}
}

func (c *CachePolicy) Set(params map[string]interface{}, data interface{}) {
	// cache disabled
	if c.Type == None {
		return
	}

	key := c.KeyFor(params)
	entry := &Cached{
		At:   time.Now(),
		Data: data,
	}

	log.Debug("cache[%s] = %v", key, entry)
	c.cache.Store(key, entry)
}
