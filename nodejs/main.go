// Node.js
//
// Install Node.js in containers based on Red Hat Universal Base Images.

// Copyright Camptocamp SA
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

const (
	// Location of npm cache
	CacheDir string = "/var/cache/node"
)

// Node.js
type Nodejs struct {
	// +private
	Npmrc *Secret
}

// Node.js constructor
func New(
	// npm configuration file (can be used to pass registry credentials)
	// +optional
	npmrc *Secret,
) *Nodejs {
	nodejs := &Nodejs{
		Npmrc: npmrc,
	}

	return nodejs
}

// Configure Node.js in a container
func (nodejs *Nodejs) Configuration(
	// Container in which to configure Node.js
	container *Container,
) *Container {
	container = container.
		WithMountedCache(CacheDir, dag.CacheVolume("nodejs")).
		WithEnvVariable("NPM_CONFIG_CACHE", CacheDir+"/npm")

	if nodejs.Npmrc != nil {
		container = container.
			WithMountedSecret("/root/.npmrc", nodejs.Npmrc)
	}

	return container
}

// Install Node.js in a Red Hat Universal Base Image container from packages
func (nodejs *Nodejs) RedhatInstallation(
	// Container in which to install Node.js
	container *Container,
) *Container {
	container = container.
		With(dag.Redhat().Module("nodejs:20").Enabled).
		With(dag.Redhat().Packages([]string{
			"npm",
		}).Installed).
		With(nodejs.Configuration)

	return container
}

// Get a Red Hat Universal Base Image container with Node.js
func (nodejs *Nodejs) RedhatContainer() *Container {
	return dag.Redhat().Container().With(nodejs.RedhatInstallation)
}

// Install Node.js in a Red Hat Minimal Universal Base Image container from packages
func (nodejs *Nodejs) RedhatMinimalInstallation(
	// Container in which to install Node.js
	container *Container,
) *Container {
	container = container.
		With(dag.Redhat().Minimal().Module("nodejs:20").Enabled).
		With(dag.Redhat().Minimal().Packages([]string{
			"npm",
		}).Installed).
		With(nodejs.Configuration)

	return container
}

// Get a Red Hat Minimal Universal Base Image container with Node.js
func (nodejs *Nodejs) RedhatMinimalContainer() *Container {
	return dag.Redhat().Minimal().Container().With(nodejs.RedhatMinimalInstallation)
}
