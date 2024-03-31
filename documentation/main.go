package main

import (
	"context"
	"encoding/json"
	"fmt"
)

type Documentation struct{}

func (*Documentation) Init() *Directory {
	return dag.CurrentModule().Source().Directory("template")
}

type DocumentationBuilder struct {
	*Container
}

func (*Documentation) Builder(
	ctx context.Context,
	directory *Directory,
) (*DocumentationBuilder, error) {
	const packageJsonFilename string = "package.json"

	packageJsonString, err := directory.File(packageJsonFilename).Contents(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to read %q file: %w", packageJsonFilename, err)
	}

	var configuration struct {
		Hugo struct {
			Version string
		}
	}

	err = json.Unmarshal([]byte(packageJsonString), &configuration)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %q file: %w", packageJsonFilename, err)
	}

	if configuration.Hugo.Version == "" {
		return nil, fmt.Errorf("Hugo version is not set in %q file", packageJsonFilename)
	}

	builder := &DocumentationBuilder{}

	kroki := dag.Kroki()

	builder.Container = dag.Redhat().Minimal().Container().
		With(dag.Nodejs().Installation).
		With(dag.Golang().Installation).
		With(dag.Hugo(configuration.Hugo.Version).Installation).
		WithMountedDirectory(".", directory).
		WithExec([]string{"npm clean-install"}).
		WithEntrypoint([]string{"npm", "run", "all", "--"}).
		WithoutDefaultArgs()

	// FIXME
	// Doesn’t work at the time.
	//WithServiceBinding("kroki", kroki.Server())
	_ = kroki

	return builder, nil
}

type DocumentationBuild struct {
	*Directory
}

func (builder *DocumentationBuilder) Build(
	// +optional
	args []string,
) *DocumentationBuild {
	build := &DocumentationBuild{
		Directory: builder.WithExec(args).Directory("public"),
	}

	return build
}

func (build *DocumentationBuild) Container() *Container {
	return dag.Caddy(build.Directory).Container()
}

func (build *DocumentationBuild) Server() *Service {
	return build.Container().AsService()
}
