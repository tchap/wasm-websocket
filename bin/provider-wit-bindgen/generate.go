package main

import (
	"fmt"
	"path/filepath"

	"github.com/patrickhuber/go-wasm/wit/ast"
)

func generatePackage(
	goPackageFragments []string,
	outputPathFragments []string,
	pkg *ast.Ast,
) error {
	goPackageFragments = append(goPackageFragments, pkg.PackageDeclaration.Unwrap().Name)

	for _, item := range pkg.Items {
		// Only interfaces are currently supported.
		if item.Interface == nil {
			continue
		}

		if err := generateInterfaceFile(goPackageFragments, item.Interface, outputPathFragments); err != nil {
			return err
		}
	}
	return nil
}

func generateInterfaceFile(
	goPackageFragments []string,
	iface *ast.Interface,
	outPathFragments []string,
) error {
	goPackageFragments = append(goPackageFragments, iface.Name)
	outPathFragments = append(outPathFragments, iface.Name, iface.Name+".go")

	fmt.Println("INTERFACE")
	fmt.Println("PKG", filepath.Join(goPackageFragments...))
	fmt.Println("FILE", filepath.Join(outPathFragments...))
	fmt.Println()
	return nil
}
