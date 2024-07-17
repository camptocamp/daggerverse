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
}

// Sass constructor
func New(
	// Sass version to get
	version string,
) (*Sass, error) {
	sass := &Sass{
		Version: version,
	}

	return sass, nil
}

// Get Sass binaries (Dart runtime and Sass snapshot)
func (sass *Sass) Binaries(
	ctx context.Context,
	// Platform to get Sass for
	// +optional
	platform dagger.Platform,
) (*dagger.Directory, error) {
	if platform == "" {
		defaultPlatform, err := dag.DefaultPlatform(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed to get platform: %s", err)
		}

		platform = defaultPlatform
	}

	platformElements := strings.Split(string(platform), "/")

	os := platformElements[0]
	arch := map[string]string{
		"amd64": "x64",
		"386":   "ia32",
		"arm":   "arm",
		"arm64": "arm64",
	}[platformElements[1]]

	downloadURL := "https://github.com/sass/dart-sass/releases/download/" + sass.Version

	tarballName := fmt.Sprintf("dart-sass-%s-%s-%s.tar.gz", sass.Version, os, arch)

	tarball := dag.HTTP(downloadURL + "/" + tarballName)

	container := dag.Redhat().Container().
		WithMountedFile("sass.tar.gz", tarball).
		WithExec([]string{"tar", "--extract", "--strip-components", "1", "--file", "sass.tar.gz"})

	binaries := dag.Directory().
		WithFile("dart", container.File("src/dart")).
		WithFile("sass.snapshot", container.File("src/sass.snapshot"))

	return binaries, nil
}

// Get a root filesystem overlay with Sass
func (sass *Sass) Overlay(
	ctx context.Context,
	// Platform to get Hugo for
	// +optional
	platform dagger.Platform,
	// Filesystem prefix under which to install Sass
	// +optional
	prefix string,
) (*dagger.Directory, error) {
	if prefix == "" {
		prefix = "/usr/local"
	}

	binaries, err := sass.Binaries(ctx, platform)

	if err != nil {
		return nil, fmt.Errorf("failed to get Sass binaries: %s", err)
	}

	overlay := dag.Directory().
		WithDirectory(prefix, dag.Directory().
			WithDirectory("libexec/sass", binaries).
			WithFile("bin/sass", dag.CurrentModule().Source().File("bin/sass"), dagger.DirectoryWithFileOpts{Permissions: 0o755}),
		)

	return overlay, nil
}

// Install Sass in a container
func (sass *Sass) Installation(
	ctx context.Context,
	// Container in which to install Sass
	container *dagger.Container,
) (*dagger.Container, error) {
	platform, err := container.Platform(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get container platform: %s", err)
	}

	overlay, err := sass.Overlay(ctx, platform, "")

	if err != nil {
		return nil, fmt.Errorf("failed to get Sass overlay: %s", err)
	}

	container = container.
		WithDirectory("/", overlay)

	return container, nil
}

// Get a Sass container from a base container
func (sass *Sass) Container(
	ctx context.Context,
	// Base container
	container *dagger.Container,
) (*dagger.Container, error) {
	container, err := sass.Installation(ctx, container)

	if err != nil {
		return nil, fmt.Errorf("failed to install Sass: %s", err)
	}

	container = container.
		WithEntrypoint([]string{"sass"}).
		WithoutDefaultArgs()

	return container, nil
}

// Get a Red Hat Universal Base Image container with Sass
func (sass *Sass) RedhatContainer(
	ctx context.Context,
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) (*dagger.Container, error) {
	container := dag.Redhat().Container(dagger.RedhatContainerOpts{Platform: platform})

	return sass.Container(ctx, container)
}

// Get a Red Hat Minimal Universal Base Image container with Sass
func (sass *Sass) RedhatMinimalContainer(
	ctx context.Context,
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) (*dagger.Container, error) {
	container := dag.Redhat().Minimal().Container(dagger.RedhatMinimalContainerOpts{Platform: platform})

	return sass.Container(ctx, container)
}

// Get a Red Hat Micro Universal Base Image container with Sass
func (sass *Sass) RedhatMicroContainer(
	ctx context.Context,
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) (*dagger.Container, error) {
	container := dag.Redhat().Micro().Container(dagger.RedhatMicroContainerOpts{Platform: platform})

	return sass.Container(ctx, container)
}
