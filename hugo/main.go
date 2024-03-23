package main

import (
	"context"
	"fmt"
	"strings"
)

const (
	CacheDir string = "/var/cache/hugo"
)

type Hugo struct {
	Version     string
	SassVersion string
	Platform    Platform
}

func New(
	ctx context.Context,
	version string,
	// +optional
	sassVersion string,
	// +optional
	platform Platform,
) (*Hugo, error) {
	if platform == "" {
		defaultPlatform, err := dag.DefaultPlatform(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed to get platform: %w", err)
		}

		platform = defaultPlatform
	}

	hugo := &Hugo{
		Version:     version,
		SassVersion: sassVersion,
		Platform:    platform,
	}

	return hugo, nil
}

func (hugo *Hugo) File() *File {
	platform := strings.Split(string(hugo.Platform), "/")

	os := platform[0]
	arch := platform[1]

	downloadURL := "https://github.com/gohugoio/hugo/releases/download/v" + hugo.Version

	tarballName := fmt.Sprintf("hugo_extended_%s_%s-%s.tar.gz", hugo.Version, os, arch)
	checksumsName := fmt.Sprintf("hugo_%s_checksums.txt", hugo.Version)

	tarball := dag.HTTP(downloadURL + "/" + tarballName)
	checksums := dag.HTTP(downloadURL + "/" + checksumsName)

	container := dag.Busybox().Container().
		WithMountedFile(tarballName, tarball).
		WithMountedFile(checksumsName, checksums).
		WithExec([]string{"grep -w " + tarballName + " " + checksumsName + " | sha256sum -c"}).
		WithExec([]string{"tar --extract --file " + tarballName})

	file := container.File("hugo")

	return file
}

func (hugo *Hugo) Directory(
	// +optional
	prefix string,
) *Directory {
	if prefix == "" {
		prefix = "/usr/local"
	}

	directory := dag.Directory().
		WithDirectory(prefix, dag.Directory().
			WithFile("bin/hugo", hugo.File()),
		)

	if hugo.SassVersion != "" {
		sass := dag.Sass(hugo.SassVersion, SassOpts{
			Platform: string(hugo.Platform),
		})

		directory = directory.
			WithDirectory("/", sass.Directory(SassDirectoryOpts{Prefix: prefix}))
	}

	return directory
}

func (hugo *Hugo) Configuration(container *Container) *Container {
	container = container.
		With(dag.Golang().Configuration).
		WithDirectory("/", hugo.Directory("")).
		WithMountedCache(CacheDir, dag.CacheVolume("hugo")).
		WithEnvVariable("HUGO_CACHEDIR", CacheDir)

	return container
}

func (hugo *Hugo) Container() *Container {
	container := dag.Redhat().Container().
		With(hugo.Configuration).
		WithEntrypoint([]string{"hugo"}).
		WithoutDefaultArgs().
		WithExposedPort(1313)

	return container
}
