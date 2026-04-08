---
layout: default
title: MCP Tools
nav_order: 7
has_children: true
---

# MCP Tools

Canopy embeds a [Model Context Protocol](https://modelcontextprotocol.io/) (MCP) server on stdin/stdout, giving you 51 tools to control your app programmatically. Inspect components, click buttons, fill forms, read and write data, capture screenshots, manage processes and channels, and install packages -- all without touching the mouse.

## Connection Modes

| Mode | How to use | Best for |
|------|-----------|----------|
| **stdin/stdout** | Launch Canopy normally; MCP is always available | Claude Code integration via `.mcp.json` |
| **Dedicated** | `canopy mcp` | Standalone MCP server for other clients |
| **HTTP** | `canopy --mcp-http :8080` | Remote access, web-based tooling |

## Tool Categories

| Category | Tools | Purpose |
|----------|-------|---------|
| [Inspecting Your App](query-tools) | 6 | Read component trees, data models, layout frames, styles |
| [Interacting with Components](interaction-tools) | 4 | Click, fill, toggle, and send generic events |
| [Managing Data](data-tools) | 2 | Write to the data model, wait for conditions |
| [Sending Messages](transport-tools) | 1 | Inject raw A2UI JSONL messages |
| [Screenshots & Capture](capture-tools) | 1 | Capture window contents as PNG |
| [System Capabilities](system-tools) | 6 | Notifications, clipboard, file dialogs, URLs |
| [Media Capture](media-tools) | 6 | Camera, microphone, screen capture |
| [Process Management](process-tools) | 4 | Spawn, stop, and message child processes |
| [Channel Communication](channel-tools) | 6 | Pub/sub messaging between processes |
| [Package Management](package-tools) | 9 | Search, install, publish packages from GitHub |
| [Logging & Actions](logging-tools) | 3 | Query logs, poll actions, send AppKit selectors |

## Using with Claude Code

Canopy's `.mcp.json` is pre-configured. When you open the project in Claude Code, all 51 tools appear as `mcp__canopy__*` tools automatically. Use `ToolSearch` with query `"canopy"` to load them.

## Quick Example

```
# List all open windows
mcp__canopy__list_surfaces

# Get the full component tree of window "main"
mcp__canopy__get_tree(surface_id: "main")

# Click a button
mcp__canopy__click(surface_id: "main", component_id: "submit-btn")

# Read a value from the data model
mcp__canopy__get_data_model(surface_id: "main", path: "/form/name")

# Take a screenshot
mcp__canopy__take_screenshot(surface_id: "main")
```
