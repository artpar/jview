---
layout: default
title: Components
nav_order: 4
has_children: true
---

# Components

Canopy provides 25 native macOS components rendered as real AppKit views. No webview, no Electron -- every component maps directly to a native Cocoa widget.

## Layout

Containers that arrange child components spatially.

| Component | Description | AppKit Class |
|-----------|-------------|--------------|
| [Row](row) | Horizontal stack layout | NSStackView |
| [Column](column) | Vertical stack layout | NSStackView |
| [Card](card) | Titled container with optional collapse | NSBox |
| [SplitView](splitview) | Resizable panes | NSSplitView |
| [Tabs](tabs) | Tabbed container | NSTabView |
| [List](list) | Scrollable container | NSScrollView |
| [Modal](modal) | Dialog overlay | NSPanel |

## Input

Interactive components that accept user input and support data binding.

| Component | Description | AppKit Class |
|-----------|-------------|--------------|
| [TextField](textfield) | Text input field | NSTextField |
| [CheckBox](checkbox) | Toggle switch | NSButton (checkBox) |
| [Slider](slider) | Range input | NSSlider |
| [ChoicePicker](choicepicker) | Dropdown or segmented selector | NSPopUpButton |
| [DateTimeInput](datetimeinput) | Date and time picker | NSDatePicker |
| [SearchField](searchfield) | Search input with clear button | NSSearchField |
| [Button](button) | Clickable action trigger | NSButton |

## Display

Read-only components for presenting content.

| Component | Description | AppKit Class |
|-----------|-------------|--------------|
| [Text](text) | Static text with variant styling | NSTextField |
| [Icon](icon) | SF Symbol icon | NSImageView |
| [Image](image) | Image from URL or file | NSImageView |
| [Divider](divider) | Horizontal separator line | NSBox |
| [ProgressBar](progressbar) | Progress indicator | NSProgressIndicator |

## Rich

Components for structured or rich content editing and display.

| Component | Description | AppKit Class |
|-----------|-------------|--------------|
| [RichTextEditor](richtexteditor) | Markdown-capable text editor | NSTextView |
| [OutlineView](outlineview) | Hierarchical tree list | NSOutlineView |
| [Video](video) | Video player | AVPlayerView |
| [AudioPlayer](audioplayer) | Audio playback controls | AVPlayer |

## Media

Components for camera and microphone access.

| Component | Description | AppKit Class |
|-----------|-------------|--------------|
| [CameraView](cameraview) | Live camera preview and capture | AVCaptureVideoPreviewLayer |
| [AudioRecorder](audiorecorder) | Microphone recording with level meter | AVAudioRecorder |
