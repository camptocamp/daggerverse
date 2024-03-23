package main

const (
	CacheDir string = "/var/cache/go"
)

type Golang struct{}

func New() *Golang {
	golang := &Golang{}

	return golang
}

func (golang *Golang) Configuration(container *Container) *Container {
	container = container.
		With(dag.Redhat().Packages([]string{
			"go",
			"git",
		}).Installed).
		WithMountedCache(CacheDir, dag.CacheVolume("golang")).
		WithEnvVariable("GOPATH", CacheDir).
		WithEnvVariable("GOCACHE", CacheDir+"/build")

	return container
}

func (golang *Golang) Container() *Container {
	container := dag.Redhat().Container().
		With(golang.Configuration)

	return container
}
