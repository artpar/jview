---
layout: default
title: Drag and Drop
parent: Building Apps
nav_order: 12
---

# Drag and Drop

Any component can accept dropped files and text by adding the `onDrop` prop.

## Adding a Drop Target

Set `onDrop` to an action that fires when files or text are dropped:

```json
{"componentId":"dropZone","type":"Card","props":{
  "title": "Drop files here",
  "onDrop": {
    "action": {
      "event": {
        "name": "fileDrop"
      }
    }
  }
},"children":["dropText"]}
```

## Drop Data

When a drop occurs, the action receives a data object with:

| Field | Type | Description |
|-------|------|-------------|
| `paths` | string array | File paths of dropped files |
| `text` | string | Dropped text (if text was dropped) |

For an event action, this data is sent to the server alongside any `dataRefs`. For a functionCall action, the data is available in the action context.

## Example: File Drop Zone

```json
{"type":"createSurface","surfaceId":"main","title":"Drop Zone","width":400,"height":300}
{"type":"updateDataModel","surfaceId":"main","ops":[
  {"op":"add","path":"/droppedFile","value":"No file dropped yet"}
]}
{"type":"updateComponents","surfaceId":"main","components":[
  {"componentId":"root","type":"Column","props":{"gap":16,"padding":20},"children":["zone","result"]},
  {"componentId":"zone","type":"Card","props":{
    "title":"Drop Zone",
    "onDrop":{"action":{"event":{"name":"fileDrop"}}}
  },"children":["instructions"],"style":{"height":150}},
  {"componentId":"instructions","type":"Text","props":{
    "content":"Drag a file here","variant":"body"
  },"style":{"textColor":"#888888","textAlign":"center"}},
  {"componentId":"result","type":"Text","props":{
    "content":{"path":"/droppedFile"},"variant":"caption"
  }}
]}
```

## Accepted Drop Types

Canopy accepts two types of drops:

- **File URLs** -- files dragged from Finder or other apps. The `paths` array contains absolute file paths.
- **Plain text** -- text dragged from other apps or selected text. The `text` field contains the dropped string.

Both can be present in the same drop if the source provides both types.

## Drop on Any Component

The `onDrop` prop works on any component type, not just Card:

```json
{"componentId":"imageTarget","type":"Image","props":{
  "src": {"path": "/imagePath"},
  "onDrop": {"action": {"event": {"name": "imageDropped"}}}
},"style":{"width":200,"height":200}}
```

## Implementation Note

Canopy adds a transparent overlay view as an `NSDraggingDestination` on any component with `onDrop`. This overlay intercepts drag operations without affecting the component's normal appearance or behavior.
