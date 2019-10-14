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

	output.WriteLine("import { list, notNull, scalar, field, obj, inline, fragment, args, fieldValue, refValue, intValue, floatValue, stringValue, boolValue, listValue, objectValue } from 'openland-graphql/spacex/types';")
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

	output.WriteLine("export const Operations = {")
	output.IndentAdd()
	for _, f := range model.Queries {
		output.WriteLine("" + f.Name + ": {")
		output.IndentAdd()
		output.WriteLine("type: 'query',")
		output.WriteLine("name: '" + f.Name + "',")
		output.WriteLine("body: '" + f.Body + "',")
		output.WriteLine("selector: " + f.Name + "Selector")
		output.IndentRemove()
		output.WriteLine("},")
	}
	for _, f := range model.Mutations {
		output.WriteLine("" + f.Name + ": {")
		output.IndentAdd()
		output.WriteLine("type: 'mutation',")
		output.WriteLine("name: '" + f.Name + "',")
		output.WriteLine("body: '" + f.Body + "',")
		output.WriteLine("selector: " + f.Name + "Selector")
		output.IndentRemove()
		output.WriteLine("},")
	}
	for _, f := range model.Subscriptions {
		output.WriteLine("" + f.Name + ": {")
		output.IndentAdd()
		output.WriteLine("type: 'subscription',")
		output.WriteLine("name: '" + f.Name + "',")
		output.WriteLine("body: '" + f.Body + "',")
		output.WriteLine("selector: " + f.Name + "Selector")
		output.IndentRemove()
		output.WriteLine("},")
	}
	output.IndentRemove()
	output.WriteLine("};");

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