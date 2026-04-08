---
layout: default
title: Channels
parent: Building Apps
nav_order: 8
---

# Channels

Channels provide publish/subscribe communication between processes and the main session. A channel is a named pipe that delivers values to subscribers.

## Creating a Channel

```json
{
  "type": "createChannel",
  "channelId": "notifications",
  "mode": "broadcast"
}
```

Fields:
- `channelId` -- unique name for the channel
- `mode` -- delivery mode: `"broadcast"` (default) or `"queue"`

### Broadcast Mode

Every subscriber receives every published value. Use this for notifications, status updates, or any case where all listeners need the same data.

### Queue Mode

Published values are delivered round-robin to subscribers. Each value goes to exactly one subscriber. Use this for work distribution or load balancing.

## Subscribing

Register a subscription that writes received values to a data model path:

```json
{
  "type": "subscribe",
  "channelId": "notifications",
  "targetPath": "/lastNotification"
}
```

Fields:
- `channelId` -- the channel to subscribe to
- `processId` -- (optional) subscribe on behalf of a specific process; omit for session-level
- `targetPath` -- JSON Pointer path where received values are written in the data model

When a value is published to the channel, it gets written to `/lastNotification` on every surface's data model. Components bound to that path re-render automatically.

## Publishing

Send a value to all subscribers:

```json
{
  "type": "publish",
  "channelId": "notifications",
  "value": {
    "title": "New message",
    "body": "Hello!"
  }
}
```

The value can be any JSON type -- string, number, object, or array.

The last published value is also stored at `/channels/{channelId}/value` in the data model.

## Unsubscribing

Remove a specific subscription:

```json
{
  "type": "unsubscribe",
  "channelId": "notifications",
  "targetPath": "/lastNotification"
}
```

Or remove all subscriptions for a process:

```json
{
  "type": "unsubscribe",
  "channelId": "notifications",
  "processId": "worker1"
}
```

## Deleting a Channel

```json
{
  "type": "deleteChannel",
  "channelId": "notifications"
}
```

This removes the channel and all its subscriptions.

## Example: Notification System

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "Notifications",
  "width": 400,
  "height": 300
}

{
  "type": "updateDataModel",
  "surfaceId": "main",
  "ops": [
    {
      "op": "add",
      "path": "/lastNotification",
      "value": ""
    }
  ]
}

{
  "type": "createChannel",
  "channelId": "alerts",
  "mode": "broadcast"
}

{
  "type": "subscribe",
  "channelId": "alerts",
  "targetPath": "/lastNotification"
}

{
  "type": "updateComponents",
  "surfaceId": "main",
  "components": [
    {
      "componentId": "root",
      "type": "Column",
      "props": {
        "gap": 12
      },
      "children": [
        "display",
        "sendBtn"
      ]
    },
    {
      "componentId": "display",
      "type": "Text",
      "props": {
        "content": {
          "path": "/lastNotification"
        }
      }
    },
    {
      "componentId": "sendBtn",
      "type": "Button",
      "props": {
        "label": "Send Alert",
        "onClick": {
          "action": {
            "event": {
              "name": "sendAlert"
            }
          }
        }
      }
    }
  ]
}
```

The server can then publish to the channel:

```json
{
  "type": "publish",
  "channelId": "alerts",
  "value": "System update available"
}
```

Every surface subscribed to `alerts` with `targetPath` of `/lastNotification` will see the update.
