---
layout: default
title: AudioPlayer
parent: Components
nav_order: 23
---

# AudioPlayer

Audio playback controls. Maps to **AVPlayer** with custom controls.

Plays audio from a URL or local file. Renders as a compact 40-point bar with play/pause button, scrubber, and time display.

![AudioPlayer component]({{ site.baseurl }}/screenshots/audio_player_app.png){: .screenshot}

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `src` | DynamicString | | Audio URL or local file path |
| `autoplay` | DynamicBoolean | `false` | Start playing automatically |
| `loop` | DynamicBoolean | `false` | Restart when playback ends |
| `onEnded` | EventAction | | Action triggered when playback completes |

## Example

Podcast player:

```json
{"type":"createSurface","surfaceId":"main","title":"AudioPlayer Example"}
{"type":"updateComponents","surfaceId":"main","components":[
  {"componentId":"root","type":"Column","props":{"padding":16,"gap":8},"children":["title","player"]},
  {"componentId":"title","type":"Text","props":{"content":"Episode 42: Native macOS Apps","variant":"h3"}},
  {"componentId":"player","type":"AudioPlayer","props":{
    "src":"https://example.com/podcast-ep42.mp3",
    "autoplay":false
  }}
]}
```

## Notes

- The player bar includes a play/pause button, a time scrubber, and elapsed/remaining time labels.
- Supports MP3, M4A, WAV, AAC, and any format AVFoundation can decode.
- Local file paths (e.g., from an audio recording) work with `src`.
- Use `onEnded` to advance to the next track in a playlist.
- The player is always 40 points tall and stretches to fill the width of its parent.
