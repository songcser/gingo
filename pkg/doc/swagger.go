package doc

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/spec"
	"github.com/songcser/gingo/pkg/api"
	"go/build"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

var swagger Swagger

type SwaggerConfig struct {
	Description string `json:"description,omitempty"`
	Title       string `json:"title,omitempty"`
}

type Swagger struct {
	apis    []SwaggerApi
	swagger spec.Swagger
}

type SwaggerApi struct {
	typ    string
	path   string
	api    api.Api
	handle gin.HandlerFunc
}

func AddSwaggerApi(a api.Api, path string) {
	s := SwaggerApi{typ: "api", api: a, path: path}
	swagger.apis = append(swagger.apis, s)
}

func AddSwaggerHandle(h gin.HandlerFunc, path string) {
	s := SwaggerApi{typ: "handler", handle: h, path: path}
	swagger.apis = append(swagger.apis, s)
}

func InitSwagger(conf SwaggerConfig) {
	swagger = Swagger{
		apis: make([]SwaggerApi, 0),
		swagger: spec.Swagger{
			SwaggerProps: spec.SwaggerProps{
				Info: &spec.Info{
					InfoProps: spec.InfoProps{
						Contact:     &spec.ContactInfo{},
						License:     nil,
						Title:       conf.Title,
						Description: conf.Description,
					},
					VendorExtensible: spec.VendorExtensible{
						Extensions: spec.Extensions{},
					},
				},
				Paths: &spec.Paths{
					Paths: make(map[string]spec.PathItem),
					VendorExtensible: spec.VendorExtensible{
						Extensions: nil,
					},
				},
				Definitions:         make(map[string]spec.Schema),
				SecurityDefinitions: make(map[string]*spec.SecurityScheme),
			},
		},
	}
}

func GenerateSwagger() {
	apis := swagger.apis
	wd, _ := os.Getwd()
	fmt.Println(wd)
	pkgName, _ := getPkgName(wd)
	fmt.Println(pkgName)
	for _, a := range apis {
		typ := a.typ
		if typ == "api" {
			pkg := reflect.TypeOf(a.api).PkgPath()
			reflect.TypeOf(a.api).PkgPath()
			fmt.Println(pkg)
		}
	}
}

func getPkgName(searchDir string) (string, error) {
	cmd := exec.Command("go", "list", "-f={{.ImportPath}}")
	cmd.Dir = searchDir
	var stdout, stderr strings.Builder

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("execute go list command, %s, stdout:%s, stderr:%s", err, stdout.String(), stderr.String())
	}

	outStr, _ := stdout.String(), stderr.String()

	if outStr[0] == '_' { // will shown like _/{GOPATH}/src/{YOUR_PACKAGE} when NOT enable GO MODULE.
		outStr = strings.TrimPrefix(outStr, "_"+build.Default.GOPATH+"/src/")
	}

	f := strings.Split(outStr, "\n")

	outStr = f[0]

	return outStr, nil
}
