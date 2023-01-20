package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pjmd89/gogql/lib/generate"
	"github.com/pjmd89/gogql/lib/generate/gqltypes"
	"golang.org/x/exp/slices"
	"golang.org/x/mod/modfile"
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
	//module-name
	goModBytes, err := ioutil.ReadFile("go.mod")
	if err != nil {
		//exitf(func() {}, 1, "%+v\n", err)
	}

	modName := modfile.ModulePath(goModBytes)
	fmt.Fprintf(os.Stdout, "modName=%+v\n", modName)

	//module-path
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)
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

	if render.Schema.Query != nil {
		generate.OmitObject = append(generate.OmitObject, render.Schema.Query.Name)
	}
	if render.Schema.Mutation != nil {
		generate.OmitObject = append(generate.OmitObject, render.Schema.Mutation.Name)
	}
	if render.Schema.Subscription != nil {
		generate.OmitObject = append(generate.OmitObject, render.Schema.Subscription.Name)
	}
	for k, v := range render.Schema.Types {
		if !slices.Contains(generate.OmitObject, k) {
			switch v.Kind {
			case "OBJECT":
				types.ModelType = append(types.ModelType, gqltypes.NewModel(render, k, v, render.Schema.Types))
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
		types.ObjectType = append(types.ObjectType, gqltypes.NewObjectType(render, v, render.Schema))
	}
	types.ScalarPath = render.ModuleName + "/" + render.ResolverPath + "/" + render.ScalarPath
	gqltypes.ModelTmpl(types)
	gqltypes.ObjectTypeTmpl(types)
	gqltypes.EnumTmpl(types)
	gqltypes.ScalarTmpl(types)
	gqltypes.UnionTmpl(types)
	gqltypes.Maintmpl(types)
}
