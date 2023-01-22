package newproject

import (
	"log"

	"github.com/pjmd89/gogql/lib/generate"
	"github.com/pjmd89/gogql/lib/generate/gqltypes"
	"golang.org/x/exp/slices"
)

func Generate(gqlGenerate generate.GqlGenerate, driverDB gqltypes.DriverDB) {

	if gqlGenerate.ModulePath != "" && gqlGenerate.ModuleName != "" {
		generateSchema(gqlGenerate, driverDB)
	} else {
		log.Fatal("debes indicar el path del schema, la carpeta raiz del proyecto y el nombre del modulo")
	}
}
func generateSchema(render generate.GqlGenerate, driverDB gqltypes.DriverDB) {
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
				types.ModelType = append(types.ModelType, gqltypes.NewModel(render, k, v, render.Schema.Types, driverDB))
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
