---
layout: default
title: Calculator
parent: Showcase
nav_order: 1
---

# Calculator

A calculator app with a display, digit buttons, operators, clear, and equals. Dark background with orange accent for operator buttons.

![Calculator](../screenshots/calculator.png)

## Prompt

> Build a calculator app with display, digit buttons (0-9), operators (+, -, *, /), clear, equals. Dark background, orange accent for operators.

## Key Features

- **defineComponent** -- reusable `DigitButton` and `OpButton` components with different styles
- **defineFunction** -- `appendDigit`, `applyOperator`, `calculate`, and `clearDisplay` functions that manage calculator state
- **Grid layout** -- buttons arranged in a calculator grid using nested Row and Column components
- **Data binding** -- display updates in real time as you press buttons

## How to Run

```bash
build/canopy sample_apps/calculator
```

## What to Look For

- The display shows the current input and updates as you click digit buttons
- Operator buttons use the orange accent color
- Clear resets the entire state
- Equals computes the result using the stored operator
