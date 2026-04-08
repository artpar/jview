---
layout: default
title: Row
parent: Components
nav_order: 1
---

# Row

Horizontal stack layout. Maps to **NSStackView** (horizontal orientation).

Arranges child components side by side from left to right. Use Row for toolbars, button groups, or any horizontal arrangement of elements.

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `justify` | string | `"start"` | Horizontal alignment: `start`, `center`, `end`, `spaceBetween`, `spaceAround` |
| `align` | string | `"start"` | Vertical alignment: `start`, `center`, `end`, `stretch` |
| `gap` | int | `0` | Spacing in points between children |
| `padding` | int | `0` | Internal padding in points |

## Example

Two buttons side by side with spacing:

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "Row Example"
}

{
  "type": "updateComponents",
  "surfaceId": "main",
  "components": [
    {
      "componentId": "root",
      "type": "Row",
      "props": {
        "gap": 12,
        "align": "center"
      },
      "children": [
        "cancel",
        "save"
      ]
    },
    {
      "componentId": "cancel",
      "type": "Button",
      "props": {
        "label": "Cancel",
        "style": "secondary"
      }
    },
    {
      "componentId": "save",
      "type": "Button",
      "props": {
        "label": "Save",
        "style": "primary"
      }
    }
  ]
}
```

## Notes

- Children are laid out in the order they appear in the `children` array.
- `spaceBetween` distributes equal space between children, pushing the first and last to the edges.
- `spaceAround` distributes equal space around each child.
- `stretch` alignment makes all children match the tallest child's height.
- Nest Rows inside Columns (and vice versa) to build complex layouts.
