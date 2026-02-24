# jview

Native macOS renderer for the [A2UI](https://a2ui.org) JSONL protocol. No webview, no Electron — real AppKit widgets driven by declarative JSON.

## What It Does

jview reads A2UI JSONL messages (from a file, and eventually from SSE/WebSocket streams) and renders them as native Cocoa widgets. Each JSON line describes a UI operation — create a window, add components, update data — and jview turns that into live NSViews on screen.

```
JSONL messages  -->  Engine (Go)  -->  CGo bridge  -->  AppKit (Obj-C)  -->  Native UI
```

## Quick Start

**Requirements:** macOS, Go 1.24+

```bash
# Build
make build

# Run a fixture
build/jview testdata/hello.jsonl

# Run all tests
make test

# Full gate (tests + screenshot verification)
make check
```

## How It Works

A2UI JSONL defines surfaces (windows), components (widgets), and a data model (reactive state). jview processes these through a layered architecture:

```
Transport (goroutine)          <- reads JSONL from file/SSE/WS
    |
engine.Session (goroutine)     <- routes messages to surfaces
    |
engine.Surface                 <- manages tree, data model, bindings
    |
Dispatcher.RunOnMain()         <- batches render ops to main thread
    |
darwin.Renderer (main thread)  <- CGo -> ObjC -> NSView creation/updates
    |
Native Cocoa widgets           <- visible on screen
```

### Supported Components

| Component | Description |
|-----------|-------------|
| Text | Labels and headings (h1-h5, body, caption) |
| Row | Horizontal stack layout |
| Column | Vertical stack layout |
| Card | NSBox container with title |
| Button | Clickable button with server actions |
| TextField | Text input with two-way data binding |
| CheckBox | Toggle with two-way data binding |

### Data Binding

Components can bind to the data model using JSON Pointers. When a user types in a TextField bound to `/name`, any Text component displaying `{"path": "/name"}` updates automatically.

## Example

`testdata/hello.jsonl`:
```json
{"type":"createSurface","surfaceId":"main","title":"Hello jview","width":600,"height":400}
{"type":"updateComponents","surfaceId":"main","components":[
  {"componentId":"card1","type":"Card","props":{"title":"Welcome"},"children":["heading","body"]},
  {"componentId":"heading","type":"Text","props":{"content":"Hello, jview!","variant":"h1"}},
  {"componentId":"body","type":"Text","props":{"content":"Native macOS rendering.","variant":"body"}}
]}
```

## Project Structure

```
protocol/          JSONL parsing, message types, dynamic values
engine/            Session, Surface, DataModel, BindingTracker, Resolver
renderer/          Platform-agnostic Renderer interface + mock for tests
platform/darwin/   CGo + Objective-C AppKit implementation
transport/         Message sources (file, future: SSE, WebSocket)
testdata/          JSONL fixtures for testing and demos
```

## Testing

Three layers:

- **Unit tests** — pure Go, no display needed: protocol parsing, data model, bindings, resolver
- **Integration tests** — engine with mock renderer: component creation, data binding, callbacks
- **Screenshot verification** — builds real binary, launches fixtures, captures screenshots

All tests run with `-race` detection enabled.

```bash
make test          # Headless unit + integration tests
make verify        # Build + screenshot capture for all fixtures
make check         # Both (the gate)
```

## License

MIT
