package jlog

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Level represents a log severity level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel converts a string to a Level. Returns LevelInfo for unknown strings.
func ParseLevel(s string) Level {
	switch strings.ToLower(s) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

// Entry is a single log entry.
type Entry struct {
	Seq       uint64    `json:"seq"`
	Time      time.Time `json:"time"`
	Level     Level     `json:"level"`
	LevelStr  string    `json:"levelStr"`
	Component string    `json:"component"`
	Surface   string    `json:"surface,omitempty"`
	Message   string    `json:"message"`
}

// QueryOpts controls log query filtering and pagination.
type QueryOpts struct {
	MinLevel  Level
	Component string
	Surface   string
	Pattern   string // regex pattern on message text
	Limit     int
	Offset    int
}

// QueryResult contains matching entries and total count for pagination.
type QueryResult struct {
	Entries []Entry `json:"entries"`
	Total   int     `json:"total"`
}

// Config configures the logger at startup.
type Config struct {
	MaxEntries int
	MinLevel   Level
	LogDir     string // directory for log files; empty disables file output
}

// Logger is the central log sink: ring buffer + file output + stderr.
type Logger struct {
	mu       sync.RWMutex
	entries  []Entry
	head     int // next write position in ring buffer
	count    int // number of entries currently stored
	maxSize  int
	file     *os.File
	minLevel Level
	seq      atomic.Uint64
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// Init initializes the global logger. Safe to call multiple times; only the first call takes effect.
func Init(cfg Config) {
	once.Do(func() {
		if cfg.MaxEntries <= 0 {
			cfg.MaxEntries = 10000
		}
		l := &Logger{
			entries:  make([]Entry, cfg.MaxEntries),
			maxSize:  cfg.MaxEntries,
			minLevel: cfg.MinLevel,
		}
		if cfg.LogDir != "" {
			dir := expandHome(cfg.LogDir)
			os.MkdirAll(dir, 0755)
			path := filepath.Join(dir, fmt.Sprintf("jview-%d.log", os.Getpid()))
			f, err := os.Create(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "jlog: failed to open log file %s: %v\n", path, err)
			} else {
				l.file = f
			}
		}
		defaultLogger = l
	})
}

// Close flushes and closes the log file.
func Close() {
	if defaultLogger != nil && defaultLogger.file != nil {
		defaultLogger.file.Close()
		defaultLogger.file = nil
	}
}

// SetLevel changes the minimum capture level at runtime.
func SetLevel(level Level) {
	if defaultLogger != nil {
		defaultLogger.mu.Lock()
		defaultLogger.minLevel = level
		defaultLogger.mu.Unlock()
	}
}

// Log writes a log entry to the global logger.
func Log(level Level, component, surface, msg string) {
	if defaultLogger == nil {
		// Fallback: write to stderr if Init hasn't been called
		fmt.Fprintf(os.Stderr, "%s %s [%s/%s] %s\n",
			time.Now().Format("2006-01-02T15:04:05.000"), level, component, surface, msg)
		return
	}
	defaultLogger.log(level, component, surface, msg)
}

func (l *Logger) log(level Level, component, surface, msg string) {
	l.mu.RLock()
	minLevel := l.minLevel
	l.mu.RUnlock()

	if level < minLevel {
		return
	}

	entry := Entry{
		Seq:       l.seq.Add(1),
		Time:      time.Now(),
		Level:     level,
		LevelStr:  level.String(),
		Component: component,
		Surface:   surface,
		Message:   msg,
	}

	// Format line once for file + stderr
	var line string
	if surface != "" {
		line = fmt.Sprintf("%s %s [%s/%s] %s",
			entry.Time.Format("2006-01-02T15:04:05.000"), entry.LevelStr, component, surface, msg)
	} else {
		line = fmt.Sprintf("%s %s [%s] %s",
			entry.Time.Format("2006-01-02T15:04:05.000"), entry.LevelStr, component, msg)
	}

	// Write to stderr
	fmt.Fprintln(os.Stderr, line)

	l.mu.Lock()
	// Ring buffer write
	l.entries[l.head] = entry
	l.head = (l.head + 1) % l.maxSize
	if l.count < l.maxSize {
		l.count++
	}
	f := l.file
	l.mu.Unlock()

	// Write to file (outside lock — file writes are thread-safe enough with a single writer)
	if f != nil {
		fmt.Fprintln(f, line)
	}
}

// Query returns log entries matching the given filters.
func Query(opts QueryOpts) QueryResult {
	if defaultLogger == nil {
		return QueryResult{}
	}
	return defaultLogger.query(opts)
}

func (l *Logger) query(opts QueryOpts) QueryResult {
	if opts.Limit <= 0 {
		opts.Limit = 50
	}
	if opts.Limit > 500 {
		opts.Limit = 500
	}

	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return QueryResult{}
		}
	}

	l.mu.RLock()
	count := l.count
	head := l.head
	entries := l.entries
	l.mu.RUnlock()

	// Iterate oldest to newest
	var matched []Entry
	start := 0
	if count == l.maxSize {
		start = head // ring buffer wrapped; head is the oldest
	}

	for i := 0; i < count; i++ {
		idx := (start + i) % l.maxSize
		e := entries[idx]

		if e.Level < opts.MinLevel {
			continue
		}
		if opts.Component != "" && e.Component != opts.Component {
			continue
		}
		if opts.Surface != "" && e.Surface != opts.Surface {
			continue
		}
		if re != nil && !re.MatchString(e.Message) {
			continue
		}
		matched = append(matched, e)
	}

	total := len(matched)

	// Apply offset + limit
	if opts.Offset >= len(matched) {
		return QueryResult{Total: total}
	}
	matched = matched[opts.Offset:]
	if len(matched) > opts.Limit {
		matched = matched[:opts.Limit]
	}

	return QueryResult{
		Entries: matched,
		Total:   total,
	}
}

// Convenience functions

func Debug(component, surface, msg string) { Log(LevelDebug, component, surface, msg) }
func Info(component, surface, msg string)  { Log(LevelInfo, component, surface, msg) }
func Warn(component, surface, msg string)  { Log(LevelWarn, component, surface, msg) }
func Error(component, surface, msg string) { Log(LevelError, component, surface, msg) }

func Debugf(component, surface, format string, args ...any) {
	Log(LevelDebug, component, surface, fmt.Sprintf(format, args...))
}
func Infof(component, surface, format string, args ...any) {
	Log(LevelInfo, component, surface, fmt.Sprintf(format, args...))
}
func Warnf(component, surface, format string, args ...any) {
	Log(LevelWarn, component, surface, fmt.Sprintf(format, args...))
}
func Errorf(component, surface, format string, args ...any) {
	Log(LevelError, component, surface, fmt.Sprintf(format, args...))
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}
