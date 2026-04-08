---
layout: default
title: CheckBox
parent: Components
nav_order: 9
---

# CheckBox

Toggle switch. Maps to **NSButton** (checkBox style).

A labeled checkbox that toggles between checked and unchecked states. Supports two-way data binding for reactive UIs.

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `label` | DynamicString | | Text label displayed next to the checkbox |
| `checked` | DynamicBoolean | `false` | Whether the checkbox is checked |
| `dataBinding` | string | | JSON Pointer path for two-way binding |
| `onToggle` | EventAction | | Action triggered when the checkbox is toggled |

## Example

Terms agreement checkbox:

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "CheckBox Example"
}

{
  "type": "updateDataModel",
  "surfaceId": "main",
  "operations": [
    {
      "op": "replace",
      "path": "/agreed",
      "value": false
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
        "padding": 16,
        "gap": 12
      },
      "children": [
        "terms",
        "submit"
      ]
    },
    {
      "componentId": "terms",
      "type": "CheckBox",
      "props": {
        "label": "I agree to the terms and conditions",
        "dataBinding": "/agreed"
      }
    },
    {
      "componentId": "submit",
      "type": "Button",
      "props": {
        "label": "Continue",
        "style": "primary",
        "disabled": {
          "$not": {
            "$ref": "/agreed"
          }
        }
      }
    }
  ]
}
```

## Notes

- When `dataBinding` is set, toggling the checkbox writes `true` or `false` to the data model.
- Other components can read the bound value using `{"$ref":"/path"}` to react to checkbox state.
- The `onToggle` action fires every time the checkbox changes state.
- Use `checked` to set the initial state when not using data binding.
