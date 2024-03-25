package main

import (
	"context"
	"encoding/json"
	"fmt"
)

type Presentation struct{}

func New() *Presentation {
	presentation := &Presentation{}

	return presentation
}

func (presentation *Presentation) Init() *Directory {
	template := dag.CurrentModule().Source().Directory("template")

	return template
}

type PresentationBuilder struct {
	// +private
	Directory *Directory
	// +private
	Npmrc *Secret
	// +private
	Configuration PresentationBuilderConfiguration
}

type PresentationBuilderConfiguration struct {
}

func (presentation *Presentation) Builder(
	ctx context.Context,
	directory *Directory,
	npmrc *Secret,
) (*PresentationBuilder, error) {
	const packageJsonFilename string = "package.json"

	builder := &PresentationBuilder{
		Directory: directory,
		Npmrc:     npmrc,
	}

	packageJsonString, err := directory.File(packageJsonFilename).Contents(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to read %q file: %w", packageJsonFilename, err)
	}

	err = json.Unmarshal([]byte(packageJsonString), &builder.Configuration)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %q file: %w", packageJsonFilename, err)
	}

	return builder, nil
}

func (builder *PresentationBuilder) Container() *Container {
	kroki := dag.Kroki()

	container := dag.Redhat().Minimal().Container().
		With(dag.Nodejs(NodejsOpts{
			Npmrc: builder.Npmrc,
		}).Configuration).
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

type PresentationBuild struct {
	// +private
	Builder *Container
}

func (builder *PresentationBuilder) Build() *PresentationBuild {
	build := &PresentationBuild{
		Builder: builder.Container().WithExec(nil),
	}

	return build
}

func (build *PresentationBuild) Directory() *Directory {
	directory := build.Builder.Directory("dist")

	return directory
}

func (build *PresentationBuild) Container() *Container {
	directory := build.Directory()
	container := dag.Caddy(directory).Container()

	return container
}

func (build *PresentationBuild) Server() *Service {
	server := build.Container().AsService()

	return server
}
