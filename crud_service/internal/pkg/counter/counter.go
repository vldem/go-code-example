package counter

import (
	"expvar"
	"runtime"
	"strconv"
	"sync"
)

var cErr *counter
var cInRequest *counter
var cOutRequest *counter
var cSuccessRequest *counter
var cFailedRequest *counter
var cCacheMis *counter
var cCacheHit *counter

type counter struct {
	cnt int
	m   *sync.RWMutex
}

func (c *counter) Inc() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt++
}

func (c *counter) String() string {
	c.m.RLock()
	defer c.m.RUnlock()
	return strconv.FormatInt(int64(c.cnt), 10)
}

func ErrorCounterInc() {
	cErr.Inc()
}
func InRequestInc() {
	cInRequest.Inc()
}
func OutRequestInc() {
	cOutRequest.Inc()
}
func SuccessRequestInc() {
	cSuccessRequest.Inc()
}
func FailedRequestInc() {
	cFailedRequest.Inc()
}
func CacheMisInc() {
	cCacheMis.Inc()
}
func CacheHitInc() {
	cCacheHit.Inc()
}

type Goroutines struct {
}

func (g *Goroutines) String() string {
	return strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
}

func init() {
	cErr = &counter{m: &sync.RWMutex{}}
	cInRequest = &counter{m: &sync.RWMutex{}}
	cOutRequest = &counter{m: &sync.RWMutex{}}
	cSuccessRequest = &counter{m: &sync.RWMutex{}}
	cFailedRequest = &counter{m: &sync.RWMutex{}}
	cCacheMis = &counter{m: &sync.RWMutex{}}
	cCacheHit = &counter{m: &sync.RWMutex{}}
	expvar.Publish("Errors", cErr)
	expvar.Publish("In requests", cInRequest)
	expvar.Publish("Out requests", cOutRequest)
	expvar.Publish("Success requests", cSuccessRequest)
	expvar.Publish("Failed requests", cFailedRequest)
	expvar.Publish("Cache miss", cCacheMis)
	expvar.Publish("Cache hit", cCacheHit)
	g := &Goroutines{}
	expvar.Publish("Goroutines", g)
}
