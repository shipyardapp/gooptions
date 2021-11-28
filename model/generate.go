package model

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"go/token"
	"os"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

var variableNameCounter = 0

//go:embed generate.gotemplate
var GenerateTemplate string

func Generate(m *Model) error {
	t := template.New("generator")
	t = t.Funcs(
		map[string]interface{}{
			"ArgumentName": ArgumentName,
			"ReceiverName": ReceiverName,
			"Title":        strings.Title,
		},
	)

	var err error
	t, err = t.Parse(GenerateTemplate)
	if err != nil {
		return err
	}

	b := &bytes.Buffer{}
	templateData := m
	if err := t.Execute(b, templateData); err != nil {
		return err
	}

	result, err := format.Source(b.Bytes())
	if err != nil {
		return err
	}

	if _, err := os.Stdout.Write(result); err != nil {
		return err
	}

	return nil
}

func ArgumentName(name string) string {
	r, size := utf8.DecodeRuneInString(name)
	result := string(unicode.ToLower(r)) + string(name[size:])
	return SanitizeName(result)
}

func ReceiverName(name string) string {
	r, _ := utf8.DecodeRuneInString(name)
	return string(unicode.ToLower(r))
}

func SanitizeName(name string) string {
	result := name
	if token.IsKeyword(name) || name == "byte" || name == "rune" {
		result = fmt.Sprintf("%s%d", name, variableNameCounter)
		variableNameCounter++
	}
	return result
}
