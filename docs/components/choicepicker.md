---
layout: default
title: ChoicePicker
parent: Components
nav_order: 11
---

# ChoicePicker

Dropdown or segmented selector. Maps to **NSPopUpButton**.

Presents a list of options for the user to choose from. Renders as a native macOS popup button (dropdown menu).

![ChoicePicker in a theme switcher]({{ site.baseurl }}/screenshots/theme_switcher.png){: .screenshot}

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `options` | array | | Array of `{"value":"...","label":"..."}` objects |
| `dataBinding` | string | | JSON Pointer path for two-way binding |
| `mutuallyExclusive` | DynamicBoolean | `true` | Whether only one option can be selected |

## Example

Theme selector:

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "ChoicePicker Example"
}

{
  "type": "updateDataModel",
  "surfaceId": "main",
  "operations": [
    {
      "op": "replace",
      "path": "/theme",
      "value": "system"
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
        "gap": 8
      },
      "children": [
        "label",
        "picker"
      ]
    },
    {
      "componentId": "label",
      "type": "Text",
      "props": {
        "content": "Theme",
        "variant": "h3"
      }
    },
    {
      "componentId": "picker",
      "type": "ChoicePicker",
      "props": {
        "options": [
          {
            "value": "light",
            "label": "Light"
          },
          {
            "value": "dark",
            "label": "Dark"
          },
          {
            "value": "system",
            "label": "System"
          }
        ],
        "dataBinding": "/theme"
      }
    }
  ]
}
```

## Notes

- The `value` field is stored in the data model; the `label` field is displayed to the user.
- When `dataBinding` is set, selecting an option writes its `value` to the data model.
- The initial selection matches the current data model value.
- Set `mutuallyExclusive` to `false` for multi-select behavior (less common).
