package goevent

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/google/uuid"
)

var (
	globEmitter *Emitter
)

// New : New event emitter to be used locally
func New(handlePanic bool) *Emitter {
	return &Emitter{
		listeners:    make(map[string]map[string]Handler),
		handlePanics: handlePanic,
	}
}

// Global : global event listener singleton , run init on your main before this
func Global() *Emitter {
	return globEmitter
}

// InitGlobal : initialize the global event emitter
// put this in your main function
func InitGlobal(handlePanic bool) {
	globEmitter = New(handlePanic)
}

const hashSplit = "::"

// Handler : event handler function type
type Handler = func(d interface{})

// Emitter : event emitter
type Emitter struct {
	mu           sync.RWMutex
	listeners    map[string]map[string]Handler
	handlePanics bool
}

// Emit : emit an event to the system
func (e *Emitter) Emit(name string, data interface{}) {
	var wg sync.WaitGroup

	e.mu.RLock()
	defer e.mu.RUnlock()
	listeners := e.listeners[name]
	for _, fn := range listeners {
		if fn == nil {
			continue
		}
		wg.Add(1)
		go func(data interface{}, fn Handler, handlePanics bool) {
			defer wg.Done()
			defer func(handlepanics bool) {
				if handlepanics {
					if err := recover(); err != nil {
						log.Println(": govevent : panic occurred :", err)
					}
				}
			}(handlePanics)
			fn(data)
		}(data, fn, e.handlePanics)
	}
	wg.Wait()
}

// AddListener : add a listener to the system
func (e *Emitter) AddListener(eventName string, fn Handler) (hash string) {
	id := uuid.New().String()
	e.mu.Lock()
	defer e.mu.Unlock()
	val := e.listeners[eventName]
	if val == nil {
		val = make(map[string]Handler)

	}
	val[id] = fn
	e.listeners[eventName] = val

	return fmt.Sprintf("%s%s%s", eventName, hashSplit, id)
}

func (e *Emitter) parseHash(hash string) (eventName string, listenerId string) {
	slc := strings.Split(hash, hashSplit)

	if len(slc) != 2 {
		return "", ""
	}
	return slc[0], slc[1]
}

// RemoveListener : remove a listener from the system
func (e *Emitter) RemoveListener(hash string) (hasBeenRemoved bool) {

	eventName, listenerId := e.parseHash(hash)
	isEmpty := !(eventName != "" && listenerId != "")
	if isEmpty {
		return false
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	eventListeners := e.listeners[eventName]

	if eventListeners == nil {
		return false
	} else if eventListeners[listenerId] == nil {
		return false
	}

	delete(eventListeners, listenerId)
	return true
}
