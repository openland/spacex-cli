package codegen

import (
	"github.com/openland/spacex-cli/il"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GenerateTypescript(model *il.Model, to string) {
	output := NewOutput()

	setComma("'")
	setArgFunc("args")
	setNeedSemicolon(true)
	setExplicitArgs(true)

	output.WriteLine("/* tslint:disable */")
	output.WriteLine("/* eslint-disable */")
	output.WriteLine("// @ts-ignore")
	output.WriteLine("import { WebDefinitions, OperationDefinition, Definitions as AllDefinitions } from '@openland/spacex';")
	output.WriteLine("// @ts-ignore")
	output.WriteLine("const { list, notNull, scalar, field, obj, inline, fragment, args, fieldValue, refValue, intValue, floatValue, stringValue, boolValue, listValue, objectValue } = WebDefinitions;")
	output.WriteLine("")

	for _, f := range model.Fragments {
		output.NextScope()
		output.WriteLine("const " + f.Name + "Selector = ")
		output.IndentAdd()
		generateSelector(f.SelectionSet, output)
		output.IndentRemove()
		output.WriteLine("")
	}

	for _, f := range model.Queries {
		output.NextScope()
		output.WriteLine("const " + f.Name + "Selector = ")
		output.IndentAdd()
		generateSelector(f.SelectionSet, output)
		output.IndentRemove()
	}

	for _, f := range model.Mutations {
		output.NextScope()
		output.WriteLine("const " + f.Name + "Selector = ")
		output.IndentAdd()
		generateSelector(f.SelectionSet, output)
		output.IndentRemove()
	}

	for _, f := range model.Subscriptions {
		output.NextScope()
		output.WriteLine("const " + f.Name + "Selector = ")
		output.IndentAdd()
		generateSelector(f.SelectionSet, output)
		output.IndentRemove()
	}

	output.WriteLine("export const Operations: { [key: string]: OperationDefinition } = {")
	output.IndentAdd()
	for _, f := range model.Queries {
		output.WriteLine("" + f.Name + ": {")
		output.IndentAdd()
		output.WriteLine("kind: 'query',")
		output.WriteLine("name: '" + f.Name + "',")
		output.WriteLine("body: '" + f.Body + "',")
		output.WriteLine("selector: " + f.Name + "Selector")
		output.IndentRemove()
		output.WriteLine("},")
	}
	for _, f := range model.Mutations {
		output.WriteLine("" + f.Name + ": {")
		output.IndentAdd()
		output.WriteLine("kind: 'mutation',")
		output.WriteLine("name: '" + f.Name + "',")
		output.WriteLine("body: '" + f.Body + "',")
		output.WriteLine("selector: " + f.Name + "Selector")
		output.IndentRemove()
		output.WriteLine("},")
	}
	for _, f := range model.Subscriptions {
		output.WriteLine("" + f.Name + ": {")
		output.IndentAdd()
		output.WriteLine("kind: 'subscription',")
		output.WriteLine("name: '" + f.Name + "',")
		output.WriteLine("body: '" + f.Body + "',")
		output.WriteLine("selector: " + f.Name + "Selector")
		output.IndentRemove()
		output.WriteLine("},")
	}
	output.IndentRemove()
	output.WriteLine("};")

	output.WriteLine("export const Definitions: AllDefinitions = { operations: Operations };")

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