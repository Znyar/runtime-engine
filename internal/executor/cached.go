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
	cache map[string]runners.RunnerResult
	mu    sync.RWMutex
	ttl   time.Duration
}

func NewCachedExecutor(ttl time.Duration) *CachedExecutor {
	return &CachedExecutor{
		cache: make(map[string]runners.RunnerResult),
		ttl:   ttl,
	}
}

func (e *CachedExecutor) Run(lang runners.Language, code []byte, log *logger.Logger) (runners.RunnerResult, error) {
	op := "executor.Run"

	key := fmt.Sprintf("%s:%x", lang, md5.Sum(code))

	log.Debug(op, "lock check cache job %s", key)

	e.mu.RLock()
	item, ok := e.cache[key]
	e.mu.RUnlock()
	log.Debug(op, "unlock check cache job %s", key)

	if ok {
		log.Debug(op, "found cached job result %s", key)
		return item, nil
	}

	r, err := runners.GetRunner(lang)

	if err != nil {
		log.Error("get runner error: %s", err)
		return runners.RunnerResult{}, err
	}

	result, err := r.Execute(code)

	log.Debug(op, "lock runner job %s", key)
	e.mu.Lock()
	defer func() {
		log.Debug(op, "unlock runner job %s", key)
		e.mu.Unlock()
	}()

	e.cache[key] = result

	log.Debug(op, "job result cached %s", key)

	time.AfterFunc(e.ttl, func() {
		log.Debug(op, "job cache clean scheduled %s", key)
		e.mu.Lock()
		delete(e.cache, key)
		e.mu.Unlock()
		log.Debug(op, "cached job result cleaned %s\n", key)
	})

	return result, err
}
