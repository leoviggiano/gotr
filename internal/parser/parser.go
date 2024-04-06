package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	ErrInvalidPath = errors.New("invalid path or invalid json")
	ErrInvalidType = errors.New("not a json")
	ErrEmptyPath   = errors.New("empty path")
)

func Parse(currentJSON any, mapPath string) ([]byte, error) {
	currentPath := strings.Split(mapPath, ".")
	if len(currentPath) == 0 || currentPath[0] == "" {
		return nil, ErrEmptyPath
	}

	return parse(currentJSON, mapPath, currentPath[0])
}

func parse(currentJSON any, mapPath, currentPath string) ([]byte, error) {
	if currentJSON == nil {
		return nil, fmt.Errorf("%w: mapPath: %s, currentPath: %s", ErrInvalidPath, mapPath, currentPath)
	}

	v, ok := currentJSON.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrInvalidType, currentJSON)
	}

	pathSlice := strings.Split(mapPath, ".")
	idx := slices.IndexFunc(pathSlice, func(i string) bool {
		return i == currentPath
	})

	if idx == len(pathSlice)-1 {
		selectedJSON := v[pathSlice[idx]]
		if selectedJSON == nil {
			return nil, fmt.Errorf("%w: mapPath: %s, currentPath: %s", ErrInvalidPath, mapPath, currentPath)
		}

		switch v := selectedJSON.(type) {
		case map[string]any:
			return json.Marshal(v)
		default:
			return json.Marshal(map[string]any{pathSlice[idx]: v})
		}
	}

	return parse(v[pathSlice[idx]], mapPath, pathSlice[idx+1])
}
