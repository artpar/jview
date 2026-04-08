---
layout: default
title: DateTimeInput
parent: Components
nav_order: 12
---

# DateTimeInput

Date and time picker. Maps to **NSDatePicker**.

A native macOS date picker that can show date, time, or both. Uses the system locale for formatting.

![DateTimeInput component]({{ site.baseurl }}/screenshots/datetimeinput.png){: .screenshot}

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `enableDate` | DynamicBoolean | `true` | Show date selection |
| `enableTime` | DynamicBoolean | `false` | Show time selection |
| `dataBinding` | string | | JSON Pointer path for two-way binding |

## Example

Event date picker:

```json
{"type":"createSurface","surfaceId":"main","title":"DateTimeInput Example"}
{"type":"updateDataModel","surfaceId":"main","operations":[
  {"op":"replace","path":"/eventDate","value":""}
]}
{"type":"updateComponents","surfaceId":"main","components":[
  {"componentId":"root","type":"Column","props":{"padding":16,"gap":8},"children":["label","picker"]},
  {"componentId":"label","type":"Text","props":{"content":"Event Date","variant":"h3"}},
  {"componentId":"picker","type":"DateTimeInput","props":{"enableDate":true,"enableTime":true,"dataBinding":"/eventDate"}}
]}
```

## Notes

- Set both `enableDate` and `enableTime` to `true` for a full date-time picker.
- Set only `enableTime` to `true` for a time-only picker.
- The bound value is stored as an ISO 8601 string in the data model.
- The picker follows the user's system locale and date format preferences.
