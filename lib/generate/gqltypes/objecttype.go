package gqltypes

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pjmd89/gogql/lib/generate"
	"github.com/pjmd89/gqlparser/v2/ast"
)

func NewObjectType(render generate.GqlGenerate, value generate.ModelDef, schema *ast.Schema) (oType generate.ObjectTypeDef) {
	oType.FilePath = value.ObjectTypePath
	oType.QueryPath = value.QueryPath
	oType.MutationPath = value.MutationPath
	oType.SubscriptionPath = value.SubscriptionPath
	oType.Name = strings.Title(value.RealName)
	oType.RealName = value.RealName
	oType.PackageName = strings.ToLower(value.RealName)
	oType.DefinitionPath = render.ModuleName + "/resolvers/" + render.ObjecttypePath + "/" + oType.PackageName
	oType.ModelPath = render.ModelPath
	oType.ModuleName = render.ModuleName
	if schema.Query != nil {
		for _, v := range schema.Query.Fields {
			typeName := v.Type.NamedType
			if typeName == "" {
				typeName = v.Type.Elem.NamedType
			}
			if value.RealName == typeName {
				oType.HasQueries = true

				oType.QueryResolvers = append(oType.QueryResolvers, map[string]string{"Name": v.Name, "Resolver": v.Name + "Query"})
			}
		}
	}
	if schema.Mutation != nil {
		for _, v := range schema.Mutation.Fields {
			typeName := v.Type.NamedType
			if typeName == "" {
				typeName = v.Type.Elem.NamedType
			}
			if value.RealName == typeName {
				oType.HasMutations = true
				oType.MutationResolvers = append(oType.MutationResolvers, map[string]string{"Name": v.Name, "Resolver": v.Name + "Mutation"})
			}
		}
	}
	if schema.Subscription != nil {
		for _, v := range schema.Subscription.Fields {
			typeName := v.Type.NamedType
			if typeName == "" {
				typeName = v.Type.Elem.NamedType
			}
			if value.RealName == typeName {
				oType.HasSubscriptions = true
				oType.SubscriptionResolvers = append(oType.SubscriptionResolvers, map[string]string{"Name": v.Name, "Resolver": v.Name + "Subscription"})
			}
		}
	}

	return
}
func ObjectTypeTmpl(types generate.RenderTypes) {
	mt, err := template.New("objecttypes.tmpl").Parse(string(generate.Objecttypetmpl))
	if err != nil {
		panic(err)
	}
	qt, err := template.New("queries.tmpl").Parse(string(generate.Queriestmpl))
	if err != nil {
		panic(err)
	}
	mmt, err := template.New("mutations.tmpl").Parse(string(generate.Mutationstmpl))
	if err != nil {
		panic(err)
	}
	st, err := template.New("mutations.tmpl").Parse(string(generate.Subscriptionstmpl))
	if err != nil {
		panic(err)
	}
	for _, v := range types.ObjectType {
		//definitions
		dir := filepath.Dir(v.FilePath)
		os.MkdirAll(dir, 0770)
		file, err := os.Create(v.FilePath)
		if err != nil {
			panic(err)
		}
		var tmpl bytes.Buffer
		err = mt.Execute(&tmpl, v)
		if err != nil {
			panic(err)
		}
		/*
			x, err := format.Source(tmpl.Bytes())
			if err != nil {
				log.Fatal(err.Error())
			}
			file.Write(x)
			//*/
		file.Write(tmpl.Bytes())

		//queries
		if len(v.QueryResolvers) > 0 {
			dir := filepath.Dir(v.QueryPath)
			os.MkdirAll(dir, 0770)
			file, err := os.Create(v.QueryPath)
			if err != nil {
				panic(err)
			}
			var tmpl bytes.Buffer
			err = qt.Execute(&tmpl, v)
			if err != nil {
				panic(err)
			}
			x, err := format.Source(tmpl.Bytes())
			if err != nil {
				log.Fatal(err.Error())
			}
			file.Write(x)
		}
		//mutations
		if len(v.MutationResolvers) > 0 {
			dir := filepath.Dir(v.MutationPath)
			os.MkdirAll(dir, 0770)
			file, err := os.Create(v.MutationPath)
			if err != nil {
				panic(err)
			}
			var tmpl bytes.Buffer
			err = mmt.Execute(&tmpl, v)
			if err != nil {
				panic(err)
			}
			x, err := format.Source(tmpl.Bytes())
			if err != nil {
				log.Fatal(err.Error())
			}
			file.Write(x)
		}
		//subscriptions
		if len(v.SubscriptionResolvers) > 0 {
			dir := filepath.Dir(v.SubscriptionPath)
			os.MkdirAll(dir, 0770)
			file, err := os.Create(v.SubscriptionPath)
			if err != nil {
				panic(err)
			}
			var tmpl bytes.Buffer
			err = st.Execute(&tmpl, v)
			if err != nil {
				panic(err)
			}
			x, err := format.Source(tmpl.Bytes())
			if err != nil {
				log.Fatal(err.Error())
			}
			file.Write(x)
		}

	}
}
