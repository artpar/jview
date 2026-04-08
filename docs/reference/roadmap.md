---
layout: default
title: Roadmap
parent: Reference
nav_order: 2
---

# Roadmap

Canopy's development history, organized by phase. All phases are complete.

## Phase 1: Core MVP

The foundation -- protocol parsing, engine core, and the first 7 component bridges.

| Feature | Description |
|---------|-------------|
| A2UI JSONL protocol parser | Parse all message types from JSONL streams |
| Engine core | Session, Surface, Tree, DataModel, BindingTracker |
| Text | Styled text with typography variants (h1-h6, body, caption) |
| Row / Column | Horizontal and vertical stack layouts |
| Card | Container with title, border, shadow |
| Button | Clickable button with label and onClick action |
| TextField | Text input with label, placeholder, and data binding |
| CheckBox | Toggle with label and onToggle binding |
| File transport | Read JSONL from files |
| Makefile verify pipeline | Build, launch fixtures, capture screenshots |

## Phase 2: Full Interactivity

FunctionCall evaluator with 17 built-in functions, validation, template expansion, and 7 new components.

| Feature | Description |
|---------|-------------|
| FunctionCall evaluator | 17 functions including array ops (append, removeLast, slice) |
| Validation | Field-level and form-level validation |
| Template expansion | Dynamic string interpolation from data model |
| Divider | Horizontal/vertical separator |
| Icon | SF Symbol rendering |
| Image | Image display from URL or file path |
| Slider | Continuous value slider |
| ChoicePicker | Segmented control / dropdown |
| DateTimeInput | Date and time picker |
| List | Dynamic list with data-bound items |

## Phase 3: Media + Live Transport + Polish

LLM connectivity, native testing, visual styling, and FFI.

| Feature | Description |
|---------|-------------|
| LLM transport | Live agent connectivity via any-llm-go |
| Action response pipeline | Two-way communication with LLM agents |
| Native e2e test framework | Real AppKit assertions on component layout and style |
| Tabs | Tabbed container with tab switching |
| Modal | Modal dialog windows |
| Video | AVPlayerView for video playback |
| AudioPlayer | Audio playback controls |
| Theme support | NSAppearance light/dark theme switching |
| Scroll view | Automatic scroll for overflow content |
| Embedded MCP server | 51 tools on stdin/stdout |

## Phase 4: Production Hardening

Reliability, multi-process architecture, and always-on MCP.

| Feature | Description |
|---------|-------------|
| CGo memory cleanup | Proper cleanup of Objective-C objects and CGo handles |
| Error recovery | Graceful degradation on protocol errors and rendering failures |
| Process model | createProcess / stopProcess for multi-process apps |
| Channel primitives | Pub/sub with broadcast and queue modes |
| Always-on MCP server | MCP available at all times, not just during testing |

## Media Capture

Camera, microphone, and screen capture -- 2 new components, 6 evaluator functions, 8 MCP tools.

| Feature | Description |
|---------|-------------|
| CameraView | Live camera preview using AVCaptureSession |
| AudioRecorder | Recording UI with level meter |
| Headless camera capture | Take photos without a visible CameraView |
| Headless audio recording | Record audio without an AudioRecorder component |
| Screen capture | Screenshot via ScreenCaptureKit |
| Info.plist privacy descriptions | Camera and microphone permission prompts |

## Notes Clone

4 new native components and a full Apple Notes sample app.

| Feature | Description |
|---------|-------------|
| SplitView | Resizable multi-pane layout (NSSplitView) |
| OutlineView | Hierarchical tree view (NSOutlineView) |
| SearchField | Native search input (NSSearchField) |
| RichTextEditor | Rich text editing (NSTextView) |
| filter / find / getField functions | Array and object query functions |
| Notes sample app | 3-pane layout with folders, note list, and editor |

## Package Management

GitHub-based package ecosystem.

| Feature | Description |
|---------|-------------|
| Package manifest | `canopy.json` format for apps, components, themes, FFI configs |
| Install / uninstall | Download and manage packages from GitHub Releases |
| Publish | Create tags and releases, set `canopy-package` topic |
| Search / browse | Discover packages via GitHub topic search |
| CLI | Full `canopy pkg` command suite |
| MCP tools | 9 package management tools for programmatic access |
