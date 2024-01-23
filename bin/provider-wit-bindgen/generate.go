package main

import (
	"fmt"
	"path/filepath"

	"github.com/patrickhuber/go-wasm/wit/ast"
)

func generate(goPackageFragments []string, inputPath string, outputRoot string, witAST *ast.Ast) error {
	goPackageFragments = append(goPackageFragments, witAST.PackageDeclaration.Unwrap().Name)
	goPackage := filepath.Join(goPackageFragments...)
	fmt.Println("PACKAGE", goPackage)

	for _, item := range witAST.Items {
		// Only interfaces are currently supported.
		if item.Interface == nil {
			continue
		}

		fmt.Println("INTERFACE", item.Interface.Name)
	}
	return nil
}
