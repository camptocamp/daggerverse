package main

import (
	"context"
	"encoding/json"
	"fmt"
)

type Presentation struct{}

func (*Presentation) Init() *Directory {
	return dag.CurrentModule().Source().Directory("template")
}

type PresentationBuilder struct {
	*Container
}

func (*Presentation) Builder(
	ctx context.Context,
	directory *Directory,
	npmrc *Secret,
) (*PresentationBuilder, error) {
	const packageJsonFilename string = "package.json"

	packageJsonString, err := directory.File(packageJsonFilename).Contents(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to read %q file: %w", packageJsonFilename, err)
	}

	var configuration struct{}

	err = json.Unmarshal([]byte(packageJsonString), &configuration)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %q file: %w", packageJsonFilename, err)
	}

	builder := &PresentationBuilder{}

	kroki := dag.Kroki()

	builder.Container = dag.Redhat().Minimal().Container().
		With(dag.Nodejs(NodejsOpts{
			Npmrc: npmrc,
		}).Installation).
		WithMountedDirectory(".", directory).
		WithExec([]string{"npm clean-install"}).
		WithEntrypoint([]string{"npm", "run", "all"}).
		WithoutDefaultArgs()

	// FIXME
	// Doesn’t work at the time.
	//WithServiceBinding("kroki", kroki.Server())
	_ = kroki

	return builder, nil
}

type PresentationBuild struct {
	*Directory
}

func (builder *PresentationBuilder) Build() *PresentationBuild {
	build := &PresentationBuild{
		Directory: builder.WithExec(nil).Directory("dist"),
	}

	return build
}

func (build *PresentationBuild) Container() *Container {
	return dag.Caddy(build.Directory).Container()
}

func (build *PresentationBuild) Server() *Service {
	return build.Container().AsService()
}
