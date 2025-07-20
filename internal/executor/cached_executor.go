package executor

import (
	"crypto/md5"
	"fmt"
	"log/slog"
	"runtime-engine/internal/runners"
	"runtime-engine/pkg/semaphore"
	"sync"
	"time"
)

type CachedExecutor struct {
	cache      map[string]runners.RunnerResult
	cacheMutex sync.RWMutex
	ttl        time.Duration
	semaphore  *semaphore.Semaphore
}

func NewCachedExecutor(ttl time.Duration, maxParallel int) *CachedExecutor {
	return &CachedExecutor{
		cache:     make(map[string]runners.RunnerResult),
		ttl:       ttl,
		semaphore: semaphore.New(maxParallel),
	}
}

func (e *CachedExecutor) Run(lang runners.Language, code []byte, log *slog.Logger) (runners.RunnerResult, error) {
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

	r, err := runners.GetRunner(lang)
	if err != nil {
		log.Error(op, "get runner error: %s", err)
		return runners.RunnerResult{}, err
	}

	result, err := r.Execute(code)
	if err != nil {
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
