package engine

import (
	"encoding/json"
	"strings"
)

// substituteParams walks a JSON tree (interface{} from json.Unmarshal),
// replacing {"param":"name"} nodes with values from args.
func substituteParams(val interface{}, args map[string]interface{}) interface{} {
	switch v := val.(type) {
	case map[string]interface{}:
		// Check if this is a {"param":"name"} node
		if paramName, ok := v["param"].(string); ok && len(v) == 1 {
			if replacement, exists := args[paramName]; exists {
				return deepCopyJSON(replacement)
			}
			return val
		}
		result := make(map[string]interface{}, len(v))
		for k, child := range v {
			result[k] = substituteParams(child, args)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, child := range v {
			result[i] = substituteParams(child, args)
		}
		return result
	default:
		return val
	}
}

// rewriteComponentIDs rewrites componentId, parentId, and children.static
// references in a slice of raw JSON component maps using the provided ID map.
func rewriteComponentIDs(trees []map[string]interface{}, idMap map[string]string) {
	for _, tree := range trees {
		if cid, ok := tree["componentId"].(string); ok {
			if mapped, ok := idMap[cid]; ok {
				tree["componentId"] = mapped
			}
		}
		if pid, ok := tree["parentId"].(string); ok {
			if mapped, ok := idMap[pid]; ok {
				tree["parentId"] = mapped
			}
		}
		// Rewrite children: can be array ["a","b"] or {"static":["a","b"]}
		if children, ok := tree["children"]; ok {
			tree["children"] = rewriteChildrenIDs(children, idMap)
		}
	}
}

// rewriteChildrenIDs rewrites child ID references in either array or object format.
func rewriteChildrenIDs(children interface{}, idMap map[string]string) interface{} {
	switch c := children.(type) {
	case []interface{}:
		result := make([]interface{}, len(c))
		for i, id := range c {
			if s, ok := id.(string); ok {
				if mapped, ok := idMap[s]; ok {
					result[i] = mapped
				} else {
					result[i] = s
				}
			} else {
				result[i] = id
			}
		}
		return result
	case map[string]interface{}:
		result := make(map[string]interface{}, len(c))
		for k, v := range c {
			if k == "static" {
				result[k] = rewriteChildrenIDs(v, idMap)
			} else {
				result[k] = v
			}
		}
		return result
	default:
		return children
	}
}

// rewriteScopedPaths replaces paths starting with "$/" with the scope prefix.
// Walks the entire JSON tree, handling:
// - {"path":"$/foo"} → {"path":"/scope/foo"}
// - "dataBinding":"$/foo" → "dataBinding":"/scope/foo"
// - "forEach":"$/items" → "forEach":"/scope/items"
// - Any string value starting with "$/" in known keys
func rewriteScopedPaths(val interface{}, scope string) interface{} {
	switch v := val.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(v))
		for k, child := range v {
			result[k] = rewriteScopedPaths(child, scope)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, child := range v {
			result[i] = rewriteScopedPaths(child, scope)
		}
		return result
	case string:
		if strings.HasPrefix(v, "$/") {
			return scope + v[1:]
		}
		return v
	default:
		return val
	}
}

// deepCopyJSON creates a deep copy of a JSON-compatible value tree.
func deepCopyJSON(val interface{}) interface{} {
	switch v := val.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(v))
		for k, child := range v {
			result[k] = deepCopyJSON(child)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, child := range v {
			result[i] = deepCopyJSON(child)
		}
		return result
	default:
		return val
	}
}

// jsonToMaps parses raw JSON messages into maps for manipulation.
func jsonToMaps(raw []json.RawMessage) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, len(raw))
	for i, r := range raw {
		var m map[string]interface{}
		if err := json.Unmarshal(r, &m); err != nil {
			return nil, err
		}
		result[i] = m
	}
	return result, nil
}
