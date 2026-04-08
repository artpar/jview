---
layout: default
title: List
parent: Components
nav_order: 6
---

# List

Scrollable container. Renders children in a vertical stack inside an **NSScrollView**.

Use List when you have many items that may exceed the visible area. List provides the same layout props as Column but adds automatic scrolling.

![List component]({{ site.baseurl }}/screenshots/list.png){: .screenshot}

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `justify` | string | `"start"` | Vertical alignment: `start`, `center`, `end`, `spaceBetween`, `spaceAround` |
| `align` | string | `"start"` | Horizontal alignment: `start`, `center`, `end`, `stretch` |
| `gap` | int | `0` | Spacing in points between children |
| `padding` | int | `0` | Internal padding in points |

## Example

A scrollable list of cards:

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "List Example"
}

{
  "type": "updateComponents",
  "surfaceId": "main",
  "components": [
    {
      "componentId": "root",
      "type": "List",
      "props": {
        "gap": 8,
        "padding": 12
      },
      "children": [
        "item1",
        "item2",
        "item3"
      ]
    },
    {
      "componentId": "item1",
      "type": "Card",
      "props": {
        "title": "Item 1"
      },
      "children": [
        "t1"
      ]
    },
    {
      "componentId": "t1",
      "type": "Text",
      "props": {
        "content": "First item content"
      }
    },
    {
      "componentId": "item2",
      "type": "Card",
      "props": {
        "title": "Item 2"
      },
      "children": [
        "t2"
      ]
    },
    {
      "componentId": "t2",
      "type": "Text",
      "props": {
        "content": "Second item content"
      }
    },
    {
      "componentId": "item3",
      "type": "Card",
      "props": {
        "title": "Item 3"
      },
      "children": [
        "t3"
      ]
    },
    {
      "componentId": "t3",
      "type": "Text",
      "props": {
        "content": "Third item content"
      }
    }
  ]
}
```

## Notes

- List scrolls automatically when content exceeds the available height.
- Children are laid out vertically in order, just like Column.
- Use List instead of Column when you expect the content to grow beyond the visible area.
- List supports dynamic children -- add or remove items by updating the `children` array.
