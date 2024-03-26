package main

import (
	"context"
	"encoding/json"
	"fmt"
)

type Documentation struct{}

func New() *Documentation {
	documentation := &Documentation{}

	return documentation
}

func (documentation *Documentation) Init() *Directory {
	template := dag.CurrentModule().Source().Directory("template")

	return template
}

type DocumentationBuilder struct {
	// +private
	Directory *Directory
	// +private
	Configuration DocumentationBuilderConfiguration
}

type DocumentationBuilderConfiguration struct {
	Hugo struct {
		Version string
	}
}

func (documentation *Documentation) Builder(
	ctx context.Context,
	directory *Directory,
) (*DocumentationBuilder, error) {
	const packageJsonFilename string = "package.json"

	builder := &DocumentationBuilder{
		Directory: directory,
	}

	packageJsonString, err := directory.File(packageJsonFilename).Contents(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to read %q file: %w", packageJsonFilename, err)
	}

	err = json.Unmarshal([]byte(packageJsonString), &builder.Configuration)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %q file: %w", packageJsonFilename, err)
	}

	if builder.Configuration.Hugo.Version == "" {
		return nil, fmt.Errorf("Hugo version is not set in %q file", packageJsonFilename)
	}

	return builder, nil
}

func (builder *DocumentationBuilder) Container() *Container {
	kroki := dag.Kroki()

	container := dag.Redhat().Minimal().Container().
		With(dag.Nodejs().Configuration).
		With(dag.Golang().Configuration).
		With(dag.Hugo(builder.Configuration.Hugo.Version).Configuration).
		WithMountedDirectory(".", builder.Directory).
		WithExec([]string{"npm clean-install"}).
		WithEntrypoint([]string{"npm", "run", "all", "--"}).
		WithoutDefaultArgs()

	// FIXME
	// Doesn’t work at the time.
	//WithServiceBinding("kroki", kroki.Server())
	_ = kroki

	return container
}

type DocumentationBuild struct {
	// +private
	Builder *Container
}

func (builder *DocumentationBuilder) Build(
	// +optional
	args []string,
) *DocumentationBuild {
	build := &DocumentationBuild{
		Builder: builder.Container().WithExec(args),
	}

	return build
}

func (build *DocumentationBuild) Directory() *Directory {
	directory := build.Builder.Directory("public")

	return directory
}

func (build *DocumentationBuild) Container() *Container {
	directory := build.Directory()
	container := dag.Caddy(directory).Container()

	return container
}

func (build *DocumentationBuild) Server() *Service {
	server := build.Container().AsService()

	return server
}
