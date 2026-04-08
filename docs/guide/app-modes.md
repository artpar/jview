---
layout: default
title: App Modes
parent: Building Apps
nav_order: 10
---

# App Modes

Canopy apps can run in three activation modes, controlled by the `setAppMode` message.

## Modes

| Mode | Dock Icon | Menu Bar | Windows | Behavior |
|------|-----------|----------|---------|----------|
| `normal` | Yes | No | Yes | Standard macOS app (default) |
| `menubar` | No | Yes | Yes | Status bar item; clicking toggles window visibility |
| `accessory` | No | No | Yes | Background process; no dock or menu bar presence |

## Normal Mode

The default. Your app appears in the Dock and behaves like a standard macOS application.

## Menubar Mode

Creates an NSStatusItem in the system menu bar. The app disappears from the Dock. Clicking the status item toggles the window.

```json
{
  "type": "setAppMode",
  "mode": "menubar",
  "icon": "bolt.fill",
  "title": "My App"
}
```

Fields:
- `icon` -- an SF Symbol name for the status bar icon
- `title` -- text to show in the status bar (used if no icon, or alongside it)
- `menuItems` -- (optional) array of menu items for the status bar dropdown

The app stays alive even when all windows are closed. Users click the menu bar icon to show/hide the window.

### Menu Items

Add dropdown items to the status bar menu:

```json
{
  "type": "setAppMode",
  "mode": "menubar",
  "icon": "bolt.fill",
  "menuItems": [
    {
      "id": "show",
      "label": "Show Window"
    },
    {
      "separator": true
    },
    {
      "id": "quit",
      "label": "Quit",
      "keyEquivalent": "q"
    }
  ]
}
```

## Accessory Mode

The app has no Dock icon and no menu bar item. It runs entirely in the background. Windows are still visible if created.

```json
{
  "type": "setAppMode",
  "mode": "accessory"
}
```

Use this for apps that are driven entirely by MCP tools, processes, or channels.

## Switching Modes

You can switch modes at any time:

```json
{
  "type": "setAppMode",
  "mode": "normal"
}
```

## Example: Menubar Timer

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "Timer",
  "width": 300,
  "height": 150
}

{
  "type": "setAppMode",
  "mode": "menubar",
  "icon": "timer",
  "title": "Timer"
}

{
  "type": "updateDataModel",
  "surfaceId": "main",
  "ops": [
    {
      "op": "add",
      "path": "/seconds",
      "value": 0
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
        "gap": 12,
        "align": "center"
      },
      "children": [
        "time"
      ]
    },
    {
      "componentId": "time",
      "type": "Text",
      "props": {
        "content": {
          "functionCall": {
            "name": "concat",
            "args": [
              {
                "path": "/seconds"
              },
              "s"
            ]
          }
        },
        "variant": "h1"
      }
    }
  ]
}

{
  "type": "createProcess",
  "processId": "tick",
  "transport": {
    "type": "interval",
    "interval": 1000,
    "message": {
      "type": "updateDataModel",
      "surfaceId": "main",
      "ops": [
        {
          "op": "replace",
          "path": "/seconds",
          "value": {
            "functionCall": {
              "name": "add",
              "args": [
                {
                  "path": "/seconds"
                },
                1
              ]
            }
          }
        }
      ]
    }
  }
}
```

The timer icon appears in the menu bar. Clicking it toggles the timer window.
