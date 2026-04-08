---
layout: default
title: Todo
parent: Showcase
nav_order: 3
---

# Todo

A todo list app with add/remove functionality and a live count of remaining items.

## Prompt

> Build a todo list with add/remove and a count of remaining items.

## Key Features

- **List** -- dynamic list that grows and shrinks as items are added and removed
- **Data binding** -- the item count text updates automatically when the todo array changes
- **Dynamic children** -- new todo items are created as components bound to array elements
- **Array functions** -- uses `append` and `removeLast` evaluator functions to manage the todo array

## How to Run

```bash
build/canopy sample_apps/todo
```

## What to Look For

- Type a task name and click Add to create a new item
- Each item has a checkbox to mark it complete and a delete button
- The "X items remaining" counter updates as you check or remove items
- The list scrolls when there are more items than fit in the window
