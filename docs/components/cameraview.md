---
layout: default
title: CameraView
parent: Components
nav_order: 24
---

# CameraView

Live camera preview and capture. Maps to **AVCaptureVideoPreviewLayer**.

Displays a live feed from the Mac's camera with an optional capture button. Requires camera permission (granted via system prompt on first use).

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `devicePosition` | string | `"front"` | Camera to use: `"front"` or `"back"` |
| `mirrored` | DynamicBoolean | `true` | Mirror the preview horizontally |
| `onCapture` | EventAction | | Action triggered when a photo is captured (receives file path) |
| `onError` | EventAction | | Action triggered on camera errors |

Size is controlled via `style.width` and `style.height` on the component.

## Example

Photo capture view:

```json
{"type":"createSurface","surfaceId":"main","title":"CameraView Example"}
{"type":"updateDataModel","surfaceId":"main","operations":[
  {"op":"replace","path":"/photoPath","value":""}
]}
{"type":"updateComponents","surfaceId":"main","components":[
  {"componentId":"root","type":"Column","props":{"padding":16,"gap":12,"align":"center"},"children":["camera","captureBtn"]},
  {"componentId":"camera","type":"CameraView","props":{
    "devicePosition":"front",
    "mirrored":true,
    "onCapture":{"action":{"setValues":[{"path":"/photoPath","value":"$event.path"}]}}
  },"style":{"width":320,"height":240}},
  {"componentId":"captureBtn","type":"Button","props":{"label":"Take Photo","style":"primary"}}
]}
```

## Notes

- The camera session starts automatically when the component is created.
- macOS prompts for camera permission on first use. Add `NSCameraUsageDescription` to Info.plist.
- Captured photos are saved as JPEG files to a temporary directory. The file path is passed to the `onCapture` action.
- Set `mirrored` to `false` for document scanning or rear-camera use cases.
- Most Macs only have a front-facing camera; `"back"` is primarily for external cameras.
