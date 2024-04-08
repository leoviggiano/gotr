# GOTR

[![workflow](https://github.com/leoviggiano/gotr/actions/workflows/ci.yml/badge.svg?event=push)](https://github.com/leoviggiano/gotr/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/leoviggiano/gotr/graph/badge.svg?token=KUC1RFVHET)](https://codecov.io/gh/leoviggiano/gotr)

This library provides functionality for translation in Go, similar to i18n. It allows you to manage translation files in JSON format and easily retrieve translated texts within your Go applications.

## JSON Translation Template

Translation files follow a JSON structure as shown below:

### English
```json
{
    "hello_world": "Hello World",
    "texts": {
        "welcome": "Welcome to my translator!",
        "goodbye": "Goodbye!"
    },
    "items": {
        "equipments": {
            "armor": {
                "description": "Armor text",
                "singular": "{{.Name}} has {{.Count}} Armor.",
                "plural": "{{.Name}} has {{.Count}} Armors.",
                "none": "{{.Name}} has no Armor."
            }
        }
    }
}
```

### Portuguese
```json
{
    "hello_world": "Olá Mundo",
    "texts": {
        "welcome": "Bem-vindo ao meu tradutor!",
        "goodbye": "Até logo!"
    },
    "items": {
        "equipments": {
            "armor": {
                "description": "Armadura text",
                "singular": "{{.Name}} tem {{.Count}} Armadura.",
                "plural": "{{.Name}} tem {{.Count}} Armaduras.",
                "none": "{{.Name}} não tem Armadura."
            }
        }
    }
}
```

## Installation
To install `gotr`, use `go get`:

```sh
go get github.com/leoviggiano/gotr
```

## How to Use

1. **Import the Library:**

```go
import "github.com/leoviggiano/gotr"
```

2. **Create a New Translator Instance:**

```go
translator, err := gotr.NewTranslator(
    translator.WithDefault("en", "path/to/json"),
)
if err != nil {
    fmt.Println(err)
}
```

3. **Register Additional Translation Files:**

```go
translator.Register("pt", "path/to/json")
```

4. **Define Arguments for Translation:**

```go
// You can use either the text to translate or the json path

argsPTText := gotr.Args{
    Identifier: "pt",
    Localizer:  "{{.Name}} has {{.Count}} Armor.",
    Args: map[string]interface{}{
        "Name":  "John",
        "Count": 10,
    },
    Count: 10,
}

argsPTJSONPath := gotr.Args{
    Identifier: "pt",
    Localizer:  "items.equipments.armor",
    Args: map[string]interface{}{
        "Name": "John",
    },
    Count: 0,
}

// When the translation does not exist on your template, it'll just replace the variables
nonExistentText := gotr.Args{
    Identifier: "pt",
    Localizer:  "{{.Name}} has {{.Count}} Apples",
    Args: map[string]interface{}{
        "Name":  "John",
        "Count": 3,
    },
    Count: 3,
}
```

5. **Retrieve Translated Texts:**

```go
fmt.Println(translator.Get(argsPTText)) // John tem 10 Armaduras.
fmt.Println(translator.Get(argsPTJSONPath)) // John não tem Armadura.
fmt.Println(translator.Get(nonExistentText)) // John has 3 Apples
```

## Notes

- `Args` struct is used to pass arguments for translation.
- Ensure that the translation files are stored in the correct directory and follow the specified JSON format.
- Use appropriate error handling to manage translation file loading and retrieval errors.
- The translation texts can be accessed using their respective keys or paths as defined in the translation JSON files.

Enjoy translating your applications with ease using this library!
