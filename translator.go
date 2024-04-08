package gotr

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/leoviggiano/gotr/internal/parser"
	"github.com/leoviggiano/gotr/internal/scanner"
)

type Translator interface {
	Register(identifier, jsonPath string) error
	Get(args Args) string
}

type translator struct {
	defaultIdentifier string
	templates         map[string]map[string]template
}

type option func(*translator) error

var (
	errDefaultAlreadyRegistered = errors.New("default identifier already registered")
)

func WithDefault(identifier, jsonPath string) option {
	return func(t *translator) error {
		if t.defaultIdentifier != "" {
			return errDefaultAlreadyRegistered
		}

		t.defaultIdentifier = identifier
		return t.Register(identifier, jsonPath)
	}
}

func NewTranslator(options ...option) (Translator, error) {
	t := &translator{
		templates: make(map[string]map[string]template),
	}

	for _, option := range options {
		err := option(t)
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}

func (t *translator) Register(identifier, jsonPath string) error {
	file, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	var v map[string]any
	err = json.Unmarshal(file, &v)
	if err != nil {
		return err
	}

	jsonTree := scanner.Scan(v)
	translator := make(map[string]template)

	for _, path := range jsonTree {
		k, err := parser.Parse(v, path)
		if err != nil {
			return err
		}

		tpl, err := newTemplate(k)
		if err != nil {
			return err
		}

		translator[path] = tpl

		if t.defaultIdentifier == identifier {
			translator[tpl.Singular] = tpl
		}

		defaultTemplate, ok := t.templates[t.defaultIdentifier][path]
		if ok {
			translator[defaultTemplate.Singular] = tpl
		}
	}

	t.templates[identifier] = translator
	return nil
}

// Get the translation by the given path or text and identifier.
func (t *translator) Get(args Args) string {
	identifiedTranslator, ok := t.templates[args.Identifier]
	if !ok {
		return t.defaultGet(args)
	}

	template, ok := identifiedTranslator[args.Localizer]
	if !ok {
		return t.defaultGet(args)
	}

	return template.apply(args)
}

func (t *translator) defaultGet(args Args) string {
	identifiedTranslator, ok := t.templates[t.defaultIdentifier]
	if !ok {
		return args.apply(args.Localizer)
	}

	template, ok := identifiedTranslator[args.Localizer]
	if !ok {
		return args.apply(args.Localizer)
	}

	return template.apply(args)
}
