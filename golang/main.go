// Go
//
// Install Go in containers based on Red Hat Universal Base Images.

// Copyright Camptocamp SA
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

const (
	// Location of Go download and Go build caches
	CacheDir string = "/var/cache/go"
)

// Go
type Golang struct{}

// Go constructor
func New() *Golang {
	return &Golang{}
}

// Configure Go in a container
func (*Golang) Configuration(
	// Container in which to configure Go
	container *Container,
) *Container {
	container = container.
		WithMountedCache(CacheDir, dag.CacheVolume("golang")).
		WithEnvVariable("GOPATH", CacheDir).
		WithEnvVariable("GOCACHE", CacheDir+"/build")

	return container
}

// Install Go in a Red Hat Universal Base Image container from packages
func (golang *Golang) RedhatInstallation(
	// Container in which to install Go
	container *Container,
) *Container {
	container = container.
		With(dag.Redhat().Packages([]string{
			"go",
			"git",
		}).Installed).
		With(golang.Configuration)

	return container
}

// Get a Red Hat Universal Base Image container with Go
func (golang *Golang) RedhatContainer() *Container {
	return dag.Redhat().Container().With(golang.RedhatInstallation)
}

// Install Go in a Red Hat Minimal Universal Base Image container from packages
func (golang *Golang) RedhatMinimalInstallation(
	// Container in which to install Go
	container *Container,
) *Container {
	container = container.
		With(dag.Redhat().Minimal().Packages([]string{
			"go",
			"git",
		}).Installed).
		With(golang.Configuration)

	return container
}

// Get a Red Hat Minimal Universal Base Image container with Go
func (golang *Golang) RedhatMinimalContainer() *Container {
	return dag.Redhat().Minimal().Container().With(golang.RedhatMinimalInstallation)
}
