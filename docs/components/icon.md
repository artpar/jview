---
layout: default
title: Icon
parent: Components
nav_order: 16
---

# Icon

SF Symbol icon. Maps to **NSImageView**.

Displays a system icon from Apple's SF Symbols library. Over 5,000 symbols are available covering common UI metaphors.

![Icon component]({{ site.baseurl }}/screenshots/icon.png){: .screenshot}

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `name` | DynamicString | | SF Symbol name (e.g., `"star.fill"`, `"folder"`, `"gear"`) |
| `size` | int | `16` | Icon size in points |

## Example

Icons in a toolbar:

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "Icon Example"
}

{
  "type": "updateComponents",
  "surfaceId": "main",
  "components": [
    {
      "componentId": "root",
      "type": "Row",
      "props": {
        "gap": 16,
        "padding": 12,
        "align": "center"
      },
      "children": [
        "home",
        "search",
        "settings"
      ]
    },
    {
      "componentId": "home",
      "type": "Icon",
      "props": {
        "name": "house.fill",
        "size": 20
      }
    },
    {
      "componentId": "search",
      "type": "Icon",
      "props": {
        "name": "magnifyingglass",
        "size": 20
      }
    },
    {
      "componentId": "settings",
      "type": "Icon",
      "props": {
        "name": "gear",
        "size": 20
      }
    }
  ]
}
```

## Common SF Symbols

| Symbol | Name |
|--------|------|
| House | `house.fill` |
| Gear | `gear` |
| Search | `magnifyingglass` |
| Plus | `plus` |
| Trash | `trash` |
| Star | `star.fill` |
| Folder | `folder` |
| Document | `doc.text` |
| Person | `person.fill` |
| Arrow right | `arrow.right` |

## Notes

- Browse all available symbols in Apple's [SF Symbols app](https://developer.apple.com/sf-symbols/).
- The `name` prop supports dynamic values for icons that change based on state.
- Icons inherit the system accent color by default.
- Pair Icon with Text inside a Row for labeled icon buttons.
