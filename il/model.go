package il

// Model

type Model struct {
	Fragments     []*Fragment
	FragmentsMap  map[string]*Fragment
	Queries       []*Operation
	Mutations     []*Operation
	Subscriptions []*Operation
	Enums         []*EnumType
	Unions        []*UnionType
	Interfaces    []*InterfaceType
	InputTypes    []*InputType
	Schema        *Schema
}

func NewModel() *Model {
	return &Model{
		Fragments:     make([]*Fragment, 0),
		FragmentsMap:  make(map[string]*Fragment),
		Queries:       make([]*Operation, 0),
		Mutations:     make([]*Operation, 0),
		Subscriptions: make([]*Operation, 0),
		Enums:         make([]*EnumType, 0),
		Unions:        make([]*UnionType, 0),
		Interfaces:    make([]*InterfaceType, 0),
		InputTypes:    make([]*InputType, 0),
	}
}

// Fragments

type Fragment struct {
	Name         string
	TypeName     string
	SelectionSet *SelectionSet
	Uses         []*Fragment
	UsedBy       []*Fragment
}

func NewFragment(name string, typeName string) *Fragment {
	return &Fragment{
		Name:         name,
		TypeName:     typeName,
		SelectionSet: nil,
		Uses:         make([]*Fragment, 0),
		UsedBy:       make([]*Fragment, 0),
	}
}

type InlineFragment struct {
	TypeName  string
	Selection *SelectionSet
}

// Selection

type SelectionSet struct {
	Fields          []*SelectionField
	Fragments       []*Fragment
	InlineFragments []*InlineFragment
}

type SelectionField struct {
	Name      string
	Alias     string
	Type      Type
	Selection *SelectionSet
	Arguments []*Argument
}

// Operation
type Operation struct {
	Type         string
	Name         string
	Body         string
	SelectionSet *SelectionSet
	Variables    *Variables
}

// Variables
type Variables struct {
	Variables []*Variable
}

type Variable struct {
	Name         string
	Type         Type
	DefaultValue *Value
}

// Arguments
type Argument struct {
	Name  string
	Value Value
}

// Enums
type EnumType struct {
	Name   string
	Values []string
}

// Unions
type UnionType struct {
	Name   string
	Values []string
}

// Interfaces
type InterfaceType struct {
	Name   string
	Values []string
}

// Input Types
type InputType struct {
	Name   string
	Fields []InputTypeField
}

type InputTypeField struct {
	Name string
	Type Type
}