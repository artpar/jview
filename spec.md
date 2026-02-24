# jview A2UI Protocol Specification

This document describes the A2UI JSONL protocol subset implemented by jview and the rendering rules applied by the engine.

## Wire Format

Messages are newline-delimited JSON (JSONL). Each line is a self-contained JSON object. Blank lines are skipped. Maximum line size: 10MB.

Every message has a `type` field and most have a `surfaceId` field:

```json
{"type":"<messageType>","surfaceId":"<id>", ...}
```

## Message Types

### createSurface

Opens a new native window.

```json  
{
  "type": "createSurface",
  "surfaceId": "main",
  "title": "My App",
  "width": 800,
  "height": 600
}
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| surfaceId | string | yes | | Unique identifier for this surface |
| title | string | yes | | Window title |
| width | int | no | 800 | Window width in points |
| height | int | no | 600 | Window height in points |

Duplicate `createSurface` for the same `surfaceId` is silently ignored.

### deleteSurface

Closes and removes a surface.

```json
{
  "type": "deleteSurface",
  "surfaceId": "main"
}
```

Deleting a non-existent surface is a no-op.

### updateComponents

Sends a batch of component definitions. Components are created or replaced atomically.

```json
{
  "type": "updateComponents",
  "surfaceId": "main",
  "components": [
    {"componentId": "t1", "type": "Text", "props": {"content": "Hello"}},
    {"componentId": "col", "type": "Column", "children": ["t1"]}
  ]
}
```

| Field | Type | Description |
|-------|------|-------------|
| components | Component[] | Array of component definitions |

Components within a batch may reference each other as children. The engine topologically sorts them (leaves first) before rendering.

### updateDataModel

Applies JSON Patch-style operations to the surface's data model.

```json
{
  "type": "updateDataModel",
  "surfaceId": "main",
  "ops": [
    {"op": "add", "path": "/name", "value": "Alice"},
    {"op": "replace", "path": "/count", "value": 42},
    {"op": "remove", "path": "/temp"}
  ]
}
```

| Op | Description |
|----|-------------|
| add | Set value at path, creating intermediate objects |
| replace | Same as add (idempotent) |
| remove | Delete value at path |

Paths use JSON Pointer syntax (RFC 6901): `/foo/bar/0` addresses `root.foo.bar[0]`.

After all ops execute, the engine finds components bound to affected paths and re-renders them.

### test

Defines an inline test case with assertions and event simulations. Test messages are interleaved with app messages in the same JSONL file. `jview <file.jsonl>` ignores test messages. `jview test <file.jsonl>` executes them against real AppKit rendering.

```json
{
  "type": "test",
  "surfaceId": "main",
  "name": "initial state",
  "steps": [
    {"assert": "component", "componentId": "heading", "props": {"content": "Welcome", "variant": "h1"}},
    {"assert": "dataModel", "path": "/name", "value": ""},
    {"assert": "children", "componentId": "root", "children": ["heading", "body"]},
    {"assert": "count", "componentId": "root", "count": 2},
    {"assert": "notExists", "componentId": "ghost"},
    {"assert": "layout", "componentId": "heading", "layout": {"width": 200}},
    {"assert": "style", "componentId": "heading", "style": {"fontSize": 24, "bold": true}}
  ]
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| surfaceId | string | yes | Surface to test against |
| name | string | yes | Human-readable test case name |
| steps | TestStep[] | yes | Sequence of assertions and simulations |

#### Assertion Types

| Assert | Fields | Description |
|--------|--------|-------------|
| `component` | componentId, props, componentType | Subset match on resolved props. Optionally check component type. |
| `dataModel` | path, value | Check data model value at JSON Pointer path |
| `children` | componentId, children | Check ordered list of child component IDs |
| `count` | componentId, count | Check number of children |
| `notExists` | componentId | Verify component does not exist |
| `action` | name, data | Check that a server action was fired with matching name and data |
| `layout` | componentId, layout | Check computed layout (x, y, width, height) from real NSView frames |
| `style` | componentId, style | Check computed style (fontName, fontSize, bold, italic, textColor, bgColor, hidden, opacity) from real NSView properties |

#### Event Simulation

```json
{"simulate": "event", "componentId": "nameField", "event": "change", "eventData": "Alice"}
```

| Event | Component | Description |
|-------|-----------|-------------|
| `change` | TextField | Set text value |
| `click` | Button | Trigger onClick action |
| `toggle` | CheckBox | Toggle checked state |
| `slide` | Slider | Set slider value |
| `select` | ChoicePicker | Select option |
| `datechange` | DateTimeInput | Set date value |

#### Test Runner Behavior

- Tests execute sequentially in file order
- Side effects from simulations persist across tests (shared session state)
- Captured actions reset at the start of each test
- `jview test` uses real AppKit rendering (not mocked) with synchronous dispatch
- Exit code 0 if all pass, 1 if any fail

### setTheme

Changes the visual theme. *Not yet implemented — reserved for Phase 3.*

```json
{
  "type": "setTheme",
  "surfaceId": "main",
  "theme": "dark"
}
```

Values: `"light"`, `"dark"`, `"system"`.

---

## Component Model

### Component Definition

```json
{
  "componentId": "unique_id",
  "type": "Text",
  "props": { ... },
  "children": ["child1", "child2"]
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| componentId | string | yes | Unique ID within the surface |
| type | string | yes | Component type name |
| props | object | no | Component-specific properties |
| children | ChildList | no | Child component references |

### ChildList

Either a static array of component IDs:

```json
"children": ["a", "b", "c"]
```

Or a template for dynamic expansion (Phase 2):

```json
"children": {
  "forEach": "/items",
  "templateId": "item_tmpl",
  "itemVariable": "item"
}
```

### Dynamic Values

Any string, number, or boolean property can be either a literal or a data model reference.

**Literal:**
```json
"content": "Hello"
```

**Path reference (resolves from data model):**
```json
"content": {"path": "/user/name"}
```

**Function call (Phase 2):**
```json
"content": {"functionCall": {"name": "concat", "args": ["Hello, ", {"path": "/name"}]}}
```

When a path reference is used, the engine registers a binding so the component re-renders when that data model path changes.

---

## Component Types

### Text

Displays read-only text.

| Prop | Type | Default | Values |
|------|------|---------|--------|
| content | DynamicString | `""` | Text to display |
| variant | string | `"body"` | `h1`, `h2`, `h3`, `h4`, `h5`, `body`, `caption` |

### Row

Horizontal stack layout (NSStackView, horizontal).

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| justify | string | | `start`, `center`, `end`, `spaceBetween`, `spaceAround` |
| align | string | | `start`, `center`, `end`, `stretch` |
| gap | int | 0 | Spacing between children in points |
| padding | int | 0 | Internal padding in points |

### Column

Vertical stack layout (NSStackView, vertical).

Same props as Row.

### Card

Titled container (NSBox).

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| title | DynamicString | `""` | Box title |
| subtitle | DynamicString | `""` | Subtitle text |
| collapsible | DynamicBoolean | `false` | Whether card can collapse |
| collapsed | DynamicBoolean | `false` | Initial collapsed state |
| padding | int | 0 | Internal padding |

### Button

Clickable button with server action.

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| label | DynamicString | `""` | Button text |
| style | string | `"secondary"` | `primary`, `secondary`, `destructive` |
| disabled | DynamicBoolean | `false` | Whether button is disabled |
| onClick | EventAction | | Action to fire on click |

### TextField

Text input with optional two-way data binding.

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| placeholder | DynamicString | `""` | Placeholder text |
| value | DynamicString | `""` | Current value |
| inputType | string | `"shortText"` | `shortText`, `longText`, `number`, `obscured` |
| readOnly | DynamicBoolean | `false` | Whether editing is disabled |
| dataBinding | string | | JSON Pointer for two-way binding |
| onChange | EventAction | | Action to fire on change |

When `dataBinding` is set, typing in the field writes the value to the data model at that path. Any components bound to overlapping paths re-render automatically (excluding the source field).

### CheckBox

Toggle with optional two-way data binding.

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| label | DynamicString | `""` | Label text |
| checked | DynamicBoolean | `false` | Current state |
| dataBinding | string | | JSON Pointer for two-way binding |
| onToggle | EventAction | | Action to fire on toggle |

---

## Actions

```json
"onClick": {
  "action": {
    "type": "serverAction",
    "name": "submitForm",
    "params": {"key": "value"},
    "dataRefs": ["/name", "/email"]
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| type | string | `"serverAction"` or `"clientAction"` |
| name | string | Action identifier |
| params | object | Static parameters |
| dataRefs | string[] | Data model paths to include with the action |

When a button with `dataRefs` is clicked, the engine resolves each path from the surface's data model and includes the values in the action payload. In LLM transport mode, the action is formatted as a user message and triggers a new conversation turn, allowing the LLM to respond by updating the UI.

---

## Data Model

Each surface has an independent data model — a JSON document addressed by JSON Pointers.

### Operations

- **Get** `pointer` — returns value at path, or not-found
- **Set** `pointer, value` — creates intermediate objects as needed, returns changed paths
- **Delete** `pointer` — removes value, shrinks arrays properly, returns changed paths

### Binding Propagation

1. Data model changes (from `updateDataModel` or user input via `dataBinding`)
2. Engine collects all changed paths
3. `BindingTracker.Affected(changedPaths)` finds components bound to overlapping paths
4. Those components are re-resolved and re-rendered
5. For user-input bindings, the source component is excluded from re-render

Path overlap: `/a` and `/a/b` overlap (parent-child). `/a` and `/b` do not.

---

## Rendering Pipeline

1. **Parse** — JSONL line to typed `Message`
2. **Route** — `Session` routes to the correct `Surface` by `surfaceId`
3. **Tree update** — `Surface.tree.Update()` stores components, returns changed IDs
4. **Resolve** — `Resolver` evaluates dynamic values against data model, registers bindings
5. **Callback registration** — interactive components get callbacks registered (old ones unregistered first)
6. **Topological sort** — changed components sorted leaves-first
7. **Dispatch to main thread** — two-pass render:
   - Pass 1: create or update each view (leaves first ensures children exist)
   - Pass 2: set children on containers
   - Set root view (single root directly, multiple roots wrapped in Column)

---

## Reserved Component Types (Not Yet Implemented)

| Type | Phase | Description |
|------|-------|-------------|
| Tabs | 3 | Tabbed container |
| Modal | 3 | Modal dialog overlay |
| Video | 3 | AVPlayerView video playback |
| AudioPlayer | 3 | Audio playback controls |

Props for these types are parsed but not rendered. The protocol types and JSON structs are already defined.
