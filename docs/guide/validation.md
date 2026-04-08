---
layout: default
title: Validation
parent: Building Apps
nav_order: 11
---

# Validation

TextField components support validation rules that check user input and display error messages.

## Adding Validation Rules

Set the `validations` prop to an array of rule objects:

```json
{
  "componentId": "emailField",
  "type": "TextField",
  "props": {
    "placeholder": "you@example.com",
    "dataBinding": "/email",
    "value": {
      "path": "/email"
    },
    "validations": [
      {
        "type": "required",
        "message": "Email is required"
      },
      {
        "type": "email",
        "message": "Enter a valid email address"
      }
    ]
  }
}
```

Validation runs every time the field value changes. Error messages appear below the field.

## Validation Types

| Type | Value Field | Default Message | Description |
|------|-------------|-----------------|-------------|
| `required` | -- | "This field is required" | Field must not be empty or whitespace |
| `minLength` | number | "Must be at least N characters" | Minimum character count |
| `maxLength` | number | "Must be at most N characters" | Maximum character count |
| `pattern` | string (regex) | "Invalid format" | Must match the regular expression |
| `email` | -- | "Invalid email address" | Must be a valid email format |

## Custom Error Messages

Every rule accepts an optional `message` field that overrides the default:

```json
{
  "type": "minLength",
  "value": 8,
  "message": "Password must be at least 8 characters"
}
```

## Example: Registration Form

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "Register",
  "width": 400,
  "height": 400
}

{
  "type": "updateDataModel",
  "surfaceId": "main",
  "ops": [
    {
      "op": "add",
      "path": "/username",
      "value": ""
    },
    {
      "op": "add",
      "path": "/email",
      "value": ""
    },
    {
      "op": "add",
      "path": "/password",
      "value": ""
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
        "gap": 12,
        "padding": 20
      },
      "children": [
        "title",
        "usernameField",
        "emailField",
        "passwordField",
        "submitBtn"
      ]
    },
    {
      "componentId": "title",
      "type": "Text",
      "props": {
        "content": "Create Account",
        "variant": "h2"
      }
    },
    {
      "componentId": "usernameField",
      "type": "TextField",
      "props": {
        "placeholder": "Username",
        "dataBinding": "/username",
        "value": {
          "path": "/username"
        },
        "validations": [
          {
            "type": "required",
            "message": "Username is required"
          },
          {
            "type": "minLength",
            "value": 3,
            "message": "At least 3 characters"
          },
          {
            "type": "maxLength",
            "value": 20
          }
        ]
      }
    },
    {
      "componentId": "emailField",
      "type": "TextField",
      "props": {
        "placeholder": "Email",
        "dataBinding": "/email",
        "value": {
          "path": "/email"
        },
        "validations": [
          {
            "type": "required"
          },
          {
            "type": "email"
          }
        ]
      }
    },
    {
      "componentId": "passwordField",
      "type": "TextField",
      "props": {
        "placeholder": "Password",
        "inputType": "obscured",
        "dataBinding": "/password",
        "value": {
          "path": "/password"
        },
        "validations": [
          {
            "type": "required"
          },
          {
            "type": "minLength",
            "value": 8
          },
          {
            "type": "pattern",
            "value": "[A-Z]",
            "message": "Must contain an uppercase letter"
          }
        ]
      }
    },
    {
      "componentId": "submitBtn",
      "type": "Button",
      "props": {
        "label": "Create Account",
        "style": "primary",
        "onClick": {
          "action": {
            "event": {
              "name": "register",
              "dataRefs": [
                "/username",
                "/email",
                "/password"
              ]
            }
          }
        }
      }
    }
  ]
}
```

## Multiple Rules

Rules are checked in order. All failing rules produce error messages -- not just the first one. This means a user sees all problems at once rather than fixing them one at a time.
