package main

import (
	"bytes"
	_ "embed"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/shipyardapp/gooptions/model"
	"golang.org/x/tools/go/packages"
)

func NewModelPackageReflect(cwd, pattern string) (*model.Package, error) {
	packages, err := packages.Load(
		&packages.Config{
			Mode:  packages.NeedName,
			Tests: false,
			Dir:   cwd,
		},
		pattern,
	)
	if err != nil {
		return nil, err
	}

	if len(packages) != 1 {
		return nil, fmt.Errorf("wrong number of packages found %v %v", len(packages), packages)
	}

	p := packages[0]
	return &model.Package{
		Name: p.Name,
		Path: p.PkgPath,
	}, nil
}

//go:embed modelreflect/main.go
var ModelReflectMainGo []byte

//go:embed modelreflect/variable.gotemplate
var ModelReflectVariableGoTemplate string

func BuildRunReflectProgram(mp *model.Package, typeName string) (*model.StructType, error) {
	t, err := template.New("modelreflect").Parse(ModelReflectVariableGoTemplate)
	if err != nil {
		return nil, err
	}

	templateData := map[string]string{
		"Path":     mp.Path,
		"TypeName": typeName,
	}

	variableGo := &bytes.Buffer{}
	if err := t.Execute(variableGo, templateData); err != nil {
		return nil, err
	}

	return BuildRunInDirectory(variableGo.Bytes(), ".")
}

func BuildRunInDirectory(variableGo []byte, dir string) (*model.StructType, error) {
	mainDir, err := GenerateModelReflectMainDirectory(variableGo, dir)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := os.RemoveAll(mainDir); err != nil {
			log.Printf("failed to remove temp directory: %v", mainDir)
		}
	}()

	programName := "modelreflect.bin"
	programBinary := filepath.Join(mainDir, programName)

	cmdArgs := []string{"build", "-o", programName, "."}
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = mainDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return RunProgram(programBinary)
}

func GenerateModelReflectMainDirectory(variableGo []byte, dir string) (string, error) {
	tempDir, err := ioutil.TempDir(dir, "modelreflect")
	if err != nil {
		return "", err
	}

	if err := ioutil.WriteFile(
		filepath.Join(tempDir, "variable.go"),
		variableGo,
		0666,
	); err != nil {
		return "", os.RemoveAll(tempDir)
	}

	if err := ioutil.WriteFile(
		filepath.Join(tempDir, "main.go"),
		ModelReflectMainGo,
		0666,
	); err != nil {
		return "", os.RemoveAll(tempDir)
	}

	return tempDir, nil
}

func RunProgram(programBinary string) (*model.StructType, error) {
	outputFile, err := ioutil.TempFile("", "modelreflectoutput")
	if err != nil {
		return nil, err
	}

	defer os.Remove(outputFile.Name())
	if err := outputFile.Close(); err != nil {
		return nil, err
	}

	cmd := exec.Command(programBinary, "-output", outputFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	f, err := os.Open(outputFile.Name())
	if err != nil {
		return nil, err
	}

	var modelStructType model.StructType
	if err := gob.NewDecoder(f).Decode(&modelStructType); err != nil {
		return nil, err
	}

	if err := f.Close(); err != nil {
		return nil, err
	}

	return &modelStructType, nil
}
