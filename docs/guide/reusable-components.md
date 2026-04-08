---
layout: default
title: Reusable Components
parent: Building Apps
nav_order: 5
---

# Reusable Components

Use `defineComponent` to create a reusable component template, then instantiate it with `useComponent`.

## Defining a Component

A `defineComponent` message registers a template with a name, parameters, and a component tree:

```json
{"type":"defineComponent","name":"DigitButton","params":["digit"],
  "components":[
    {"componentId":"_root","type":"Button","props":{
      "label":{"param":"digit"},
      "onClick":{"action":{"functionCall":{
        "call":"updateDataModel",
        "args":{"ops":[
          {"op":"replace","path":"/display","value":{"functionCall":{"name":"concat","args":[{"path":"/display"},{"param":"digit"}]}}}
        ]}
      }}}
    },"style":{"width":60,"height":60,"fontSize":24}}
  ]
}
```

Key rules:

- **`_root`** is required. It becomes the top-level view of each instance.
- **`{"param": "name"}`** placeholders are replaced with actual values at instantiation.
- Other component IDs should start with `_` (e.g., `_label`, `_icon`). They get rewritten to avoid collisions between instances.

## Using a Component

Reference the template with `useComponent` and pass arguments:

```json
{"componentId":"btn7","useComponent":"DigitButton","args":{"digit":"7"}}
```

This creates a Button instance with:
- `componentId` "btn7" (replaces `_root`)
- Label "7"
- Click action that appends "7" to the display

## ID Rewriting

When a template has multiple internal components, IDs are rewritten to prevent collisions:

- `_root` becomes the instance's `componentId`
- `_label` becomes `instanceId__label`
- `_icon` becomes `instanceId__icon`

For example, with a template containing `_root`, `_label`, and `_icon`, and an instance `btn7`:

| Template ID | Instance ID |
|-------------|-------------|
| `_root` | `btn7` |
| `_label` | `btn7__label` |
| `_icon` | `btn7__icon` |

## State Scoping with $ Prefix

Parameters prefixed with `$` create scoped data model paths for each instance. This lets multiple instances of the same template have independent state:

```json
{"type":"defineComponent","name":"Counter","params":["$count","label"],
  "components":[
    {"componentId":"_root","type":"Column","children":["_text","_btn"]},
    {"componentId":"_text","type":"Text","props":{
      "content":{"functionCall":{"name":"concat","args":[{"param":"label"},": ",{"path":{"param":"$count"}}]}}
    }},
    {"componentId":"_btn","type":"Button","props":{
      "label":"+1",
      "onClick":{"action":{"functionCall":{"call":"updateDataModel","args":{"ops":[
        {"op":"replace","path":{"param":"$count"},"value":{"functionCall":{"name":"add","args":[{"path":{"param":"$count"}},1]}}}
      ]}}}}
    }}
  ]
}
```

Each instance passes its own data model path:

```json
{"componentId":"counterA","useComponent":"Counter","args":{"$count":"/counters/a","label":"Counter A"}}
{"componentId":"counterB","useComponent":"Counter","args":{"$count":"/counters/b","label":"Counter B"}}
```

## Full Example: Calculator Digit Buttons

```json
{"type":"defineComponent","name":"DigitButton","params":["digit"],
  "components":[
    {"componentId":"_root","type":"Button","props":{
      "label":{"param":"digit"},
      "onClick":{"action":{"functionCall":{"call":"updateDataModel","args":{"ops":[
        {"op":"replace","path":"/display","value":{"functionCall":{"name":"concat","args":[{"path":"/display"},{"param":"digit"}]}}}
      ]}}}}
    },"style":{"width":60,"height":60,"fontSize":24}}
  ]
}

{"type":"updateComponents","surfaceId":"main","components":[
  {"componentId":"row1","type":"Row","props":{"gap":4},"children":["btn7","btn8","btn9"]},
  {"componentId":"btn7","useComponent":"DigitButton","args":{"digit":"7"}},
  {"componentId":"btn8","useComponent":"DigitButton","args":{"digit":"8"}},
  {"componentId":"btn9","useComponent":"DigitButton","args":{"digit":"9"}}
]}
```
