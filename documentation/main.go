// Camptocamp branded documentation
//
// Build a statically generated documentation directory, container or service using Hugo and Camptocamp branded Docsy theme.

// Copyright Camptocamp SA
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"context"
	"encoding/json"
	"fmt"
)

// Documentation
type Documentation struct{}

// Get a directory containing a newly initialized documentation
func (*Documentation) Init() *Directory {
	return dag.CurrentModule().Source().Directory("template")
}

// Documentation builder
type DocumentationBuilder struct {
	// Get a container ready to build the documentation
	*Container
}

// Documentation builder constructor
func (*Documentation) Builder(
	ctx context.Context,
	// Directory containing documentation to build
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
		With(dag.Nodejs().RedhatMinimalInstallation).
		With(dag.Golang().RedhatMinimalInstallation).
		With(dag.Hugo(configuration.Hugo.Version).Installation).
		WithServiceBinding("kroki", kroki.Server()).
		WithMountedDirectory(".", directory).
		WithExec([]string{"npm clean-install"}).
		WithEntrypoint([]string{"npm", "run", "build", "--"}).
		WithoutDefaultArgs()

	return builder, nil
}

// Documentation build result
type DocumentationBuildResult struct {
	// Get a directory containing the documentation build result
	*Directory
}

// Build the documentation
func (builder *DocumentationBuilder) Build(
	// Documentation builder arguments (Hugo arguments)
	// +optional
	args []string,
) *DocumentationBuildResult {
	build := &DocumentationBuildResult{
		Directory: builder.WithExec(args).Directory("public"),
	}

	return build
}

// Get a container ready to serve the documentation
//
// Container exposes port 8080.
func (build *DocumentationBuildResult) Container() *Container {
	return dag.Caddy(build.Directory).Container()
}

// Get a service serving the documentation
//
// See `container()` for details.
func (build *DocumentationBuildResult) Server() *Service {
	return build.Container().AsService()
}
