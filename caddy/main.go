// Simply serve static content using Caddy HTTP server
//
// Get a container or service running Caddy HTTP server to serve static content without TLS.
//
// Copyright Camptocamp SA
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

const (
	// Caddy container image registry
	ImageRegistry string = "docker.io"
	// Caddy container image repository
	ImageRepository string = "caddy"
	// Caddy container image tag
	ImageTag string = "2.7.6"
	// Caddy container image digest
	ImageDigest string = "sha256:d8d3637a26f50bf0bd27a6151d2bd4f7a9f0455936fe7ca2498abbc2e26c841e"
)

// Caddy module
type Caddy struct {
	// Directory containing static content to be served
	// +private
	Directory *Directory
}

// Caddy module constructor
func New(
	// Directory containing static content to be served
	directory *Directory,
) *Caddy {
	caddy := &Caddy{
		Directory: directory,
	}

	return caddy
}

// Get a Caddy container ready to serve the static content
//
// Static content is mounted under `/usr/share/caddy` and container exposes port 8080.
func (caddy *Caddy) Container() *Container {
	caddyfile := dag.CurrentModule().Source().File("Caddyfile")

	container := dag.Container().
		From(ImageRegistry+"/"+ImageRepository+":"+ImageTag+"@"+ImageDigest).
		WithEntrypoint([]string{"sh", "-c"}).
		WithExec([]string{"chown 65535:65535 /config/caddy"}).
		WithExec([]string{"chown 65535:65535 /data/caddy"}).
		WithEntrypoint([]string{"caddy"}).
		WithDefaultArgs([]string{"run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"}).
		WithFile("/etc/caddy/Caddyfile", caddyfile).
		WithMountedDirectory("/usr/share/caddy", caddy.Directory).
		WithUser("65535").
		WithExposedPort(8080).
		WithExec(nil)

	return container
}

// Get a Caddy container serving the static content
func (caddy *Caddy) Server() *Service {
	return caddy.Container().AsService()
}
