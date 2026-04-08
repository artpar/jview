---
layout: default
title: Expressions
parent: Building Apps
nav_order: 2
---

# Expressions

Any component prop that accepts a dynamic value can use three kinds of expressions: **literals**, **path references**, and **function calls**.

## Literals

Plain strings, numbers, and booleans pass through unchanged:

```json
{
  "props": {
    "content": "Hello"
  }
}

{
  "props": {
    "min": 0,
    "max": 100
  }
}

{
  "props": {
    "checked": true
  }
}
```

## Path References

Read a value from the data model:

```json
{
  "props": {
    "content": {
      "path": "/user/name"
    }
  }
}
```

If the path doesn't exist, the value resolves to an empty string.

## Function Calls

Compute a value using built-in or user-defined functions:

```json
{
  "props": {
    "content": {
      "functionCall": {
        "name": "concat",
        "args": [
          "Count: ",
          {
            "path": "/count"
          }
        ]
      }
    }
  }
}
```

A function call has two fields:
- `name` -- the function name (see [Functions](../functions/))
- `args` -- an array of arguments, each of which can be a literal, path reference, or nested function call

## Nesting

Arguments can themselves be function calls, allowing arbitrary computation:

```json
{
  "functionCall": {
    "name": "concat",
    "args": [
      "Total: $",
      {
        "functionCall": {
          "name": "multiply",
          "args": [
            {
              "path": "/quantity"
            },
            {
              "path": "/price"
            }
          ]
        }
      }
    ]
  }
}
```

## Conditional Values

Use the `if` function to choose between values:

```json
{
  "functionCall": {
    "name": "if",
    "args": [
      {
        "functionCall": {
          "name": "greaterThan",
          "args": [
            {
              "path": "/count"
            },
            0
          ]
        }
      },
      {
        "functionCall": {
          "name": "concat",
          "args": [
            {
              "path": "/count"
            },
            " items"
          ]
        }
      },
      "No items"
    ]
  }
}
```

The `if` function is lazy -- it only evaluates the branch that matches the condition.

## Where Expressions Work

Expressions work in any prop that accepts `DynamicString`, `DynamicNumber`, or `DynamicBoolean`. This includes:

- Text `content`, Card `title` and `subtitle`
- Button `label`, TextField `placeholder` and `value`
- CheckBox `checked`, Slider `sliderValue`, `min`, `max`
- All `style` properties (backgroundColor, textColor, fontSize, etc.)
- `disabled`, `visible`, `readOnly`, and other boolean props

Expressions are resolved every time the component renders, so they always reflect the current data model state.
