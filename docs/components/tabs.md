---
layout: default
title: Tabs
parent: Components
nav_order: 5
---

# Tabs

Tabbed container. Maps to **NSTabView**.

Displays multiple panels with tab labels at the top. Only one tab is visible at a time. Supports data binding for the active tab.

![Tabs component]({{ site.baseurl }}/screenshots/tabs.png){: .screenshot}

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `tabLabels` | string[] | | Array of tab label strings, one per child |
| `activeTab` | DynamicString | | Label of the currently active tab |
| `dataBinding` | string | | JSON Pointer path to bind the active tab label |

## Example

Settings panel with two tabs:

```json
{"type":"createSurface","surfaceId":"main","title":"Tabs Example"}
{"type":"updateComponents","surfaceId":"main","components":[
  {"componentId":"root","type":"Tabs","props":{"tabLabels":["General","Advanced"]},"children":["general","advanced"]},
  {"componentId":"general","type":"Column","props":{"padding":12,"gap":8},"children":["g1"]},
  {"componentId":"g1","type":"Text","props":{"content":"General settings go here."}},
  {"componentId":"advanced","type":"Column","props":{"padding":12,"gap":8},"children":["a1"]},
  {"componentId":"a1","type":"Text","props":{"content":"Advanced settings go here."}}
]}
```

## Notes

- The number of `tabLabels` entries must match the number of children.
- Each child component becomes the content of the corresponding tab.
- When `dataBinding` is set, switching tabs writes the active tab label to the data model, and updating the data model switches the visible tab.
- Tab content is preserved when switching between tabs.
