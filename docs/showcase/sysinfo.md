---
layout: default
title: System Info
parent: Showcase
nav_order: 6
---

# System Info

A system information viewer that calls native C libraries via FFI to gather data about the running system.

## Key Features

- **loadLibrary** -- loads native shared libraries (`libcurl`, `libsqlite3`, `libz`) at runtime
- **FFI function calls** -- calls C functions from loaded libraries and displays their results
- **Native integration** -- demonstrates how Canopy apps can access system-level functionality beyond what the built-in evaluator functions provide

## How to Run

```bash
build/canopy sample_apps/sysinfo
```

## What to Look For

- Library version strings retrieved by calling into `libcurl`, `libsqlite3`, and `libz`
- System information gathered through native API calls
- Results displayed in Cards with proper formatting

## How FFI Works

Canopy's FFI system lets you call functions in any native shared library:

1. **Load a library** using `loadLibrary` with the library name (e.g., `libcurl`)
2. **Call functions** using `callFunction` with the library handle, function name, and arguments
3. **Use results** -- return values are available as strings or numbers in the data model

This is configured via `ffi-config` packages (see [Packages](../packages/)) or inline in the JSONL.
