// Hugo
//
// Get Hugo static site generator.

// Copyright Camptocamp SA
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"context"
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
	// +private
	Platform Platform
}

// Hugo constructor
func New(
	ctx context.Context,
	// Hugo version to get
	version string,
	// Hugo edition to get
	// +optional
	// +default=true
	extended bool,
	// Platform to get Hugo for
	// +optional
	platform Platform,
) (*Hugo, error) {
	if platform == "" {
		defaultPlatform, err := dag.DefaultPlatform(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed to get platform: %w", err)
		}

		platform = defaultPlatform
	}

	hugo := &Hugo{
		Version:  version,
		Extended: extended,
		Platform: platform,
	}

	return hugo, nil
}

// Get Hugo executable binary
func (hugo *Hugo) Binary() *File {
	platform := strings.Split(string(hugo.Platform), "/")

	os := platform[0]
	arch := platform[1]

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
		WithExec([]string{"grep -w " + tarballName + " " + checksumsName + " | sha256sum -c"}).
		WithExec([]string{"tar --extract --file " + tarballName})

	file := container.File("hugo")

	return file
}

// Get a root filesystem overlay with Hugo
func (hugo *Hugo) Overlay(
	// Filesystem prefix under which to install Hugo
	// +optional
	prefix string,
) *Directory {
	if prefix == "" {
		prefix = "/usr/local"
	}

	directory := dag.Directory().
		WithDirectory(prefix, dag.Directory().
			WithFile("bin/hugo", hugo.Binary()),
		)

	return directory
}

// Install Hugo in a container
func (hugo *Hugo) Installation(
	// Container in which to install Hugo
	container *Container,
) *Container {
	container = container.
		WithDirectory("/", hugo.Overlay("")).
		WithMountedCache(CacheDir, dag.CacheVolume("hugo")).
		WithEnvVariable("HUGO_CACHEDIR", CacheDir)

	return container
}

// Get a Hugo container from a base container
func (hugo *Hugo) Container(
	// Base container
	container *Container,
) *Container {
	container = container.
		With(hugo.Installation).
		WithEntrypoint([]string{"hugo"}).
		WithoutDefaultArgs().
		WithExposedPort(1313)

	return container
}

// Get a Red Hat Universal Base Image container with Hugo
func (hugo *Hugo) RedhatContainer() *Container {
	container := hugo.Container(
		dag.Redhat().Container().
			With(dag.Golang().RedhatInstallation),
	)

	return container
}

// Get a Red Hat Minimal Universal Base Image container with Hugo
func (hugo *Hugo) RedhatMinimalContainer() *Container {
	container := hugo.Container(
		dag.Redhat().Minimal().Container().
			With(dag.Golang().RedhatMinimalInstallation),
	)

	return container
}

// Get a Red Hat Micro Universal Base Image container with Hugo
//
// Hugo extended edition and Hugo modules cannot be used.
func (hugo *Hugo) RedhatMicroContainer() *Container {
	container := hugo.Container(dag.Redhat().Micro().Container())

	return container
}
