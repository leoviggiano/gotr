package main

import (
	"fmt"

	"github.com/leoviggiano/gotr"
)

func main() {
	translator, err := gotr.NewTranslator(
		gotr.WithDefault("en", "./translations/en_US.json"),
		gotr.WithDefault("en", "./translations/en_US_items.json"),
	)
	if err != nil {
		fmt.Println(err)
	}

	err = translator.Register("pt", "./translations/pt_BR.json")
	if err != nil {
		fmt.Println(err)
	}

	argsPTArmorDescription := gotr.Args{
		Identifier: "pt",
		Localizer:  "items.equipments.armor.description",
		Args: map[string]any{
			"Name":  "John",
			"Count": 10,
		},
		Count: 10,
	}

	argsPTArmorText := gotr.Args{
		Identifier: "pt",
		Localizer:  "items.equipments.armor",
		Args: map[string]any{
			"Name":  "John",
			"Count": 10,
		},
		Count: 10,
	}

	argsPTArmorFullText := gotr.Args{
		Identifier: "pt",
		Localizer:  "{{.Name}} has {{.Count}} Armor.",
		Args: map[string]any{
			"Name":  "John",
			"Count": 10,
		},
		Count: 10,
	}

	argsPTText2 := gotr.Args{
		Identifier: "pt",
		Localizer:  "hello_world",
		Count:      1,
	}

	fmt.Println(translator.Get(argsPTArmorText))
	fmt.Println(translator.Get(argsPTArmorDescription))
	fmt.Println(translator.Get(argsPTArmorFullText))
	fmt.Println(translator.Get(argsPTText2))
}
