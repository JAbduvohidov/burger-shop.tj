package app

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
)

type exactMux struct {
	mutex           sync.RWMutex
	routes          map[string]map[string]exactMuxEntry
	routesSorted    map[string][]exactMuxEntry
	notFoundHandler http.Handler
}

func NewExactMux() *exactMux {
	return &exactMux{}
}

func (m *exactMux) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if handler, err := m.handler(request.Method, request.URL.Path); err == nil {
		handler.ServeHTTP(writer, request)
	}

	if m.notFoundHandler != nil {
		m.notFoundHandler.ServeHTTP(writer, request)
	}
}

func (m *exactMux) GET(pattern string, handlerFunc func(responseWriter http.ResponseWriter, request *http.Request)) {
	m.HandleFunc(http.MethodGet, pattern, handlerFunc)
}

func (m *exactMux) POST(pattern string, handlerFunc func(responseWriter http.ResponseWriter, request *http.Request)) {
	m.HandleFunc(http.MethodPost, pattern, handlerFunc)
}

func (m *exactMux) HandleFunc(method string, pattern string, handlerFunc func(responseWriter http.ResponseWriter, request *http.Request)) {

	if !strings.HasPrefix(pattern, "/") {
		panic(fmt.Errorf("pattern must start with /: %s", pattern))
	}

	if handlerFunc == nil { // ?
		panic(errors.New("handler can't be empty"))
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	entry := exactMuxEntry{
		pattern: pattern,
		handler: http.HandlerFunc(handlerFunc),
		weight:  calculateWeight(pattern),
	}

	if _, exists := m.routes[method][pattern]; exists {
		panic(fmt.Errorf("ambigious mapping: %s", pattern))
	}

	if m.routes == nil {
		m.routes = make(map[string]map[string]exactMuxEntry)
	}

	if m.routes[method] == nil {
		m.routes[method] = make(map[string]exactMuxEntry)
	}

	m.routes[method][pattern] = entry
	m.appendSorted(method, entry)
}

func (m *exactMux) appendSorted(method string, entry exactMuxEntry) {
	if m.routesSorted == nil {
		m.routesSorted = make(map[string][]exactMuxEntry)
	}

	if m.routesSorted[method] == nil {
		m.routesSorted[method] = make([]exactMuxEntry, 0)
	}

	routes := append(m.routesSorted[method], entry)
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].weight > routes[j].weight
	})
	m.routesSorted[method] = routes
}

func (m *exactMux) handler(method string, path string) (handler http.Handler, err error) {
	entries, exists := m.routes[method]
	if !exists {
		return nil, fmt.Errorf("can't find handler for: %s, %s", method, path)
	}

	if entry, ok := entries[path]; ok {
		return entry.handler, nil
	}

	return nil, fmt.Errorf("can't find handler for: %s, %s", method, path)
}

type exactMuxEntry struct {
	pattern string
	handler http.Handler
	weight  int
}

func calculateWeight(pattern string) int {
	if pattern == "/" {
		return 0
	}

	count := (strings.Count(pattern, "/") - 1) * 2
	if !strings.HasSuffix(pattern, "/") {
		return count + 1
	}
	return count
}
