---
layout: default
title: Channel Demo
parent: Showcase
nav_order: 5
---

# Channel Demo

A multi-process application demonstrating cross-process communication via channels.

![Channel Demo](../screenshots/channel_demo.png)

## Key Features

- **createProcess** -- spawns multiple child processes, each with its own transport and surfaces
- **createChannel** -- sets up named message channels between processes
- **Broadcast mode** -- one channel broadcasts status updates to all subscribers
- **Queue mode** -- another channel distributes work items round-robin to worker processes
- **Cross-process data flow** -- processes publish results that update the UI in real time

## How to Run

```bash
build/canopy sample_apps/channel_demo
```

## What to Look For

- Multiple processes running simultaneously, each with its own section in the UI
- Messages flowing through broadcast channels appear in all subscriber panels
- Work items distributed through queue channels go to one worker at a time
- The UI updates in real time as messages are published and consumed
