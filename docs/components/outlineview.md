---
layout: default
title: OutlineView
parent: Components
nav_order: 21
---

# OutlineView

Hierarchical tree list. Maps to **NSOutlineView**.

Displays a tree structure with expandable/collapsible nodes. Ideal for file browsers, navigation sidebars, and any hierarchical data.

![OutlineView in a notes app]({{ site.baseurl }}/screenshots/notes.png){: .screenshot}

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `outlineData` | DynamicString | | JSON string containing the tree data |
| `labelKey` | string | `"name"` | Key in each node for the display label |
| `childrenKey` | string | `"children"` | Key in each node for the children array |
| `iconKey` | string | | Key in each node for an SF Symbol icon name |
| `idKey` | string | `"id"` | Key in each node for the unique identifier |
| `selectedId` | DynamicString | | ID of the currently selected node |
| `badgeKey` | string | | Key in each node for a badge count |
| `dataBinding` | string | | JSON Pointer path to bind the selected ID |
| `onSelect` | EventAction | | Action triggered when a node is selected |

## Example

Folder tree:

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "OutlineView Example"
}

{
  "type": "updateDataModel",
  "surfaceId": "main",
  "operations": [
    {
      "op": "replace",
      "path": "/selectedFolder",
      "value": ""
    },
    {
      "op": "replace",
      "path": "/folders",
      "value": "[{\"id\":\"inbox\",\"name\":\"Inbox\",\"icon\":\"tray\",\"badge\":3},{\"id\":\"projects\",\"name\":\"Projects\",\"icon\":\"folder\",\"children\":[{\"id\":\"p1\",\"name\":\"Website Redesign\",\"icon\":\"doc\"},{\"id\":\"p2\",\"name\":\"Mobile App\",\"icon\":\"doc\"}]},{\"id\":\"archive\",\"name\":\"Archive\",\"icon\":\"archivebox\"}]"
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
        "padding": 0
      },
      "children": [
        "tree"
      ]
    },
    {
      "componentId": "tree",
      "type": "OutlineView",
      "props": {
        "outlineData": {
          "$ref": "/folders"
        },
        "labelKey": "name",
        "childrenKey": "children",
        "iconKey": "icon",
        "idKey": "id",
        "badgeKey": "badge",
        "dataBinding": "/selectedFolder"
      }
    }
  ]
}
```

## Notes

- The `outlineData` prop expects a JSON string containing an array of node objects.
- Nodes with a `children` array are expandable. Leaf nodes have no children key.
- When `dataBinding` is set, selecting a node writes its ID to the data model.
- The `badgeKey` displays a count badge on the right side of the row (e.g., unread count).
- The `iconKey` displays an SF Symbol icon to the left of the label.
- Nodes are expanded by default. Click the disclosure triangle to collapse.
