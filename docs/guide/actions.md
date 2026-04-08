---
layout: default
title: Actions
parent: Building Apps
nav_order: 3
---

# Actions

Actions define what happens when a user interacts with a component -- clicking a button, toggling a checkbox, or changing a text field. There are two kinds.

## Event Actions (Server-Bound)

An event action sends a named event to the LLM or transport that created the surface. Use events when you want the server to decide what happens next.

```json
{
  "componentId": "submitBtn",
  "type": "Button",
  "props": {
    "label": "Submit",
    "onClick": {
      "action": {
        "event": {
          "name": "submitForm",
          "dataRefs": [
            "/name",
            "/email"
          ]
        }
      }
    }
  }
}
```

Fields:
- `name` -- the event name the server receives
- `dataRefs` -- an array of JSON Pointer paths whose current values are sent with the event

When the user clicks Submit, the server receives an event named `submitForm` with the current values at `/name` and `/email`.

## FunctionCall Actions (Client-Bound)

A functionCall action runs a function locally, without contacting the server. Use these for client-side state updates.

```json
{
  "componentId": "addBtn",
  "type": "Button",
  "props": {
    "label": "Add Item",
    "onClick": {
      "action": {
        "functionCall": {
          "call": "updateDataModel",
          "args": {
            "ops": [
              {
                "op": "replace",
                "path": "/count",
                "value": {
                  "functionCall": {
                    "name": "add",
                    "args": [
                      {
                        "path": "/count"
                      },
                      1
                    ]
                  }
                }
              }
            ]
          }
        }
      }
    }
  }
}
```

## Built-in FunctionCall Actions

### updateDataModel

The most common client-side action. Applies JSON Patch operations to the data model:

```json
{
  "action": {
    "functionCall": {
      "call": "updateDataModel",
      "args": {
        "ops": [
          {
            "op": "add",
            "path": "/items/-",
            "value": {
              "title": "New"
            }
          },
          {
            "op": "replace",
            "path": "/count",
            "value": 5
          },
          {
            "op": "remove",
            "path": "/temp"
          }
        ]
      }
    }
  }
}
```

Operation types:
- `add` -- set a value (creates path if needed; use `/-` to append to array)
- `replace` -- update an existing value
- `remove` -- delete a value

Values in ops can be expressions (path references or function calls), which are resolved before applying.

### setTheme

Switch the app theme:

```json
{
  "action": {
    "functionCall": {
      "call": "setTheme",
      "args": {
        "theme": "dark"
      }
    }
  }
}
```

Valid themes: `"light"`, `"dark"`, `"system"`.

## When to Use Which

| Scenario | Use |
|----------|-----|
| Form submission to LLM | Event |
| Increment a counter | FunctionCall (updateDataModel) |
| Navigate between views | FunctionCall (updateDataModel) |
| Toggle dark mode | FunctionCall (setTheme) |
| Ask the server for new data | Event |
| Client-side filtering | FunctionCall (updateDataModel with expressions) |

## Action Props by Component

| Component | Action Prop | Trigger |
|-----------|------------|---------|
| Button | `onClick` | Click |
| TextField | `onChange` | Text changes |
| CheckBox | `onToggle` | Toggle |
| Slider | `onSlide` | Value changes |
| ChoicePicker | `onSelect` | Selection changes |
| DateTimeInput | `onDateChange` | Date changes |
| SearchField | `onSearch` | Search text changes |
| Modal | `onDismiss` | Modal dismissed |
| Video | `onEnded` | Playback ends |
| Any component | `onDrop` | File/text dropped |
| **Any component** | **`on`** | **Any native event** (see [Events guide](events.md)) |

## DataPath Shorthand

Most event handlers just write a value to the data model. Instead of a verbose `functionCall.updateDataModel`, use `dataPath`:

```json
"on": {
  "mouseEnter": {"dataPath": "/hovered", "dataValue": true},
  "mouseLeave": {"dataPath": "/hovered", "dataValue": false}
}
```

When `dataValue` is omitted, the native event data is written (e.g. `{"x": 150, "y": 200}` for mouse events).

## Event Filtering

Filter which events trigger the handler. Useful for keyboard shortcuts:

```json
"on": {
  "keyDown": {
    "filter": {"key": "s", "modifiers": ["cmd"]},
    "action": {"event": {"name": "save"}}
  }
}
```

Filter fields: `key` (key name), `modifiers` (`"cmd"`, `"shift"`, `"option"`, `"ctrl"`), `button` (0=left, 1=right, 2=middle).

## Throttle and Debounce

Rate-limit high-frequency event handlers:

```json
"on": {
  "scrollWheel": {"dataPath": "/scroll", "throttle": 100},
  "change": {"dataPath": "/searchQuery", "debounce": 300}
}
```

- **throttle** (ms): fire at most once per interval
- **debounce** (ms): wait for quiet period, then fire with the last event's data

See the [Events guide](events.md) for the full event catalog, window events, and system events.
