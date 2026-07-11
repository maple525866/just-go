package main

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestCoreLayersDoNotImportOuterLayers(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate test file")
	}
	chapter := filepath.Dir(file)
	for _, layer := range []string{"domain", "usecase"} {
		packages, err := parser.ParseDir(token.NewFileSet(), filepath.Join(chapter, layer), nil, parser.ImportsOnly)
		if err != nil {
			t.Fatal(err)
		}
		for _, pkg := range packages {
			for filename, parsed := range pkg.Files {
				for _, spec := range parsed.Imports {
					path, err := strconv.Unquote(spec.Path.Value)
					if err != nil {
						t.Fatal(err)
					}
					if strings.Contains(path, "/12-clean-architecture/interface/") || strings.Contains(path, "/12-clean-architecture/infrastructure/") {
						position := parsed.Pos()
						if spec.Pos().IsValid() {
							position = spec.Pos()
						}
						t.Errorf("%s (%v) imports outer layer %q", filename, position, path)
					}
				}
			}
		}
	}
}
