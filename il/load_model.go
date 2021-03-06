package il

import (
	"encoding/json"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/printer"
	"github.com/graphql-go/graphql/language/source"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type ClientModel struct {
	Schema        Schema
	Fragments     map[string]*ast.FragmentDefinition
	Subscriptions map[string]*ast.OperationDefinition
	Queries       map[string]*ast.OperationDefinition
	Mutations     map[string]*ast.OperationDefinition
	Sources       map[string]string
}

type SchemaType struct {
	Kind          string             `json:"kind"`
	Name          string             `json:"name"`
	Fields        []*SchemaField     `json:"fields"`
	OfType        *SchemaType        `json:"ofType"`
	EnumValues    []*SchemaEnumValue `json:"enumValues"`
	PossibleTypes []*SchemaType      `json:"possibleTypes"`
	InputFields   []*InputField      `json:"inputFields"`
}

type InputField struct {
	Name string     `json:"name"`
	Type SchemaType `json:"type"`
}

type SchemaEnumValue struct {
	Name string `json:"name"`
}

type SchemaField struct {
	Name string     `json:"name"`
	Type SchemaType `json:"type"`
}

type Schema struct {
	QueryType        *SchemaRootType `json:"queryType"`
	MutationType     *SchemaRootType `json:"mutationType"`
	SubscriptionType *SchemaRootType `json:"subscriptionType"`
	Types            []SchemaType    `json:"types"`
}

type SchemaRoot struct {
	Schema Schema `json:"__schema"`
}

type SchemaRootType struct {
	Name string `json:"name"`
}

// Process All Fragment Dependencies

func prepareSelectionSet(selectionSet *ast.SelectionSet, model *Model, clModel *ClientModel) {
	if selectionSet == nil {
		return
	}
	for s := range selectionSet.Selections {
		ss := selectionSet.Selections[s].(ast.Node)
		if ss.GetKind() == "FragmentSpread" {
			fs := ss.(*ast.FragmentSpread)
			prepareFragment(clModel.Fragments[fs.Name.Value], model, clModel)
		} else if ss.GetKind() == "Field" {
			f := ss.(*ast.Field)
			prepareSelectionSet(f.SelectionSet, model, clModel)
		} else if ss.GetKind() == "InlineFragment" {
			fs := ss.(*ast.InlineFragment)
			prepareSelectionSet(fs.SelectionSet, model, clModel)
		} else {
			panic("Unknown selection: " + ss.GetKind())
		}
	}
}

func prepareFragment(fragment *ast.FragmentDefinition, model *Model, clModel *ClientModel) {
	if _, ok := model.FragmentsMap[fragment.Name.Value]; ok {
		return
	}
	fr := NewFragment(fragment.Name.Value, fragment.TypeCondition.Name.Value)
	model.FragmentsMap[fragment.Name.Value] = fr
	prepareSelectionSet(fragment.SelectionSet, model, clModel)
	fr.SelectionSet = convertSelection(fr.TypeName, fragment.SelectionSet, model, clModel)
	deps := collectDependencies(fragment.SelectionSet, model, clModel)
	for i := range deps {
		fr2 := model.FragmentsMap[deps[i]]
		fr.Uses = append(fr.Uses, fr2)
		fr2.UsedBy = append(fr2.UsedBy, fr)
	}

	model.Fragments = append(model.Fragments, fr)
}

// Operations

func convertValue(value ast.Value) Value {
	if value.GetKind() == "Variable" {
		return VariableValue{(value.(*ast.Variable)).Name.Value}
	} else if value.GetKind() == "IntValue" {
		iv := value.(*ast.IntValue)
		v, e := strconv.ParseInt(iv.Value, 10, 32)
		if e != nil {
			panic(e)
		}
		return IntValue{int32(v)}
	} else if value.GetKind() == "FloatValue" {
		iv := value.(*ast.FloatValue)
		v, e := strconv.ParseFloat(iv.Value, 64)
		if e != nil {
			panic(e)
		}
		return FloatValue{v}
	} else if value.GetKind() == "StringValue" {
		iv := value.(*ast.StringValue)
		return StringValue{iv.Value}
	} else if value.GetKind() == "BooleanValue" {
		iv := value.(*ast.BooleanValue)
		return BooleanValue{iv.Value}
	} else if value.GetKind() == "EnumValue" {
		iv := value.(*ast.EnumValue)
		return EnumValue{iv.Value}
	} else if value.GetKind() == "ListValue" {
		iv := value.(*ast.ListValue)
		lv := make([]Value, 0)
		for i := range iv.Values {
			lv = append(lv, convertValue(iv.Values[i]))
		}
		return ListValue{lv}
	} else if value.GetKind() == "ObjectValue" {
		ov := value.(*ast.ObjectValue)
		fld := make([]ObjectValueField, 0)
		for i := range ov.Fields {
			fld = append(fld, ObjectValueField{ov.Fields[i].Name.Value, convertValue(ov.Fields[i].Value)})
		}
		return ObjectValue{fld}
	} else {
		panic("Unknown value kind: " + value.GetKind())
	}
}

func convertVariables(variables []*ast.VariableDefinition, model *Model, clModel *ClientModel) *Variables {
	res := make([]*Variable, 0)
	for i := range variables {
		tp := convertInputType(variables[i].Type, model, clModel)
		var vl *Value
		if variables[i].DefaultValue != nil {
			vl2 := convertValue(variables[i].DefaultValue)
			vl = &vl2
		}
		vr := &Variable{Name: variables[i].Variable.Name.Value, Type: tp, DefaultValue: vl}
		res = append(res, vr)
	}
	return &Variables{res}
}

func normalizeSelector(selectionSet *ast.SelectionSet, root bool) {
	hasTypename := false
	for s := range selectionSet.Selections {
		ss := selectionSet.Selections[s].(ast.Node)
		if ss.GetKind() == "FragmentSpread" {
			// Nothing to do
		} else if ss.GetKind() == "Field" {
			f := ss.(*ast.Field)
			if f.Name.Value == "__typename" {
				hasTypename = true
			}
			if f.SelectionSet != nil {
				normalizeSelector(f.SelectionSet, false)
			}
		} else if ss.GetKind() == "InlineFragment" {
			f := ss.(*ast.InlineFragment)
			normalizeSelector(f.SelectionSet, false)
		} else {
			panic("Unknown selection: " + ss.GetKind())
		}
	}

	if !hasTypename && !root {
		selectionSet.Selections = append([]ast.Selection{ast.NewField(&ast.Field{
			Name: ast.NewName(&ast.Name{
				Value: "__typename",
			}),
		})}, selectionSet.Selections...)
	}
}

func printQuery(src ast.Node) string {
	return printer.Print(src).(string)
}

func prepareOperation(definition *ast.OperationDefinition, model *Model, clModel *ClientModel) {
	root := clModel.Schema.QueryType.Name
	if definition.Operation == "mutation" {
		root = clModel.Schema.MutationType.Name
	} else if definition.Operation == "subscription" {
		root = clModel.Schema.SubscriptionType.Name
	}
	variables := &Variables{}
	if definition.VariableDefinitions != nil {
		variables = convertVariables(definition.VariableDefinitions, model, clModel)
	}
	selection := convertSelection(root, definition.SelectionSet, model, clModel)

	dependencies := collectDependencies(definition.SelectionSet, model, clModel)

	body := clModel.Sources[definition.Name.Value]
	for _, d := range dependencies {
		body = body + " " + clModel.Sources[d]
	}
	// body = normalizeQuery(body)
	body = strings.Replace(body, "\n", " ", -1)
	for strings.Contains(body, "  ") {
		body = strings.Replace(body, "  ", " ", -1)
	}
	body = strings.Replace(body, ": ", ":", -1)
	body = strings.Replace(body, ", ", ",", -1)
	body = strings.Replace(body, "( ", "(", -1)
	body = strings.Replace(body, " (", "(", -1)
	body = strings.Replace(body, ") ", ")", -1)
	body = strings.Replace(body, " ) ", ")", -1)
	body = strings.Replace(body, "{ ", "{", -1)
	body = strings.Replace(body, " {", "{", -1)
	body = strings.Replace(body, "} ", "}", -1)
	body = strings.Replace(body, " }", "}", -1)

	if definition.Operation == "mutation" {
		model.Mutations = append(model.Mutations, &Operation{
			Type:         definition.Operation,
			Name:         definition.Name.Value,
			Body:         body,
			SelectionSet: selection,
			Variables:    variables,
		})
	} else if definition.Operation == "subscription" {
		model.Subscriptions = append(model.Subscriptions, &Operation{
			Type:         definition.Operation,
			Name:         definition.Name.Value,
			Body:         body,
			SelectionSet: selection,
			Variables:    variables,
		})
	} else {
		model.Queries = append(model.Queries, &Operation{
			Type:         definition.Operation,
			Name:         definition.Name.Value,
			Body:         body,
			SelectionSet: selection,
			Variables:    variables,
		})
	}

}

// Collect Dependencies

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func collectDependencies(selectionSet *ast.SelectionSet, model *Model, clModel *ClientModel) []string {
	res := make([]string, 0)
	if selectionSet == nil {
		return res
	}
	for s := range selectionSet.Selections {
		ss := selectionSet.Selections[s].(ast.Node)
		if ss.GetKind() == "FragmentSpread" {
			fs := ss.(*ast.FragmentSpread)
			r := collectDependencies(clModel.Fragments[fs.Name.Value].SelectionSet, model, clModel)
			if !contains(res, fs.Name.Value) {
				res = append(res, fs.Name.Value)
			}
			for i := range r {
				if !contains(res, r[i]) {
					res = append(res, r[i])
				}
			}
		} else if ss.GetKind() == "Field" {
			f := ss.(*ast.Field)
			r := collectDependencies(f.SelectionSet, model, clModel)
			for i := range r {
				if !contains(res, r[i]) {
					res = append(res, r[i])
				}
			}
		} else if ss.GetKind() == "InlineFragment" {
			fs := ss.(*ast.InlineFragment)
			r := collectDependencies(fs.SelectionSet, model, clModel)
			for i := range r {
				if !contains(res, r[i]) {
					res = append(res, r[i])
				}
			}
		} else {
			panic("Unknown selection: " + ss.GetKind())
		}
	}
	return res
}

// Convert Selections

func Find(a []SchemaType, name string) *SchemaType {
	for i, n := range a {
		if name == n.Name {
			return &a[i]
		}
	}
	panic("Unknown type: " + name)
}

func FindField(a []*SchemaField, name string) *SchemaField {
	for i, n := range a {
		if name == n.Name {
			return a[i]
		}
	}
	return nil
}

func convertType(field *ast.Field, ff SchemaType, model *Model, clModel *ClientModel) Type {
	if ff.Kind == "NON_NULL" {
		// Not null
		return NotNull{convertType(field, *ff.OfType, model, clModel)}
	} else if ff.Kind == "SCALAR" {
		// Scalar
		return Scalar{ff.Name}
	} else if ff.Kind == "OBJECT" {
		// Object
		return Object{ff.Name, convertSelection(ff.Name, field.SelectionSet, model, clModel)}
	} else if ff.Kind == "UNION" {
		// Union
		return Union{ff.Name, convertSelection(ff.Name, field.SelectionSet, model, clModel)}
	} else if ff.Kind == "LIST" {
		// List
		return List{convertType(field, *ff.OfType, model, clModel)}
	} else if ff.Kind == "INTERFACE" {
		// Interface
		return Interface{ff.Name, convertSelection(ff.Name, field.SelectionSet, model, clModel)}
	} else if ff.Kind == "ENUM" {
		// Enum
		return Enum{ff.Name}
	} else {
		panic("Unexpected type kind: " + ff.Kind)
	}
}

func convertInputType2(ff SchemaType, model *Model, clModel *ClientModel) Type {
	if ff.Kind == "NON_NULL" {
		return NotNull{convertInputType2(*ff.OfType, model, clModel)}
	} else if ff.Kind == "SCALAR" {
		return Scalar{ff.Name}
	} else if ff.Kind == "OBJECT" {
		panic("Can't convert object for input")
	} else if ff.Kind == "UNION" {
		panic("Can't convert object for input")
	} else if ff.Kind == "LIST" {
		// List
		return List{convertInputType2(*ff.OfType, model, clModel)}
	} else if ff.Kind == "INTERFACE" {
		panic("Can't convert object for input")
	} else if ff.Kind == "ENUM" {
		return Enum{ff.Name}
	} else if ff.Kind == "INPUT_OBJECT" {
		return Input{ff.Name}
	} else {
		panic("Unexpected type kind: " + ff.Kind)
	}
}

func convertInputType(ff ast.Type, model *Model, clModel *ClientModel) Type {
	if ff.GetKind() == "NonNull" {
		nn := ff.(*ast.NonNull)
		return NotNull{convertInputType(nn.Type, model, clModel)}
	} else if ff.GetKind() == "Named" {
		nn := ff.(*ast.Named)
		tp := Find(clModel.Schema.Types, nn.Name.Value)
		if tp.Kind == "SCALAR" {
			return Scalar{tp.Name}
		} else if tp.Kind == "ENUM" {
			return Enum{tp.Name}
		} else if tp.Kind == "INPUT_OBJECT" {
			return Input{tp.Name}
		} else {
			panic("Unexpected Named type: " + tp.Kind)
		}
	} else if ff.GetKind() == "List" {
		ln := ff.(*ast.List)
		return List{convertInputType(ln.Type, model, clModel)}
	} else {
		panic("Unexpected input type kind: " + ff.GetKind())
	}
}

func isIdenityType(name string, clModel *ClientModel) bool {
	f := Find(clModel.Schema.Types, name)
	for _, ff := range f.Fields {
		if ff.Name == "id" {
			return true
		}
	}
	return false
}

func hasField(name string, selection *ast.SelectionSet, clModel *ClientModel) bool {
	if selection == nil {
		return false
	}
	for s := range selection.Selections {
		ss := selection.Selections[s].(ast.Node)
		if ss.GetKind() == "FragmentSpread" {
			fs := ss.(*ast.FragmentSpread)
			if hasField(name, clModel.Fragments[fs.Name.Value].SelectionSet, clModel) {
				return true
			}
		} else if ss.GetKind() == "Field" {
			f := ss.(*ast.Field)
			if f.Name.Value == name {
				return true
			}
		} else if ss.GetKind() == "InlineFragment" {
			fs := ss.(*ast.InlineFragment)
			if hasField(name, fs.SelectionSet, clModel) {
				return true
			}
		} else {
			panic("Unknown selection: " + ss.GetKind())
		}
	}
	return false
}

func convertSelection(typeName string, selection *ast.SelectionSet, model *Model, clModel *ClientModel) *SelectionSet {
	if selection == nil {
		return nil
	}
	tp := Find(clModel.Schema.Types, typeName)
	fields := make([]*SelectionField, 0)
	fragments := make([]*Fragment, 0)
	inlineFragments := make([]*InlineFragment, 0)
	if isIdenityType(typeName, clModel) && !hasField("id", selection, clModel) {
		panic("Type " + typeName + " is missing id field in " + string(selection.Loc.Source.Body))
	}
	for s := range selection.Selections {
		ss := selection.Selections[s].(ast.Node)
		if ss.GetKind() == "FragmentSpread" {
			fs := ss.(*ast.FragmentSpread)
			fragments = append(fragments, model.FragmentsMap[fs.Name.Value])
		} else if ss.GetKind() == "Field" {
			f := ss.(*ast.Field)
			alias := f.Name.Value
			if f.Alias != nil {
				alias = f.Alias.Value
			}
			var arguments []*Argument
			var typ Type
			if f.Name.Value == "__typename" {
				// Special Case
				typ = NotNull{Scalar{"String"}}
			} else {
				ff := FindField(tp.Fields, f.Name.Value)
				typ = convertType(f, ff.Type, model, clModel)
				if f.Arguments != nil {
					arguments = make([]*Argument, 0)
					for i := range f.Arguments {
						arg := f.Arguments[i]
						vl := convertValue(arg.Value)
						arguments = append(arguments, &Argument{Name: arg.Name.Value, Value: vl})
					}
				}
			}
			fld := &SelectionField{Name: f.Name.Value, Alias: alias, Type: typ, Arguments: arguments}
			fields = append(fields, fld)
		} else if ss.GetKind() == "InlineFragment" {
			fs := ss.(*ast.InlineFragment)
			fr := convertSelection(fs.TypeCondition.Name.Value, fs.SelectionSet, model, clModel)
			tp := Find(clModel.Schema.Types, fs.TypeCondition.Name.Value)
			if tp.Kind != "OBJECT" {
				panic("Inline fragments are possible only on explicit types")
			}
			ifr := &InlineFragment{
				TypeName:  fs.TypeCondition.Name.Value,
				Selection: fr,
			}
			inlineFragments = append(inlineFragments, ifr)
		} else {
			panic("Unknown selection: " + ss.GetKind())
		}
	}
	return &SelectionSet{Fields: fields, Fragments: fragments, InlineFragments: inlineFragments}
}

// Load Model from files

func LoadModel(schemaPath string, files []string) *Model {
	sort.Strings(files)

	// Read Schema and Queries
	schemaBody, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		panic(err)
	}
	var schemaRoot SchemaRoot
	err = json.Unmarshal(schemaBody, &schemaRoot)
	if err != nil {
		panic(err)
	}
	schema := schemaRoot.Schema
	model := &ClientModel{
		Schema:        schema,
		Fragments:     make(map[string]*ast.FragmentDefinition),
		Queries:       make(map[string]*ast.OperationDefinition),
		Mutations:     make(map[string]*ast.OperationDefinition),
		Subscriptions: make(map[string]*ast.OperationDefinition),
		Sources:       make(map[string]string),
	}
	for i := 0; i < len(files); i++ {
		path := files[i]
		body, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		src := source.NewSource(&source.Source{Body: body})
		parsed, err := parser.Parse(parser.ParseParams{Source: src})
		if err != nil {
			panic(err)
		}
		for i := 0; i < len(parsed.Definitions); i++ {
			node := parsed.Definitions[i]
			if node.GetKind() == "OperationDefinition" {
				op := node.(*ast.OperationDefinition)
				normalizeSelector(op.SelectionSet, true)

				if op.Operation == "query" {
					model.Queries[op.Name.Value] = op
				} else if op.Operation == "mutation" {
					model.Mutations[op.Name.Value] = op
				} else if op.Operation == "subscription" {
					model.Subscriptions[op.Name.Value] = op
				} else {
					panic("Unknown operation: " + op.Operation)
				}
				model.Sources[op.Name.Value] = printQuery(op)
			} else if node.GetKind() == "FragmentDefinition" {
				fr := node.(*ast.FragmentDefinition)
				normalizeSelector(fr.SelectionSet, false)
				model.Fragments[fr.Name.Value] = fr
				model.Sources[fr.Name.Value] = printQuery(fr)
			} else {
				panic("Unknown node: " + node.GetKind())
			}
		}
	}

	// Build IL model
	ilModel := NewModel()
	ilModel.Schema = &schema

	// Fragments
	keys := make([]string, 0)
	for k := range model.Fragments {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := model.Fragments[k]
		prepareFragment(v, ilModel, model)
	}

	// Queries
	keys = make([]string, 0)
	for k := range model.Queries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := model.Queries[k]
		prepareOperation(v, ilModel, model)
	}

	// Mutations
	keys = make([]string, 0)
	for k := range model.Mutations {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := model.Mutations[k]
		prepareOperation(v, ilModel, model)
	}

	// Subscriptions
	keys = make([]string, 0)
	for k := range model.Subscriptions {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := model.Subscriptions[k]
		prepareOperation(v, ilModel, model)
	}

	// Enums
	for _, e := range model.Schema.Types {
		if e.Kind == "ENUM" && !strings.HasPrefix(e.Name, "__") {
			values := make([]string, 0)
			for _, v := range e.EnumValues {
				values = append(values, v.Name)
			}
			ilModel.Enums = append(ilModel.Enums, &EnumType{Name: e.Name, Values: values})
		}
	}

	// Unions
	for _, e := range model.Schema.Types {
		if e.Kind == "UNION" && !strings.HasPrefix(e.Name, "__") {
			values := make([]string, 0)
			for _, p := range e.PossibleTypes {
				if p.Kind != "OBJECT" {
					panic("Unions of unions are not supported yet")
				}
				values = append(values, p.Name)
			}
			ilModel.Unions = append(ilModel.Unions, &UnionType{Name: e.Name, Values: values})
		}
	}

	for _, e := range model.Schema.Types {
		if e.Kind == "INTERFACE" && !strings.HasPrefix(e.Name, "__") {
			values := make([]string, 0)
			for _, p := range e.PossibleTypes {
				if p.Kind != "OBJECT" {
					panic("Interface of interfaces are not supported yet")
				}
				values = append(values, p.Name)
			}
			ilModel.Interfaces = append(ilModel.Interfaces, &InterfaceType{Name: e.Name, Values: values})
		}
	}

	for _, e := range model.Schema.Types {
		if e.Kind == "INPUT_OBJECT" && !strings.HasPrefix(e.Name, "__") {
			fields := make([]InputTypeField, 0)
			for _, v := range e.InputFields {
				fields = append(fields, InputTypeField{Name: v.Name, Type: convertInputType2(v.Type, ilModel, model)})
			}
			ilModel.InputTypes = append(ilModel.InputTypes, &InputType{Name: e.Name, Fields: fields})
		}
	}

	return ilModel
}
