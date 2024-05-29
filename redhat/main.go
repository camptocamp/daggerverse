package main

import (
	"strings"
)

const (
	ImageRegistry string = "registry.access.redhat.com"

	RedhatImageRepository string = "ubi9"
	RedhatImageTag        string = "9.3-1610"
	RedhatImageDigest     string = "sha256:66233eebd72bb5baa25190d4f55e1dc3fff3a9b77186c1f91a0abdb274452072"

	RedhatMinimalImageRepository string = "ubi9-minimal"
	RedhatMinimalImageTag        string = "9.3-1612"
	RedhatMinimalImageDigest     string = "sha256:bc552efb4966aaa44b02532be3168ac1ff18e2af299d0fe89502a1d9fabafbc5"

	RedhatMicroImageRepository string = "ubi9-micro"
	RedhatMicroImageTag        string = "9.3-15"
	RedhatMicroImageDigest     string = "sha256:8e33df2832f039b4b1adc53efd783f9404449994b46ae321ee4a0bf4499d5c42"
)

type Redhat struct{}

func (*Redhat) Container() *Container {
	container := dag.Container().
		From(ImageRegistry + "/" + RedhatImageRepository + ":" + RedhatImageTag + "@" + RedhatImageDigest).
		WithEntrypoint([]string{"sh", "-c"}).
		WithoutDefaultArgs().
		WithWorkdir("/home")

	return container
}

type RedhatModule struct {
	Name string
}

func (*Redhat) Module(name string) *RedhatModule {
	module := &RedhatModule{
		Name: name,
	}

	return module
}

func (module *RedhatModule) Enabled(container *Container) *Container {
	return container.WithExec([]string{"dnf module enable --assumeyes " + module.Name + " && dnf clean all"})
}

func (module *RedhatModule) Disabled(container *Container) *Container {
	return container.WithExec([]string{"dnf module disable --assumeyes " + module.Name + " && dnf clean all"})
}

type RedhatPackages struct {
	Names []string
}

func (*Redhat) Packages(names []string) *RedhatPackages {
	packages := &RedhatPackages{
		Names: names,
	}

	return packages
}

func (packages *RedhatPackages) Installed(container *Container) *Container {
	return container.WithExec([]string{"dnf install --nodocs --setopt install_weak_deps=0 --assumeyes " + strings.Join(packages.Names, " ") + " && dnf clean all"})
}

func (packages *RedhatPackages) Removed(container *Container) *Container {
	return container.WithExec([]string{"dnf remove --assumeyes " + strings.Join(packages.Names, " ") + " && dnf clean all"})
}

func (redhat *Redhat) CaCertificates() *Directory {
	const installroot string = "/tmp/rootfs"

	caCertificates := redhat.Container().
		WithExec([]string{"mkdir " + installroot + " && dnf --installroot " + installroot + " install --nodocs --setopt install_weak_deps=0 --assumeyes ca-certificates && dnf --installroot " + installroot + " clean all"}).
		Directory(installroot + "/etc/pki/ca-trust")

	return caCertificates
}

type RedhatMinimal struct{}

func (*Redhat) Minimal() *RedhatMinimal {
	return &RedhatMinimal{}
}

func (*RedhatMinimal) Container() *Container {
	container := dag.Container().
		From(ImageRegistry + "/" + RedhatMinimalImageRepository + ":" + RedhatMinimalImageTag + "@" + RedhatMinimalImageDigest).
		WithEntrypoint([]string{"sh", "-c"}).
		WithoutDefaultArgs().
		WithWorkdir("/home")

	return container
}

type RedhatMinimalModule struct {
	Name string
}

func (*RedhatMinimal) Module(name string) *RedhatMinimalModule {
	module := &RedhatMinimalModule{
		Name: name,
	}

	return module
}

func (module *RedhatMinimalModule) Enabled(container *Container) *Container {
	return container.WithExec([]string{"microdnf module enable --assumeyes " + module.Name + " && microdnf clean all"})
}

func (module *RedhatMinimalModule) Disabled(container *Container) *Container {
	return container.WithExec([]string{"microdnf module disable --assumeyes " + module.Name + " && microdnf clean all"})
}

type RedhatMinimalPackages struct {
	Names []string
}

func (*RedhatMinimal) Packages(names []string) *RedhatMinimalPackages {
	packages := &RedhatMinimalPackages{
		Names: names,
	}

	return packages
}

func (packages *RedhatMinimalPackages) Installed(container *Container) *Container {
	return container.WithExec([]string{"microdnf install --nodocs --setopt install_weak_deps=0 --assumeyes " + strings.Join(packages.Names, " ") + " && microdnf clean all"})
}

func (packages *RedhatMinimalPackages) Removed(container *Container) *Container {
	return container.WithExec([]string{"microdnf remove --assumeyes " + strings.Join(packages.Names, " ") + " && microdnf clean all"})
}

type RedhatMicro struct{}

func (*Redhat) Micro() *RedhatMicro {
	return &RedhatMicro{}
}

func (*RedhatMicro) Container() *Container {
	container := dag.Container().
		From(ImageRegistry + "/" + RedhatMicroImageRepository + ":" + RedhatMicroImageTag + "@" + RedhatMicroImageDigest).
		WithEntrypoint([]string{"sh", "-c"}).
		WithoutDefaultArgs().
		WithWorkdir("/home")

	return container
}
