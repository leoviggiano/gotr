package scanner

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScan(t *testing.T) {
	testJSON := []byte(`
	{
		"value": "ok",
		"test": {
			"test1-2": {
				"value": "ok",
				"test3": {
					"value": "ok"
				}
			}
		},
		"test2": {
			"value": "ok",
			"test3": {
				"value": "ok"
			}
		}
	}`)

	expected := []string{
		"value",
		"test.test1-2.value",
		"test.test1-2.test3",
		"test2.value",
		"test2.test3.value",
	}

	t.Run("success", func(t *testing.T) {
		var currentJSON map[string]any
		err := json.Unmarshal(testJSON, &currentJSON)
		require.NoError(t, err)

		paths := Scan(currentJSON)
		require.ElementsMatch(t, expected, paths)
	})
}
