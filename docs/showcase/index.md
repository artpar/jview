---
layout: default
title: Showcase
nav_order: 9
has_children: true
---

# Showcase

A gallery of sample apps built with Canopy. Each app was generated from a text prompt and runs as a native macOS application -- no webview, no Electron.

These samples demonstrate the range of what you can build: from simple calculators to full productivity apps with rich text editing, multi-process architectures, and native library integration.

## Running a Sample

Every sample app lives in the `sample_apps/` directory. Run any of them with:

```bash
build/canopy sample_apps/<name>
```

## The Apps

| App | What it demonstrates |
|-----|---------------------|
| [Calculator](calculator) | Custom components, functions, grid layout |
| [Notes](notes) | SplitView, OutlineView, RichTextEditor, SearchField |
| [Todo](todo) | List, data binding, dynamic children |
| [Dashboard](dashboard) | Row/Column layout, Card, nested components |
| [Channel Demo](channel-demo) | Multi-process communication via channels |
| [System Info](sysinfo) | FFI -- calling native C libraries |
