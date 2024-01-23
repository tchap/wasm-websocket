package generator

type File struct {
	PackageName string
	Imports     []ImportNode
	Enums       []EnumNode
	Records     []RecordNode
	Structs     []StructNode
}

type ImportNode struct {
	Alias string
	Path  string
}

type EnumNode struct {
	InterfaceName string
	Name          string
	Cases         []string
}

type StructNode struct {
	Name    string
	Methods []MethodNode
}

type RecordNode struct {
	InterfaceName string
	Name          string
	Fields        []ArgNode
}

type ArgNode struct {
	Name string
	Type string
}

type MethodNode struct {
	TypeName    string
	Name        string
	Args        []ArgNode
	ReturnTypes []string
}
