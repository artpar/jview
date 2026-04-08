---
layout: default
title: Data Binding
parent: Building Apps
nav_order: 1
---

# Data Binding

Every surface has a **data model** -- a JSON document that stores your app's state. Components read from it, write to it, and re-render automatically when it changes.

## Setting Up the Data Model

Use `updateDataModel` to initialize values before creating components:

```json
{
  "type": "updateDataModel",
  "surfaceId": "main",
  "ops": [
    {
      "op": "add",
      "path": "/name",
      "value": ""
    },
    {
      "op": "add",
      "path": "/email",
      "value": ""
    },
    {
      "op": "add",
      "path": "/count",
      "value": 0
    }
  ]
}
```

## JSON Pointer Syntax

Paths use [JSON Pointer](https://datatracker.ietf.org/doc/html/rfc6901) notation:

| Path | Points to |
|------|-----------|
| `/name` | Top-level "name" field |
| `/user/name` | Nested field: `{"user":{"name":"..."}}` |
| `/items/0` | First element of "items" array |
| `/items/0/title` | "title" field of the first item |

## Reading from the Data Model

Use a **path reference** anywhere a prop accepts a dynamic value:

```json
{
  "componentId": "greeting",
  "type": "Text",
  "props": {
    "content": {
      "path": "/name"
    }
  }
}
```

When `/name` changes, this Text re-renders with the new value.

## Two-Way Binding with `dataBinding`

Input components (TextField, CheckBox, Slider, SearchField) support `dataBinding` -- a JSON Pointer that creates a two-way link between the component and the data model:

```json
{
  "componentId": "nameField",
  "type": "TextField",
  "props": {
    "placeholder": "Enter your name",
    "value": {
      "path": "/name"
    },
    "dataBinding": "/name"
  }
}
```

When the user types, the data model updates at `/name`. When `/name` changes from any source, the field updates.

## Binding Propagation

Changes propagate automatically. Here, typing in the TextField updates the Text in real time:

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "Binding Demo",
  "width": 400,
  "height": 200
}

{
  "type": "updateDataModel",
  "surfaceId": "main",
  "ops": [
    {
      "op": "add",
      "path": "/name",
      "value": ""
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
        "gap": 12
      },
      "children": [
        "field",
        "display"
      ]
    },
    {
      "componentId": "field",
      "type": "TextField",
      "props": {
        "placeholder": "Type your name",
        "value": {
          "path": "/name"
        },
        "dataBinding": "/name"
      }
    },
    {
      "componentId": "display",
      "type": "Text",
      "props": {
        "content": {
          "functionCall": {
            "name": "concat",
            "args": [
              "Hello, ",
              {
                "path": "/name"
              },
              "!"
            ]
          }
        }
      }
    }
  ]
}
```

The flow:

1. User types "Alice" into the TextField.
2. `dataBinding` writes "Alice" to `/name` in the data model.
3. The binding tracker finds all components that reference `/name`.
4. The Text re-renders, showing "Hello, Alice!".

## CheckBox Binding

CheckBox works the same way with boolean values:

```json
{
  "componentId": "check",
  "type": "CheckBox",
  "props": {
    "label": "Subscribe",
    "checked": {
      "path": "/subscribe"
    },
    "dataBinding": "/subscribe"
  }
}
```

## Nested Paths

Bind to any depth:

```json
{
  "type": "updateDataModel",
  "surfaceId": "main",
  "ops": [
    {
      "op": "add",
      "path": "/user",
      "value": {
        "name": "",
        "preferences": {
          "darkMode": false
        }
      }
    }
  ]
}
```

```json
{
  "componentId": "darkToggle",
  "type": "CheckBox",
  "props": {
    "label": "Dark Mode",
    "checked": {
      "path": "/user/preferences/darkMode"
    },
    "dataBinding": "/user/preferences/darkMode"
  }
}
```
