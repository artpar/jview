---
layout: default
title: Styling
parent: Building Apps
nav_order: 4
---

# Styling

Every component accepts a `style` object that controls its visual appearance. All style properties accept dynamic values (literals, path references, or function calls).

## Style Properties

| Property | Type | Description |
|----------|------|-------------|
| `backgroundColor` | string | Background color (hex `#RRGGBB`, named colors) |
| `textColor` | string | Text/foreground color |
| `cornerRadius` | number | Corner rounding in points |
| `width` | number | Fixed width in points |
| `height` | number | Fixed height in points |
| `fontSize` | number | Font size in points |
| `fontWeight` | string | `"bold"`, `"medium"`, `"light"` |
| `textAlign` | string | `"left"`, `"center"`, `"right"` |
| `opacity` | number | 0.0 (transparent) to 1.0 (opaque) |
| `flexGrow` | number | Expand to fill available space in parent stack |
| `fontFamily` | string | Font family name |

## Example: Styled Button

```json
{"componentId":"submit","type":"Button","props":{
  "label": "Submit",
  "style": "primary"
},"style":{
  "backgroundColor": "#007AFF",
  "cornerRadius": 8,
  "fontSize": 16,
  "fontWeight": "bold"
}}
```

## Dynamic Styles

Style values can reference the data model:

```json
{"componentId":"status","type":"Text","props":{
  "content": {"path": "/message"}
},"style":{
  "textColor": {"functionCall": {
    "name": "if",
    "args": [
      {"functionCall": {"name": "equals", "args": [{"path": "/status"}, "error"]}},
      "#FF3B30",
      "#34C759"
    ]
  }}
}}
```

## Surface-Level Styling

The `createSurface` message accepts `backgroundColor` and `padding` for the window itself:

```json
{"type":"createSurface","surfaceId":"main","title":"Styled App",
  "width":600,"height":400,
  "backgroundColor":"#1E1E1E",
  "padding":20
}
```

## Flexible Layouts with flexGrow

Use `flexGrow` to make a component expand within its parent Row or Column:

```json
{"type":"updateComponents","surfaceId":"main","components":[
  {"componentId":"row","type":"Row","props":{"gap":8},"children":["sidebar","content"]},
  {"componentId":"sidebar","type":"Column","children":["nav"],"style":{"width":200}},
  {"componentId":"content","type":"Column","children":["main"],"style":{"flexGrow":1}}
]}
```

The sidebar has a fixed 200pt width. The content area expands to fill the remaining space.

## Layout Props vs Style

Layout properties like `gap`, `padding`, `justify`, and `align` live in `props`, not `style`. The `style` object is for visual properties only.

```json
{
  "componentId": "container",
  "type": "Column",
  "props": {"gap": 16, "padding": 20, "align": "center"},
  "style": {"backgroundColor": "#F5F5F5", "cornerRadius": 12}
}
```
