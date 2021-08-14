package main

import (
	"fmt"
	"os"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		os.Exit(2)
	}

	ginGqlgen := NewGinGglgen()

	err = api.Generate(cfg,
		api.AddPlugin(ginGqlgen),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}

type gqlgenBinding struct {
}

func NewGinGglgen() plugin.Plugin {
	return &gqlgenBinding{}
}

func (g *gqlgenBinding) Name() string {
	return "GqlgenBinding"
}

func (g *gqlgenBinding) MutateConfig(cfg *config.Config) error {
	cfg.Directives["binding"] = config.DirectiveConfig{SkipRuntime: false}

	return nil
}

func (g *gqlgenBinding) InjectSourceEarly() *ast.Source {
	return &ast.Source{
		Name:    "gqlgenBinding/directives.graphql",
		BuiltIn: false,
		Input: `
            directive @binding(constraint: String!, trim: Boolean) on INPUT_FIELD_DEFINITION
        `,
	}
}
