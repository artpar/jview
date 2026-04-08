---
layout: default
title: Image
parent: Components
nav_order: 17
---

# Image

Image display. Maps to **NSImageView**.

Loads and displays an image from a URL or local file path. Supports sizing and accessibility labels.

![Image component]({{ site.baseurl }}/screenshots/image.png){: .screenshot}

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `src` | DynamicString | | Image URL or local file path |
| `alt` | DynamicString | | Accessibility description |
| `width` | int | | Image width in points |
| `height` | int | | Image height in points |

## Example

Profile picture:

```json
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "Image Example"
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
        "gap": 8,
        "align": "center"
      },
      "children": [
        "avatar",
        "name"
      ]
    },
    {
      "componentId": "avatar",
      "type": "Image",
      "props": {
        "src": "https://example.com/avatar.jpg",
        "alt": "User avatar",
        "width": 80,
        "height": 80
      }
    },
    {
      "componentId": "name",
      "type": "Text",
      "props": {
        "content": "Jane Smith",
        "variant": "h3"
      }
    }
  ]
}
```

## Notes

- Images load asynchronously. The view reserves space based on `width` and `height` while loading.
- If only `width` or `height` is set, the image scales proportionally.
- If neither is set, the image displays at its natural size.
- Local file paths (e.g., from a camera capture or file dialog) are also supported via `src`.
- The `alt` prop is used for VoiceOver accessibility.
