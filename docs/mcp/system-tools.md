---
layout: default
title: System Capabilities
parent: MCP Tools
nav_order: 6
---

# System Capabilities

These tools access native macOS features -- notifications, clipboard, file dialogs, and URL opening. They work the same way whether called from MCP tools or from evaluator functions inside your JSONL app.

---

## notify

Show a macOS notification via Notification Center.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `title` | string | yes | Notification title |
| `body` | string | yes | Notification body text |
| `subtitle` | string | no | Optional subtitle |

**Example:**
```
mcp__canopy__notify(title: "Download Complete", body: "Your file has been saved.")
```

---

## alert

Show a native macOS alert dialog. The calling thread waits for the user to dismiss it; the UI stays responsive.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `title` | string | yes | Alert title |
| `msg` | string | yes | Alert message |
| `style` | string | no | `"warning"`, `"critical"`, or `"informational"` (default) |
| `buttons` | array | no | Button labels (default: `["OK"]`) |

**Returns:** The index of the button the user clicked (0-based).

**Example:**
```
mcp__canopy__alert(
  title: "Delete item?",
  msg: "This action cannot be undone.",
  style: "warning",
  buttons: ["Delete", "Cancel"]
)
```

Returns `0` if the user clicked "Delete", `1` for "Cancel".

---

## clipboard_read

Read the current contents of the system clipboard.

**Parameters:** None

**Returns:** The clipboard text.

**Example:**
```
mcp__canopy__clipboard_read()
```

---

## clipboard_write

Write text to the system clipboard.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `text` | string | yes | The text to copy |

**Example:**
```
mcp__canopy__clipboard_write(text: "Copied from Canopy!")
```

---

## open_url

Open a URL or file path in the default application using NSWorkspace.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `url` | string | yes | URL or file path to open |

**Example:**
```
mcp__canopy__open_url(url: "https://github.com")
```

Opens the URL in the default browser. Also works with file paths to open them in their default app.

---

## file_open

Show a native file open dialog (NSOpenPanel). The main thread stays unblocked while the dialog is open.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `title` | string | no | Dialog title |
| `types` | array | no | Allowed file extensions (e.g., `["png", "jpg"]`) |
| `multi` | boolean | no | Allow multiple selection (default: false) |

**Returns:** Selected file path(s), or empty string if cancelled.

**Example:**
```
mcp__canopy__file_open(title: "Choose an image", types: ["png", "jpg", "gif"])
```

---

## file_save

Show a native file save dialog (NSSavePanel).

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `title` | string | no | Dialog title |
| `name` | string | no | Default file name |
| `types` | array | no | Allowed file extensions |

**Returns:** Selected save path, or empty string if cancelled.

**Example:**
```
mcp__canopy__file_save(title: "Export PDF", name: "report.pdf", types: ["pdf"])
```
