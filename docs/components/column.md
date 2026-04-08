---
layout: default
title: Column
parent: Components
nav_order: 2
---

# Column

Vertical stack layout. Maps to **NSStackView** (vertical orientation).

Arranges child components top to bottom. Column is the most common root container and the default choice for page-level layout.

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `justify` | string | `"start"` | Vertical alignment: `start`, `center`, `end`, `spaceBetween`, `spaceAround` |
| `align` | string | `"start"` | Horizontal alignment: `start`, `center`, `end`, `stretch` |
| `gap` | int | `0` | Spacing in points between children |
| `padding` | int | `0` | Internal padding in points |

## Example

A heading followed by body text:

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "Column Example"
}

{
  "type": "updateComponents",
  "surfaceId": "main",
  "components": [
    {
      "componentId": "root",
      "type": "Column",
      "props": {
        "gap": 8,
        "padding": 16
      },
      "children": [
        "heading",
        "body"
      ]
    },
    {
      "componentId": "heading",
      "type": "Text",
      "props": {
        "content": "Welcome",
        "variant": "h1"
      }
    },
    {
      "componentId": "body",
      "type": "Text",
      "props": {
        "content": "This is a vertical layout with a heading and paragraph."
      }
    }
  ]
}
```

## Notes

- Column is typically the root component of a surface.
- Use `stretch` alignment to make all children fill the full width.
- Columns scroll automatically when content exceeds the window height.
- Combine with Row for grid-like layouts.
