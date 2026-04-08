---
layout: default
title: RichTextEditor
parent: Components
nav_order: 20
---

# RichTextEditor

Rich text editor with markdown support. Maps to **NSTextView**.

A full-featured text editor that renders markdown formatting in real time. Supports headings, bold, italic, strikethrough, code, checklists, bullets, and numbered lists.

![RichTextEditor in a notes app]({{ site.baseurl }}/screenshots/notes.png){: .screenshot}

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `content` | DynamicString | | Initial markdown content |
| `editable` | DynamicBoolean | `true` | Whether the text is editable |
| `dataBinding` | string | | JSON Pointer path to bind the markdown text |
| `formatBinding` | string | | JSON Pointer path to bind the attributed/formatted text |
| `onChange` | EventAction | | Action triggered when content changes |

## Supported Markdown

| Syntax | Result |
|--------|--------|
| `# Heading` | Heading levels 1-6 |
| `**bold**` | **Bold** text |
| `*italic*` | *Italic* text |
| `~~strike~~` | ~~Strikethrough~~ text |
| `` `code` `` | Inline code |
| `- item` | Bullet list |
| `1. item` | Numbered list |
| `- [ ] task` | Checklist item |
| `- [x] done` | Checked checklist item |

## Example

Note editor:

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "RichTextEditor Example"
}

{
  "type": "updateDataModel",
  "surfaceId": "main",
  "operations": [
    {
      "op": "replace",
      "path": "/noteContent",
      "value": "# My Note\n\nThis is **bold** and *italic* text.\n\n- [ ] First task\n- [x] Completed task"
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
        "padding": 16
      },
      "children": [
        "editor"
      ]
    },
    {
      "componentId": "editor",
      "type": "RichTextEditor",
      "props": {
        "dataBinding": "/noteContent",
        "editable": true
      }
    }
  ]
}
```

## Notes

- Markdown is rendered as styled attributed text in real time as you type.
- The `dataBinding` stores the raw markdown string.
- The `formatBinding` stores the rich/attributed text representation separately from the raw markdown.
- Set `editable` to `false` for a read-only markdown viewer.
- Checklists are interactive -- click the checkbox to toggle completion state.
