package codegen

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/openland/spacex-cli/il"
)

func write(output *Output, to string) {
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

func fileHeader(output *Output, pkg string) {
	output.WriteLine("package " + pkg)
	output.WriteLine("")
	output.WriteLine("import com.openland.spacex.*")
	output.WriteLine("import com.openland.spacex.gen.*")
	output.WriteLine("import org.json.*")
	output.WriteLine("")
}

func GenerateKotlin(model *il.Model, to string, pkg string) {

	// Clean output folder

	err := os.RemoveAll(to + "/operations")
	if err != nil {
		panic(err)
	}

	//
	// Selectors
	//
	output := NewOutput()
	fileHeader(output, pkg)
	for _, f := range model.Fragments {
		output.NextScope()
		output.WriteLine("internal val " + f.Name + "Selector = ")
		output.IndentAdd()
		generateSelector(f.SelectionSet, output)
		output.IndentRemove()
		output.WriteLine("")
	}
	write(output, to+"/operations/Fragments.kt")

	for _, f := range model.Queries {
		output = NewOutput()
		fileHeader(output, pkg)

		output.NextScope()
		output.WriteLine("internal val " + f.Name + "Selector = ")
		output.IndentAdd()
		generateSelector(f.SelectionSet, output)
		output.IndentRemove()

		output.WriteLine("val " + f.Name + " = object: OperationDefinition {")
		output.IndentAdd()
		output.WriteLine("override val name = \"" + f.Name + "\"")
		output.WriteLine("override val kind = OperationKind.QUERY")
		output.WriteLine("override val body = \"" + strings.Replace(f.Body, "$", "\\$", -1) + "\"")
		output.WriteLine("override val selector = " + f.Name + "Selector")
		output.IndentRemove()
		output.WriteLine("}")

		write(output, to+"/operations/queries/"+f.Name+".kt")
	}

	for _, f := range model.Mutations {
		output = NewOutput()
		fileHeader(output, pkg)

		output.NextScope()
		output.WriteLine("internal val " + f.Name + "Selector = ")
		output.IndentAdd()
		generateSelector(f.SelectionSet, output)
		output.IndentRemove()

		output.WriteLine("val " + f.Name + " = object: OperationDefinition {")
		output.IndentAdd()
		output.WriteLine("override val name = \"" + f.Name + "\"")
		output.WriteLine("override val kind = OperationKind.MUTATION")
		output.WriteLine("override val body = \"" + strings.Replace(f.Body, "$", "\\$", -1) + "\"")
		output.WriteLine("override val selector = " + f.Name + "Selector")
		output.IndentRemove()
		output.WriteLine("}")

		write(output, to+"/operations/mutations/"+f.Name+".kt")
	}

	for _, f := range model.Subscriptions {
		output = NewOutput()
		fileHeader(output, pkg)

		output.NextScope()
		output.WriteLine("internal val " + f.Name + "Selector = ")
		output.IndentAdd()
		generateSelector(f.SelectionSet, output)
		output.IndentRemove()

		output.WriteLine("val " + f.Name + " = object: OperationDefinition {")
		output.IndentAdd()
		output.WriteLine("override val name = \"" + f.Name + "\"")
		output.WriteLine("override val kind = OperationKind.SUBSCRIPTION")
		output.WriteLine("override val body = \"" + strings.Replace(f.Body, "$", "\\$", -1) + "\"")
		output.WriteLine("override val selector = " + f.Name + "Selector")
		output.IndentRemove()
		output.WriteLine("}")

		write(output, to+"/operations/subscriptions/"+f.Name+".kt")
	}

	output = NewOutput()
	output.WriteLine("package " + pkg)
	output.WriteLine("")
	output.WriteLine("import com.openland.spacex.*")
	output.WriteLine("import com.openland.spacex.gen.*")
	output.WriteLine("import org.json.*")
	output.WriteLine("")

	output.WriteLine("")
	output.WriteLine("object Operations {")
	output.IndentAdd()

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

	write(output, to+"/Operations.kt")
}
