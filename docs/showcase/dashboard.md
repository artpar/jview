---
layout: default
title: Dashboard
parent: Showcase
nav_order: 4
---

# Dashboard

A data dashboard with summary cards, an activity list, and quick action buttons.

## Prompt

> Build a data dashboard with summary cards across the top, activity list, and quick action buttons.

## Key Features

- **Row/Column layout** -- summary cards arranged horizontally across the top, with vertical stacking below
- **Card** -- each summary metric gets its own Card with a title and value
- **Nested components** -- cards contain Text components with different typography styles (h2 for values, body for labels)
- **Button actions** -- quick action buttons in the sidebar trigger different operations

## How to Run

```bash
build/canopy sample_apps/dashboard
```

## What to Look For

- Summary cards across the top show key metrics (users, revenue, etc.)
- The activity list in the main area shows recent events
- Quick action buttons on the side provide shortcuts to common tasks
- The layout adapts to the window size with proper spacing and alignment
