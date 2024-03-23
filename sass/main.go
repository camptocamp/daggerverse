package main

import (
	"context"
	"fmt"
	"strings"
)

type Sass struct {
	Version  string
	Platform Platform
}

func New(
	ctx context.Context,
	version string,
	// +optional
	platform Platform,
) (*Sass, error) {
	if platform == "" {
		defaultPlatform, err := dag.DefaultPlatform(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed to get platform: %w", err)
		}

		platform = defaultPlatform
	}

	sass := &Sass{
		Version:  version,
		Platform: platform,
	}

	return sass, nil
}

func (sass *Sass) Tarball() *Directory {
	platform := strings.Split(string(sass.Platform), "/")

	os := platform[0]
	arch := map[string]string{
		"amd64": "x64",
		"386":   "ia32",
		"arm":   "arm",
		"arm64": "arm64",
	}[platform[1]]

	downloadURL := "https://github.com/sass/dart-sass/releases/download/" + sass.Version

	tarballName := fmt.Sprintf("dart-sass-%s-%s-%s.tar.gz", sass.Version, os, arch)

	tarball := dag.HTTP(downloadURL + "/" + tarballName)

	container := dag.Busybox().Container().
		WithMountedFile("sass.tar.gz", tarball).
		WithExec([]string{"tar --extract --strip-components 1 --file sass.tar.gz"})

	directory := dag.Directory().
		WithFile("dart", container.File("src/dart")).
		WithFile("sass.snapshot", container.File("src/sass.snapshot"))

	return directory
}

func (sass *Sass) Directory(
	// +optional
	prefix string,
) *Directory {
	if prefix == "" {
		prefix = "/usr/local"
	}

	directory := dag.Directory().
		WithDirectory(prefix, dag.Directory().
			WithDirectory("libexec/sass", sass.Tarball()).
			WithFile("bin/sass", dag.CurrentModule().Source().File("bin/sass"), DirectoryWithFileOpts{Permissions: 0o755}),
		)

	return directory
}

func (sass *Sass) Configuration(container *Container) *Container {
	container = container.
		WithDirectory("/", sass.Directory(""))

	return container
}

func (sass *Sass) Container() *Container {
	container := dag.Redhat().Container().
		With(sass.Configuration).
		WithEntrypoint([]string{"sass"}).
		WithoutDefaultArgs()

	return container
}
