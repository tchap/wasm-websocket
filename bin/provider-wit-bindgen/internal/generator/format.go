package generator

import (
	"bytes"
	"fmt"
)

func FormatFile(file File) []byte {
	var b bytes.Buffer
	formatFile(&b, file)
	return b.Bytes()
}

func formatFile(b *bytes.Buffer, file File) {
	// Package
	fmt.Fprintf(b, "package %s\n\n", file.PackageName)

	// Imports
	if len(file.Imports) > 0 {
		fmt.Fprintln(b, "import (")
		for _, imp := range file.Imports {
			if imp.Alias != "" {
				fmt.Fprintf(b, "\t%s \"%s\"\n", imp.Alias, imp.Path)
			} else {
				fmt.Fprintf(b, "\t\"%s\"\n", imp.Path)
			}
		}
		fmt.Fprintln(b, ")")
	}

	// Enums
	for _, enum := range file.Enums {
		if len(enum.Cases) == 0 {
			continue
		}

		fmt.Fprintf(b, "\ntype %s_%s int\n\n", enum.InterfaceName, enum.Name)
		fmt.Fprintln(b, "const (")
		fmt.Fprintf(b, "\t%s%s_%s %s_%s = iota + 1\n", enum.InterfaceName, enum.Name, enum.Cases[0], enum.InterfaceName, enum.Name)
		for _, c := range enum.Cases[1:] {
			fmt.Fprintf(b, "\t%s%s_%s\n", enum.InterfaceName, enum.Name, c)
		}
		fmt.Fprintln(b, ")")
	}

	// Records
	for _, record := range file.Records {
		fmt.Fprintf(b, "\ntype %s_%s struct {\n", record.InterfaceName, record.Name)
		for _, field := range record.Fields {
			fmt.Fprintf(b, "\t%s %s\n", field.Name, field.Type)
		}
		fmt.Fprintln(b, "}")
	}

	// Structs
	for _, node := range file.Structs {
		fmt.Fprintf(b, "\ntype %s struct {\n", node.Name)
		fmt.Fprintf(b, "\tp *provider.WasmcloudProvider\n")
		fmt.Fprintln(b, "}")

		fmt.Fprintf(b, "\nfunc New%s(p *provider.WasmcloudProvider) *%s {\n", node.Name, node.Name)
		fmt.Fprintf(b, "\treturn &%s{p: p}\n", node.Name)
		fmt.Fprintln(b, "}")

		for _, method := range node.Methods {
			fmt.Fprintf(b, "\nfunc (self *%s) %s(", node.Name, method.Name)
			fmt.Fprintf(b, ") %s {\n", method.ReturnType)
			fmt.Fprintln(b, "}")
		}
	}
}
