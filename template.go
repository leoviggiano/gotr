package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type template struct {
	Singular string `json:"singular"`
	Plural   string `json:"plural"`
	None     string `json:"none"`
}

var (
	errInvalidJSON  = errors.New("invalid json")
	errInvalidValue = errors.New("invalid value")
)

func newTemplate(data []byte) (template, error) {
	var tpl template
	err := json.Unmarshal(data, &tpl)
	if err != nil {
		return template{}, fmt.Errorf("%w: %v", errInvalidJSON, err)
	}

	if tpl.empty() {
		value, err := tpl.extractValue(data)
		if err != nil {
			return template{}, err
		}

		if value == "" {
			return template{}, fmt.Errorf("%w: value is empty: %s", errInvalidValue, string(data))
		}

		return template{
			Singular: value,
			Plural:   value,
			None:     value,
		}, nil
	}

	return tpl, nil
}

func (t template) extractValue(data []byte) (string, error) {
	var dataJSON map[string]any
	err := json.Unmarshal(data, &dataJSON)
	if err != nil {
		return "", fmt.Errorf("%w: %v", errInvalidJSON, err)
	}

	for _, v := range dataJSON {
		switch v := v.(type) {
		case map[string]any:
			return "", fmt.Errorf("%w: given wrong value, must be a string, bool or number: value: %v", errInvalidValue, v)
		default:
			return fmt.Sprintf("%v", v), nil
		}
	}

	return "", nil
}

func (t template) apply(args Args) string {
	var str string

	switch args.Count {
	case 0:
		str = fmt.Sprintf(t.None)
	case 1:
		str = fmt.Sprintf(t.Singular)
	default:
		str = fmt.Sprintf(t.Plural)
	}

	return args.apply(str)
}

func (t template) empty() bool {
	return (t.Singular == "" && t.Plural == "" && t.None == "")
}
