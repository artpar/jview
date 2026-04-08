---
layout: default
title: Publishing Packages
parent: Packages
nav_order: 3
---

# Publishing Packages

Share your Canopy apps and components with others by publishing them to GitHub.

## Prerequisites

- A GitHub account
- A GitHub repo containing your package
- A `canopy.json` manifest at the repo root (see [Package Manifest](manifest))

## Step-by-Step

### 1. Create your repo

Your repo should have at minimum:
- `canopy.json` -- the package manifest
- Your entry file (e.g., `prompt.jsonl` for apps)

### 2. Tag a release

Version your package using git tags:

```bash
git tag v1.0.0
git push --tags
```

### 3. Authenticate

Login to GitHub via the OAuth device flow. This opens your browser for authorization.

```bash
canopy pkg login
```

### 4. Publish

From your package directory:

```bash
canopy pkg publish .
```

Or specify the repo explicitly:

```bash
canopy pkg publish . --repo=yourname/your-app
```

## What Publish Does

1. Reads `canopy.json` to get the package name and version
2. Creates a git tag matching the version (if not already tagged)
3. Creates a GitHub Release from that tag
4. Sets the `canopy-package` topic on the repo

The `canopy-package` topic is how other users discover your package through `canopy pkg search`.

## Updating Your Package

To publish a new version:

1. Update the `version` field in `canopy.json`
2. Commit your changes
3. Run `canopy pkg publish .`

This creates a new tag and release for the updated version.

## Discovery

Published packages are found through GitHub's topic search. Users find your package with:

```bash
canopy pkg search your-keyword
```

Good keywords in your `canopy.json` help users find your package. Choose terms that describe what your app does, not how it works.
