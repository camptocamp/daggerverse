package main

import (
	"context"
	"fmt"
	"strings"
)

const (
	CacheDir string = "/var/cache/terraform/plugins"
)

type Terraform struct {
	Version  string
	Platform Platform
}

func New(
	ctx context.Context,
	version string,
	// +optional
	platform Platform,
) (*Terraform, error) {
	if platform == "" {
		defaultPlatform, err := dag.DefaultPlatform(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed to get platform: %w", err)
		}

		platform = defaultPlatform
	}

	terraform := &Terraform{
		Version:  version,
		Platform: platform,
	}

	return terraform, nil
}

func (terraform *Terraform) File() *File {
	platform := strings.Split(string(terraform.Platform), "/")

	os := platform[0]
	arch := platform[1]

	downloadURL := "https://releases.hashicorp.com/terraform/" + terraform.Version

	tarballName := fmt.Sprintf("terraform_%s_%s_%s.zip", terraform.Version, os, arch)
	checksumsName := fmt.Sprintf("terraform_%s_SHA256SUMS", terraform.Version)
	checksumsSignatureName := fmt.Sprintf("terraform_%s_SHA256SUMS.sig", terraform.Version)

	tarball := dag.HTTP(downloadURL + "/" + tarballName)
	checksums := dag.HTTP(downloadURL + "/" + checksumsName)
	checksumsSignature := dag.HTTP(downloadURL + "/" + checksumsSignatureName)

	container := dag.Redhat().Container().
		With(dag.Redhat().Packages([]string{"unzip"}).Installed).
		WithFile("pgp-key.txt", dag.CurrentModule().Source().File("pgp-key.txt")).
		WithExec([]string{"gpg --import pgp-key.txt"}).
		WithMountedFile(tarballName, tarball).
		WithMountedFile(checksumsName, checksums).
		WithMountedFile(checksumsSignatureName, checksumsSignature).
		WithExec([]string{"gpg --verify " + checksumsSignatureName + " " + checksumsName}).
		WithExec([]string{"grep -w " + tarballName + " " + checksumsName + " | sha256sum -c"}).
		WithExec([]string{"unzip " + tarballName})

	file := container.File("terraform")

	return file
}

func (terraform *Terraform) Directory(
	// +optional
	prefix string,
) *Directory {
	if prefix == "" {
		prefix = "/usr/local"
	}

	directory := dag.Directory().
		WithDirectory(prefix, dag.Directory().
			WithFile("bin/terraform", terraform.File()),
		)

	return directory
}

func (terraform *Terraform) Installation(container *Container) *Container {
	container = container.
		WithDirectory("/", terraform.Directory("")).
		WithMountedCache(CacheDir, dag.CacheVolume("terraform")).
		WithEnvVariable("TF_PLUGIN_CACHE_DIR", CacheDir)

	return container
}

func (terraform *Terraform) Container() *Container {
	container := dag.Redhat().Micro().Container().
		WithDirectory("/etc/pki/ca-trust", dag.Redhat().CaCertificates()).
		With(terraform.Installation).
		WithEntrypoint([]string{"terraform"}).
		WithoutDefaultArgs()

	return container
}
