---
layout: default
title: Video
parent: Components
nav_order: 22
---

# Video

Video player. Maps to **AVPlayerView**.

Plays video from a URL or local file with native macOS playback controls.

![Video player]({{ site.baseurl }}/screenshots/video_player_app.png){: .screenshot}

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `src` | DynamicString | | Video URL or local file path |
| `width` | int | | Player width in points |
| `height` | int | | Player height in points |
| `autoplay` | DynamicBoolean | `false` | Start playing automatically |
| `loop` | DynamicBoolean | `false` | Restart when playback ends |
| `controls` | DynamicBoolean | `true` | Show playback controls |
| `muted` | DynamicBoolean | `false` | Mute audio |
| `onEnded` | EventAction | | Action triggered when playback completes |

## Example

Video player:

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "Video Example"
}

{
  "type": "updateComponents",
  "surfaceId": "main",
  "components": [
    {
      "componentId": "root",
      "type": "Column",
      "props": {
        "padding": 16
      },
      "children": [
        "player"
      ]
    },
    {
      "componentId": "player",
      "type": "Video",
      "props": {
        "src": "https://example.com/demo.mp4",
        "width": 640,
        "height": 360,
        "controls": true,
        "autoplay": false
      }
    }
  ]
}
```

## Notes

- Supports any format that AVFoundation can play: MP4, MOV, M4V, and HLS streams.
- Set `controls` to `false` to hide the native playback bar (useful for background video).
- The `onEnded` action fires when the video reaches the end. Combine with `loop: false` to trigger a next action.
- Local file paths (e.g., from a file dialog or screen recording) work with `src`.
- If `width` and `height` are omitted, the player sizes to the video's natural dimensions.
