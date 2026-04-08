package engine

import (
	"canopy/jlog"
	"canopy/protocol"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

// EventSubscription represents a single event subscription from an "on" message.
type EventSubscription struct {
	ID        string
	Event     string
	SurfaceID string
	Handler   protocol.EventAction
	Cancel    func() // stops timer, watcher, etc. Nil for passive subscriptions.
}

// EventManager manages event subscriptions from "on"/"off" protocol messages.
// It handles subscription lifecycle and cleanup.
type EventManager struct {
	mu   sync.Mutex
	subs map[string]*EventSubscription // subscriptionID → subscription
	sess *Session
	seq  int // auto-increment for unnamed subscriptions
}

// NewEventManager creates an EventManager attached to the given session.
func NewEventManager(sess *Session) *EventManager {
	return &EventManager{
		subs: make(map[string]*EventSubscription),
		sess: sess,
	}
}

// Subscribe registers an event subscription. If no ID is provided, one is generated.
func (em *EventManager) Subscribe(msg protocol.OnMessage) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	id := msg.ID
	if id == "" {
		em.seq++
		id = fmt.Sprintf("__auto_%d", em.seq)
	}

	// Remove existing subscription with same ID
	if old, exists := em.subs[id]; exists {
		if old.Cancel != nil {
			old.Cancel()
		}
		delete(em.subs, id)
	}

	sub := &EventSubscription{
		ID:        id,
		Event:     msg.Event,
		SurfaceID: msg.SurfaceID,
		Handler:   msg.Handler,
	}

	// Start system event source if applicable
	em.startEventSource(sub, msg.Config)

	em.subs[id] = sub

	jlog.Infof("events", msg.SurfaceID, "subscribed: id=%s event=%s", id, msg.Event)
	return nil
}

// startEventSource starts the background process for system events that need active monitoring.
// Sets sub.Cancel to stop the source when the subscription is removed.
func (em *EventManager) startEventSource(sub *EventSubscription, config map[string]interface{}) {
	if !strings.HasPrefix(sub.Event, "system.") {
		return
	}

	switch sub.Event {
	case "system.timer":
		em.startTimer(sub, config)
	}
}

// startTimer starts a periodic timer that fires the subscription's event.
func (em *EventManager) startTimer(sub *EventSubscription, config map[string]interface{}) {
	intervalMs := 1000 // default 1 second
	if config != nil {
		if v, ok := config["interval"]; ok {
			switch iv := v.(type) {
			case float64:
				intervalMs = int(iv)
			case json.Number:
				if n, err := iv.Int64(); err == nil {
					intervalMs = int(n)
				}
			}
		}
	}
	if intervalMs < 10 {
		intervalMs = 10 // minimum 10ms to prevent spin
	}

	done := make(chan struct{})
	sub.Cancel = func() { close(done) }

	surfaceID := sub.SurfaceID
	event := sub.Event

	go func() {
		ticker := time.NewTicker(time.Duration(intervalMs) * time.Millisecond)
		defer ticker.Stop()
		tick := 0
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				tick++
				data := fmt.Sprintf(`{"tick":%d,"elapsed":%d}`, tick, tick*intervalMs)
				em.Fire(event, surfaceID, data)
			}
		}
	}()
}

// Unsubscribe removes an event subscription by ID.
func (em *EventManager) Unsubscribe(id string) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	sub, exists := em.subs[id]
	if !exists {
		return fmt.Errorf("subscription %q not found", id)
	}
	if sub.Cancel != nil {
		sub.Cancel()
	}
	delete(em.subs, id)
	jlog.Infof("events", "", "unsubscribed: id=%s", id)
	return nil
}

// Fire invokes all subscriptions matching the given event name and optional surfaceID.
// Called by native event sources when an event occurs.
func (em *EventManager) Fire(event string, surfaceID string, data string) {
	em.mu.Lock()
	// Collect matching subscriptions while holding the lock
	var matches []*EventSubscription
	for _, sub := range em.subs {
		if sub.Event != event {
			continue
		}
		if sub.SurfaceID != "" && sub.SurfaceID != surfaceID {
			continue
		}
		matches = append(matches, sub)
	}
	em.mu.Unlock()

	// Execute handlers outside the lock
	for _, sub := range matches {
		sid := sub.SurfaceID
		if sid == "" {
			sid = surfaceID
		}
		em.sess.mu.Lock()
		if surf, ok := em.sess.surfaces[sid]; ok {
			surf.executeEventAction(&sub.Handler, data)
		}
		em.sess.mu.Unlock()
	}
}

// CleanupSurface removes all subscriptions scoped to the given surface.
func (em *EventManager) CleanupSurface(surfaceID string) {
	em.mu.Lock()
	defer em.mu.Unlock()

	for id, sub := range em.subs {
		if sub.SurfaceID == surfaceID {
			if sub.Cancel != nil {
				sub.Cancel()
			}
			delete(em.subs, id)
		}
	}
}

// CleanupAll removes all subscriptions.
func (em *EventManager) CleanupAll() {
	em.mu.Lock()
	defer em.mu.Unlock()

	for id, sub := range em.subs {
		if sub.Cancel != nil {
			sub.Cancel()
		}
		delete(em.subs, id)
	}
}
