package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidPath = errors.New("invalid path or invalid json")
	ErrInvalidType = errors.New("not a json")
	ErrEmptyPath   = errors.New("empty path")
)

func Parse(currentJSON any, mapPath string) ([]byte, error) {
	if mapPath == "" {
		return nil, ErrEmptyPath
	}

	pathSlice := strings.Split(mapPath, ".")

	return parse(currentJSON, pathSlice, 0)
}

func parse(currentJSON any, pathSlice []string, idx int) ([]byte, error) {
	if currentJSON == nil {
		return nil, fmt.Errorf("%w: mapPath: %s, currentPath: %s", ErrInvalidPath, strings.Join(pathSlice, "."), pathSlice[idx])
	}

	v, ok := currentJSON.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrInvalidType, currentJSON)
	}

	if idx == len(pathSlice)-1 {
		selectedJSON := v[pathSlice[idx]]
		if selectedJSON == nil {
			return nil, fmt.Errorf("%w: mapPath: %s, currentPath: %s", ErrInvalidPath, strings.Join(pathSlice, "."), pathSlice[idx])
		}

		switch v := selectedJSON.(type) {
		case map[string]any:
			return json.Marshal(v)
		default:
			return json.Marshal(map[string]any{pathSlice[idx]: v})
		}
	}

	return parse(v[pathSlice[idx]], pathSlice, idx+1)
}
