package codegen

import (
	"github.com/openland/spacex-cli/il"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GenerateSwift(model *il.Model, to string) {
	output := NewOutput()

	for _, f := range model.Fragments {
		output.NextScope()
		output.WriteLine("private let " + f.Name + "Selector = ")
		output.IndentAdd()
		generateSelector(f.SelectionSet, output)
		output.IndentRemove()
		output.WriteLine("")
	}

	for _, f := range model.Queries {
		output.NextScope()
		output.WriteLine("private let " + f.Name + "Selector = ")
		output.IndentAdd()
		generateSelector(f.SelectionSet, output)
		output.IndentRemove()
	}

	for _, f := range model.Mutations {
		output.NextScope()
		output.WriteLine("private let " + f.Name + "Selector = ")
		output.IndentAdd()
		generateSelector(f.SelectionSet, output)
		output.IndentRemove()
	}

	for _, f := range model.Subscriptions {
		output.NextScope()
		output.WriteLine("private let " + f.Name + "Selector = ")
		output.IndentAdd()
		generateSelector(f.SelectionSet, output)
		output.IndentRemove()
	}

	output.WriteLine("")
	output.WriteLine("class Operations {")
	output.IndentAdd()
	output.WriteLine("static let shared = Operations()")
	output.WriteLine("")
	output.WriteLine("private init() { }")
	for _, f := range model.Queries {
		output.WriteLine("let " + f.Name + " = OperationDefinition(")
		output.IndentAdd()
		output.WriteLine("\"" + f.Name + "\",")
		output.WriteLine(".query, ")
		output.WriteLine("\"" + f.Body + "\",")
		output.WriteLine(f.Name + "Selector")
		output.IndentRemove()
		output.WriteLine(")")
	}
	for _, f := range model.Mutations {
		output.WriteLine("let " + f.Name + " = OperationDefinition(")
		output.IndentAdd()
		output.WriteLine("\"" + f.Name + "\",")
		output.WriteLine(".mutation, ")
		output.WriteLine("\"" + f.Body + "\",")
		output.WriteLine(f.Name + "Selector")
		output.IndentRemove()
		output.WriteLine(")")
	}
	for _, f := range model.Subscriptions {
		output.WriteLine("let " + f.Name + " = OperationDefinition(")
		output.IndentAdd()
		output.WriteLine("\"" + f.Name + "\",")
		output.WriteLine(".subscription, ")
		output.WriteLine("\"" + f.Body + "\",")
		output.WriteLine(f.Name + "Selector")
		output.IndentRemove()
		output.WriteLine(")")
	}
	output.WriteLine("")
	output.WriteLine("func operationByName(_ name: String) -> OperationDefinition {")
	output.IndentAdd()
	for _, f := range model.Queries {
		output.WriteLine("if name == \"" + f.Name + "\" { return " + f.Name + " }")
	}
	for _, f := range model.Mutations {
		output.WriteLine("if name == \"" + f.Name + "\" { return " + f.Name + " }")
	}
	for _, f := range model.Subscriptions {
		output.WriteLine("if name == \"" + f.Name + "\" { return " + f.Name + " }")
	}
	output.WriteLine("fatalError(\"Unknown operation: \" + name)")
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
