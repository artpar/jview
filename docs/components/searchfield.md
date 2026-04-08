---
layout: default
title: SearchField
parent: Components
nav_order: 13
---

# SearchField

Search input with clear button. Maps to **NSSearchField**.

A text field styled for search with a magnifying glass icon and a clear button. Ideal for filtering lists or triggering search queries.

![SearchField in a notes app]({{ site.baseurl }}/screenshots/notes.png){: .screenshot}

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `placeholder` | DynamicString | | Placeholder text shown when empty |
| `value` | DynamicString | | Current search text |
| `dataBinding` | string | | JSON Pointer path for two-way binding |
| `onChange` | EventAction | | Action triggered when the text changes |

## Example

Search field that filters a list:

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "SearchField Example"
}

{
  "type": "updateDataModel",
  "surfaceId": "main",
  "operations": [
    {
      "op": "replace",
      "path": "/query",
      "value": ""
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
        "search",
        "results"
      ]
    },
    {
      "componentId": "search",
      "type": "SearchField",
      "props": {
        "placeholder": "Search items...",
        "dataBinding": "/query"
      }
    },
    {
      "componentId": "results",
      "type": "List",
      "props": {
        "gap": 4
      },
      "children": [
        "r1",
        "r2"
      ]
    },
    {
      "componentId": "r1",
      "type": "Text",
      "props": {
        "content": "Result 1"
      }
    },
    {
      "componentId": "r2",
      "type": "Text",
      "props": {
        "content": "Result 2"
      }
    }
  ]
}
```

## Notes

- SearchField behaves like TextField but with native search styling.
- The clear button (x) appears when the field has text and clears it when clicked.
- Use `onChange` with a `filter` function call to dynamically filter displayed items.
- When `dataBinding` is set, typing updates the data model on every keystroke.
