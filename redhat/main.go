package main

import (
	"strings"
)

const (
	ImageRegistry string = "registry.access.redhat.com"

	MicroImageRepository string = "ubi9-micro"
	MicroImageTag        string = "9.3-15"
	MicroImageDigest     string = "sha256:8e33df2832f039b4b1adc53efd783f9404449994b46ae321ee4a0bf4499d5c42"

	MinimalImageRepository string = "ubi9-minimal"
	MinimalImageTag        string = "9.3-1552"
	MinimalImageDigest     string = "sha256:582e18f13291d7c686ec4e6e92d20b24c62ae0fc72767c46f30a69b1a6198055"
)

type Redhat struct{}

type RedhatMicro struct{}

func (*Redhat) Micro() *RedhatMicro {
	redhatMicro := &RedhatMicro{}

	return redhatMicro
}

func (*RedhatMicro) Container() *Container {
	container := dag.Container().
		From(ImageRegistry + "/" + MicroImageRepository + ":" + MicroImageTag + "@" + MicroImageDigest).
		WithEntrypoint([]string{"sh", "-c"}).
		WithoutDefaultArgs().
		WithWorkdir("/home")

	return container
}

type RedhatMinimal struct{}

func (*Redhat) Minimal() *RedhatMinimal {
	redhatMinimal := &RedhatMinimal{}

	return redhatMinimal
}

func (*RedhatMinimal) Container() *Container {
	container := dag.Container().
		From(ImageRegistry + "/" + MinimalImageRepository + ":" + MinimalImageTag + "@" + MinimalImageDigest).
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
