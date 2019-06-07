package codegen

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/openland/spacex-cli/il"
)

func inputValue(value il.Value) string {
	if value.GetKind() == "IntValue" {
		return "intValue(" + strconv.FormatInt(int64(value.(il.IntValue).Int), 10) + ")"
	} else if value.GetKind() == "StringValue" {
		return "stringValue(\"" + value.(il.StringValue).String + "\")"
	} else if value.GetKind() == "BooleanValue" {
		return "boolValue(" + strconv.FormatBool(value.(il.BooleanValue).Boolean) + ")"
	} else if value.GetKind() == "EnumValue" {
		return "stringValue(\"" + value.(il.EnumValue).String + "\")"
	} else if value.GetKind() == "VariableValue" {
		return "refValue(\"" + value.(il.VariableValue).Name + "\")"
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
			inner = append(inner, "fieldValue(\""+f.Name+"\", "+inputValue(f.Value)+")")
		}
		return "objectValue(" + strings.Join(inner, ",") + ")"
	}

	panic("Unexpected input type: " + value.GetKind())
}

func outputSelectors(set *il.SelectionSet, output *Output) {
	output.Append("obj(")
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
				args = append(args, "fieldValue(\""+a.Name+"\", "+inputValue(a.Value)+")")
			}
			output.WriteLine("field(\"" + s.Name + "\",\"" + s.Alias + "\", arguments(" + strings.Join(args, ", ") + "), ")
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
	output.WriteLine(")")
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

func GenerateKotlin(model *il.Model, to string, pgk string) {
	output := NewOutput()
	output.WriteLine("package " + pgk)
	output.WriteLine("")
	output.WriteLine("import com.openland.spacex.*")
	output.WriteLine("import com.openland.spacex.gen.*")
	output.WriteLine("import org.json.*")
	output.WriteLine("")

	//
	// Selectors
	//

	for _, f := range model.Fragments {
		output.NextScope()
		output.WriteLine("private val " + f.Name + "Selector = ")
		output.IndentAdd()
		generateSelector(f.SelectionSet, output)
		output.IndentRemove()
		output.WriteLine("")
	}

	for _, f := range model.Queries {
		output.NextScope()
		output.WriteLine("private val " + f.Name + "Selector = ")
		output.IndentAdd()
		generateSelector(f.SelectionSet, output)
		output.IndentRemove()
	}

	for _, f := range model.Mutations {
		output.NextScope()
		output.WriteLine("private val " + f.Name + "Selector = ")
		output.IndentAdd()
		generateSelector(f.SelectionSet, output)
		output.IndentRemove()
	}

	for _, f := range model.Subscriptions {
		output.NextScope()
		output.WriteLine("private val " + f.Name + "Selector = ")
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
		output.WriteLine("override val body = \"" + strings.Replace(f.Body, "$", "\\$", -1) + "\"")
		output.WriteLine("override val selector = " + f.Name + "Selector")
		output.IndentRemove()
		output.WriteLine("}")
	}
	for _, f := range model.Mutations {
		output.WriteLine("val " + f.Name + " = object: OperationDefinition {")
		output.IndentAdd()
		output.WriteLine("override val name = \"" + f.Name + "\"")
		output.WriteLine("override val kind = OperationKind.MUTATION")
		output.WriteLine("override val body = \"" + strings.Replace(f.Body, "$", "\\$", -1) + "\"")
		output.WriteLine("override val selector = " + f.Name + "Selector")
		output.IndentRemove()
		output.WriteLine("}")
	}
	for _, f := range model.Subscriptions {

		output.WriteLine("val " + f.Name + " = object: OperationDefinition {")
		output.IndentAdd()
		output.WriteLine("override val name = \"" + f.Name + "\"")
		output.WriteLine("override val kind = OperationKind.SUBSCRIPTION")
		output.WriteLine("override val body = \"" + strings.Replace(f.Body, "$", "\\$", -1) + "\"")
		output.WriteLine("override val selector = " + f.Name + "Selector")
		output.IndentRemove()
		output.WriteLine("}")
	}

	output.WriteLine("fun operationByName(name: String): OperationDefinition {")
	output.IndentAdd()
	for _, f := range model.Queries {
		output.WriteLine("if (name == \"" + f.Name + "\") return " + f.Name)
	}
	for _, f := range model.Mutations {
		output.WriteLine("if (name == \"" + f.Name + "\") return " + f.Name)
	}
	for _, f := range model.Subscriptions {
		output.WriteLine("if (name == \"" + f.Name + "\") return " + f.Name)
	}
	output.WriteLine("error(\"Unknown operation: $name\")")
	output.IndentRemove()
	output.WriteLine("}")

	output.IndentRemove()
	output.WriteLine("}")

	// Write result
	err := os.MkdirAll(filepath.Dir(to), os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(to, []byte(output.String()), 0644)
	if err != nil {
		panic(err)
	}
}
