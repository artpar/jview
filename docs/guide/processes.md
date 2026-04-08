---
layout: default
title: Processes
parent: Building Apps
nav_order: 7
---

# Processes

A process is a named background worker with its own message transport. Processes let you run timers, load additional JSONL files, or connect to an LLM -- all while the main app keeps running.

## Creating a Process

```json
{
  "type": "createProcess",
  "processId": "timer",
  "transport": {
    "type": "interval",
    "interval": 1000,
    "message": {
      "type": "updateDataModel",
      "surfaceId": "main",
      "ops": [
        {
          "op": "replace",
          "path": "/tick",
          "value": {
            "functionCall": {
              "name": "add",
              "args": [
                {
                  "path": "/tick"
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
```

The `transport` object configures how the process receives messages:

### Transport Types

| Type | Fields | Description |
|------|--------|-------------|
| `file` | `path` | Reads JSONL from a file, one message per line |
| `interval` | `interval`, `message` | Sends `message` every `interval` milliseconds |
| `llm` | `provider`, `model`, `prompt` | Connects to an LLM provider |

### File Transport

Loads and processes a JSONL file:

```json
{
  "type": "createProcess",
  "processId": "loader",
  "transport": {
    "type": "file",
    "path": "components/sidebar.jsonl"
  }
}
```

### Interval Transport

Sends the same message on a timer:

```json
{
  "type": "createProcess",
  "processId": "clock",
  "transport": {
    "type": "interval",
    "interval": 1000,
    "message": {
      "type": "updateDataModel",
      "surfaceId": "main",
      "ops": [
        {
          "op": "replace",
          "path": "/time",
          "value": {
            "functionCall": {
              "name": "now",
              "args": []
            }
          }
        }
      ]
    }
  }
}
```

### LLM Transport

Connects to an LLM that sends A2UI messages:

```json
{
  "type": "createProcess",
  "processId": "agent",
  "transport": {
    "type": "llm",
    "provider": "anthropic",
    "model": "claude-sonnet-4-20250514",
    "prompt": "Build a todo list app"
  }
}
```

## Process Status

Each process publishes its status to the data model at `/processes/{processId}/status`:

| Status | Meaning |
|--------|---------|
| `running` | Process is active |
| `stopped` | Process was stopped normally |
| `error` | Process encountered an error |

Display it in your UI:

```json
{
  "componentId": "status",
  "type": "Text",
  "props": {
    "content": {
      "functionCall": {
        "name": "concat",
        "args": [
          "Timer: ",
          {
            "path": "/processes/timer/status"
          }
        ]
      }
    }
  }
}
```

## Stopping a Process

```json
{
  "type": "stopProcess",
  "processId": "timer"
}
```

## Sending Messages to a Process

Route a message to a running process:

```json
{
  "type": "sendToProcess",
  "processId": "agent",
  "message": {
    "type": "updateDataModel",
    "surfaceId": "main",
    "ops": [
      {
        "op": "replace",
        "path": "/prompt",
        "value": "Add a delete button"
      }
    ]
  }
}
```

## Example: Auto-Incrementing Counter

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "Counter",
  "width": 300,
  "height": 200
}

{
  "type": "updateDataModel",
  "surfaceId": "main",
  "ops": [
    {
      "op": "add",
      "path": "/tick",
      "value": 0
    }
  ]
}

{
  "type": "updateComponents",
  "surfaceId": "main",
  "components": [
    {
      "componentId": "root",
      "type": "Column",
      "props": {
        "gap": 12,
        "align": "center"
      },
      "children": [
        "count",
        "controls"
      ]
    },
    {
      "componentId": "count",
      "type": "Text",
      "props": {
        "content": {
          "path": "/tick"
        },
        "variant": "h1"
      }
    },
    {
      "componentId": "controls",
      "type": "Row",
      "props": {
        "gap": 8
      },
      "children": [
        "startBtn",
        "stopBtn"
      ]
    },
    {
      "componentId": "startBtn",
      "type": "Button",
      "props": {
        "label": "Start"
      }
    },
    {
      "componentId": "stopBtn",
      "type": "Button",
      "props": {
        "label": "Stop"
      }
    }
  ]
}

{
  "type": "createProcess",
  "processId": "counter",
  "transport": {
    "type": "interval",
    "interval": 1000,
    "message": {
      "type": "updateDataModel",
      "surfaceId": "main",
      "ops": [
        {
          "op": "replace",
          "path": "/tick",
          "value": {
            "functionCall": {
              "name": "add",
              "args": [
                {
                  "path": "/tick"
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
```
