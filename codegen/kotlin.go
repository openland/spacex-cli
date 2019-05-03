package codegen

import (
	"github.com/openland/spacex-cli/il"
	"io/ioutil"
	"strconv"
	"strings"
)

func formatValue(arg il.Value, output *Output) string {
	scope := "scope" + strconv.FormatInt(output.GetScope(), 10)
	if arg.GetKind() == "IntValue" {
		return "\"" + strconv.FormatInt(int64(arg.(il.IntValue).Int), 10) + "\""
	} else if arg.GetKind() == "FloatValue" {
		return "\"" + strconv.FormatFloat(arg.(il.FloatValue).Float, 'E', -1, 64) + "\""
	} else if arg.GetKind() == "StringValue" {
		return "\"" + arg.(il.StringValue).String + "\""
	} else if arg.GetKind() == "BooleanValue" {
		return "\"" + strconv.FormatBool(arg.(il.BooleanValue).Boolean) + "\""
	} else if arg.GetKind() == "VariableValue" {
		return scope + ".argumentKey(\"" + arg.(il.VariableValue).Name + "\")"
	} else if arg.GetKind() == "ListValue" {
		inner := make([]string, 0)
		for _, a := range arg.(il.ListValue).Values {
			inner = append(inner, formatValue(a, output))
		}
		if len(inner) == 0 {
			return "\"[]\""
		}
		return "\"[\"+" + strings.Join(inner, "+\",\"+") + "+\"]\""
	} else if arg.GetKind() == "EnumValue" {
		return "\"" + arg.(il.EnumValue).String + "\""
	} else if arg.GetKind() == "ObjectValue" {
		inner := make([]string, 0)
		for _, f := range arg.(il.ObjectValue).Fields {
			inner = append(inner, "\""+f.Name+"\" to "+formatValue(f.Value, output))
		}
		return scope + ".formatObjectKey(" + strings.Join(inner, ",") + ")"
	} else {
		panic("Unsupported variable value: " + arg.GetKind())
	}
}

func storeKey(field *il.SelectionField, output *Output) string {
	scope := "scope" + strconv.FormatInt(output.GetScope(), 10)
	if len(field.Arguments) > 0 {
		argsString := make([]string, 0)
		for _, a := range field.Arguments {
			argsString = append(argsString, "\""+a.Name+"\" to "+formatValue(a.Value, output))
		}
		args := strings.Join(argsString, ",")
		return "\"" + field.Name + "\" + " + scope + ".formatArguments(" + args + ")"
	} else {
		return "\"" + field.Name + "\""
	}
}

func generateReadScalar(field *il.SelectionField, tp string, output *Output) {
	scope := "scope" + strconv.FormatInt(output.GetScope(), 10)
	storeKey := storeKey(field, output)
	requestKey := field.Alias
	output.WriteLine(scope + ".set(" + storeKey + ", " + scope + ".read" + tp + "(\"" + requestKey + "\"))")
}

func generateReadOptionalScalar(field *il.SelectionField, tp string, output *Output) {
	scope := "scope" + strconv.FormatInt(output.GetScope(), 10)
	storeKey := storeKey(field, output)
	requestKey := field.Alias
	output.WriteLine(scope + ".set(" + storeKey + ", " + scope + ".read" + tp + "Optional(\"" + requestKey + "\"))")
}

func generateReadOptionalListScalar(tp string, output *Output) {
	scope := "scope" + strconv.FormatInt(output.GetScope(), 10)
	output.WriteLine(scope + ".next(" + scope + ".read" + tp + "Optional(i))")
}

func generateReadListScalar(tp string, output *Output) {
	scope := "scope" + strconv.FormatInt(output.GetScope(), 10)
	output.WriteLine(scope + ".next(" + scope + ".read" + tp + "(i))")
}

func newScope(field *il.SelectionField, output *Output) {
	storeKey := storeKey(field, output)
	requestKey := field.Alias
	output.NextScope()
	output.WriteLine("val scope" + strconv.FormatInt(output.GetScope(), 10) + " = scope" + strconv.FormatInt(output.ParentScope(), 10) +
		".child(\"" + requestKey + "\", " + storeKey + ")")
}

func newScopeInList(output *Output) {
	output.NextScope()
	output.WriteLine("val scope" + strconv.FormatInt(output.GetScope(), 10) + " = scope" + strconv.FormatInt(output.ParentScope(), 10) +
		".child(i)")
}

func newListScope(field *il.SelectionField, output *Output) {
	storeKey := storeKey(field, output)
	requestKey := field.Alias
	output.NextScope()
	output.WriteLine("val scope" + strconv.FormatInt(output.GetScope(), 10) + " = scope" + strconv.FormatInt(output.ParentScope(), 10) +
		".childList(\"" + requestKey + "\", " + storeKey + ")")
}

func newListScopeInList(output *Output) {
	output.NextScope()
	output.WriteLine("val scope" + strconv.FormatInt(output.GetScope(), 10) + " = scope" + strconv.FormatInt(output.ParentScope(), 10) +
		".childList(i)")
}

var nextLevel int64

func generateListNormalizer(level int64, fld *il.SelectionField, list il.List, output *Output) {
	scope := "scope" + strconv.FormatInt(output.GetScope(), 10)
	output.WriteLine("for (i in 0 until " + scope + ".size) {")
	output.IndentAdd()
	if list.Inner.GetKind() == "NotNull" {
		inner := list.Inner.(il.NotNull).Inner
		if inner.GetKind() == "Scalar" {
			scalar := inner.(il.Scalar)
			if scalar.Name == "String" {
				generateReadListScalar("String", output)
			} else if scalar.Name == "Int" {
				generateReadListScalar("Int", output)
			} else if scalar.Name == "Float" {
				generateReadListScalar("Float", output)
			} else if scalar.Name == "ID" {
				generateReadListScalar("String", output)
			} else if scalar.Name == "Boolean" {
				generateReadListScalar("Boolean", output)
			} else if scalar.Name == "Date" {
				generateReadListScalar("String", output)
			} else {
				panic("Unsupported scalar: " + scalar.Name)
			}
		} else if inner.GetKind() == "Object" || inner.GetKind() == "Union" || inner.GetKind() == "Interface" {
			output.WriteLine(scope + ".assertObject(i)")
			newScopeInList(output)
			if inner.GetKind() == "Object" {
				obj := inner.(il.Object)
				generateNormalizer(obj.SelectionSet, output)
			} else if inner.GetKind() == "Interface" {
				obj := inner.(il.Interface)
				generateNormalizer(obj.SelectionSet, output)
			} else {
				obj := inner.(il.Union)
				generateNormalizer(obj.SelectionSet, output)
			}
			output.ScopePop()
		} else if inner.GetKind() == "Enum" {
			generateReadListScalar("String", output)
		} else {
			panic("Unsupported list inner type " + inner.GetKind())
		}
	} else {
		if list.Inner.GetKind() == "Scalar" {
			scalar := list.Inner.(il.Scalar)
			if scalar.Name == "String" {
				generateReadOptionalListScalar("String", output)
			} else if scalar.Name == "Int" {
				generateReadOptionalListScalar("Int", output)
			} else if scalar.Name == "Float" {
				generateReadOptionalListScalar("Float", output)
			} else if scalar.Name == "ID" {
				generateReadOptionalListScalar("String", output)
			} else if scalar.Name == "Boolean" {
				generateReadOptionalListScalar("Boolean", output)
			} else if scalar.Name == "Date" {
				generateReadOptionalListScalar("String", output)
			} else {
				panic("Unsupported scalar: " + scalar.Name)
			}
		} else if list.Inner.GetKind() == "Object" || list.Inner.GetKind() == "Union" || list.Inner.GetKind() == "Interface" {

		} else if list.Inner.GetKind() == "Enum" {
			generateReadOptionalListScalar("String", output)
		} else if list.Inner.GetKind() == "List" {
			output.WriteLine("if (" + scope + ".isNotNull(i)) {")
			output.IndentAdd()
			newListScopeInList(output)
			generateListNormalizer(nextLevel, fld, list.Inner.(il.List), output)
			output.ScopePop()
			output.IndentRemove()
			output.WriteLine("}")
		} else {
			panic("Unsupported list inner type " + list.Inner.GetKind())
		}
	}
	output.IndentRemove()
	output.WriteLine("}")
}

func generateNormalizer(set *il.SelectionSet, output *Output) {
	scope := "scope" + strconv.FormatInt(output.GetScope(), 10)
	for _, fld := range set.Fields {
		if fld.Type.GetKind() == "NotNull" {
			inner := fld.Type.(il.NotNull).Inner
			if inner.GetKind() == "Scalar" {
				scalar := inner.(il.Scalar)
				if scalar.Name == "String" {
					generateReadScalar(fld, "String", output)
				} else if scalar.Name == "Int" {
					generateReadScalar(fld, "Int", output)
				} else if scalar.Name == "Float" {
					generateReadScalar(fld, "Float", output)
				} else if scalar.Name == "ID" {
					generateReadScalar(fld, "String", output)
				} else if scalar.Name == "Boolean" {
					generateReadScalar(fld, "Boolean", output)
				} else if scalar.Name == "Date" {
					generateReadScalar(fld, "String", output)
				} else {
					panic("Unsupported scalar: " + scalar.Name)
				}
			} else if inner.GetKind() == "Object" || inner.GetKind() == "Union" || inner.GetKind() == "Interface" {
				output.WriteLine(scope + ".assertObject(\"" + fld.Alias + "\")")
				newScope(fld, output)
				if inner.GetKind() == "Object" {
					obj := inner.(il.Object)
					generateNormalizer(obj.SelectionSet, output)
				} else if inner.GetKind() == "Interface" {
					obj := inner.(il.Interface)
					generateNormalizer(obj.SelectionSet, output)
				} else {
					obj := inner.(il.Union)
					generateNormalizer(obj.SelectionSet, output)
				}
				output.ScopePop()
			} else if inner.GetKind() == "List" {
				output.WriteLine("if (" + scope + ".assertList(\"" + fld.Alias + "\")) {")
				output.IndentAdd()
				storeKey := storeKey(fld, output)
				newListScope(fld, output)
				nextLevel++
				generateListNormalizer(nextLevel, fld, inner.(il.List), output)
				listScope := "scope" + strconv.FormatInt(output.GetScope(), 10)
				output.WriteLine(scope + ".set(" + storeKey + ", " + listScope + ".completed())")
				output.ScopePop()
				output.IndentRemove()
				output.WriteLine("}")
			} else if inner.GetKind() == "Enum" {
				generateReadScalar(fld, "String", output)
			} else {
				panic("Unsupported type: " + inner.GetKind())
			}
			//
		} else {
			if fld.Type.GetKind() == "Scalar" {
				scalar := fld.Type.(il.Scalar)
				if scalar.Name == "String" {
					generateReadOptionalScalar(fld, "String", output)
				} else if scalar.Name == "Int" {
					generateReadOptionalScalar(fld, "Int", output)
				} else if scalar.Name == "Float" {
					generateReadOptionalScalar(fld, "Float", output)
				} else if scalar.Name == "ID" {
					generateReadOptionalScalar(fld, "String", output)
				} else if scalar.Name == "Boolean" {
					generateReadOptionalScalar(fld, "Boolean", output)
				} else if scalar.Name == "Date" {
					generateReadOptionalScalar(fld, "String", output)
				} else {
					panic("Unsupported scalar: " + scalar.Name)
				}
			} else if fld.Type.GetKind() == "Object" || fld.Type.GetKind() == "Union" || fld.Type.GetKind() == "Interface" {
				output.WriteLine("if (" + scope + ".hasKey(\"" + fld.Alias + "\")) {")
				output.IndentAdd()
				storeKey := storeKey(fld, output)
				newScope(fld, output)
				if fld.Type.GetKind() == "Object" {
					obj := fld.Type.(il.Object)
					generateNormalizer(obj.SelectionSet, output)
				} else if fld.Type.GetKind() == "Interface" {
					obj := fld.Type.(il.Interface)
					generateNormalizer(obj.SelectionSet, output)
				} else {
					obj := fld.Type.(il.Union)
					generateNormalizer(obj.SelectionSet, output)
				}
				output.IndentRemove()
				output.ScopePop()
				output.WriteLine("} else {")
				output.IndentAdd()
				output.WriteLine(scope + ".setNull(" + storeKey + ")")
				output.IndentRemove()
				output.WriteLine("}")
			} else if fld.Type.GetKind() == "List" {
				output.WriteLine("if (" + scope + ".hasKey(\"" + fld.Alias + "\")) {")
				output.IndentAdd()
				storeKey := storeKey(fld, output)
				newListScope(fld, output)
				nextLevel++
				generateListNormalizer(nextLevel, fld, fld.Type.(il.List), output)
				output.ScopePop()
				output.IndentRemove()
				output.WriteLine("} else {")
				output.IndentAdd()
				output.WriteLine(scope + ".setNull(" + storeKey + ")")
				output.IndentRemove()
				output.WriteLine("}")
			} else if fld.Type.GetKind() == "Enum" {
				generateReadOptionalScalar(fld, "String", output)
			} else {
				panic("Unsupported type: " + fld.Type.GetKind())
			}
		}
	}
	for _, inf := range set.InlineFragments {
		output.WriteLine("if (" + scope + ".isType(\"" + inf.TypeName + "\")) {")
		output.IndentAdd()
		generateNormalizer(inf.Selection, output)
		output.IndentRemove()
		output.WriteLine("}")
	}
	for _, fr := range set.Fragments {
		output.WriteLine("normalize" + fr.Name + "(scope" + strconv.FormatInt(output.GetScope(), 10) + ")")
	}
}

func inputValue(value il.Value) string {
	if value.GetKind() == "IntValue" {
		return "i(" + strconv.FormatInt(int64(value.(il.IntValue).Int), 10) + ")"
	} else if value.GetKind() == "StringValue" {
		return "s(\"" + value.(il.StringValue).String + "\")"
	} else if value.GetKind() == "VariableValue" {
		return "reference(\"" + value.(il.VariableValue).Name + "\")"
	}

	// TODO: Implement all types
	panic("Unexpected input type: " + value.GetKind())
}

func outputSelectors(set *il.SelectionSet, output *Output) {
	output.Append("obj(listOf(")
	output.IndentAdd()
	isFirst := true
	for _, s := range set.Fields {
		if isFirst {
			isFirst = false
		} else {
			output.Append(",")
		}
		if len(s.Arguments) > 0 {
			args := make([]string, 0)
			for _, a := range s.Arguments {
				args = append(args, "\""+a.Name+"\" to "+inputValue(a.Value)+"")
			}
			output.WriteLine("field(\"" + s.Name + "\",\"" + s.Alias + "\", mapOf(" + strings.Join(args, ", ") + "), ")
		} else {
			output.WriteLine("field(\"" + s.Name + "\",\"" + s.Alias + "\", ")
		}
		outputType(s.Type, output)
		output.Append(")")
	}
	for _, fr := range set.Fragments {
		if isFirst {
			isFirst = false
		} else {
			output.Append(",")
		}
		output.WriteLine("fragment(\"" + fr.TypeName + "\", " + fr.Name + "Selector)")
	}
	for _, fr := range set.InlineFragments {
		if isFirst {
			isFirst = false
		} else {
			output.Append(",")
		}
		output.WriteLine("inline(\"" + fr.TypeName + "\", ")
		outputSelectors(fr.Selection, output)
		output.Append(")")
	}
	output.IndentRemove()
	output.WriteLine("))")
}

func outputType(tp il.Type, output *Output) {
	if tp.GetKind() == "NotNull" {
		inner := tp.(il.NotNull)
		output.Append("notNull(")
		outputType(inner.Inner, output)
		output.Append(")")
	} else if tp.GetKind() == "Scalar" {
		scalar := tp.(il.Scalar)
		output.Append("scalar(\"" + scalar.Name + "\")")
	} else if tp.GetKind() == "Enum" {
		output.Append("scalar(\"String\")")
	} else if tp.GetKind() == "Object" || tp.GetKind() == "Union" || tp.GetKind() == "Interface" {
		var set *il.SelectionSet
		if tp.GetKind() == "Object" {
			set = tp.(il.Object).SelectionSet
		} else if tp.GetKind() == "Union" {
			set = tp.(il.Union).SelectionSet
		} else {
			set = tp.(il.Interface).SelectionSet
		}
		output.IndentAdd()
		outputSelectors(set, output)
		output.IndentRemove()
	} else if tp.GetKind() == "List" {
		output.Append("list(")
		outputType((tp.(il.List)).Inner, output)
		output.Append(")")
	} else {
		panic("Unsupported output type: " + tp.GetKind())
	}
}

func generateSelector(set *il.SelectionSet, output *Output) {
	output.IndentAdd()
	outputSelectors(set, output)
	output.IndentRemove()
}

func GenerateKotlin(model *il.Model, to string) {
	output := NewOutput()
	output.WriteLine("package com.openland.soyuz.gen")
	output.WriteLine("")
	output.WriteLine("import com.openland.soyuz.store.RecordSet")
	output.WriteLine("import kotlinx.serialization.json.JsonObject")
	output.WriteLine("")

	//
	// Normalizers
	//

	for _, f := range model.Fragments {
		output.NextScope()
		output.WriteLine("fun normalize" + f.Name + "(scope" + strconv.FormatInt(output.GetScope(), 10) + ": Scope) {")
		output.IndentAdd()
		generateNormalizer(f.SelectionSet, output)
		output.IndentRemove()
		output.WriteLine("}")
	}

	for _, f := range model.Queries {
		output.NextScope()
		output.WriteLine("fun normalize" + f.Name + "(scope" + strconv.FormatInt(output.GetScope(), 10) + ": Scope) {")
		output.IndentAdd()
		generateNormalizer(f.SelectionSet, output)
		output.IndentRemove()
		output.WriteLine("}")
	}

	for _, f := range model.Mutations {
		output.NextScope()
		output.WriteLine("fun normalize" + f.Name + "(scope" + strconv.FormatInt(output.GetScope(), 10) + ": Scope) {")
		output.IndentAdd()
		generateNormalizer(f.SelectionSet, output)
		output.IndentRemove()
		output.WriteLine("}")
	}

	for _, f := range model.Subscriptions {
		output.NextScope()
		output.WriteLine("fun normalize" + f.Name + "(scope" + strconv.FormatInt(output.GetScope(), 10) + ": Scope) {")
		output.IndentAdd()
		generateNormalizer(f.SelectionSet, output)
		output.IndentRemove()
		output.WriteLine("}")
	}

	//
	// Selectors
	//

	for _, f := range model.Fragments {
		output.NextScope()
		output.WriteLine("val " + f.Name + "Selector = ")
		output.IndentAdd()
		generateSelector(f.SelectionSet, output)
		output.IndentRemove()
		output.WriteLine("")
	}

	for _, f := range model.Queries {
		output.NextScope()
		output.WriteLine("val " + f.Name + "Selector = ")
		output.IndentAdd()
		generateSelector(f.SelectionSet, output)
		output.IndentRemove()
	}

	//
	// Operations
	//

	output.WriteLine("")
	output.WriteLine("object Operations {")
	output.IndentAdd()
	for _, f := range model.Queries {
		output.WriteLine("val " + f.Name + " = object: OperationDefinition {")
		output.IndentAdd()
		output.WriteLine("override val name = \"" + f.Name + "\"")
		output.WriteLine("override val kind = OperationKind.QUERY")
		output.WriteLine("override val body = \"" + f.Body + "\"")
		output.WriteLine("override val selector = " + f.Name + "Selector")
		output.WriteLine("override fun normalizeResponse(response: JsonObject): RecordSet {")
		output.IndentAdd()
		output.WriteLine("val collection = NormalizedCollection()")
		output.WriteLine("normalize" + f.Name + "(Scope(\"ROOT_QUERY\", collection, response))")
		output.WriteLine("return collection.build()")
		output.IndentRemove()
		output.WriteLine("}")
		output.IndentRemove()
		output.WriteLine("}")
	}
	for _, f := range model.Mutations {
		output.WriteLine("val " + f.Name + " = object: OperationDefinition {")
		output.IndentAdd()
		output.WriteLine("override val name = \"" + f.Name + "\"")
		output.WriteLine("override val kind = OperationKind.MUTATION")
		output.WriteLine("override val body = \"" + f.Body + "\"")
		output.WriteLine("override val selector = null")
		output.WriteLine("override fun normalizeResponse(response: JsonObject): RecordSet {")
		output.IndentAdd()
		output.WriteLine("val collection = NormalizedCollection()")
		output.WriteLine("normalize" + f.Name + "(Scope(\"ROOT_QUERY\", collection, response))")
		output.WriteLine("return collection.build()")
		output.IndentRemove()
		output.WriteLine("}")
		output.IndentRemove()
		output.WriteLine("}")
	}
	for _, f := range model.Subscriptions {
		output.WriteLine("val " + f.Name + " = object: OperationDefinition {")
		output.IndentAdd()
		output.WriteLine("override val name = \"" + f.Name + "\"")
		output.WriteLine("override val kind = OperationKind.SUBSCRIPTION")
		output.WriteLine("override val body = \"" + f.Body + "\"")
		output.WriteLine("override val selector = null")
		output.WriteLine("override fun normalizeResponse(response: JsonObject): RecordSet {")
		output.IndentAdd()
		output.WriteLine("val collection = NormalizedCollection()")
		output.WriteLine("normalize" + f.Name + "(Scope(\"ROOT_QUERY\", collection, response))")
		output.WriteLine("return collection.build()")
		output.IndentRemove()
		output.WriteLine("}")
		output.IndentRemove()
		output.WriteLine("}")
	}
	output.IndentRemove()
	output.WriteLine("}")

	// Write result
	err := ioutil.WriteFile(to, []byte(output.String()), 0644)
	if err != nil {
		panic(err)
	}
}
