package parser

import (
	_ "embed"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	//go:embed 1_level_path_test.json
	oneLevelPath []byte

	//go:embed 2_level_path_test.json
	twoLevelPath []byte

	//go:embed 3_level_path_test.json
	threeLevelPath []byte

	//go:embed multiple_levels_path_test.json
	multipleLevelsPath []byte
)

func TestParse(t *testing.T) {
	var (
		oneLevelJSON       any
		twoLevelJSON       any
		threeLevelJSON     any
		multipleLevelsJSON any
	)

	err := json.Unmarshal(oneLevelPath, &oneLevelJSON)
	require.NoError(t, err)

	err = json.Unmarshal(twoLevelPath, &twoLevelJSON)
	require.NoError(t, err)

	err = json.Unmarshal(threeLevelPath, &threeLevelJSON)
	require.NoError(t, err)

	err = json.Unmarshal(multipleLevelsPath, &multipleLevelsJSON)
	require.NoError(t, err)

	expected := []byte(`{"test":true}`)

	tt := []struct {
		name        string
		currentJSON any
		mapPath     string
		err         error
	}{
		{name: "1 level path", currentJSON: oneLevelJSON, mapPath: "test"},
		{name: "2 level path", currentJSON: twoLevelJSON, mapPath: "test2.test"},
		{name: "3 level path", currentJSON: threeLevelJSON, mapPath: "test3.test2.test"},
		{name: "1 - multiple levels path", currentJSON: multipleLevelsJSON, mapPath: "test"},
		{name: "2 - multiple levels path", currentJSON: multipleLevelsJSON, mapPath: "test2.test"},
		{name: "3 - multiple levels path", currentJSON: multipleLevelsJSON, mapPath: "test3.test2"},
		{name: "invalid path", currentJSON: oneLevelJSON, mapPath: "invalid", err: ErrInvalidPath},
		{name: "invalid type", currentJSON: "invalid", mapPath: "test", err: ErrInvalidType},
		{name: "empty path", err: ErrEmptyPath},
		{name: "empty json", mapPath: "test", err: ErrInvalidPath},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Parse(tc.currentJSON, tc.mapPath)

			if tc.err != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tc.err))
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, result)
			require.Equal(t, expected, result)
		})
	}
}
