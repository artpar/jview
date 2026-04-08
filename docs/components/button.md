---
layout: default
title: Button
parent: Components
nav_order: 14
---

# Button

Clickable action trigger. Maps to **NSButton**.

A standard macOS button that triggers an action when clicked. Supports primary, secondary, and destructive styles.

![Button in a validation form]({{ site.baseurl }}/screenshots/validation.png){: .screenshot}

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `label` | DynamicString | | Button text |
| `style` | string | `"secondary"` | Visual style: `primary`, `secondary`, `destructive` |
| `disabled` | DynamicBoolean | `false` | Whether the button is grayed out and non-clickable |
| `onClick` | EventAction | | Action triggered when clicked |

## Example

Submit button that sets a value:

```json
{"type":"createSurface","surfaceId":"main","title":"Button Example"}
{"type":"updateDataModel","surfaceId":"main","operations":[
  {"op":"replace","path":"/submitted","value":false}
]}
{"type":"updateComponents","surfaceId":"main","components":[
  {"componentId":"root","type":"Column","props":{"padding":16,"gap":12},"children":["heading","submitBtn"]},
  {"componentId":"heading","type":"Text","props":{"content":"Ready to submit?","variant":"h2"}},
  {"componentId":"submitBtn","type":"Button","props":{
    "label":"Submit",
    "style":"primary",
    "onClick":{"action":{"setValues":[{"path":"/submitted","value":true}]}}
  }}
]}
```

## Notes

- `primary` renders as a prominent blue button (the default action).
- `secondary` renders as a standard gray button.
- `destructive` renders as a red button for dangerous actions.
- The `disabled` prop can reference a data model value using `{"$ref":"/path"}` for dynamic enable/disable.
- `onClick` can trigger multiple actions: `setValues` to update data, `event` to send actions to the transport, or `functionCall` to invoke evaluator functions.
