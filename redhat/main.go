package main

import (
	"strings"
)

const (
	ImageRegistry   string = "registry.access.redhat.com"
	ImageRepository string = "ubi9-minimal"
	ImageTag        string = "9.3-1552"
	ImageDigest     string = "sha256:582e18f13291d7c686ec4e6e92d20b24c62ae0fc72767c46f30a69b1a6198055"
)

type Redhat struct{}

func (redhat *Redhat) Container() *Container {
	container := dag.Container().
		From(ImageRegistry + "/" + ImageRepository + ":" + ImageTag + "@" + ImageDigest).
		WithEntrypoint([]string{"sh", "-c"}).
		WithoutDefaultArgs().
		WithWorkdir("/home")

	return container
}

type RedHatModule struct {
	Name string
}

func (redhat *Redhat) Module(name string) *RedHatModule {
	module := &RedHatModule{
		Name: name,
	}

	return module
}

func (module *RedHatModule) Enabled(container *Container) *Container {
	container = container.
		WithExec([]string{"microdnf module enable --assumeyes " + module.Name + " && microdnf clean all"})

	return container
}

type RedHatPackages struct {
	Names []string
}

func (redhat *Redhat) Packages(names []string) *RedHatPackages {
	packages := &RedHatPackages{
		Names: names,
	}

	return packages
}

func (packages *RedHatPackages) Installed(container *Container) *Container {
	container = container.
		WithExec([]string{"microdnf install --nodocs --setopt install_weak_deps=0 --assumeyes " + strings.Join(packages.Names, " ") + " && microdnf clean all"})

	return container
}
