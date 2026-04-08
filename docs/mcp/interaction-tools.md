---
layout: default
title: Interacting with Components
parent: MCP Tools
nav_order: 2
---

# Interacting with Components

These tools simulate user interactions -- clicking buttons, typing into fields, toggling checkboxes, and sending custom events.

---

## click

Click a button or any component with an `onClick` handler.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `surface_id` | string | yes | The window containing the component |
| `component_id` | string | yes | The component to click |

**Example:**
```
mcp__canopy__click(surface_id: "main", component_id: "submit-btn")
```

This triggers the component's `onClick` action, exactly as if the user had clicked it with the mouse.

---

## fill

Type text into a TextField or SearchField component.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `surface_id` | string | yes | The window containing the field |
| `component_id` | string | yes | The text field to fill |
| `value` | string | yes | The text to enter |

**Example:**
```
mcp__canopy__fill(surface_id: "main", component_id: "name-field", value: "Alice")
```

This replaces the field's current content and triggers any `onChange` bindings, updating the data model just as if the user had typed the value.

---

## toggle

Toggle a CheckBox component on or off.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `surface_id` | string | yes | The window containing the checkbox |
| `component_id` | string | yes | The checkbox to toggle |

**Example:**
```
mcp__canopy__toggle(surface_id: "main", component_id: "agree-checkbox")
```

Flips the checkbox state and triggers its `onToggle` binding, updating the data model.

---

## interact

Send a generic interaction event to a component. Use this for component types that have specialized events beyond click, fill, and toggle.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `surface_id` | string | yes | The window containing the component |
| `component_id` | string | yes | The component to interact with |
| `event_type` | string | yes | The type of event to send |
| `data` | object | no | Event-specific data |

**Supported event types:**

| Event | Components | Data |
|-------|-----------|------|
| `slide` | Slider | `{ "value": 0.75 }` |
| `select` | ChoicePicker | `{ "value": "option-2" }` |
| `datechange` | DateTimeInput | `{ "value": "2026-04-08T10:30:00" }` |

**Examples:**

Move a slider to 75%:
```
mcp__canopy__interact(
  surface_id: "main",
  component_id: "volume-slider",
  event_type: "slide",
  data: { "value": 0.75 }
)
```

Select an option from a picker:
```
mcp__canopy__interact(
  surface_id: "main",
  component_id: "color-picker",
  event_type: "select",
  data: { "value": "blue" }
)
```

Change a date:
```
mcp__canopy__interact(
  surface_id: "main",
  component_id: "deadline-input",
  event_type: "datechange",
  data: { "value": "2026-12-31T23:59:00" }
)
```
