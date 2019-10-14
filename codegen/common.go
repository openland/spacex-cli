package codegen

import (
	"github.com/openland/spacex-cli/il"
	"strconv"
	"strings"
)

var comma = "\""
var argfunc = "arguments"
var needSemicolon = false
var explicitArgs = false

func setComma(src string) {
	comma = src
}

func setArgFunc(src string) {
	argfunc = src
}

func setNeedSemicolon(src bool) {
	needSemicolon = src;
}

func setExplicitArgs(src bool) {
	explicitArgs = src
}

func inputValue(value il.Value) string {
	if value.GetKind() == "IntValue" {
		return "intValue(" + strconv.FormatInt(int64(value.(il.IntValue).Int), 10) + ")"
	} else if value.GetKind() == "StringValue" {
		return "stringValue(" + comma + value.(il.StringValue).String + comma + ")"
	} else if value.GetKind() == "BooleanValue" {
		return "boolValue(" + strconv.FormatBool(value.(il.BooleanValue).Boolean) + ")"
	} else if value.GetKind() == "EnumValue" {
		return "stringValue(\"" + value.(il.EnumValue).String + "\")"
	} else if value.GetKind() == "VariableValue" {
		return "refValue(" + comma + value.(il.VariableValue).Name + comma + ")"
	} else if value.GetKind() == "ListValue" {
		list := value.(il.ListValue)
		inner := make([]string, 0)
		for _, v := range list.Values {
			inner = append(inner, inputValue(v))
		}
		return "listValue(" + strings.Join(inner, ",") + ")"
	} else if value.GetKind() == "ObjectValue" {
		obj := value.(il.ObjectValue)
		inner := make([]string, 0)
		for _, f := range obj.Fields {
			inner = append(inner, "fieldValue("+comma+f.Name+comma+", "+inputValue(f.Value)+")")
		}
		return "objectValue(" + strings.Join(inner, ",") + ")"
	}

	panic("Unexpected input type: " + value.GetKind())
}

func outputSelectors(set *il.SelectionSet, output *Output, first bool) {
	output.Append("obj(")
	output.IndentAdd()
	isFirst := true
	for _, s := range set.Fields {
		if isFirst {
			isFirst = false
		} else {
			output.Append(",")
		}
		if len(s.Arguments) > 0 || explicitArgs {
			args := make([]string, 0)
			for _, a := range s.Arguments {
				args = append(args, "fieldValue(\""+a.Name+"\", "+inputValue(a.Value)+")")
			}
			output.WriteLine("field(" + comma + s.Name + comma + ", " + comma + s.Alias + comma + ", " + argfunc + "(" + strings.Join(args, ", ") + "), ")
		} else {
			output.WriteLine("field(" + comma + s.Name + comma + ", " + comma + s.Alias + comma + ", ")
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
		output.WriteLine("fragment(" + comma + fr.TypeName + comma + ", " + fr.Name + "Selector)")
	}
	for _, fr := range set.InlineFragments {
		if isFirst {
			isFirst = false
		} else {
			output.Append(",")
		}
		output.WriteLine("inline(" + comma + fr.TypeName + comma + ", ")
		outputSelectors(fr.Selection, output, false)
		output.Append(")")
	}
	output.IndentRemove()
	if needSemicolon && first {
		output.WriteLine(");")
	} else {
		output.WriteLine(")")
	}
}

func outputType(tp il.Type, output *Output) {
	if tp.GetKind() == "NotNull" {
		inner := tp.(il.NotNull)
		output.Append("notNull(")
		outputType(inner.Inner, output)
		output.Append(")")
	} else if tp.GetKind() == "Scalar" {
		scalar := tp.(il.Scalar)
		output.Append("scalar(" + comma + scalar.Name + comma + ")")
	} else if tp.GetKind() == "Enum" {
		output.Append("scalar(" + comma + "String" + comma + ")")
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
		outputSelectors(set, output, false)
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
	outputSelectors(set, output, true)
	output.IndentRemove()
}
