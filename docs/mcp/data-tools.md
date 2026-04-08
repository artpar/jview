---
layout: default
title: Managing Data
parent: MCP Tools
nav_order: 3
---

# Managing Data

These tools let you write to the data model and wait for conditions to be met. Every surface has a data model -- a JSON document that components bind to. Changing a value in the data model automatically updates any component bound to that path.

---

## set_data_model

Write one or more values to the surface's data model.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `surface_id` | string | yes | The window whose data model to update |
| `ops` | array | yes | Array of set operations |

Each operation in `ops`:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `path` | string | yes | JSON Pointer path (e.g., `/form/name`) |
| `value` | any | yes | The value to set (string, number, boolean, object, array) |

**Example -- set a single value:**
```
mcp__canopy__set_data_model(
  surface_id: "main",
  ops: [{ "path": "/form/name", "value": "Alice" }]
)
```

**Example -- set multiple values at once:**
```
mcp__canopy__set_data_model(
  surface_id: "main",
  ops: [
    { "path": "/form/name", "value": "Alice" },
    { "path": "/form/email", "value": "alice@example.com" },
    { "path": "/form/agree", "value": true }
  ]
)
```

All operations are applied atomically. Components bound to the affected paths re-render with the new values.

---

## wait_for

Wait for a condition to be met before continuing. Useful for testing flows where you need to wait for an animation, a network response, or a data model update.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `surface_id` | string | yes | The window to watch |
| `component_id` | string | no | Wait for this component to exist |
| `path` | string | no | Wait for a non-null value at this data model path |
| `timeout` | number | no | Maximum wait time in milliseconds (default: 5000) |

Provide either `component_id` or `path`, not both.

**Example -- wait for a component to appear:**
```
mcp__canopy__wait_for(
  surface_id: "main",
  component_id: "results-list",
  timeout: 10000
)
```

**Example -- wait for data to load:**
```
mcp__canopy__wait_for(
  surface_id: "main",
  path: "/api/response",
  timeout: 5000
)
```

Returns when the condition is met, or an error if the timeout is reached.
