// Hugo
//
// Get Hugo static site generator.

// Copyright Camptocamp SA
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"context"
	"dagger/hugo/internal/dagger"
	"errors"
	"fmt"
	"strings"
)

const (
	// Location of Hugo cache
	CacheDir string = "/var/cache/hugo"
)

// Hugo
type Hugo struct {
	// +private
	Version string
	// +private
	Extended bool
}

// Hugo constructor
func New(
	// Hugo version to get
	version string,
	// Hugo edition to get
	// +optional
	// +default=true
	extended bool,
) (*Hugo, error) {
	hugo := &Hugo{
		Version:  version,
		Extended: extended,
	}

	return hugo, nil
}

// Get Hugo executable binary
func (hugo *Hugo) Binary(
	ctx context.Context,
	// Platform to get Hugo for
	// +optional
	platform dagger.Platform,
) (*dagger.File, error) {
	if platform == "" {
		defaultPlatform, err := dag.DefaultPlatform(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed to get platform: %s", err)
		}

		platform = defaultPlatform
	}

	platformElements := strings.Split(string(platform), "/")

	os := platformElements[0]
	arch := platformElements[1]

	downloadURL := "https://github.com/gohugoio/hugo/releases/download/v" + hugo.Version

	tarballBaseName := "hugo"

	if hugo.Extended {
		tarballBaseName += "_extended"
	}

	tarballName := fmt.Sprintf("%s_%s_%s-%s.tar.gz", tarballBaseName, hugo.Version, os, arch)
	checksumsName := fmt.Sprintf("hugo_%s_checksums.txt", hugo.Version)

	tarball := dag.HTTP(downloadURL + "/" + tarballName)
	checksums := dag.HTTP(downloadURL + "/" + checksumsName)

	container := dag.Redhat().Container().
		WithMountedFile(tarballName, tarball).
		WithMountedFile(checksumsName, checksums).
		WithExec([]string{"sh", "-c", "grep -w " + tarballName + " " + checksumsName + " | sha256sum -c"}).
		WithExec([]string{"tar", "--extract", "--file", tarballName})

	binary := container.File("hugo")

	return binary, nil
}

// Get a root filesystem overlay with Hugo
func (hugo *Hugo) Overlay(
	ctx context.Context,
	// Platform to get Hugo for
	// +optional
	platform dagger.Platform,
	// Filesystem prefix under which to install Hugo
	// +optional
	prefix string,
) (*dagger.Directory, error) {
	if prefix == "" {
		prefix = "/usr/local"
	}

	binary, err := hugo.Binary(ctx, platform)

	if err != nil {
		return nil, fmt.Errorf("failed to get Hugo binary: %s", err)
	}

	overlay := dag.Directory().
		WithDirectory(prefix, dag.Directory().
			WithFile("bin/hugo", binary),
		)

	return overlay, nil
}

// Install Hugo in a container
func (hugo *Hugo) Installation(
	ctx context.Context,
	// Container in which to install Hugo
	container *dagger.Container,
) (*dagger.Container, error) {
	platform, err := container.Platform(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get container platform: %s", err)
	}

	overlay, err := hugo.Overlay(ctx, platform, "")

	if err != nil {
		return nil, fmt.Errorf("failed to get Hugo overlay: %s", err)
	}

	container = container.
		WithDirectory("/", overlay).
		WithMountedCache(CacheDir, dag.CacheVolume("hugo")).
		WithEnvVariable("HUGO_CACHEDIR", CacheDir)

	return container, nil
}

// Get a Hugo container from a base container
func (hugo *Hugo) Container(
	ctx context.Context,
	// Base container
	container *dagger.Container,
) (*dagger.Container, error) {
	container, err := hugo.Installation(ctx, container)

	if err != nil {
		return nil, fmt.Errorf("failed to install Hugo: %s", err)
	}

	container = container.
		WithEntrypoint([]string{"hugo"}).
		WithoutDefaultArgs().
		WithExposedPort(1313)

	return container, nil
}

// Get a Red Hat Universal Base Image container with Hugo
func (hugo *Hugo) RedhatContainer(
	ctx context.Context,
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) (*dagger.Container, error) {
	container := dag.Redhat().Container(dagger.RedhatContainerOpts{Platform: platform}).
		With(dag.Golang().RedhatInstallation)

	return hugo.Container(ctx, container)
}

// Get a Red Hat Minimal Universal Base Image container with Hugo
func (hugo *Hugo) RedhatMinimalContainer(
	ctx context.Context,
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) (*dagger.Container, error) {
	container := dag.Redhat().Minimal().Container(dagger.RedhatMinimalContainerOpts{Platform: platform}).
		With(dag.Golang().RedhatMinimalInstallation)

	return hugo.Container(ctx, container)
}

// Get a Red Hat Micro Universal Base Image container with Hugo
//
// Hugo extended edition and Hugo modules cannot be used.
func (hugo *Hugo) RedhatMicroContainer(
	ctx context.Context,
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) (*dagger.Container, error) {
	if hugo.Extended {
		return nil, errors.New("extended version is not compatible with Red Hat micro container")
	}

	container := dag.Redhat().Micro().Container(dagger.RedhatMicroContainerOpts{Platform: platform})

	return hugo.Container(ctx, container)
}
