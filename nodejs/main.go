package main

const (
	CacheDir string = "/var/cache/node"
)

type Nodejs struct {
	// +private
	Npmrc *Secret
}

func New(
	// +optional
	npmrc *Secret,
) *Nodejs {
	nodejs := &Nodejs{
		Npmrc: npmrc,
	}

	return nodejs
}

func (nodejs *Nodejs) Installation(container *Container) *Container {
	container = container.
		With(dag.Redhat().Minimal().Module("nodejs:20").Enabled).
		With(dag.Redhat().Minimal().Packages([]string{
			"npm",
		}).Installed).
		WithMountedCache(CacheDir, dag.CacheVolume("nodejs")).
		WithEnvVariable("NPM_CONFIG_CACHE", CacheDir+"/npm")

	if nodejs.Npmrc != nil {
		container = container.
			WithMountedSecret("/root/.npmrc", nodejs.Npmrc)
	}

	return container
}

func (nodejs *Nodejs) Container() *Container {
	return dag.Redhat().Minimal().Container().With(nodejs.Installation)
}
