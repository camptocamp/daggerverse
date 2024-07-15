// Sass
//
// Get Sass CSS preprocessor.

// Copyright Camptocamp SA
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"context"
	"dagger/sass/internal/dagger"
	"fmt"
	"strings"
)

// Sass
type Sass struct {
	// +private
	Version string
	// +private
	Platform dagger.Platform
}

// Sass constructor
func New(
	ctx context.Context,
	// Sass version to get
	version string,
	// Platform to get Sass for
	// +optional
	platform dagger.Platform,
) (*Sass, error) {
	if platform == "" {
		defaultPlatform, err := dag.DefaultPlatform(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed to get platform: %w", err)
		}

		platform = defaultPlatform
	}

	sass := &Sass{
		Version:  version,
		Platform: platform,
	}

	return sass, nil
}

// Get Sass binaries (Dart runtime and Sass snapshot)
func (sass *Sass) Binaries() *dagger.Directory {
	platform := strings.Split(string(sass.Platform), "/")

	os := platform[0]
	arch := map[string]string{
		"amd64": "x64",
		"386":   "ia32",
		"arm":   "arm",
		"arm64": "arm64",
	}[platform[1]]

	downloadURL := "https://github.com/sass/dart-sass/releases/download/" + sass.Version

	tarballName := fmt.Sprintf("dart-sass-%s-%s-%s.tar.gz", sass.Version, os, arch)

	tarball := dag.HTTP(downloadURL + "/" + tarballName)

	container := dag.Redhat().Container().
		WithMountedFile("sass.tar.gz", tarball).
		WithExec([]string{"tar", "--extract", "--strip-components", "1", "--file", "sass.tar.gz"})

	directory := dag.Directory().
		WithFile("dart", container.File("src/dart")).
		WithFile("sass.snapshot", container.File("src/sass.snapshot"))

	return directory
}

// Get a root filesystem overlay with Sass
func (sass *Sass) Overlay(
	// Filesystem prefix under which to install Sass
	// +optional
	prefix string,
) *dagger.Directory {
	if prefix == "" {
		prefix = "/usr/local"
	}

	directory := dag.Directory().
		WithDirectory(prefix, dag.Directory().
			WithDirectory("libexec/sass", sass.Binaries()).
			WithFile("bin/sass", dag.CurrentModule().Source().File("bin/sass"), dagger.DirectoryWithFileOpts{Permissions: 0o755}),
		)

	return directory
}

// Install Sass in a container
func (sass *Sass) Installation(
	// Container in which to install Sass
	container *dagger.Container,
) *dagger.Container {
	container = container.
		WithDirectory("/", sass.Overlay(""))

	return container
}

// Get a Sass container from a base container
func (sass *Sass) Container(
	// Base container
	container *dagger.Container,
) *dagger.Container {
	container = container.
		With(sass.Installation).
		WithEntrypoint([]string{"sass"}).
		WithoutDefaultArgs()

	return container
}

// Get a Red Hat Universal Base Image container with Sass
func (sass *Sass) RedhatContainer() *dagger.Container {
	container := sass.Container(dag.Redhat().Container())

	return container
}

// Get a Red Hat Minimal Universal Base Image container with Sass
func (sass *Sass) RedhatMinimalContainer() *dagger.Container {
	container := sass.Container(dag.Redhat().Minimal().Container())

	return container
}

// Get a Red Hat Micro Universal Base Image container with Sass
func (sass *Sass) RedhatMicroContainer() *dagger.Container {
	container := sass.Container(dag.Redhat().Micro().Container())

	return container
}
