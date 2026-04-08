---
layout: default
title: Card
parent: Components
nav_order: 3
---

# Card

Titled container with optional collapse. Maps to **NSBox**.

Wraps child components in a bordered box with a title bar. Cards visually group related content and can be made collapsible for dense interfaces.

![Card in a contact form]({{ site.baseurl }}/screenshots/contact_form.png){: .screenshot}

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `title` | DynamicString | | Title displayed in the box header |
| `subtitle` | DynamicString | | Subtitle displayed below the title |
| `collapsible` | DynamicBoolean | `false` | Whether the card can be collapsed |
| `collapsed` | DynamicBoolean | `false` | Whether the card starts collapsed |
| `padding` | int | `0` | Internal padding in points |

## Example

A card containing user profile info:

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "Card Example"
}

{
  "type": "updateComponents",
  "surfaceId": "main",
  "components": [
    {
      "componentId": "root",
      "type": "Column",
      "props": {
        "padding": 16
      },
      "children": [
        "profile"
      ]
    },
    {
      "componentId": "profile",
      "type": "Card",
      "props": {
        "title": "User Profile",
        "padding": 12
      },
      "children": [
        "name",
        "email"
      ]
    },
    {
      "componentId": "name",
      "type": "Text",
      "props": {
        "content": "Jane Smith",
        "variant": "h3"
      }
    },
    {
      "componentId": "email",
      "type": "Text",
      "props": {
        "content": "jane@example.com"
      }
    }
  ]
}
```

## Notes

- Card titles support dynamic values bound to the data model.
- When `collapsible` is true, clicking the title bar toggles visibility of the card body.
- The `collapsed` prop can be bound to the data model to programmatically control collapse state.
- Cards can be nested inside other layout containers (Row, Column, List).
