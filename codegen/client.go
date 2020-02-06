package codegen

import (
	"github.com/openland/spacex-cli/il"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//func findUnion(name string, model *il.Model) []string {
//	res := make([]string, 0)
//	for _, u := range model.Unions {
//		if u.Name == name {
//			for _, c := range u.Values {
//				res = append(res, c)
//			}
//		}
//	}
//	return res
//}
//
//func findInterface(name string, model *il.Model) []string {
//	res := make([]string, 0)
//	for _, u := range model.Interfaces {
//		if u.Name == name {
//			for _, c := range u.Values {
//				res = append(res, c)
//			}
//		}
//	}
//	return res
//}

//
//func generateClientInnerObjectTypes(field string, tp il.Type, output *Output, model *il.Model) {
//	if tp.GetKind() == "NotNull" {
//		inner := tp.(il.NotNull)
//		generateClientInnerObjectTypes(field, inner.Inner, output, model)
//	} else if tp.GetKind() == "Scalar" {
//		return
//	} else if tp.GetKind() == "Enum" {
//		return
//	} else if tp.GetKind() == "Object" || tp.GetKind() == "Union" || tp.GetKind() == "Interface" {
//		if tp.GetKind() == "Object" {
//			generateClientObjectFragment(field, tp.(il.Object).Name, tp.(il.Object).SelectionSet, output, model)
//		} else if tp.GetKind() == "Union" {
//			unions := findUnion(tp.(il.Union).Name, model)
//
//			// Concrete types
//			for _, u := range unions {
//				generateClientObjectFragment(field+"_"+u, u, tp.(il.Union).SelectionSet, output, model)
//			}
//
//			// Union Type
//			fields := make([]string, 0)
//			for _, u := range unions {
//				fields = append(fields, field+"_"+u)
//			}
//			output.WriteLine("export type " + field + " = " + strings.Join(fields, " | ") + ";")
//		} else {
//			interfaces := findInterface(tp.(il.Interface).Name, model)
//
//			// Concrete types
//			for _, u := range interfaces {
//				generateClientObjectFragment(field+"_"+u, u, tp.(il.Interface).SelectionSet, output, model)
//			}
//
//			// Union type
//			fields := make([]string, 0)
//			for _, u := range interfaces {
//				fields = append(fields, field+"_"+u)
//			}
//			output.WriteLine("export type " + field + " = " + strings.Join(fields, " | ") + ";")
//		}
//	} else if tp.GetKind() == "List" {
//		generateClientInnerObjectTypes(field, (tp.(il.List)).Inner, output, model)
//	} else {
//		panic("Unsupported output type: " + tp.GetKind())
//	}
//}

//func generateClientInnerObjectTypesSelector(name string, typename string, set *il.SelectionSet, output *Output, model *il.Model) {
//	for _, sf := range set.Fields {
//		generateClientInnerObjectTypes(name+"_"+sf.Alias, sf.Type, output, model)
//	}
//	for _, sf := range set.InlineFragments {
//		if sf.TypeName == typename {
//			generateClientInnerObjectTypesSelector(name, typename, sf.Selection, output, model)
//		}
//	}
//}

//func generateClientObjectInnerFragment(name string, typename string, set *il.SelectionSet, output *Output, model *il.Model) {
//	output.WriteLine("export type " + name + "_" + typename + " = (")
//
//	//sset := &il.SelectionSet{Fields: make([]*il.SelectionField, 0),Fr}
//
//	generateClientObjectFragmentBody(name, typename, set, output, model)
//	//output.WriteLine("__typename: '" + typename + "';")
//	output.WriteLine(");")
//}

//func generateClientObjectFragmentBody(name string, typename string, set *il.SelectionSet, output *Output, model *il.Model) {
//	tp := findType(typename, model)
//	var u []string
//	polymorphic := false
//	if tp.Kind == "UNION" || tp.Kind == "INTERFACE" {
//		polymorphic = true
//		if tp.Kind == "UNION" {
//			u = findUnion(typename, model)
//		} else {
//			u = findInterface(typename, model)
//		}
//	} else {
//		u = make([]string, 1)
//		u[0] = typename
//	}
//	u2 := make([]string, 0)
//	for _, k := range u {
//		u2 = append(u2, "'"+k+"'")
//	}
//
//	output.WriteLine("{")
//	output.IndentAdd()
//	if polymorphic {
//		output.WriteLine("__typename: " + strings.Join(u2, " | ") + " | string;")
//	} else {
//		output.WriteLine("__typename: " + strings.Join(u2, " | ") + ";")
//	}
//	processed := make(map[string]string)
//	processed["__typename"] = "__typename"
//	for _, sf := range set.Fields {
//		if _, ok := processed[sf.Alias]; ok {
//			// Ignore
//		} else {
//			processed[sf.Alias] = sf.Alias
//			output.WriteLine(sf.Alias + ": " + clientOutputType(sf.Type, false) + ";")
//		}
//	}
//	output.IndentRemove()
//	output.WriteLine("}")
//
//	if len(set.Fragments) > 0 || len(set.InlineFragments) > 0 {
//
//		for _, sf := range set.Fragments {
//			output.WriteLine("& " + sf.Name)
//		}
//
//		di := make([]string, 0)
//		processedInlineTypes := make(map[string]string)
//		for _, sf := range set.InlineFragments {
//			if _, ok := processedInlineTypes[sf.TypeName]; !ok {
//				di = append(di, name+"_"+sf.TypeName)
//				processedInlineTypes[sf.TypeName] = sf.TypeName
//			}
//		}
//		if len(di) > 0 {
//			output.WriteLine("& ({} | " + strings.Join(di, " | ") + ")")
//		}
//	}
//}

//func generateClientObjectFragment(name string, typename string, set *il.SelectionSet, output *Output, model *il.Model) {
//	for _, sf := range set.Fields {
//		generateClientInnerObjectTypes(name+"_"+sf.Alias, sf.Type, output, model)
//	}
//	processedInlineTypes := make(map[string]string)
//	for _, sf := range set.InlineFragments {
//		if _, ok := processedInlineTypes[sf.TypeName]; !ok {
//			generateClientObjectInnerFragment(name, sf.TypeName, set, output, model)
//			processedInlineTypes[sf.TypeName] = sf.TypeName
//		}
//	}
//	output.WriteLine("export type " + name + " = (")
//	output.IndentAdd()
//	generateClientObjectFragmentBody(name, typename, set, output, model)
//	output.IndentRemove()
//	output.WriteLine(");")
//}

//func generateClientSelectionSetOp(name string, set *il.SelectionSet, output *Output, model *il.Model) {
//	// Generate object fields
//	for _, sf := range set.Fields {
//		generateClientInnerObjectTypes(name+"_"+sf.Alias, sf.Type, output, model)
//	}
//	output.WriteLine("export interface " + name + " {")
//	output.IndentAdd()
//	for _, sf := range set.Fields {
//		if sf.Name != "__typename" {
//			output.WriteLine(sf.Alias + ": " + clientOutputType(name+"_"+sf.Alias, sf.Type, false) + ";")
//		}
//	}
//	if len(set.InlineFragments) > 0 {
//		panic("Inline fragments on root types are prohibited")
//	}
//	if len(set.Fragments) > 0 {
//		panic("Fragments on root types are prohibited")
//	}
//	output.IndentRemove()
//	output.WriteLine("}")
//}

//func generateClientSelectionSetBody(name string, typename string, set *il.SelectionSet, output *Output, processed map[string]string) {
//	for _, sf := range set.Fields {
//		if sf.Name != "__typename" {
//			if _, ok := processed[sf.Alias]; ok {
//				// Ignore
//			} else {
//				processed[sf.Alias] = sf.Alias
//				output.WriteLine(sf.Alias + ": " + clientOutputType(name+"_"+sf.Alias, sf.Type, false) + ";")
//			}
//		}
//	}
//	for _, sf := range set.InlineFragments {
//		if sf.TypeName == typename {
//			generateClientSelectionSetBody(name, typename, sf.Selection, output, processed)
//		}
//	}
//}
//func contains(s []string, e string) bool {
//	for _, a := range s {
//		if a == e {
//			return true
//		}
//	}
//	return false
//}
//func collectFragments(typename string, set *il.SelectionSet) []string {
//	res := make([]string, 0)
//	for _, sf := range set.InlineFragments {
//		if sf.TypeName == typename {
//			r := collectFragments(typename, sf.Selection)
//			for _, i := range r {
//				if !contains(res, i) {
//					res = append(res, i)
//				}
//			}
//		}
//	}
//	for _, sf := range set.Fragments {
//		if sf.TypeName == typename {
//			if !contains(res, sf.Name) {
//				res = append(res, sf.Name)
//			}
//		}
//	}
//	return res
//}

//func generateVariables(name string, variables *il.Variables, output *Output) {
//	output.WriteLine("export interface " + name + "Variables {")
//	output.IndentAdd()
//	for _, v := range variables.Variables {
//		if v.Type.GetKind() == "NotNull" {
//			output.WriteLine(v.Name + ": " + clientInputType(v.Type, false) + ";")
//		} else {
//			output.WriteLine(v.Name + "?: " + clientInputType(v.Type, false) + ";")
//		}
//	}
//	output.IndentRemove()
//	output.WriteLine("}")
//}

//func generateSelectorSetFields(name string, typename string, set *il.SelectionSet, output *Output, processed map[string]string) {
//	for _, sf := range set.Fields {
//		if _, ok := processed[sf.Alias]; !ok {
//			output.WriteLine(sf.Alias + ": " + clientOutputType(name+"_"+sf.Alias, sf.Type, false) + ";")
//			processed[sf.Alias] = sf.Alias
//		}
//	}
//	for _, sf := range set.InlineFragments {
//		if sf.TypeName == typename {
//			generateSelectorSetFields(name, typename, sf.Selection, output, processed)
//		}
//	}
//}

//func generateSelectorSet(name string, typename string, set *il.SelectionSet, output *Output, model *il.Model) {
//	output.WriteLine("{")
//	output.IndentAdd()
//	output.WriteLine("__typename: '" + typename + "';")
//	processed := make(map[string]string)
//	processed["__typename"] = "__typename"
//	generateSelectorSetFields(name, typename, set, output, processed)
//	output.IndentRemove()
//	output.WriteLine("}")
//	for _, sf := range set.Fragments {
//		output.WriteLine("& " + sf.Name + "_" + typename)
//	}
//}

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
			output.Append("[")
		} else {
			output.Append("Maybe<[")
		}
		generateOutputType((tp.(il.List)).Inner, false, output, model)
		if notNull {
			output.EndLine("]")
		} else {
			output.EndLine("]>")
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

func generateClientSelector(set *il.SelectionSet, output *Output, model *il.Model) {
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
		output.WriteLine("& Inline<(")
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
		subtypes2 = append(subtypes2, "string")
		output.WriteLine("& { __typename: " + strings.Join(subtypes2, " | ") + " }")
	} else {
		output.WriteLine("& { __typename: '" + typename + "' }")
	}
	generateClientSelector(set, output, model)
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
		name = clientInputType(tp.(il.List).Inner, false)
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
	generateClientSelector(t.SelectionSet, output, model)
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
	output.WriteLine("type Inline<V> = {} | V;")
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
