package gotr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTranslator(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		translator := &translator{
			templates: make(map[string]map[string]template),
		}
		identifier := "identifier"
		applyOption := WithDefault(identifier, "./translations/en_US.json")

		err := applyOption(translator)
		require.NoError(t, err)
		require.NotNil(t, translator)
		require.Equal(t, identifier, translator.defaultIdentifier)
	})
}

func TestTranslator_Register(t *testing.T) {
	t.Run("success - with default", func(t *testing.T) {
		translator := &translator{
			templates: make(map[string]map[string]template),
		}
		identifier := "identifier"
		applyOption := WithDefault(identifier, "./translations/en_US.json")

		err := applyOption(translator)

		require.NoError(t, err)
		require.NotEmpty(t, translator.templates[identifier])
	})

	t.Run("error - default already registered", func(t *testing.T) {
		translator := &translator{
			templates: make(map[string]map[string]template),
		}
		identifier := "identifier"

		applyOption := WithDefault(identifier, "./translations/en_US.json")

		err := applyOption(translator)
		require.NoError(t, err)

		translator.defaultIdentifier = "different"
		err = applyOption(translator)
		require.Equal(t, errDefaultAlreadyRegistered, err)
	})

	t.Run("error read file - with default", func(t *testing.T) {
		translator, err := NewTranslator(
			WithDefault("en", "invalid"),
		)

		require.Error(t, err)
		require.Nil(t, translator)
	})

	t.Run("error unmarshal file - with default", func(t *testing.T) {
		translator, err := NewTranslator(
			WithDefault("en", "./translator_test.go"),
		)
		require.Error(t, err)
		require.Nil(t, translator)
	})

	t.Run("success", func(t *testing.T) {
		translator := &translator{
			templates: make(map[string]map[string]template),
		}
		identifier := "identifier"
		applyOption := WithDefault(identifier, "./translations/en_US.json")

		err := applyOption(translator)
		require.NoError(t, err)

		err = translator.Register(identifier, "./translations/en_US.json")
		require.NoError(t, err)
		require.NotEmpty(t, translator.templates[identifier])
	})

	t.Run("error read file", func(t *testing.T) {
		translator := &translator{
			templates: make(map[string]map[string]template),
		}
		identifier := "identifier"

		err := translator.Register(identifier, "invalid")
		require.Error(t, err)
		require.Nil(t, translator.templates[identifier])
	})

	t.Run("error unmarshal file", func(t *testing.T) {
		translator := &translator{
			templates: make(map[string]map[string]template),
		}
		identifier := "identifier"

		err := translator.Register(identifier, "./translator_test.go")
		require.Error(t, err)
		require.Nil(t, translator.templates[identifier])
	})
}

func TestTranslator_Get(t *testing.T) {
	translator, err := NewTranslator(
		WithDefault("en", "./translations/en_US.json"),
	)
	require.NoError(t, err)

	err = translator.Register("pt", "./translations/pt_BR.json")
	require.NoError(t, err)

	tt := []struct {
		name       string
		localizer  string
		identifier string
		expected   string
	}{
		// get with json path
		{name: "success - pt", localizer: "hello_world", identifier: "pt", expected: "Olá Mundo"},
		{name: "success - en", localizer: "hello_world", identifier: "en", expected: "Hello World"},
		{name: "pt fallback to default", localizer: "hello_world2", identifier: "pt", expected: "Hello World 2"},
		{name: "empty path", localizer: "invalid.path", expected: "invalid.path"},

		// get with texts
		{name: "success - pt", localizer: "Hello World", identifier: "pt", expected: "Olá Mundo"},
		{name: "success - en", localizer: "Hello World", identifier: "en", expected: "Hello World"},
		{name: "pt fallback to default", localizer: "Hello World 2", identifier: "pt", expected: "Hello World 2"},
		{name: "empty path", localizer: "invalid.path", expected: "invalid.path"},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			args := Args{
				Localizer:  test.localizer,
				Identifier: test.identifier,
			}

			require.Equal(t, test.expected, translator.Get(args))
		})
	}
}

func TestTranslator_defaultGet(t *testing.T) {
	translator := translator{
		defaultIdentifier: "en",
		templates:         make(map[string]map[string]template),
	}

	err := translator.Register("en", "./translations/en_US.json")
	require.NoError(t, err)

	err = translator.Register("pt", "./translations/pt_BR.json")
	require.NoError(t, err)

	tt := []struct {
		name              string
		localizer         string
		expected          string
		defaultIdentifier string
	}{
		{name: "success", localizer: "hello_world", expected: "Hello World"},
		{name: "empty path", localizer: "invalid.path", expected: "invalid.path"},
		{name: "translator not found", defaultIdentifier: "invalid", localizer: "hello_world", expected: "hello_world"},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			args := Args{
				Localizer:  test.localizer,
				Identifier: "pt",
			}

			if test.defaultIdentifier != "" {
				oldIdentifier := translator.defaultIdentifier
				defer func() { translator.defaultIdentifier = oldIdentifier }()

				translator.defaultIdentifier = test.defaultIdentifier
			}

			require.Equal(t, test.expected, translator.defaultGet(args))
		})
	}
}
