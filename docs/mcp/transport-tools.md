---
layout: default
title: Sending Messages
parent: MCP Tools
nav_order: 4
---

# Sending Messages

The `send_message` tool lets you inject raw A2UI JSONL messages into a surface. This is the most powerful tool -- you can create components, update props, modify the data model, or trigger any protocol message.

---

## send_message

Inject any A2UI JSONL message into a surface's message stream.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `surface_id` | string | yes | The window to send the message to |
| `message` | string | yes | A valid A2UI JSONL message (JSON string) |

**Example -- create a new component:**
```
mcp__canopy__send_message(
  surface_id: "main",
  message: "{\"componentId\":\"greeting\",\"type\":\"Text\",\"props\":{\"content\":\"Hello, world!\"}}"
)
```

**Example -- update an existing component's props:**
```
mcp__canopy__send_message(
  surface_id: "main",
  message: "{\"componentId\":\"greeting\",\"type\":\"Text\",\"props\":{\"content\":\"Updated text\"}}"
)
```

**Example -- update the data model:**
```
mcp__canopy__send_message(
  surface_id: "main",
  message: "{\"type\":\"updateDataModel\",\"operations\":[{\"op\":\"replace\",\"path\":\"/status\",\"value\":\"active\"}]}"
)
```

**Example -- create a surface with a full layout:**
```
mcp__canopy__send_message(
  surface_id: "main",
  message: "{\"type\":\"surface\",\"surfaceId\":\"main\",\"title\":\"My App\",\"width\":800,\"height\":600}"
)
```

Then send component messages to populate the window:
```
mcp__canopy__send_message(
  surface_id: "main",
  message: "{\"componentId\":\"root\",\"type\":\"Column\",\"props\":{\"padding\":16},\"children\":[\"title\",\"content\"]}"
)
```

## Tips

- Messages are processed in order. Create child components before their parent references them, or send the parent last.
- You can build entire UIs dynamically by sending a sequence of `send_message` calls.
- Combine with `get_tree` and `get_data_model` to read state, modify it, and verify the result.
