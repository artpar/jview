---
layout: default
title: SplitView
parent: Components
nav_order: 4
---

# SplitView

Resizable panes. Maps to **NSSplitView**.

Splits the available space between two or more child components with draggable dividers. Ideal for master-detail layouts like file browsers, email clients, or note-taking apps.

![SplitView in a notes app]({{ site.baseurl }}/screenshots/notes.png){: .screenshot}

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `dividerStyle` | string | `"thin"` | Divider appearance: `thin`, `thick`, `paneSplitter` |
| `vertical` | DynamicBoolean | `true` | `true` for side-by-side panes, `false` for top/bottom |
| `collapsedPane` | DynamicNumber | `-1` | Index of the collapsed pane, or `-1` for none |

## Example

Two-pane layout with a sidebar and content area:

```json
{"type":"createSurface","surfaceId":"main","title":"SplitView Example"}
{"type":"updateComponents","surfaceId":"main","components":[
  {"componentId":"root","type":"SplitView","props":{"vertical":true,"dividerStyle":"thin"},"children":["sidebar","content"]},
  {"componentId":"sidebar","type":"Column","props":{"padding":8,"gap":4},"children":["nav1","nav2","nav3"]},
  {"componentId":"nav1","type":"Text","props":{"content":"Inbox"}},
  {"componentId":"nav2","type":"Text","props":{"content":"Drafts"}},
  {"componentId":"nav3","type":"Text","props":{"content":"Archive"}},
  {"componentId":"content","type":"Column","props":{"padding":16},"children":["detail"]},
  {"componentId":"detail","type":"Text","props":{"content":"Select an item from the sidebar.","variant":"h2"}}
]}
```

## Notes

- The first child becomes the left (or top) pane; the second child becomes the right (or bottom) pane.
- Users can drag the divider to resize panes.
- Set `collapsedPane` to `0` or `1` to programmatically collapse a pane.
- The `paneSplitter` divider style shows a thicker, more visible handle.
