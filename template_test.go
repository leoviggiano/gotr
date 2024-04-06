package main

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTemplate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		jsonTemplate := []byte(`
		{
			"description": "Armor text",
			"singular": "{{.Name}} has {{.Count}} Armor.",
			"plural": "{{.Name}} has {{.Count}} Armors.",
			"none": "{{.Name}} has no Armor."
		}`)

		tpl, err := newTemplate(jsonTemplate)
		require.NoError(t, err)

		jsonMap := map[string]any{}
		err = json.Unmarshal(jsonTemplate, &jsonMap)
		require.NoError(t, err)

		require.Equal(t, tpl.Singular, jsonMap["singular"])
		require.Equal(t, tpl.Plural, jsonMap["plural"])
		require.Equal(t, tpl.None, jsonMap["none"])
	})

	t.Run("success - only key/value data", func(t *testing.T) {
		jsonTemplate := []byte(`{"key": "value"}`)

		tpl, err := newTemplate(jsonTemplate)
		require.NoError(t, err)

		require.Equal(t, tpl.Singular, "value")
		require.Equal(t, tpl.Plural, "value")
		require.Equal(t, tpl.None, "value")
	})

	ttErrors := []struct {
		name          string
		jsonData      []byte
		expectedError error
	}{
		{name: "invalid json", jsonData: []byte("invalid"), expectedError: errInvalidJSON},
		{name: "error extract value: many levels json", jsonData: []byte(`{"value": { "test": "value" }}`), expectedError: errInvalidValue},
		{name: "empty json", jsonData: []byte(`{}`), expectedError: errInvalidValue},
		{name: "empty value", jsonData: []byte(`{"value": ""}`), expectedError: errInvalidValue},
	}

	for _, tt := range ttErrors {
		t.Run(tt.name, func(t *testing.T) {
			tpl, err := newTemplate(tt.jsonData)
			require.Error(t, err)
			require.True(t, errors.Is(err, tt.expectedError))
			require.Equal(t, template{}, tpl)
		})
	}
}

func TestTemplate_extractValue(t *testing.T) {
	tt := []struct {
		name          string
		jsonData      []byte
		expectedValue string
		expectedError error
	}{
		{name: "success - string", jsonData: []byte(`{"value": "string"}`), expectedValue: "string"},
		{name: "success - number", jsonData: []byte(`{"value": 10}`), expectedValue: "10"},
		{name: "success - bool", jsonData: []byte(`{"value": true}`), expectedValue: "true"},
		{name: "success - empty json", jsonData: []byte(`{}`), expectedValue: ""},
		{name: "error - many levels json", jsonData: []byte(`{"value": { "test": "value" }}`), expectedError: errInvalidValue},
		{name: "error - invalid json", jsonData: []byte("invalid"), expectedError: errInvalidJSON},
	}

	for _, tt := range tt {
		t.Run(tt.name, func(t *testing.T) {
			tpl := template{}
			value, err := tpl.extractValue(tt.jsonData)
			require.Equal(t, tt.expectedValue, value)
			require.True(t, errors.Is(err, tt.expectedError))
		})
	}
}

func TestTemplate_apply(t *testing.T) {
	tpl := template{
		Singular: "{{.Name}} has {{.Count}} Armor.",
		Plural:   "{{.Name}} has {{.Count}} Armors.",
		None:     "{{.Name}} has no Armor.",
	}

	tt := []struct {
		name          string
		expectedValue string
		args          Args
	}{
		{
			name:          "none",
			expectedValue: "John has no Armor.",
			args: Args{
				Args: map[string]any{
					"Name": "John",
				},
			},
		},
		{
			name:          "singular",
			expectedValue: "John has 1 Armor.",
			args: Args{
				Count: 1,
				Args: map[string]any{
					"Name":  "John",
					"Count": 1,
				},
			},
		},
		{
			name:          "plural",
			expectedValue: "John has 2 Armors.",
			args: Args{
				Count: 2,
				Args: map[string]any{
					"Name":  "John",
					"Count": 2,
				},
			},
		},
	}

	for _, tt := range tt {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedValue, tpl.apply(tt.args))
		})
	}
}
