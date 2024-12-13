// Kroki
//
// Get a container or service running Kroki to create diagrams from textual descriptions.

// Copyright Camptocamp SA
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"dagger/kroki/internal/dagger"
)

const (
	// Kroki container image registry
	ImageRegistry string = "docker.io"
	// Kroki container image repository
	ImageRepository string = "yuzutech/kroki"
	// Kroki container image tag
	ImageTag string = "0.25.0"
	// Kroki container image digest
	ImageDigest string = "sha256:a9db3ab74543b84d641d5ff32272ffd4c6a21126ea6a529248bf276367c14273"
)

// Kroki
type Kroki struct{}

// Kroki constructor
func New() *Kroki {
	return &Kroki{}
}

// Get a Kroki container ready to create diagrams
//
// Container exposes port 8080.
func (*Kroki) Container(
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) *dagger.Container {
	container := dag.Container(dagger.ContainerOpts{Platform: platform}).
		From(ImageRegistry + "/" + ImageRepository + ":" + ImageTag + "@" + ImageDigest).
		WithExposedPort(8000)

	return container
}

// Get a Kroki service creating diagrams
//
// See `container()` for details.
func (kroki *Kroki) Server() *dagger.Service {
	return kroki.Container("").AsService(dagger.ContainerAsServiceOpts{UseEntrypoint: true})
}
