---
layout: default
title: Text
parent: Components
nav_order: 15
---

# Text

Read-only text display. Maps to **NSTextField** (non-editable).

Renders static or dynamic text with variant-based styling. Use Text for headings, labels, paragraphs, and any read-only content.

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `content` | DynamicString | | The text to display |
| `variant` | string | `"body"` | Typography style: `h1`, `h2`, `h3`, `h4`, `h5`, `body`, `caption` |

## Example

Heading and body text:

```json
{"type":"createSurface","surfaceId":"main","title":"Text Example"}
{"type":"updateComponents","surfaceId":"main","components":[
  {"componentId":"root","type":"Column","props":{"padding":16,"gap":8},"children":["heading","body","caption"]},
  {"componentId":"heading","type":"Text","props":{"content":"Welcome to Canopy","variant":"h1"}},
  {"componentId":"body","type":"Text","props":{"content":"Canopy renders native macOS components from JSONL."}},
  {"componentId":"caption","type":"Text","props":{"content":"Built with AppKit","variant":"caption"}}
]}
```

## Variant Sizes

| Variant | Usage | Font |
|---------|-------|------|
| `h1` | Page title | System bold, 28pt |
| `h2` | Section heading | System bold, 22pt |
| `h3` | Subsection heading | System semibold, 18pt |
| `h4` | Group heading | System semibold, 15pt |
| `h5` | Minor heading | System medium, 13pt |
| `body` | Body text (default) | System regular, 13pt |
| `caption` | Secondary text | System regular, 11pt, gray |

## Notes

- The `content` prop supports dynamic values. Use `{"$ref":"/path"}` to display data model values.
- Use `{"$concat":["Hello, ",{"$ref":"/name"},"!"]}` to build composite strings.
- Text wraps automatically when it exceeds the available width.
