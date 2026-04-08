---
layout: default
title: Protocol Reference
nav_order: 5
has_children: true
---

# Protocol Reference

Canopy apps communicate via the **A2UI protocol** -- a stream of JSON messages, one per line (JSONL). Every message has a `type` field that determines its structure and behavior.

## Wire Format

- One JSON object per line (no pretty-printing across lines)
- Maximum 10 MB per line
- UTF-8 encoding
- Messages are processed in order

## Message Types

### Surface Management

| Message | Description |
|---------|-------------|
| [createSurface](create-surface) | Open a new window |
| [deleteSurface](delete-surface) | Close a window |
| [updateWindow](update-window) | Modify window properties |
| [setTheme](set-theme) | Switch light/dark/system theme |
| [setAppMode](set-app-mode) | Switch normal/menubar/accessory mode |

### Components

| Message | Description |
|---------|-------------|
| [updateComponents](update-components) | Create or update components in a surface |
| [defineComponent](define-component) | Register a reusable component template |

### Data

| Message | Description |
|---------|-------------|
| [updateDataModel](update-data-model) | Apply JSON Patch operations to the data model |
| [defineFunction](define-function) | Register a reusable parametric function |

### UI Chrome

| Message | Description |
|---------|-------------|
| [updateMenu](update-menu) | Define the menu bar |
| [updateToolbar](update-toolbar) | Define the window toolbar |

### Processes and Channels

| Message | Description |
|---------|-------------|
| [createProcess](create-process) | Spawn a background process |
| [stopProcess](stop-process) | Terminate a process |
| [sendToProcess](send-to-process) | Route a message to a process |
| [createChannel](create-channel) | Create a pub/sub channel |
| [deleteChannel](delete-channel) | Remove a channel |
| [subscribe](subscribe) | Subscribe to channel values |
| [unsubscribe](unsubscribe) | Remove a subscription |
| [publish](publish-message) | Send a value to a channel |

### Advanced

| Message | Description |
|---------|-------------|
| [loadLibrary](load-library) | Load a native C library via FFI |
| [include](include) | Include another JSONL file |
| [test](test-message) | Inline testing with assertions |
