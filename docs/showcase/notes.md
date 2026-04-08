---
layout: default
title: Notes
parent: Showcase
nav_order: 2
---

# Notes

A 3-pane Notes app inspired by Apple Notes. Folder sidebar, note list, and rich text editor with search and formatting toolbar.

![Notes](../screenshots/notes.png)

## Prompt

> Build a 3-pane Notes app with SplitView: folder sidebar (OutlineView), note list, rich text editor. Search field, toolbar with formatting.

## Key Features

- **SplitView** -- 3-pane resizable layout using native NSSplitView
- **OutlineView** -- hierarchical folder tree with expand/collapse in the sidebar
- **RichTextEditor** -- full rich text editing with NSTextView, supporting bold, italic, underline, and lists
- **SearchField** -- native NSSearchField that filters the note list in real time
- **Toolbar** -- formatting controls for the editor

## How to Run

```bash
build/canopy sample_apps/notes
```

## What to Look For

- The sidebar shows a folder hierarchy you can expand and collapse
- Selecting a folder filters the note list in the middle pane
- Selecting a note opens it in the rich text editor
- The search field at the top filters notes across all folders
- Formatting toolbar buttons apply styles to selected text in the editor
