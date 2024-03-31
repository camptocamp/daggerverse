package main

const (
	CacheDir string = "/var/cache/go"
)

type Golang struct{}

func (*Golang) Installation(container *Container) *Container {
	container = container.
		With(dag.Redhat().Minimal().Packages([]string{
			"go",
			"git",
		}).Installed).
		WithMountedCache(CacheDir, dag.CacheVolume("golang")).
		WithEnvVariable("GOPATH", CacheDir).
		WithEnvVariable("GOCACHE", CacheDir+"/build")

	return container
}

func (golang *Golang) Container() *Container {
	return dag.Redhat().Minimal().Container().With(golang.Installation)
}
