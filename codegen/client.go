package codegen

import (
	"github.com/openland/spacex-cli/il"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//
// Selector Generator
//

func generateOutputType(tp il.Type, notNull bool, output *Output, model *il.Model) {
	var name string
	if tp.GetKind() == "NotNull" {
		inner := tp.(il.NotNull)
		generateOutputType(inner.Inner, true, output, model)
		return
	}
	if tp.GetKind() == "Scalar" {
		scalar := tp.(il.Scalar)
		if scalar.Name == "String" || scalar.Name == "ID" || scalar.Name == "Date" {
			name = "string"
		} else if scalar.Name == "Int" || scalar.Name == "Float" {
			name = "number"
		} else if scalar.Name == "Boolean" {
			name = "boolean"
		} else {
			panic("unknown scalar: " + scalar.Name)
		}
	} else if tp.GetKind() == "Enum" {
		name = tp.(il.Enum).Name
	} else if tp.GetKind() == "Object" || tp.GetKind() == "Union" || tp.GetKind() == "Interface" {
		if notNull {
			output.EndLine("(")
		} else {
			output.EndLine("Maybe<(")
		}
		output.IndentAdd()
		var tpname string
		var sel *il.SelectionSet
		if tp.GetKind() == "Object" {
			tpname = tp.(il.Object).Name
			sel = tp.(il.Object).SelectionSet
		} else if tp.GetKind() == "Union" {
			tpname = tp.(il.Union).Name
			sel = tp.(il.Union).SelectionSet
		} else if tp.GetKind() == "Interface" {
			tpname = tp.(il.Interface).Name
			sel = tp.(il.Interface).SelectionSet
		} else {
			panic("impossible")
		}
		generateSelectorTyped(tpname, sel, output, model)
		output.IndentRemove()
		output.EndLine("")
		if notNull {
			output.WriteLine(")")
		} else {
			output.WriteLine(")>")
		}
		return
	} else if tp.GetKind() == "List" {
		if notNull {
			output.Append("(")
		} else {
			output.Append("Maybe<(")
		}
		generateOutputType((tp.(il.List)).Inner, false, output, model)
		if notNull {
			output.EndLine(")[]")
		} else {
			output.EndLine(")[]>")
		}
		return
	} else {
		panic("Unsupported output type: " + tp.GetKind())
	}

	if notNull {
		output.Append(name)
	} else {
		output.Append("Maybe<" + name + ">")
	}
}

func generateClientSelector(set *il.SelectionSet, output *Output, model *il.Model, parentType *string) {
	processed := make(map[string]string)
	processed["__typename"] = "__typename"
	for _, sf := range set.Fields {
		if _, ok := processed[sf.Alias]; !ok {
			output.BeginLine("& { " + sf.Alias + ": ")
			generateOutputType(sf.Type, false, output, model)
			output.EndLine("}")
			processed[sf.Alias] = sf.Alias
		}
	}
	for _, sf := range set.Fragments {
		output.WriteLine("& " + sf.Name)
	}
	for _, sf := range set.InlineFragments {
		if isPolymorphic(sf.TypeName, model) {
			panic("Inline polymorphic fragments are not supported")
		}

		s := make([]string, 0)
		for _, r := range findSubtypes(*parentType, model) {
			if r != sf.TypeName {
				s = append(s, "'"+r+"'")
			}
		}
		if len(s) == 0 {
			s = append(s, "never")
		}

		output.WriteLine("& Inline<" + strings.Join(s, " | ") + ",(")
		output.IndentAdd()
		generateSelectorTyped(sf.TypeName, sf.Selection, output, model)
		output.IndentRemove()
		output.WriteLine(")>")
	}
}

func generateSelectorTyped(typename string, set *il.SelectionSet, output *Output, model *il.Model) {
	if isPolymorphic(typename, model) {
		subtypes := findSubtypes(typename, model)
		subtypes2 := make([]string, 0)
		for _, s := range subtypes {
			subtypes2 = append(subtypes2, "'"+s+"'")
		}
		// subtypes2 = append(subtypes2, "string")
		output.WriteLine("& { __typename: " + strings.Join(subtypes2, " | ") + " }")
	} else {
		output.WriteLine("& { __typename: '" + typename + "' }")
	}
	generateClientSelector(set, output, model, &typename)
}

//
// Input Values
//

func clientInputType(tp il.Type, notNull bool) string {
	if tp.GetKind() == "NotNull" {
		return clientInputType(tp.(il.NotNull).Inner, true)
	}
	var name string
	if tp.GetKind() == "Scalar" {
		scalar := tp.(il.Scalar).Name
		if scalar == "String" || scalar == "ID" || scalar == "Date" {
			name = "string"
		} else if scalar == "Int" || scalar == "Float" {
			name = "number"
		} else if scalar == "Boolean" {
			name = "boolean"
		} else {
			panic("unknown scalar: " + scalar)
		}

	} else if tp.GetKind() == "Input" {
		name = tp.(il.Input).Name
	} else if tp.GetKind() == "Enum" {
		name = tp.(il.Enum).Name
	} else if tp.GetKind() == "List" {
		name = "(" + clientInputType(tp.(il.List).Inner, false) + ")[]"
	} else {
		panic("Unknown input type " + tp.GetKind())
	}
	if notNull {
		return name
	} else {
		return "MaybeInput<" + name + ">"
	}
}

//
// Type Generators
//

func generateFragment(fragment *il.Fragment, output *Output, model *il.Model) {
	output.WriteLine("export type " + fragment.Name + " = (")
	output.IndentAdd()
	generateSelectorTyped(fragment.TypeName, fragment.SelectionSet, output, model)
	output.IndentRemove()
	output.WriteLine(");")
}

func generateEnum(enum *il.EnumType, output *Output) {
	output.WriteLine("export enum " + enum.Name + " {")
	output.IndentAdd()
	for _, v := range enum.Values {
		output.WriteLine("" + v + " = '" + v + "',")
	}
	output.IndentRemove()
	output.WriteLine("}")
}

func generateInputType(t *il.InputType, output *Output) {
	output.WriteLine("export interface " + t.Name + " {")
	output.IndentAdd()
	for _, v := range t.Fields {
		if v.Type.GetKind() == "NotNull" {
			output.WriteLine(v.Name + ": " + clientInputType(v.Type, false) + ";")
		} else {
			output.WriteLine(v.Name + "?: " + clientInputType(v.Type, false) + ";")
		}
	}
	output.IndentRemove()
	output.WriteLine("}")
}

func generateVariables(name string, variables *il.Variables, output *Output) {
	output.WriteLine("export interface " + name + "Variables {")
	output.IndentAdd()
	for _, v := range variables.Variables {
		if v.Type.GetKind() == "NotNull" {
			output.WriteLine(v.Name + ": " + clientInputType(v.Type, false) + ";")
		} else {
			output.WriteLine(v.Name + "?: " + clientInputType(v.Type, false) + ";")
		}
	}
	output.IndentRemove()
	output.WriteLine("}")
}

func generateOperation(t *il.Operation, output *Output, model *il.Model) {
	// Generate object fields
	if len(t.SelectionSet.InlineFragments) > 0 {
		panic("Inline fragments on root types are prohibited")
	}
	if len(t.SelectionSet.Fragments) > 0 {
		panic("Fragments on root types are prohibited")
	}

	output.WriteLine("export type " + t.Name + " = (")
	output.IndentAdd()
	generateClientSelector(t.SelectionSet, output, model, nil)
	output.IndentRemove()
	output.WriteLine(");")
}

//
// Entry Point
//

func GenerateClient(model *il.Model, to string) {
	output := NewOutput()
	output.WriteLine("/* tslint:disable */")
	output.WriteLine("/* eslint-disable */")
	output.WriteLine("type Maybe<T> = T | null;")
	output.WriteLine("type MaybeInput<T> = T | null | undefined;")
	output.WriteLine("type Inline<E, V> =  { __typename: E; } | V")
	output.WriteLine("")
	output.WriteLine("// Enums")
	for _, f := range model.Enums {
		generateEnum(f, output)
	}
	output.WriteLine("")
	output.WriteLine("// Input Types")
	for _, f := range model.InputTypes {
		generateInputType(f, output)
	}
	output.WriteLine("")
	output.WriteLine("// Fragments")
	for _, f := range model.Fragments {
		generateFragment(f, output, model)
	}
	output.WriteLine("")
	output.WriteLine("// Queries")
	for _, f := range model.Queries {
		if len(f.Variables.Variables) > 0 {
			generateVariables(f.Name, f.Variables, output)
		}
		generateOperation(f, output, model)
	}
	output.WriteLine("")
	output.WriteLine("// Mutations")
	for _, f := range model.Mutations {
		if len(f.Variables.Variables) > 0 {
			generateVariables(f.Name, f.Variables, output)
		}
		generateOperation(f, output, model)
	}
	output.WriteLine("")
	output.WriteLine("// Subscriptions")
	for _, f := range model.Subscriptions {
		if len(f.Variables.Variables) > 0 {
			generateVariables(f.Name, f.Variables, output)
		}
		generateOperation(f, output, model)
	}

	// Result
	err := os.MkdirAll(filepath.Dir(to), os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(to+"/spacex.types.ts", []byte(output.String()), 0644)
	if err != nil {
		panic(err)
	}
}
