package executor

import (
	"crypto/md5"
	"fmt"
	"runtime-engine/internal/runners"
	"runtime-engine/pkg/logger"
	"sync"
	"time"
)

type CachedExecutor struct {
	cache      map[string]runners.RunnerResult
	cacheMutex sync.RWMutex
	ttl        time.Duration
	semaphore  chan struct{}
}

func NewCachedExecutor(ttl time.Duration, maxParallel int) *CachedExecutor {
	return &CachedExecutor{
		cache:     make(map[string]runners.RunnerResult),
		ttl:       ttl,
		semaphore: make(chan struct{}, maxParallel),
	}
}

func (e *CachedExecutor) Run(lang runners.Language, code []byte, log *logger.Logger) (runners.RunnerResult, error) {
	op := "executor.Run"
	key := fmt.Sprintf("%s:%x", lang, md5.Sum(code))

	e.cacheMutex.RLock()
	item, ok := e.cache[key]
	e.cacheMutex.RUnlock()

	if ok {
		log.Debug(op, "found cached job result %s", key)
		return item, nil
	}

	e.semaphore <- struct{}{}
	defer func() { <-e.semaphore }()

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
	e.cacheMutex.Unlock()

	time.AfterFunc(e.ttl, func() {
		e.cacheMutex.Lock()
		delete(e.cache, key)
		e.cacheMutex.Unlock()
		log.Debug(op, "cached job result cleaned %s", key)
	})

	return result, nil
}
