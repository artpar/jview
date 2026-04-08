---
layout: default
title: Modal
parent: Components
nav_order: 7
---

# Modal

Dialog overlay. Maps to **NSPanel**.

Displays a floating panel above the main window. Use Modal for confirmations, forms, or any content that requires user attention before continuing.

![Modal dialog]({{ site.baseurl }}/screenshots/modal_open.png){: .screenshot}

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `title` | DynamicString | | Title displayed in the panel title bar |
| `visible` | DynamicBoolean | `false` | Whether the modal is shown |
| `dataBinding` | string | | JSON Pointer path to bind the `visible` state |
| `width` | int | `480` | Panel width in points |
| `height` | int | `320` | Panel height in points |
| `onDismiss` | EventAction | | Action triggered when the modal is closed |

## Example

A confirmation dialog controlled by data binding:

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "Modal Example"
}

{
  "type": "updateDataModel",
  "surfaceId": "main",
  "operations": [
    {
      "op": "replace",
      "path": "/showConfirm",
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
        "openBtn",
        "dialog"
      ]
    },
    {
      "componentId": "openBtn",
      "type": "Button",
      "props": {
        "label": "Delete Item",
        "style": "destructive",
        "onClick": {
          "action": {
            "setValues": [
              {
                "path": "/showConfirm",
                "value": true
              }
            ]
          }
        }
      }
    },
    {
      "componentId": "dialog",
      "type": "Modal",
      "props": {
        "title": "Confirm Delete",
        "visible": {
          "$ref": "/showConfirm"
        },
        "dataBinding": "/showConfirm",
        "width": 400,
        "height": 200
      },
      "children": [
        "dialogContent"
      ]
    },
    {
      "componentId": "dialogContent",
      "type": "Column",
      "props": {
        "padding": 16,
        "gap": 12
      },
      "children": [
        "msg",
        "actions"
      ]
    },
    {
      "componentId": "msg",
      "type": "Text",
      "props": {
        "content": "Are you sure you want to delete this item?"
      }
    },
    {
      "componentId": "actions",
      "type": "Row",
      "props": {
        "gap": 8,
        "justify": "end"
      },
      "children": [
        "cancelBtn",
        "confirmBtn"
      ]
    },
    {
      "componentId": "cancelBtn",
      "type": "Button",
      "props": {
        "label": "Cancel",
        "onClick": {
          "action": {
            "setValues": [
              {
                "path": "/showConfirm",
                "value": false
              }
            ]
          }
        }
      }
    },
    {
      "componentId": "confirmBtn",
      "type": "Button",
      "props": {
        "label": "Delete",
        "style": "destructive"
      }
    }
  ]
}
```

## Notes

- Set `visible` to `true` to show the modal, `false` to hide it.
- When `dataBinding` is set, closing the modal (via the close button or Escape key) automatically sets the bound value to `false`.
- The `onDismiss` action fires when the user closes the modal by any means.
- Modal content is any valid component tree passed as children.
