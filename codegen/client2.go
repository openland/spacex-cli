package codegen

import (
	"github.com/openland/spacex-cli/il"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GenerateClient2(model *il.Model, name string, to string) {
	output := NewOutput()
	output.WriteLine("/* tslint:disable */")
	output.WriteLine("/* eslint-disable */")
	output.WriteLine("import * as Types from './spacex.types';")
	output.WriteLine("import { SpaceXClientParameters, GraphqlActiveSubscription, QueryParameters, MutationParameters, SubscriptionParameters, GraphqlSubscriptionHandler, BaseSpaceXClient, SpaceQueryWatchParameters } from '@openland/spacex';")
	output.WriteLine("")
	output.WriteLine("export class " + name + " extends BaseSpaceXClient {")
	output.IndentAdd()
	output.WriteLine("constructor(params: SpaceXClientParameters) {")
	output.IndentAdd()
	output.WriteLine("super(params);")
	output.IndentRemove()
	output.WriteLine("}")

	output.WriteLine("withParameters(params: Partial<SpaceXClientParameters>) {")
	output.IndentAdd()
	output.WriteLine("return new " + name + "({ ... params, engine: this.engine, globalCache: this.globalCache});")
	output.IndentRemove()
	output.WriteLine("}")

	for _, q := range model.Queries {
		if len(q.Variables.Variables) > 0 {
			output.WriteLine("query" + q.Name + "(variables: Types." + q.Name + "Variables, params?: QueryParameters): Promise<Types." + q.Name + "> {")
			output.IndentAdd()
			output.WriteLine("return this.query('" + q.Name + "', variables, params);")
			output.IndentRemove()
			output.WriteLine("}")
		} else {
			output.WriteLine("query" + q.Name + "(params?: QueryParameters): Promise<Types." + q.Name + "> {")
			output.IndentAdd()
			output.WriteLine("return this.query('" + q.Name + "', undefined, params);")
			output.IndentRemove()
			output.WriteLine("}")
		}
	}

	for _, q := range model.Queries {
		if len(q.Variables.Variables) > 0 {
			output.WriteLine("refetch" + q.Name + "(variables: Types." + q.Name + "Variables, params?: QueryParameters): Promise<Types." + q.Name + "> {")
			output.IndentAdd()
			output.WriteLine("return this.refetch('" + q.Name + "', variables, params);")
			output.IndentRemove()
			output.WriteLine("}")
		} else {
			output.WriteLine("refetch" + q.Name + "(params?: QueryParameters): Promise<Types." + q.Name + "> {")
			output.IndentAdd()
			output.WriteLine("return this.refetch('" + q.Name + "', undefined, params);")
			output.IndentRemove()
			output.WriteLine("}")
		}
	}

	for _, q := range model.Queries {
		if len(q.Variables.Variables) > 0 {
			output.WriteLine("update" + q.Name + "(variables: Types." + q.Name + "Variables, updater: (data: Types." + q.Name + ") => Types." + q.Name + " | null): Promise<boolean> {")
			output.IndentAdd()
			output.WriteLine("return this.updateQuery(updater, '" + q.Name + "', variables);")
			output.IndentRemove()
			output.WriteLine("}")
		} else {
			output.WriteLine("update" + q.Name + "(updater: (data: Types." + q.Name + ") => Types." + q.Name + " | null): Promise<boolean> {")
			output.IndentAdd()
			output.WriteLine("return this.updateQuery(updater, '" + q.Name + "', undefined);")
			output.IndentRemove()
			output.WriteLine("}")
		}
	}

	for _, q := range model.Queries {
		if len(q.Variables.Variables) > 0 {
			output.WriteLine("use" + q.Name + "(variables: Types." + q.Name + "Variables, params: SpaceQueryWatchParameters & { suspense: false }): Types." + q.Name + " | null;")
			output.WriteLine("use" + q.Name + "(variables: Types." + q.Name + "Variables, params?: SpaceQueryWatchParameters): Types." + q.Name + ";")
			output.WriteLine("use" + q.Name + "(variables: Types." + q.Name + "Variables, params?: SpaceQueryWatchParameters): Types." + q.Name + " | null {")
			output.IndentAdd()
			output.WriteLine("return this.useQuery('" + q.Name + "', variables, params);")
			output.IndentRemove()
			output.WriteLine("}")
		} else {
			output.WriteLine("use" + q.Name + "(params: SpaceQueryWatchParameters & { suspense: false }): Types." + q.Name + " | null;")
			output.WriteLine("use" + q.Name + "(params?: SpaceQueryWatchParameters): Types." + q.Name + ";")
			output.WriteLine("use" + q.Name + "(params?: SpaceQueryWatchParameters): Types." + q.Name + " | null {")
			output.IndentAdd()
			output.WriteLine("return this.useQuery('" + q.Name + "', undefined, params);")
			output.IndentRemove()
			output.WriteLine("}")
		}
	}

	for _, q := range model.Mutations {
		if len(q.Variables.Variables) > 0 {
			output.WriteLine("mutate" + q.Name + "(variables: Types." + q.Name + "Variables, params?: MutationParameters): Promise<Types." + q.Name + "> {")
			output.IndentAdd()
			output.WriteLine("return this.mutate('" + q.Name + "', variables, params);")
			output.IndentRemove()
			output.WriteLine("}")
		} else {
			output.WriteLine("mutate" + q.Name + "(params?: MutationParameters): Promise<Types." + q.Name + "> {")
			output.IndentAdd()
			output.WriteLine("return this.mutate('" + q.Name + "', undefined, params);")
			output.IndentRemove()
			output.WriteLine("}")
		}
	}

	for _, q := range model.Subscriptions {
		if len(q.Variables.Variables) > 0 {
			output.WriteLine("subscribe" + q.Name + "(variables: Types." + q.Name + "Variables, handler: GraphqlSubscriptionHandler<Types." + q.Name + ">, params?: SubscriptionParameters): GraphqlActiveSubscription<Types." + q.Name + "> {")
			output.IndentAdd()
			output.WriteLine("return this.subscribe(handler, '" + q.Name + "', variables, params);")
			output.IndentRemove()
			output.WriteLine("}")
		} else {
			output.WriteLine("subscribe" + q.Name + "(handler: GraphqlSubscriptionHandler<Types." + q.Name + ">, params?: SubscriptionParameters): GraphqlActiveSubscription<Types." + q.Name + "> {")
			output.IndentAdd()
			output.WriteLine("return this.subscribe(handler, '" + q.Name + "', undefined, params);")
			output.IndentRemove()
			output.WriteLine("}")
		}
	}

	output.IndentRemove()
	output.WriteLine("}")

	// Result
	err := os.MkdirAll(filepath.Dir(to), os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(to, []byte(output.String()), 0644)
	if err != nil {
		panic(err)
	}
}
