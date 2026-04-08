---
layout: default
title: ProgressBar
parent: Components
nav_order: 19
---

# ProgressBar

Progress indicator. Maps to **NSProgressIndicator**.

Displays a horizontal bar showing progress toward completion. Supports both determinate (known progress) and indeterminate (spinner) modes.

![ProgressBar component]({{ site.baseurl }}/screenshots/progressbar.png){: .screenshot}

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `value` | DynamicNumber | | Current progress value |
| `maxValue` | DynamicNumber | `100` | Maximum progress value |
| `indeterminate` | DynamicBoolean | `false` | Show a spinning/barberpole indicator instead of a fixed bar |

## Example

File upload progress:

```json
{"type":"createSurface","surfaceId":"main","title":"ProgressBar Example"}
{"type":"updateDataModel","surfaceId":"main","operations":[
  {"op":"replace","path":"/progress","value":35}
]}
{"type":"updateComponents","surfaceId":"main","components":[
  {"componentId":"root","type":"Column","props":{"padding":16,"gap":8},"children":["label","bar","status"]},
  {"componentId":"label","type":"Text","props":{"content":"Uploading...","variant":"h3"}},
  {"componentId":"bar","type":"ProgressBar","props":{"value":{"$ref":"/progress"},"maxValue":100}},
  {"componentId":"status","type":"Text","props":{"content":{"$concat":[{"$ref":"/progress"},"% complete"]},"variant":"caption"}}
]}
```

## Notes

- Update `value` through the data model to animate the progress bar.
- Set `indeterminate` to `true` when the total amount of work is unknown (e.g., loading data).
- The bar fills proportionally based on `value / maxValue`.
- The progress bar stretches to fill the width of its parent container.
