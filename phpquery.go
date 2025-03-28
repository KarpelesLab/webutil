package webutil

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
)

// ParsePhpQuery parses a PHP-compatible query string into a structured map.
//
// The function handles various PHP-style query formats:
//   - a=b (simple key-value)
//   - a[b]=c (nested object)
//   - a[]=c (array)
//   - a[b][]=c (multi-level nesting)
//   - a[][][]=c (multi-level arrays)
//
// The resulting map contains properly structured nested objects and arrays.
func ParsePhpQuery(query string) map[string]any {
	result := make(map[string]any)

	// Split query string on '&' and process each part
	for _, part := range strings.Split(query, "&") {
		if part != "" {
			parsePhpQ(result, part)
		}
	}

	// Process the result to clean up any pointer references
	parsePhpFix(result)
	return result
}

// ConvertPhpQuery converts standard url.Values to a structured map
// using PHP-style array/object notation in the keys.
func ConvertPhpQuery(values url.Values) map[string]any {
	result := make(map[string]any)

	// Process each key-value pair
	for key, vals := range values {
		for _, val := range vals {
			parsePhpQV(result, val, key)
		}
	}

	// Clean up the structure
	parsePhpFix(result)
	return result
}

// parsePhpFix recursively processes the parsed query structure,
// resolving pointer references and ensuring proper nesting.
func parsePhpFix(value any) any {
	switch val := value.(type) {
	case map[string]any:
		// Process each key-value pair in the map
		for k, v := range val {
			val[k] = parsePhpFix(v)
		}
		return val
	case []any:
		// Process each item in the array
		for i, v := range val {
			val[i] = parsePhpFix(v)
		}
		return val
	case *[]any:
		// Dereference and process
		return parsePhpFix(*val)
	default:
		// Other types are returned as-is
		return value
	}
}

// parsePhpQ parses a single key-value part from a query string.
func parsePhpQ(result map[string]any, part string) {
	if part == "" {
		return
	}

	// Split into key and value
	var key, value string
	if eqIdx := strings.IndexByte(part, '='); eqIdx != -1 {
		if eqIdx == 0 {
			// Ignore if key is empty
			return
		}
		// URL-decode the value
		var err error
		value, err = url.QueryUnescape(part[eqIdx+1:])
		if err != nil {
			// Skip malformed values
			return
		}
		key = part[:eqIdx]
	} else {
		// No value, just a key
		key = part
	}

	// URL-decode the key
	decodedKey, err := url.QueryUnescape(key)
	if err != nil {
		// Skip malformed keys
		return
	}

	// Process the key-value pair
	parsePhpQV(result, value, decodedKey)
}

// parsePhpQV processes a key-value pair, handling PHP-style array/object notation in the key.
func parsePhpQV(result map[string]any, value, key string) {
	// Find the first bracket, which indicates array/object syntax
	bracketIdx := strings.IndexByte(key, '[')
	if bracketIdx == -1 {
		// Simple key-value, no brackets
		result[key] = value
		return
	}
	if bracketIdx == 0 {
		// Key can't start with a bracket
		return
	}

	// Extract the base key and parse the array/object path
	baseName := key[:bracketIdx]
	path := parsePhpArrayPath(key[bracketIdx:])
	if len(path) == 0 {
		// Malformed path
		return
	}

	// Start with the base name
	currentPath := []string{baseName}
	currentPath = append(currentPath, path...)

	// Process the path
	processPhpArrayPath(result, currentPath, value)
}

// parsePhpArrayPath extracts the path components from PHP array/object notation.
// For example, "[a][b][]" becomes ["a", "b", ""].
func parsePhpArrayPath(pathStr string) []string {
	if pathStr == "" {
		return nil
	}

	var path []string
	for len(pathStr) >= 2 {
		// Each segment must start with '['
		if pathStr[0] != '[' {
			break
		}

		if pathStr[1] == ']' {
			// Empty segment "[]" means array
			path = append(path, "")
			pathStr = pathStr[2:]
			continue
		}

		// Find closing bracket
		closeBracket := strings.IndexByte(pathStr, ']')
		if closeBracket == -1 {
			// Malformed: no closing bracket
			break
		}

		// Extract key between brackets
		path = append(path, pathStr[1:closeBracket])
		pathStr = pathStr[closeBracket+1:]
	}

	return path
}

// processPhpArrayPath builds the nested structure according to the path components.
func processPhpArrayPath(result map[string]any, path []string, value string) {
	if len(path) == 0 {
		return
	}

	// Extract current and next path components
	current := path[0]
	remaining := path[1:]

	if len(remaining) == 0 {
		// We're at the leaf node, set the value
		result[current] = value
		return
	}

	// Handle based on the next path component
	next := remaining[0]
	isArray := next == ""

	if isArray {
		// Next component is array notation "[]"
		var arr *[]any

		// Try to get existing array or create new one
		if existingArr, ok := result[current].(*[]any); ok {
			arr = existingArr
		} else {
			newArr := make([]any, 0)
			arr = &newArr
			result[current] = arr
		}

		if len(remaining) == 1 {
			// Last component is array, append value
			*arr = append(*arr, value)
		} else {
			// More components after array, create nested structure
			newMap := make(map[string]any)
			*arr = append(*arr, newMap)
			processPhpArrayPath(newMap, remaining[1:], value)
		}
	} else {
		// Next component is named key, creating/updating a map
		var nestedMap map[string]any

		// Try to get existing map or create new one
		if existing, ok := result[current].(map[string]any); ok {
			nestedMap = existing
		} else {
			nestedMap = make(map[string]any)
			result[current] = nestedMap
		}

		// Process the rest of the path in the nested map
		processPhpArrayPath(nestedMap, remaining, value)
	}
}

// EncodePhpQuery converts a structured map back to a PHP-compatible query string.
func EncodePhpQuery(query map[string]any) string {
	var result []byte

	// Process each top-level key
	for key, value := range query {
		result = encodePhpQueryAppend(result, value, key)
	}

	return string(result)
}

// encodePhpQueryAppend recursively builds a query string by appending
// encoded key-value pairs to an existing byte slice.
func encodePhpQueryAppend(result []byte, value any, key string) []byte {
	switch val := value.(type) {
	case map[string]any:
		// Handle nested maps as PHP objects
		for subKey, subValue := range val {
			result = encodePhpQueryAppend(result, subValue, key+"["+subKey+"]")
		}

	case []any:
		// Handle arrays with [] notation
		for _, subValue := range val {
			result = encodePhpQueryAppend(result, subValue, key+"[]")
		}

	case string:
		// Handle string values with proper encoding
		if len(result) > 0 {
			result = append(result, '&')
		}
		result = append(result, url.QueryEscape(key)...)
		result = append(result, '=')
		result = append(result, url.QueryEscape(val)...)

	case []byte:
		// Convert byte slices to strings
		return encodePhpQueryAppend(result, string(val), key)

	case *bytes.Buffer:
		// Convert buffers to strings
		return encodePhpQueryAppend(result, string(val.Bytes()), key)

	default:
		// Convert anything else to string representation
		return encodePhpQueryAppend(result, fmt.Sprintf("%+v", val), key)
	}

	return result
}
