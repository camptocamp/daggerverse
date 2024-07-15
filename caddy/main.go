// Caddy
//
// Get a container or service running Caddy HTTP server to serve static content without TLS.

// Copyright Camptocamp SA
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"dagger/caddy/internal/dagger"
)

const (
	// Caddy container image registry
	ImageRegistry string = "docker.io"
	// Caddy container image repository
	ImageRepository string = "caddy"
	// Caddy container image tag
	ImageTag string = "2.8.4"
	// Caddy container image digest
	ImageDigest string = "sha256:4718355ff1e2592290e49950f01fb1d4b75adb920a7695aedd94b6a4590a684b"
)

// Caddy
type Caddy struct {
	// +private
	Directory *dagger.Directory
}

// Caddy constructor
func New(
	// Directory containing static content to serve
	directory *dagger.Directory,
) *Caddy {
	caddy := &Caddy{
		Directory: directory,
	}

	return caddy
}

// Get a Caddy container ready to serve the static content
//
// Static content is mounted under `/usr/share/caddy` and container exposes port 8080.
func (caddy *Caddy) Container() *dagger.Container {
	caddyfile := dag.CurrentModule().Source().File("Caddyfile")

	container := dag.Container().
		From(ImageRegistry+"/"+ImageRepository+":"+ImageTag+"@"+ImageDigest).
		WithExec([]string{"chown", "65535:65535", "/config/caddy"}).
		WithExec([]string{"chown", "65535:65535", "/data/caddy"}).
		WithEntrypoint([]string{"caddy"}).
		WithDefaultArgs([]string{"run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"}).
		WithFile("/etc/caddy/Caddyfile", caddyfile).
		WithMountedDirectory("/usr/share/caddy", caddy.Directory).
		WithUser("65535").
		WithExposedPort(8080)

	return container
}

// Get a Caddy service serving the static content
//
// See `container()` for details.
func (caddy *Caddy) Server() *dagger.Service {
	return caddy.Container().WithExec(nil, dagger.ContainerWithExecOpts{UseEntrypoint: true}).AsService()
}
