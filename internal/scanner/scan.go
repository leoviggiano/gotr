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
			newPaths := subScan(v, fmt.Sprintf("%s.%s", path, k))
			for _, p := range newPaths {
				mapPaths[p] = struct{}{}
			}

		default:
			mapPaths[newPath] = struct{}{}
		}
	}

	for k := range mapPaths {
		paths = append(paths, k)
	}

	return paths
}

func subScan(currentJSON map[string]any, path string) []string {
	mapPaths := make(map[string]struct{})
	paths := []string{}

	for k, v := range currentJSON {
		newPath := fmt.Sprintf("%s.%s", path, k)

		switch v := v.(type) {
		case map[string]any:
			newPaths := scan(v, fmt.Sprintf("%s.%s", path, k))
			for _, p := range newPaths {
				splittedPath := strings.Split(p, ".")
				splittedPath = splittedPath[:len(splittedPath)-1]
				newPath = strings.Join(splittedPath, ".")
				mapPaths[newPath] = struct{}{}
			}

		default:
			mapPaths[newPath] = struct{}{}
		}
	}

	for k := range mapPaths {
		paths = append(paths, k)
	}

	return paths
}
