---
layout: default
title: Slider
parent: Components
nav_order: 10
---

# Slider

Range input. Maps to **NSSlider**.

A horizontal slider for selecting a numeric value within a range. Supports snapping to discrete steps and data binding.

![Slider component]({{ site.baseurl }}/screenshots/slider.png){: .screenshot}

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `min` | DynamicNumber | `0` | Minimum value |
| `max` | DynamicNumber | `100` | Maximum value |
| `step` | DynamicNumber | `1` | Step increment |
| `sliderValue` | DynamicNumber | `0` | Current value |
| `dataBinding` | string | | JSON Pointer path for two-way binding |

## Example

Font size slider with live preview:

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "Slider Example"
}

{
  "type": "updateDataModel",
  "surfaceId": "main",
  "operations": [
    {
      "op": "replace",
      "path": "/fontSize",
      "value": 14
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
        "label",
        "slider",
        "preview"
      ]
    },
    {
      "componentId": "label",
      "type": "Text",
      "props": {
        "content": "Font Size",
        "variant": "h3"
      }
    },
    {
      "componentId": "slider",
      "type": "Slider",
      "props": {
        "min": 8,
        "max": 72,
        "step": 1,
        "dataBinding": "/fontSize"
      }
    },
    {
      "componentId": "preview",
      "type": "Text",
      "props": {
        "content": {
          "$concat": [
            "Size: ",
            {
              "$ref": "/fontSize"
            },
            "pt"
          ]
        }
      }
    }
  ]
}
```

## Notes

- The slider value snaps to the nearest `step` increment.
- When `dataBinding` is set, dragging the slider updates the data model in real time.
- Other components can reference the bound value to react to slider changes.
- Use `min` and `max` to constrain the selectable range.
