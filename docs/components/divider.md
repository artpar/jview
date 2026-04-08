---
layout: default
title: Divider
parent: Components
nav_order: 18
---

# Divider

Horizontal separator line. Maps to **NSBox** (separator style).

A thin line that visually separates sections of content. Has no props -- just place it between other components.

![Divider component]({{ site.baseurl }}/screenshots/divider.png){: .screenshot}

## Props

None. Divider has no configurable properties.

## Example

Divider between content sections:

```json
{"type":"createSurface","surfaceId":"main","title":"Divider Example"}
{"type":"updateComponents","surfaceId":"main","components":[
  {"componentId":"root","type":"Column","props":{"padding":16,"gap":12},"children":["heading","divider","body"]},
  {"componentId":"heading","type":"Text","props":{"content":"Section Title","variant":"h2"}},
  {"componentId":"divider","type":"Divider"},
  {"componentId":"body","type":"Text","props":{"content":"Content below the divider."}}
]}
```

## Notes

- Divider stretches to fill the width of its parent container.
- Use Divider to visually separate groups of related components.
- In a Row parent, the Divider renders as a vertical separator.
