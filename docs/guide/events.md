# Events

Canopy provides two mechanisms for handling native events. Both use the same `EventAction` handler type with identical dispatch (event, functionCall, standardAction).

## Component Events: the `on` prop

Any component can handle native events via the `on` prop — a map of event names to handlers.

```json
{
  "componentId": "card1",
  "type": "Card",
  "props": {
    "title": "Hover me",
    "on": {
      "mouseEnter": {"dataPath": "/hovered", "dataValue": "card1"},
      "mouseLeave": {"dataPath": "/hovered", "dataValue": null}
    }
  }
}
```

Named props (`onClick`, `onChange`, `onToggle`, etc.) still work and are preferred for simple cases. They are syntactic sugar — internally folded into the `on` map. When both exist, `on` entries take precedence.

### Available component events

**Mouse:** mouseEnter, mouseLeave, doubleClick, rightClick

**Keyboard:** keyDown, keyUp

**Focus:** focus, blur

**Gesture:** magnify, rotate, scrollWheel

**Existing:** click, change, toggle, slide, select, dateChange, drop, dismiss, capture, error, ended, search

## Window & System Events: `on`/`off` messages

For events not tied to any component (window lifecycle, system state, timers), use `on`/`off` protocol messages:

```json
{"type": "on", "surfaceId": "main", "id": "resize-1", "event": "window.resize", "handler": {"dataPath": "/window/size"}}
```

Remove a subscription:
```json
{"type": "off", "id": "resize-1"}
```

Subscriptions are automatically cleaned up when their surface is destroyed.

## EventAction handler fields

| Field | Type | Description |
|-------|------|-------------|
| action | Action | Execute an action (event, functionCall, or standardAction) |
| dataPath | string | JSON Pointer — write to data model when event fires |
| dataValue | any | Value to write at dataPath. Omit to write native event data. |
| filter | EventFilter | Conditions for firing (key, modifiers, button) |
| throttle | int | Maximum fire rate in milliseconds |
| debounce | int | Wait this many ms of quiet before firing |

### DataPath shorthand

80% of event handlers just flip a value in the data model. `dataPath` eliminates the need for a verbose `functionCall.updateDataModel`:

```json
"on": {
  "mouseEnter": {"dataPath": "/ui/focused", "dataValue": "card1"},
  "focus": {"dataPath": "/ui/editing", "dataValue": true}
}
```

When `dataValue` is omitted, the native event data is written instead:

```json
"on": {
  "mouseEnter": {"dataPath": "/mouse/position"}
}
```

After a mouseEnter, `/mouse/position` contains `{"x": 150.0, "y": 200.0}`.

### Event filtering

Filter keyboard events by key and modifiers:

```json
"on": {
  "keyDown": {
    "filter": {"key": "s", "modifiers": ["cmd"]},
    "action": {"event": {"name": "save"}}
  }
}
```

Filter mouse events by button:

```json
"on": {
  "rightClick": {
    "filter": {"button": 1},
    "action": {"functionCall": {"call": "updateDataModel", "args": {"ops": [...]}}}
  }
}
```

### Throttle and debounce

Rate-limit high-frequency events:

```json
"on": {
  "scrollWheel": {"dataPath": "/scroll", "throttle": 100},
  "change": {"dataPath": "/searchQuery", "debounce": 300}
}
```

- **throttle**: fires at most once per interval (first call fires immediately, subsequent calls dropped until interval passes)
- **debounce**: waits for quiet period — resets timer on each new event, fires only the last one

## Common patterns

### Hover state

```json
"on": {
  "mouseEnter": {"dataPath": "/hovered", "dataValue": "card1"},
  "mouseLeave": {"dataPath": "/hovered", "dataValue": null}
}
```

Use with dynamic styling:
```json
"style": {
  "backgroundColor": {"functionCall": {"name": "if", "args": [
    {"functionCall": {"name": "equals", "args": [{"path": "/hovered"}, "card1"]}},
    "#E8F0FE", "#FFFFFF"
  ]}}
}
```

### Keyboard shortcuts

```json
"on": {
  "keyDown": {
    "filter": {"key": "Enter", "modifiers": ["cmd"]},
    "action": {"event": {"name": "submit"}}
  }
}
```

### Auto-save timer

```json
{"type": "on", "id": "autosave", "event": "system.timer", "config": {"interval": 30000},
 "handler": {"action": {"event": {"name": "autoSave"}}}}
```

### Window resize tracking

```json
{"type": "on", "surfaceId": "main", "id": "resize", "event": "window.resize",
 "handler": {"dataPath": "/window/size"}}
```

After resize, `/window/size` contains `{"width": 1024, "height": 768}`.

### Dark mode detection

```json
{"type": "on", "surfaceId": "main", "id": "theme", "event": "system.appearance",
 "handler": {"dataPath": "/system/appearance"}}
```

### File watcher

```json
{"type": "on", "surfaceId": "main", "id": "watcher", "event": "system.fs.watch",
 "config": {"paths": ["/tmp/data"]},
 "handler": {"action": {"event": {"name": "fileChanged"}}}}
```

### Debounced search input

Combine data binding (for TextField value sync) with a debounced `on` handler (for triggering search):

```json
{
  "componentId": "search",
  "type": "SearchField",
  "props": {
    "dataBinding": "/searchText",
    "on": {
      "change": {"dataPath": "/searchQuery", "debounce": 300}
    }
  }
}
```

## Event data reference

See the [Event Catalog](../spec.md#event-catalog) in the protocol specification for the complete list of events and their data shapes.
