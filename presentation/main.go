package main

import (
	"context"
)

const (
	BaseImageRegistry   string = "registry.access.redhat.com"
	BaseImageRepository string = "ubi9-minimal"
	BaseImageTag        string = "9.3-1552"
	BaseImageDigest     string = "sha256:582e18f13291d7c686ec4e6e92d20b24c62ae0fc72767c46f30a69b1a6198055"

	KrokiImageRegistry   string = "docker.io"
	KrokiImageRepository string = "yuzutech/kroki"
	KrokiImageTag        string = "0.24.1"
	KrokiImageDigest     string = "sha256:b43be03ec8a210471d4eaaf044b44e765551730447516c26f15c0f4b27628d45"

	CaddyImageRegistry   string = "docker.io"
	CaddyImageRepository string = "caddy"
	CaddyImageTag        string = "2.7.6"
	CaddyImageDigest     string = "sha256:d8d3637a26f50bf0bd27a6151d2bd4f7a9f0455936fe7ca2498abbc2e26c841e"

	nodeCacheDir string = "/var/cache/node"
)

type Presentation struct {
	// +private
	Directory *Directory
	// +private
	Npmrc *Secret
}

func New(
	ctx context.Context,
	directory *Directory,
	npmrc *Secret,
) (*Presentation, error) {
	m := &Presentation{
		Directory: directory,
		Npmrc:     npmrc,
	}

	return m, nil
}

func (m *Presentation) Init(
	ctx context.Context,
) *Directory {
	template := dag.CurrentModule().Source().Directory("template").
		WithoutFile(".gitignore")

	return template
}

func (m *Presentation) Builder(
	ctx context.Context,
) *Container {
	kroki := dag.Container().
		From(KrokiImageRegistry + "/" + KrokiImageRepository + ":" + KrokiImageTag + "@" + KrokiImageDigest).
		WithExposedPort(8000).
		AsService()

	builder := dag.Container().
		From(BaseImageRegistry+"/"+BaseImageRepository+":"+BaseImageTag+"@"+BaseImageDigest).
		WithEntrypoint([]string{"sh", "-c"}).
		WithExec([]string{"microdnf module enable nodejs:20 --assumeyes && microdnf install --nodocs --setopt install_weak_deps=0 --assumeyes npm && microdnf clean all"}).
		WithWorkdir("/home").
		WithMountedCache(nodeCacheDir, dag.CacheVolume("node")).
		WithEnvVariable("NPM_CONFIG_CACHE", nodeCacheDir+"/npm").
		WithMountedDirectory(".", m.Directory).
		WithMountedSecret("/etc/npmrc", m.Npmrc).
		WithExec([]string{"npm clean-install"}).
		WithEntrypoint([]string{"npm", "run", "all"}).
		WithoutDefaultArgs()

	// FIXME
	// Doesn’t work at the time.
	//WithServiceBinding("kroki", kroki)
	_ = kroki

	return builder
}

func (m *Presentation) Build(
	ctx context.Context,
) *Directory {
	build := m.Builder(ctx).
		WithExec(nil).
		Directory("dist")

	return build
}

func (m *Presentation) Server(
	ctx context.Context,
) *Service {
	caddyfile := dag.CurrentModule().Source().File("Caddyfile")
	build := m.Build(ctx)

	server := dag.Container().
		From(CaddyImageRegistry+"/"+CaddyImageRepository+":"+CaddyImageTag+"@"+CaddyImageDigest).
		WithFile("/etc/caddy/Caddyfile", caddyfile).
		WithMountedDirectory("/usr/share/caddy", build).
		WithExposedPort(8080).
		AsService()

	return server
}
