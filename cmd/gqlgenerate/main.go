package main

import (
	"embed"
	"flag"
	"log"

	"github.com/pjmd89/gogql/lib/generate"
	"github.com/pjmd89/gogql/lib/generate/gqltypes"
	"github.com/pjmd89/gogql/lib/gql"
	"golang.org/x/exp/slices"
)

var (
	//go:embed schema
	embedFS embed.FS
)

func main() {
	var (
		schemaPath     string = ""
		moduleName            = ""
		modulePath            = ""
		modelPath             = "models"
		resolverPath          = "resolvers"
		unionPath             = "unions"
		scalarPath            = "scalars"
		enumPath              = "enums"
		objecttypePath        = "objectTypes"
	)

	flag.StringVar(&schemaPath, "schema", "", "Ruta de la carpeta contenedora del esquema de GraphQL")
	flag.StringVar(&modulePath, "module-path", "", "Ruta donde se guardaran los modelos generados")
	flag.StringVar(&moduleName, "module-name", "", "Ruta donde se guardaran los modelos generados")
	flag.StringVar(&unionPath, "union-path", unionPath, "Ruta donde se guardaran los modelos generados")
	flag.StringVar(&scalarPath, "scalar-path", scalarPath, "Ruta donde se guardaran los modelos generados")
	flag.StringVar(&enumPath, "enum-path", enumPath, "Ruta donde se guardaran los modelos generados")
	flag.StringVar(&modelPath, "model-path", modelPath, "Ruta donde se guardaran los modelos generados")
	flag.StringVar(&resolverPath, "resolver-path", resolverPath, "Ruta donde se guardaran los modelos generados")
	flag.StringVar(&objecttypePath, "objecttype-path", objecttypePath, "Ruta donde se guardaran los modelos generados")
	flag.Parse()
	render := generate.GqlGenerate{
		SchemaPath:     "schema",
		ModuleName:     moduleName,
		ModulePath:     modulePath,
		ModelPath:      modelPath,
		ResolverPath:   resolverPath,
		UnionPath:      unionPath,
		ScalarPath:     scalarPath,
		EnumPath:       enumPath,
		ObjecttypePath: objecttypePath,
	}
	if schemaPath != "" && modulePath != "" && moduleName != "" {
		generateSchema(render)
	} else {
		log.Fatal("debes indicar el path del schema, la carpeta raiz del proyecto y el nombre del modulo (-schema, -module-path, -module-name)")
	}
}

func generateSchema(render generate.GqlGenerate) {
	types := generate.RenderTypes{
		ModelType:  make([]generate.ModelDef, 0),
		ObjectType: make([]generate.ObjectTypeDef, 0),
		EnumType:   make([]generate.EnumDef, 0),
		ScalarType: make([]generate.ScalarDef, 0),
	}
	types.MainPath = render.ModulePath + "/generate/main.go"
	gql := gql.Init("", embedFS, render.SchemaPath)

	if gql.GetSchema().Query != nil {
		generate.OmitObject = append(generate.OmitObject, gql.GetSchema().Query.Name)
	}
	if gql.GetSchema().Mutation != nil {
		generate.OmitObject = append(generate.OmitObject, gql.GetSchema().Mutation.Name)
	}
	if gql.GetSchema().Subscription != nil {
		generate.OmitObject = append(generate.OmitObject, gql.GetSchema().Subscription.Name)
	}
	for k, v := range gql.GetSchema().Types {
		if !slices.Contains(generate.OmitObject, k) {
			switch v.Kind {
			case "OBJECT":
				types.ModelType = append(types.ModelType, gqltypes.NewModel(render, k, v, gql.GetSchema().Types))
				break
			case "ENUM":
				types.EnumType = append(types.EnumType, gqltypes.NewEnum(render, k, v))
				break
			case "SCALAR":
				scalar := gqltypes.NewScalar(render, k, v)
				if scalar != nil {
					types.ScalarType = append(types.ScalarType, *scalar)
				}
				break
			case "UNION":
				types.UnionType.PackageName = render.ModelPath
				types.UnionType.Type = append(types.UnionType.Type, gqltypes.NewUnion(k, v))
				types.UnionType.FilePath = render.ModulePath + "/generate/" + render.ModelPath + "/union_definition.go"
				break
			}
		}
	}
	for _, v := range types.ModelType {
		types.ObjectType = append(types.ObjectType, gqltypes.NewObjectType(render, v, gql.GetSchema()))
	}
	types.ScalarPath = render.ModuleName + "/" + render.ResolverPath + "/" + render.ScalarPath
	gqltypes.ModelTmpl(types)
	gqltypes.ObjectTypeTmpl(types)
	gqltypes.EnumTmpl(types)
	gqltypes.ScalarTmpl(types)
	gqltypes.UnionTmpl(types)
	gqltypes.Maintmpl(types)
}

//go:generate go run main.go -scheme=$SCHEME -model-path=$MODELPATH
