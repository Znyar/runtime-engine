package executor

import (
	"crypto/md5"
	"fmt"
	"log/slog"
	"runtime-engine/internal/runner"
	"runtime-engine/pkg/semaphore"
	"sync"
	"time"
)

type CachedExecutor struct {
	cache      map[string]runner.Result
	cacheMutex sync.RWMutex
	ttl        time.Duration
	semaphore  *semaphore.Semaphore
}

func NewCachedExecutor(ttl time.Duration, maxParallel int) *CachedExecutor {
	return &CachedExecutor{
		cache:     make(map[string]runner.Result),
		ttl:       ttl,
		semaphore: semaphore.New(maxParallel),
	}
}

func (e *CachedExecutor) Run(lang string, version string, code []byte, log *slog.Logger) (runner.Result, error) {
	op := "executor.Run"
	key := fmt.Sprintf("%s:%x", lang, md5.Sum(code))

	e.cacheMutex.RLock()
	item, ok := e.cache[key]
	e.cacheMutex.RUnlock()

	if ok {
		log.Debug(op, "found cached job result %s", key)
		return item, nil
	}

	e.semaphore.Acquire()
	defer e.semaphore.Release()

	result, err := runner.Execute(code, lang, version, log)
	if err != nil {
		log.Error(op, "failed to execute code:", err)
		return result, err
	}

	e.cacheMutex.Lock()
	e.cache[key] = result
	log.Debug(op, "job result cached %s", key)
	e.cacheMutex.Unlock()

	time.AfterFunc(e.ttl, func() {
		e.cacheMutex.Lock()
		delete(e.cache, key)
		e.cacheMutex.Unlock()
		log.Debug(op, "cached job result cleaned %s", key)
	})

	return result, nil
}
