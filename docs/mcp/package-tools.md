---
layout: default
title: Package Management
parent: MCP Tools
nav_order: 10
---

# Package Management

These MCP tools mirror the `canopy pkg` CLI commands, letting you search, install, and publish packages programmatically. Packages are GitHub repos with a `canopy.json` manifest.

See the [Packages](../packages/) section for full details on manifests, installation, and publishing.

---

## package_login

Authenticate with GitHub using OAuth device flow. Opens a browser for authorization.

**Parameters:** None

**Example:**
```
mcp__canopy__package_login()
```

---

## package_search

Search for packages on GitHub.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `query` | string | yes | Search terms |
| `type` | string | no | Filter by type: `app`, `component`, `theme`, `ffi-config` |
| `limit` | number | no | Max results (default: 20) |

**Example:**
```
mcp__canopy__package_search(query: "calculator", type: "app")
```

---

## package_browse

List popular packages, sorted by stars or recent activity.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `type` | string | no | Filter by package type |
| `sort` | string | no | `"stars"` or `"updated"` (default: `"stars"`) |
| `limit` | number | no | Max results (default: 20) |

**Example:**
```
mcp__canopy__package_browse(type: "app", sort: "stars", limit: 10)
```

---

## package_info

Get detailed information about a package.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `repo` | string | yes | GitHub package ref (e.g., `github.com/artpar/notes`) |

**Returns:** Package manifest, versions, description, and dependencies.

**Example:**
```
mcp__canopy__package_info(repo: "github.com/artpar/notes")
```

---

## package_install

Install a package from GitHub.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `repo` | string | yes | GitHub package ref (e.g., `github.com/artpar/calculator`) |
| `version` | string | no | Specific version tag (e.g., `1.0.0`). Default: latest. |

**Example:**
```
mcp__canopy__package_install(repo: "github.com/artpar/calculator")
```

```
mcp__canopy__package_install(repo: "github.com/artpar/calculator", version: "1.2.0")
```

---

## package_uninstall

Remove an installed package.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `name` | string | yes | Package name (e.g., `github.com/artpar/calculator`) |

**Example:**
```
mcp__canopy__package_uninstall(name: "github.com/artpar/calculator")
```

---

## package_update

Update installed packages to the latest version.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `name` | string | no | Specific package to update. Omit to update all. |

**Example:**
```
mcp__canopy__package_update(name: "github.com/artpar/notes")
```

```
mcp__canopy__package_update()
```

---

## package_list

List installed packages.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `type` | string | no | Filter by type: `app`, `component`, `theme`, `ffi-config` |

**Example:**
```
mcp__canopy__package_list(type: "app")
```

---

## package_publish

Publish a package to GitHub.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | yes | Path to the package directory |
| `repo` | string | no | GitHub owner/repo (inferred from git remote if omitted; bare `owner/repo` accepted as shorthand) |
| `tag` | string | no | Version tag (inferred from canopy.json if omitted) |

**Example:**
```
mcp__canopy__package_publish(path: ".", repo: "github.com/artpar/my-app")
```

Creates a git tag and GitHub Release, and sets the `canopy-package` topic on the repo for discovery.

> **Note:** Bare `owner/repo` is accepted as shorthand for `github.com/owner/repo` in all package tool parameters.
