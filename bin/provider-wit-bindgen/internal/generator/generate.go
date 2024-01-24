package generator

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/patrickhuber/go-wasm/wit/ast"
)

func BuildFile(pkg *ast.Ast) File {
	file := File{
		PackageName: pkg.PackageDeclaration.Unwrap().Name,
		Imports: []ImportNode{
			{
				Alias: "provider",
				Path:  "github.com/wasmCloud/provider-sdk-go",
			},
			{
				Path: "github.com/tchap/wasmcloud-websocket/bin/provider-wit-bindgen/internal/rpc/types",
			},
		},
	}
	for _, item := range pkg.Items {
		// Only interfaces are currently supported.
		if item.Interface == nil {
			continue
		}

		node, enums, records := buildStructNode(item.Interface)
		file.Structs = append(file.Structs, node)
		file.Enums = append(file.Enums, enums...)
		file.Records = append(file.Records, records...)
	}
	return file
}

func buildStructNode(iface *ast.Interface) (StructNode, []EnumNode, []RecordNode) {
	ifaceName := convertName(iface.Name)
	node := StructNode{
		Name: ifaceName,
	}
	var (
		enums   []EnumNode
		records []RecordNode
	)
	for _, item := range iface.Items {
		switch item := item.(type) {
		case *ast.Enum:
			enums = append(enums, buildEnumNode(ifaceName, item))

		case *ast.Record:
			records = append(records, buildRecordNode(ifaceName, item))

		case *ast.FuncItem:
			node.Methods = append(node.Methods, buildMethod(ifaceName, item))
		}
	}
	return node, enums, records
}

func buildEnumNode(interfaceName string, enum *ast.Enum) EnumNode {
	convertedID := convertName(enum.ID)
	node := EnumNode{
		InterfaceName: interfaceName,
		Name:          convertedID,
	}
	for _, c := range enum.Cases {
		node.Cases = append(node.Cases, convertName(c.Name))
	}
	return node
}

func buildRecordNode(interfaceName string, record *ast.Record) RecordNode {
	convertedID := convertName(record.ID)
	node := RecordNode{
		InterfaceName: interfaceName,
		Name:          convertedID,
	}
	for _, f := range record.Fields {
		node.Fields = append(node.Fields, ArgNode{
			Name: convertName(f.Name),
			Type: convertType(f.Type),
		})
	}
	return node
}

func buildMethod(interfaceName string, m *ast.FuncItem) MethodNode {
	return MethodNode{
		InterfaceName: interfaceName,
		Name:          convertName(m.ID),
		Args:          buildArgsFromParams(m.FuncType.Params),
		ReturnType:    convertType(m.FuncType.Results.Anonymous),
	}
}

func buildArgsFromParams(params []ast.Parameter) []ArgNode {
	out := make([]ArgNode, len(params))
	for i, param := range params {
		out[i] = ArgNode{
			Name: param.Id,
			Type: convertType(param.Type),
		}
	}
	return out
}

func convertName(name string) string {
	var b strings.Builder
	capitalize := true
	for _, c := range name {
		if c == '-' {
			capitalize = true
			continue
		}

		if capitalize {
			b.WriteRune(unicode.ToUpper(c))
			capitalize = false
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func convertType(t ast.Type) string {
	switch t := t.(type) {
	case *ast.Option:
		ot := convertType(t.Type)
		return fmt.Sprintf("types.Option[%s, *%s]", ot, ot)

	case *ast.List:
		switch it := t.Type.(type) {
		case *ast.U8:
			return fmt.Sprintf("types.Bytes")

		case *ast.Id:
			switch it.Value {
			case "u8":
				return fmt.Sprintf("types.Bytes")
			}

		default:
			return fmt.Sprintf("types.Option[%T]", reflect.TypeOf(it))
		}
	}
	return fmt.Sprintf("%T", t)
}
