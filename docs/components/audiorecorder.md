---
layout: default
title: AudioRecorder
parent: Components
nav_order: 25
---

# AudioRecorder

Microphone recording with level meter. Maps to **AVAudioRecorder**.

Records audio from the microphone. Renders as a compact 40-point bar with a record/stop button, audio level meter, and elapsed time display. Requires microphone permission.

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `format` | string | `"m4a"` | Audio format: `"m4a"` or `"wav"` |
| `sampleRate` | DynamicNumber | `44100` | Sample rate in Hz |
| `recordChannels` | int | `1` | Number of audio channels (1 = mono, 2 = stereo) |
| `onRecordingStarted` | EventAction | | Action triggered when recording begins |
| `onRecordingStopped` | EventAction | | Action triggered when recording stops (receives file path) |
| `onLevel` | EventAction | | Action triggered periodically with audio level data |
| `onError` | EventAction | | Action triggered on microphone errors |

## Example

Voice memo recorder:

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "AudioRecorder Example"
}

{
  "type": "updateDataModel",
  "surfaceId": "main",
  "operations": [
    {
      "op": "replace",
      "path": "/recordings",
      "value": []
    },
    {
      "op": "replace",
      "path": "/isRecording",
      "value": false
    }
  ]
}

{
  "type": "updateComponents",
  "surfaceId": "main",
  "components": [
    {
      "componentId": "root",
      "type": "Column",
      "props": {
        "padding": 16,
        "gap": 12
      },
      "children": [
        "heading",
        "recorder"
      ]
    },
    {
      "componentId": "heading",
      "type": "Text",
      "props": {
        "content": "Voice Memos",
        "variant": "h2"
      }
    },
    {
      "componentId": "recorder",
      "type": "AudioRecorder",
      "props": {
        "format": "m4a",
        "sampleRate": 44100,
        "onRecordingStarted": {
          "action": {
            "setValues": [
              {
                "path": "/isRecording",
                "value": true
              }
            ]
          }
        },
        "onRecordingStopped": {
          "action": {
            "setValues": [
              {
                "path": "/isRecording",
                "value": false
              }
            ]
          }
        }
      }
    }
  ]
}
```

## Notes

- macOS prompts for microphone permission on first use. Add `NSMicrophoneUsageDescription` to Info.plist.
- Click the record button to start; click again to stop. The recorded file path is passed to `onRecordingStopped`.
- The level meter updates in real time during recording, showing the current audio input level.
- The timer shows elapsed recording time in MM:SS format.
- Recordings are saved to a temporary directory. Use the file path to play back with AudioPlayer or upload.
- `m4a` (AAC) produces smaller files. `wav` produces uncompressed audio for higher quality.
