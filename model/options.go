package model

import (
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	OptionName string

	OptionPrefix string
}

func NewOptions() *Options {
	return &Options{
		OptionName:   "Option",
		OptionPrefix: "With",
	}
}

func (o *Options) OutputFile(typeName, sourceDir, destinationPath string) (string, error) {
	// TODO detect already exists.

	if destinationPath == "" {
		destinationPath = strings.ToLower(typeName + "_options.go")
	}
	if !filepath.IsAbs(destinationPath) {
		destinationPath = filepath.Join(sourceDir, destinationPath)
	}

	if err := os.MkdirAll(filepath.Dir(destinationPath), 0777); err != nil {
		return "", err
	}

	return destinationPath, nil
}
