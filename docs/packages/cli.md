---
layout: default
title: Package CLI Reference
parent: Packages
nav_order: 4
---

# Package CLI Reference

All package commands are subcommands of `canopy pkg`.

## Commands

| Command | Description |
|---------|-------------|
| `canopy pkg login` | Authenticate with GitHub via OAuth device flow |
| `canopy pkg search <query> [--type=TYPE]` | Search for packages by keyword |
| `canopy pkg info <github.com/owner/repo>` | Show package details, versions, and dependencies |
| `canopy pkg install <github.com/owner/repo> [@version]` | Install a package (latest or specific version) |
| `canopy pkg uninstall <github.com/owner/name>` | Remove an installed package |
| `canopy pkg update [github.com/owner/name]` | Update all packages, or a specific one |
| `canopy pkg list [--type=TYPE]` | List installed packages |
| `canopy pkg publish [path] [--repo=OWNER/REPO]` | Publish a package to GitHub |

## canopy pkg login

Starts the GitHub OAuth device flow. Opens your browser to authorize Canopy. The token is stored locally for future commands.

```bash
canopy pkg login
```

## canopy pkg search

Search GitHub for packages with the `canopy-package` topic.

```bash
canopy pkg search calculator
canopy pkg search notes --type=app
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--type=TYPE` | Filter by package type: `app`, `component`, `theme`, `ffi-config` |

## canopy pkg info

Display the manifest, available versions, and description for a package.

```bash
canopy pkg info github.com/artpar/notes
```

## canopy pkg install

Download and install a package.

```bash
canopy pkg install github.com/artpar/calculator
canopy pkg install github.com/artpar/calculator @1.2.0
```

The `@version` argument is optional. Without it, the latest release is installed.

> **Note:** Bare `owner/repo` is accepted as shorthand for `github.com/owner/repo`.

## canopy pkg uninstall

Remove a package from your system.

```bash
canopy pkg uninstall github.com/artpar/calculator
```

## canopy pkg update

Update packages to their latest available version.

```bash
# Update all installed packages
canopy pkg update

# Update a specific package
canopy pkg update github.com/artpar/notes
```

## canopy pkg list

Show installed packages.

```bash
# List all
canopy pkg list

# List only apps
canopy pkg list --type=app
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--type=TYPE` | Filter by package type |

## canopy pkg publish

Publish a package to GitHub. Creates a git tag and GitHub Release.

```bash
# Publish from current directory (repo inferred from git remote)
canopy pkg publish .

# Publish with explicit repo
canopy pkg publish . --repo=yourname/your-app
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--repo=OWNER/REPO` | GitHub repository (default: inferred from git remote) |
