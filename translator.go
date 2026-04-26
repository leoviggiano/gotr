package gotr

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/leoviggiano/gotr/internal/parser"
	"github.com/leoviggiano/gotr/internal/scanner"
)

type Translator interface {
	Register(identifier, jsonPath string) error
	Get(args Args) string
}

type translator struct {
	mu                sync.RWMutex
	defaultIdentifier string
	templates         map[string]map[string]template
}

type option func(*translator) error

var (
	errDefaultAlreadyRegistered = errors.New("default identifier already registered")
)

func WithDefault(identifier, jsonPath string) option {
	return func(t *translator) error {
		if t.defaultIdentifier != "" && identifier != t.defaultIdentifier {
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

	jsonTree, err := scanner.Scan(v)
	if err != nil {
		return err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	lang, ok := t.templates[identifier]
	if !ok {
		lang = make(map[string]template)
		t.templates[identifier] = lang
	}

	isDefault := t.defaultIdentifier == identifier

	for _, path := range jsonTree {
		k, err := parser.Parse(v, path)
		if err != nil {
			return err
		}

		tpl, err := newTemplate(k)
		if err != nil {
			return err
		}

		lang[path] = tpl

		if isDefault {
			lang[tpl.Singular] = tpl
			// Backfill non-default identifiers that were registered before the default.
			for otherId, otherLang := range t.templates {
				if otherId == identifier {
					continue
				}
				if otherTpl, ok := otherLang[path]; ok {
					otherLang[tpl.Singular] = otherTpl
				}
			}
		} else {
			defaultLang, ok := t.templates[t.defaultIdentifier]
			if ok {
				if defaultTpl, ok := defaultLang[path]; ok {
					lang[defaultTpl.Singular] = tpl
				}
			}
		}
	}

	return nil
}

// Get returns the translation for the given path or text and identifier.
func (t *translator) Get(args Args) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	identifiedTranslator, ok := t.templates[args.Identifier]
	if !ok {
		return t.defaultGet(args)
	}

	tpl, ok := identifiedTranslator[args.Localizer]
	if !ok {
		return t.defaultGet(args)
	}

	return tpl.apply(args)
}

func (t *translator) defaultGet(args Args) string {
	identifiedTranslator, ok := t.templates[t.defaultIdentifier]
	if !ok {
		return args.apply(args.Localizer)
	}

	tpl, ok := identifiedTranslator[args.Localizer]
	if !ok {
		return args.apply(args.Localizer)
	}

	return tpl.apply(args)
}
