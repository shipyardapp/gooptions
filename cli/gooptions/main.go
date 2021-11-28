package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shipyardapp/gooptions/model"
)

func main() {
	f, err := NewFlags(os.Args[1:])
	if err != nil {
		exit(err, 1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		exit(fmt.Errorf("failed to get cwd: %v", err), 2)
	}

	sourceDir := filepath.Join(cwd, f.SourceDir)

	modelPackage, err := NewModelPackageReflect(cwd, sourceDir)
	if err != nil {
		exit(fmt.Errorf("failed to load source package information: %v", err), 3)
	}

	modelStructType, err := BuildRunReflectProgram(modelPackage, f.Type)
	if err != nil {
		exit(fmt.Errorf("failed to generate model from reflection: %v", err), 4)
	}

	modelModel := model.NewModel(
		model.NewOptions(),
		modelPackage,
		modelStructType,
	)

	err = model.Generate(modelModel, f.Type, sourceDir, f.DestinationPath)
	if err != nil {
		exit(err, 5)
	}
}

type Flags struct {
	SourceDir       string
	Type            string
	DestinationPath string
}

func NewFlags(args []string) (*Flags, error) {
	f := &Flags{
		SourceDir:       ".",
		Type:            "",
		DestinationPath: "",
	}

	fs := flag.NewFlagSet("gooptions", flag.ExitOnError)

	fs.StringVar(&f.SourceDir, "source", f.SourceDir, "source package to generate options for types in the package")
	fs.StringVar(&f.Type, "type", f.Type, "name of struct type to generate options for")
	fs.StringVar(&f.DestinationPath, "dest", "", `destination file path to write options file to (default: empty value means "<os.Getwd()>/<strings.ToLower(type)>_options.go")`)

	err := fs.Parse(args)
	if err != nil {
		fs.Usage()
		return nil, err
	}
	return f, nil
}

func exit(err error, exitCode int) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(exitCode)
}
