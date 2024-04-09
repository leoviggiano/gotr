package scanner

import (
	"fmt"
	"strings"
)

func Scan(currentJSON any) []string {
	paths := []string{}

	for k, v := range currentJSON.(map[string]any) {
		switch v := v.(type) {
		case map[string]any:
			newPaths := scan(v, k)
			paths = append(paths, newPaths...)
		default:
			paths = append(paths, k)
		}
	}

	return paths
}

func scan(currentJSON map[string]any, path string) []string {
	mapPaths := make(map[string]struct{})
	paths := []string{}

	for k, v := range currentJSON {
		newPath := fmt.Sprintf("%s.%s", path, k)

		switch v := v.(type) {
		case map[string]any:
			newPaths := scan(v, fmt.Sprintf("%s.%s", path, k))
			for _, p := range newPaths {
				mapPaths[p] = struct{}{}
			}

		default:
			if strings.Contains(newPath, "singular") ||
				strings.Contains(newPath, "plural") ||
				strings.Contains(newPath, "none") {
				mapPaths[path] = struct{}{}
				continue
			}

			mapPaths[newPath] = struct{}{}
		}
	}

	for k := range mapPaths {
		paths = append(paths, k)
	}

	return paths
}
