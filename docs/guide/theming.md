---
layout: default
title: Theming
parent: Building Apps
nav_order: 9
---

# Theming

Canopy supports light, dark, and system-following themes. The theme controls NSAppearance for the entire application.

## Setting a Theme

Send a `setTheme` message:

```json
{"type":"setTheme","surfaceId":"main","theme":"dark"}
```

Valid values:
- `"light"` -- always light appearance
- `"dark"` -- always dark appearance
- `"system"` -- follow macOS system setting

## Theme on Surface Creation

Set the initial theme when creating a surface:

```json
{"type":"createSurface","surfaceId":"main","title":"My App","theme":"dark"}
```

## Theme Toggle Button

Use a `setTheme` functionCall action to let users switch themes:

```json
{"type":"updateComponents","surfaceId":"main","components":[
  {"componentId":"themeBtn","type":"Button","props":{
    "label":"Switch to Dark",
    "onClick":{"action":{"functionCall":{
      "call":"setTheme",
      "args":{"theme":"dark"}
    }}}
  }}
]}
```

## Example: Theme Switcher

A row of buttons that switch between themes:

```json
{"type":"createSurface","surfaceId":"main","title":"Theme Switcher","width":400,"height":200}
{"type":"updateComponents","surfaceId":"main","components":[
  {"componentId":"root","type":"Column","props":{"gap":16,"align":"center"},"children":["title","buttons"]},
  {"componentId":"title","type":"Text","props":{"content":"Choose Theme","variant":"h2"}},
  {"componentId":"buttons","type":"Row","props":{"gap":8},"children":["lightBtn","darkBtn","systemBtn"]},
  {"componentId":"lightBtn","type":"Button","props":{
    "label":"Light",
    "onClick":{"action":{"functionCall":{"call":"setTheme","args":{"theme":"light"}}}}
  }},
  {"componentId":"darkBtn","type":"Button","props":{
    "label":"Dark",
    "onClick":{"action":{"functionCall":{"call":"setTheme","args":{"theme":"dark"}}}}
  }},
  {"componentId":"systemBtn","type":"Button","props":{
    "label":"System",
    "onClick":{"action":{"functionCall":{"call":"setTheme","args":{"theme":"system"}}}}
  }}
]}
```
