package main

import (
	"context"
	"fmt"
	"strings"
)

const (
	HugoVersionFilename string = ".hugo-version"
)

type Documentation struct{}

func New() *Documentation {
	documentation := &Documentation{}

	return documentation
}

func (documentation *Documentation) Init() *Directory {
	template := dag.CurrentModule().Source().Directory("template").
		WithoutFile(".gitignore")

	return template
}

type DocumentationBuilder struct {
	// +private
	Directory *Directory
	// +private
	HugoVersion string
}

func (documentation *Documentation) Builder(
	ctx context.Context,
	directory *Directory,
) (*DocumentationBuilder, error) {
	hugoVersion, err := directory.File(HugoVersionFilename).Contents(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to read Hugo version from %q: %w", HugoVersionFilename, err)
	}

	builder := &DocumentationBuilder{
		Directory:   directory,
		HugoVersion: strings.TrimSpace(hugoVersion),
	}

	return builder, nil
}

func (builder *DocumentationBuilder) Container() *Container {
	kroki := dag.Kroki()

	container := dag.Redhat().Container().
		With(dag.Nodejs().Configuration).
		With(dag.Hugo(builder.HugoVersion).Configuration).
		WithMountedDirectory(".", builder.Directory).
		WithExec([]string{"npm clean-install"}).
		WithEntrypoint([]string{"npm", "run", "all"}).
		WithoutDefaultArgs()

	// FIXME
	// Doesn’t work at the time.
	//WithServiceBinding("kroki", kroki.Server())
	_ = kroki

	return container
}

type DocumentationBuild struct {
	Builder *Container
}

func (builder *DocumentationBuilder) Build(
	// +optional
	args []string,
) *DocumentationBuild {
	if args == nil {
		args = []string{
			"--cleanDestinationDir",
			"--minify",
		}
	}

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
