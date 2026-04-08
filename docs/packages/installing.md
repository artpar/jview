---
layout: default
title: Installing Packages
parent: Packages
nav_order: 2
---

# Installing Packages

## Searching

Find packages by keyword. Results come from GitHub repos tagged with the `canopy-package` topic.

```bash
canopy pkg search calculator
```

Filter by type:

```bash
canopy pkg search notes --type=app
```

## Viewing Package Details

Before installing, you can inspect a package:

```bash
canopy pkg info github.com/artpar/notes
```

This shows the manifest, available versions, description, and dependencies.

## Installing

Install the latest version:

```bash
canopy pkg install github.com/artpar/calculator
```

Install a specific version:

```bash
canopy pkg install github.com/artpar/calculator @1.0.0
```

The package is downloaded from the GitHub Release and placed in the appropriate directory based on its type:

| Type | Location |
|------|----------|
| app | `~/.canopy/apps/github.com/{owner}/{name}/` |
| component | `~/.canopy/library/{name}.jsonl` |
| theme | `~/.canopy/themes/{name}.jsonl` |
| ffi-config | `~/.canopy/ffi/{name}.json` |

Installed apps appear in the Canopy menubar menu and can be launched from there.

## Listing Installed Packages

See everything you have installed:

```bash
canopy pkg list
```

Filter by type:

```bash
canopy pkg list --type=app
```

## Updating

Update all packages to their latest versions:

```bash
canopy pkg update
```

Update a specific package:

```bash
canopy pkg update github.com/artpar/notes
```

## Uninstalling

Remove a package:

```bash
canopy pkg uninstall github.com/artpar/calculator
```

This removes the package files from the install directory.

> **Note:** Bare `owner/repo` is accepted as shorthand for `github.com/owner/repo`.
