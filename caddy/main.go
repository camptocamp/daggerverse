package main

const (
	ImageRegistry   string = "docker.io"
	ImageRepository string = "caddy"
	ImageTag        string = "2.7.6"
	ImageDigest     string = "sha256:d8d3637a26f50bf0bd27a6151d2bd4f7a9f0455936fe7ca2498abbc2e26c841e"
)

type Caddy struct {
	// +private
	Directory *Directory
}

func New(
	directory *Directory,
) *Caddy {
	caddy := &Caddy{
		Directory: directory,
	}

	return caddy
}

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

func (caddy *Caddy) Server() *Service {
	server := caddy.Container().
		AsService()

	return server
}
