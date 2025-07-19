package logger

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarning
	LevelError
)

type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
	color    bool
}

func New(out io.Writer, minLevel Level, color bool) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
		color:    color,
	}
}

func (l *Logger) log(level Level, msg string, args ...interface{}) {
	if level < l.minLevel {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().Format("2006-01-02 15:04:05.000")
	levelStr := ""

	switch level {
	case LevelDebug:
		levelStr = "DEBUG"
		if l.color {
			levelStr = "\033[36m" + levelStr + "\033[0m" // Голубой
		}
	case LevelInfo:
		levelStr = "INFO"
		if l.color {
			levelStr = "\033[32m" + levelStr + "\033[0m" // Зеленый
		}
	case LevelWarning:
		levelStr = "WARN"
		if l.color {
			levelStr = "\033[33m" + levelStr + "\033[0m" // Желтый
		}
	case LevelError:
		levelStr = "ERROR"
		if l.color {
			levelStr = "\033[31m" + levelStr + "\033[0m" // Красный
		}
	}

	_, err := fmt.Fprintf(l.out, "%s [%s] %s\n", now, levelStr, fmt.Sprintf(msg, args...))
	if err != nil {
		return
	}
}

func (l *Logger) logWithOp(level Level, op string, msg string, args ...interface{}) {
	if level < l.minLevel {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().Format("2006-01-02 15:04:05.000")
	levelStr := ""

	switch level {
	case LevelDebug:
		levelStr = "DEBUG"
		if l.color {
			levelStr = "\033[36m" + levelStr + "\033[0m" // Голубой
		}
	case LevelInfo:
		levelStr = "INFO"
		if l.color {
			levelStr = "\033[32m" + levelStr + "\033[0m" // Зеленый
		}
	case LevelWarning:
		levelStr = "WARN"
		if l.color {
			levelStr = "\033[33m" + levelStr + "\033[0m" // Желтый
		}
	case LevelError:
		levelStr = "ERROR"
		if l.color {
			levelStr = "\033[31m" + levelStr + "\033[0m" // Красный
		}
	}

	_, err := fmt.Fprintf(l.out, "%s [%s] %s at %s\n", now, levelStr, fmt.Sprintf(msg, args...), op)
	if err != nil {
		return
	}
}

func (l *Logger) Debug(op string, msg string, args ...interface{}) {
	l.logWithOp(LevelDebug, op, msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(LevelInfo, msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(LevelWarning, msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(LevelError, msg, args...)
}
