// Red Hat
//
// Get and customize containers based on Red Hat Universal Base Images.

// Copyright Camptocamp SA
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"strings"
)

const (
	// Red Hat Universal Base Image container registry
	ImageRegistry string = "registry.access.redhat.com"

	// Red Hat Universal Base Image container repository
	ImageRepository string = "ubi9"
	// Red Hat Universal Base Image container tag
	ImageTag string = "9.4-947.1716476138"
	// Red Hat Universal Base Image container digest
	ImageDigest string = "sha256:8f8fb7989ee757f27fae06a7d5041c903bc11ebf6ee7a1177c7e3070ce1d9e54"

	// Red Hat Minimal Universal Base Image container repository
	MinimalImageRepository string = "ubi9-minimal"
	// Red Hat Minimal Universal Base Image container tag
	MinimalImageTag string = "9.4-949.1716471857"
	// Red Hat Minimal Universal Base Image container digest
	MinimalImageDigest string = "sha256:adbac3083c2f340bee7cce4563665a1555901bee048bca6842b4fa0a1e6b875b"

	// Red Hat Micro Universal Base Image container repository
	MicroImageRepository string = "ubi9-micro"
	// Red Hat Micro Universal Base Image container tag
	MicroImageTag string = "9.4-6.1716471860"
	// Red Hat Micro Universal Base Image container digest
	MicroImageDigest string = "sha256:213fd2a0116a76eaa274fee20c86eef4dfba9f311784e8fb7d7f5fc38b32f3ef"
)

// Red Hat Universal Base Image
type Redhat struct{}

// Red Hat Universal Base Image constructor
func New() *Redhat {
	return &Redhat{}
}

// Get a Red Hat Universal Base Image container
func (*Redhat) Container() *Container {
	container := dag.Container().
		From(ImageRegistry + "/" + ImageRepository + ":" + ImageTag + "@" + ImageDigest).
		WithEntrypoint([]string{"sh", "-c"}).
		WithoutDefaultArgs().
		WithWorkdir("/home")

	return container
}

// Red Hat Universal Base Image module
type RedhatModule struct {
	// +private
	Name string
}

// Red Hat Universal Base Image module constructor
func (*Redhat) Module(
	// Module name
	name string,
) *RedhatModule {
	module := &RedhatModule{
		Name: name,
	}

	return module
}

// Enable a module in a Red Hat Universal Base Image container
func (module *RedhatModule) Enabled(
	// Container in which to enable the module
	container *Container,
) *Container {
	return container.WithExec([]string{"dnf module enable --assumeyes " + module.Name + " && dnf clean all"})
}

// Disable a module in a Red Hat Universal Base Image container
func (module *RedhatModule) Disabled(
	// Container in which to disable the module
	container *Container,
) *Container {
	return container.WithExec([]string{"dnf module disable --assumeyes " + module.Name + " && dnf clean all"})
}

// Red Hat Universal Base Image packages
type RedhatPackages struct {
	// +private
	Names []string
}

// Red Hat Universal Base Image packages constructor
func (*Redhat) Packages(
	// Packages name
	names []string,
) *RedhatPackages {
	packages := &RedhatPackages{
		Names: names,
	}

	return packages
}

// Install packages in a Red Hat Universal Base Image container
func (packages *RedhatPackages) Installed(
	// Container in which to install the packages
	container *Container,
) *Container {
	return container.WithExec([]string{"dnf install --nodocs --setopt install_weak_deps=0 --assumeyes " + strings.Join(packages.Names, " ") + " && dnf clean all"})
}

// Remove packages in a Red Hat Universal Base Image container
func (packages *RedhatPackages) Removed(
	// Container in which to remove the packages
	container *Container,
) *Container {
	return container.WithExec([]string{"dnf remove --assumeyes " + strings.Join(packages.Names, " ") + " && dnf clean all"})
}

// Get Red Hat Universal Base Image CA certificates
func (redhat *Redhat) CaCertificates() *Directory {
	const installroot string = "/tmp/rootfs"

	caCertificates := redhat.Container().
		WithExec([]string{"mkdir " + installroot + " && dnf --installroot " + installroot + " install --nodocs --setopt install_weak_deps=0 --assumeyes ca-certificates && dnf --installroot " + installroot + " clean all"}).
		Directory(installroot + "/etc/pki/ca-trust")

	return caCertificates
}

// Red Hat Minimal Universal Base Image
type RedhatMinimal struct{}

// Red Hat Minimal Universal Base Image constructor
func (*Redhat) Minimal() *RedhatMinimal {
	return &RedhatMinimal{}
}

// Get a Red Hat Minimal Universal Base Image container
func (*RedhatMinimal) Container() *Container {
	container := dag.Container().
		From(ImageRegistry + "/" + MinimalImageRepository + ":" + MinimalImageTag + "@" + MinimalImageDigest).
		WithEntrypoint([]string{"sh", "-c"}).
		WithoutDefaultArgs().
		WithWorkdir("/home")

	return container
}

// Red Hat Minimal Universal Base Image module
type RedhatMinimalModule struct {
	// +private
	Name string
}

// Red Hat Minimal Universal Base Image module constructor
func (*RedhatMinimal) Module(
	// Module name
	name string,
) *RedhatMinimalModule {
	module := &RedhatMinimalModule{
		Name: name,
	}

	return module
}

// Enable a module in a Red Hat Minimal Universal Base Image container
func (module *RedhatMinimalModule) Enabled(
	// Container in which to enable the module
	container *Container,
) *Container {
	return container.WithExec([]string{"microdnf module enable --assumeyes " + module.Name + " && microdnf clean all"})
}

// Disable a module in a Red Hat Minimal Universal Base Image container
func (module *RedhatMinimalModule) Disabled(
	// Container in which to disable the module
	container *Container,
) *Container {
	return container.WithExec([]string{"microdnf module disable --assumeyes " + module.Name + " && microdnf clean all"})
}

// Red Hat Minimal Universal Base Image packages
type RedhatMinimalPackages struct {
	// +private
	Names []string
}

// Red Hat Minimal Universal Base Image packages constructor
func (*RedhatMinimal) Packages(
	// Packages name
	names []string,
) *RedhatMinimalPackages {
	packages := &RedhatMinimalPackages{
		Names: names,
	}

	return packages
}

// Install packages in a Red Hat Minimal Universal Base Image container
func (packages *RedhatMinimalPackages) Installed(
	// Container in which to install the packages
	container *Container,
) *Container {
	return container.WithExec([]string{"microdnf install --nodocs --setopt install_weak_deps=0 --assumeyes " + strings.Join(packages.Names, " ") + " && microdnf clean all"})
}

// Remove packages in a Red Hat Minimal Universal Base Image container
func (packages *RedhatMinimalPackages) Removed(
	// Container in which to remove the packages
	container *Container,
) *Container {
	return container.WithExec([]string{"microdnf remove --assumeyes " + strings.Join(packages.Names, " ") + " && microdnf clean all"})
}

// Red Hat Micro Universal Base Image
type RedhatMicro struct{}

// Red Hat Micro Universal Base Image constructor
func (*Redhat) Micro() *RedhatMicro {
	return &RedhatMicro{}
}

// Get a Red Hat Micro Universal Base Image container
func (*RedhatMicro) Container() *Container {
	container := dag.Container().
		From(ImageRegistry + "/" + MicroImageRepository + ":" + MicroImageTag + "@" + MicroImageDigest).
		WithEntrypoint([]string{"sh", "-c"}).
		WithoutDefaultArgs().
		WithWorkdir("/home")

	return container
}
