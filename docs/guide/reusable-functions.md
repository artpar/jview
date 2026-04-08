---
layout: default
title: Reusable Functions
parent: Building Apps
nav_order: 6
---

# Reusable Functions

Use `defineFunction` to create named, parametric functions that can be called from expressions and actions, just like built-in functions.

## Defining a Function

```json
{"type":"defineFunction","name":"appendDigit","params":["digit"],
  "body":{
    "functionCall":{
      "call":"updateDataModel",
      "args":{
        "ops":[
          {"op":"replace","path":"/display","value":{
            "functionCall":{"name":"concat","args":[{"path":"/display"},{"param":"digit"}]}
          }}
        ]
      }
    }
  }
}
```

Fields:
- `name` -- the function name, usable in `{"functionCall": {"name": "appendDigit", ...}}`
- `params` -- parameter names
- `body` -- the function body, with `{"param": "name"}` placeholders

## How It Works

When you call a user-defined function:

1. Arguments are resolved (path references and nested function calls evaluate first).
2. The body is deep-copied.
3. Every `{"param": "name"}` placeholder is replaced with the resolved argument value.
4. The resulting expression is evaluated.

## Calling a Defined Function

In an expression:

```json
{"functionCall": {"name": "appendDigit", "args": ["5"]}}
```

In an action:

```json
{"props": {"onClick": {"action": {"functionCall": {"call": "appendDigit", "args": ["5"]}}}}}
```

## Example: Toggle Boolean

```json
{"type":"defineFunction","name":"toggleFlag","params":["path"],
  "body":{
    "functionCall":{
      "call":"updateDataModel",
      "args":{
        "ops":[
          {"op":"replace","path":{"param":"path"},"value":{
            "functionCall":{"name":"not","args":[{"path":{"param":"path"}}]}
          }}
        ]
      }
    }
  }
}
```

Usage:

```json
{"componentId":"toggle","type":"Button","props":{
  "label":"Toggle Dark Mode",
  "onClick":{"action":{"functionCall":{"call":"toggleFlag","args":["/settings/darkMode"]}}}
}}
```

## Resolution Order

When a function call is evaluated, Canopy looks up the name in this order:

1. **Built-in functions** (concat, add, if, etc.)
2. **User-defined functions** (from `defineFunction`)
3. **FFI functions** (from `loadLibrary`)

If you define a function with the same name as a built-in, the built-in takes precedence.
