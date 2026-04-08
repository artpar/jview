---
layout: default
title: Packages
nav_order: 8
has_children: true
---

# Packages

Canopy packages are GitHub repositories with a `canopy.json` manifest at the root. You can install apps, reusable components, themes, and FFI configurations directly from GitHub.

## How It Works

1. Every package lives in a GitHub repo with a `canopy.json` file
2. Packages are discovered via the `canopy-package` GitHub topic
3. Install with `canopy pkg install owner/repo`
4. Installed apps appear in the Canopy menubar menu

## Package Types

| Type | What it is | Install location |
|------|-----------|-----------------|
| **app** | A complete Canopy application | `~/.canopy/apps/{owner}/{name}/` |
| **component** | A reusable UI component | `~/.canopy/library/{name}.jsonl` |
| **theme** | Visual theme (colors, fonts) | `~/.canopy/themes/{name}.jsonl` |
| **ffi-config** | FFI library configuration | `~/.canopy/ffi/{name}.json` |

## Quick Start

```bash
# Search for packages
canopy pkg search calculator

# Install an app
canopy pkg install artpar/calculator

# List what you have installed
canopy pkg list

# Update everything
canopy pkg update
```

## In This Section

- [Package Manifest](manifest) -- the `canopy.json` file format
- [Installing Packages](installing) -- search, install, update, remove
- [Publishing Packages](publishing) -- share your apps on GitHub
- [CLI Reference](cli) -- complete command reference
