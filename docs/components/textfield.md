---
layout: default
title: TextField
parent: Components
nav_order: 8
---

# TextField

Text input field. Maps to **NSTextField**.

Accepts user text input with support for different input modes, placeholder text, and two-way data binding.

![TextField in a contact form]({{ site.baseurl }}/screenshots/contact_form.png){: .screenshot}

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `placeholder` | DynamicString | | Placeholder text shown when empty |
| `value` | DynamicString | | Current text value |
| `inputType` | string | `"shortText"` | Input mode: `shortText`, `longText`, `number`, `obscured` |
| `readOnly` | DynamicBoolean | `false` | Whether the field is non-editable |
| `dataBinding` | string | | JSON Pointer path for two-way binding |
| `onChange` | EventAction | | Action triggered when the value changes |

## Example

Name input with data binding:

```json
{"type":"createSurface","surfaceId":"main","title":"TextField Example"}
{"type":"updateDataModel","surfaceId":"main","operations":[
  {"op":"replace","path":"/name","value":""}
]}
{"type":"updateComponents","surfaceId":"main","components":[
  {"componentId":"root","type":"Column","props":{"padding":16,"gap":8},"children":["label","input","preview"]},
  {"componentId":"label","type":"Text","props":{"content":"Name","variant":"h3"}},
  {"componentId":"input","type":"TextField","props":{"placeholder":"Enter your name","dataBinding":"/name"}},
  {"componentId":"preview","type":"Text","props":{"content":{"$ref":"/name"}}}
]}
```

## Notes

- `shortText` renders a single-line input. `longText` renders a multi-line text area.
- `number` restricts input to numeric values. `obscured` masks input for passwords.
- When `dataBinding` is set, typing in the field writes to the data model, and data model changes update the field.
- The `onChange` action fires on every keystroke. Use it to trigger validation or live filtering.
- The `value` prop sets the initial text. If `dataBinding` is also set, the bound value takes precedence.
