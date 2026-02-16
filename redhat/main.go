// Red Hat
//
// Get and customize containers based on Red Hat Universal Base Images.

// Copyright Camptocamp SA
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"dagger/redhat/internal/dagger"
	"strings"
)

const (
	// Red Hat Universal Base Image container registry
	ImageRegistry string = "registry.access.redhat.com"

	// Red Hat Universal Base Image container repository
	ImageRepository string = "ubi10"
	// Red Hat Universal Base Image container tag
	ImageTag string = "10.1-1770180700"
	// Red Hat Universal Base Image container digest
	ImageDigest string = "sha256:b9e5730d0b6dba45e82c15fb8f49c6082e01cdcb5e4f6ba96535dab42a4d2cf0"

	// Red Hat Minimal Universal Base Image container repository
	MinimalImageRepository string = "ubi10-minimal"
	// Red Hat Minimal Universal Base Image container tag
	MinimalImageTag string = "10.1-1770180557"
	// Red Hat Minimal Universal Base Image container digest
	MinimalImageDigest string = "sha256:a74a7a92d3069bfac09c6882087771fc7db59fa9d8e16f14f4e012fe7288554c"

	// Red Hat Micro Universal Base Image container repository
	MicroImageRepository string = "ubi10-micro"
	// Red Hat Micro Universal Base Image container tag
	MicroImageTag string = "10.1-1769518576"
	// Red Hat Micro Universal Base Image container digest
	MicroImageDigest string = "sha256:551f8ee81be3dbabd45a9c197f3724b9724c1edb05d68d10bfe85a5c9e46a458"
)

// Red Hat Universal Base Image
type Redhat struct{}

// Red Hat Universal Base Image constructor
func New() *Redhat {
	return &Redhat{}
}

// Get a Red Hat Universal Base Image container
func (*Redhat) Container(
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) *dagger.Container {
	container := dag.Container(dagger.ContainerOpts{Platform: platform}).
		From(ImageRegistry + "/" + ImageRepository + ":" + ImageTag + "@" + ImageDigest).
		WithWorkdir("/home")

	return container
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
	container *dagger.Container,
) *dagger.Container {
	return container.WithExec([]string{"sh", "-c", "dnf install --nodocs --setopt install_weak_deps=0 --assumeyes " + strings.Join(packages.Names, " ") + " && dnf clean all"})
}

// Remove packages in a Red Hat Universal Base Image container
func (packages *RedhatPackages) Removed(
	// Container in which to remove the packages
	container *dagger.Container,
) *dagger.Container {
	return container.WithExec([]string{"sh", "-c", "dnf remove --assumeyes " + strings.Join(packages.Names, " ") + " && dnf clean all"})
}

// Get Red Hat Universal Base Image CA certificates
func (redhat *Redhat) CaCertificates() *dagger.Directory {
	const installroot string = "/tmp/rootfs"

	caCertificates := redhat.Container("").
		WithExec([]string{"sh", "-c", "mkdir " + installroot + " && dnf --installroot " + installroot + " install --nodocs --setopt install_weak_deps=0 --assumeyes ca-certificates && dnf --installroot " + installroot + " clean all"}).
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
func (*RedhatMinimal) Container(
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) *dagger.Container {
	container := dag.Container(dagger.ContainerOpts{Platform: platform}).
		From(ImageRegistry + "/" + MinimalImageRepository + ":" + MinimalImageTag + "@" + MinimalImageDigest).
		WithWorkdir("/home")

	return container
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
	container *dagger.Container,
) *dagger.Container {
	return container.WithExec([]string{"sh", "-c", "microdnf install --nodocs --setopt install_weak_deps=0 --assumeyes " + strings.Join(packages.Names, " ") + " && microdnf clean all"})
}

// Remove packages in a Red Hat Minimal Universal Base Image container
func (packages *RedhatMinimalPackages) Removed(
	// Container in which to remove the packages
	container *dagger.Container,
) *dagger.Container {
	return container.WithExec([]string{"sh", "-c", "microdnf remove --assumeyes " + strings.Join(packages.Names, " ") + " && microdnf clean all"})
}

// Red Hat Micro Universal Base Image
type RedhatMicro struct{}

// Red Hat Micro Universal Base Image constructor
func (*Redhat) Micro() *RedhatMicro {
	return &RedhatMicro{}
}

// Get a Red Hat Micro Universal Base Image container
func (*RedhatMicro) Container(
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) *dagger.Container {
	container := dag.Container(dagger.ContainerOpts{Platform: platform}).
		From(ImageRegistry + "/" + MicroImageRepository + ":" + MicroImageTag + "@" + MicroImageDigest).
		WithWorkdir("/home")

	return container
}
