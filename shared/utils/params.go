package utils

// GetStringSlice extracts a []string from a params map
func GetStringSlice(params map[string]interface{}, key string) []string {
	raw, ok := params[key]
	if !ok {
		return nil
	}

	switch v := raw.(type) {
	case []interface{}:
		var result []string
		for _, item := range v {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	case []string:
		return v
	default:
		return nil
	}
}

// GetString extracts a string from a params map
func GetString(params map[string]interface{}, key string) string {
	if val, ok := params[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetBool extracts a bool from a params map
func GetBool(params map[string]interface{}, key string) bool {
	if val, ok := params[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

// GetInt extracts an int from a params map
func GetInt(params map[string]interface{}, key string) int {
	if val, ok := params[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64: // YAML numbers default to float64
			return int(v)
		}
	}
	return 0
}

// GetFloat extracts a float64 from a params map
func GetFloat(params map[string]interface{}, key string) float64 {
	if val, ok := params[key]; ok {
		if f, ok := val.(float64); ok {
			return f
		}
	}
	return 0.0
}