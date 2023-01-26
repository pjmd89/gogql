package generate

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/pjmd89/gqlparser/v2/ast"
	"golang.org/x/mod/modfile"
)

var (
	schemaPath     = "./schema"
	modelPath      = "models"
	resolverPath   = "resolvers"
	unionPath      = "unions"
	scalarPath     = "scalars"
	enumPath       = "enums"
	objecttypePath = "objectTypes"
	libPath        = "lib"
)

func NewGqlGenerate(schema *ast.Schema, schemaPath string) (r GqlGenerate) {
	moduleByte, err := ioutil.ReadFile("go.mod")
	moduleName := modfile.ModulePath(moduleByte)
	if err != nil {
		log.Fatalln(err.Error())
	}
	modulePath, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}
	r = GqlGenerate{
		Schema:         schema,
		ModuleName:     string(moduleName),
		ModulePath:     modulePath,
		ModelPath:      modelPath,
		ResolverPath:   resolverPath,
		UnionPath:      unionPath,
		ScalarPath:     scalarPath,
		EnumPath:       enumPath,
		ObjecttypePath: objecttypePath,
		LibPath:        libPath,
		SchemaPath:     schemaPath,
	}
	return
}
