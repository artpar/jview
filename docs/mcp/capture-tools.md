---
layout: default
title: Screenshots & Capture
parent: MCP Tools
nav_order: 5
---

# Screenshots & Capture

Capture the visual output of any Canopy window as a PNG image.

---

## take_screenshot

Capture a window's contents as a PNG image.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `surface_id` | string | yes | The window to capture |
| `filePath` | string | no | Save to this file path instead of returning base64 |

**Example -- get as base64:**
```
mcp__canopy__take_screenshot(surface_id: "main")
```

Returns the screenshot as a base64-encoded PNG image that can be displayed inline.

**Example -- save to file:**
```
mcp__canopy__take_screenshot(surface_id: "main", filePath: "/tmp/screenshot.png")
```

Saves the PNG to the specified path and returns the file path.

## Use Cases

- **Visual verification**: After building a UI, take a screenshot to confirm layout and styling look correct.
- **Regression testing**: Capture screenshots before and after changes to compare.
- **Documentation**: Generate images of your app for docs or presentations.
- **Debugging**: See exactly what the user would see, including text rendering, colors, and spacing.
