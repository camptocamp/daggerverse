package main

const (
	ImageRegistry   string = "docker.io"
	ImageRepository string = "yuzutech/kroki"
	ImageTag        string = "0.24.1"
	ImageDigest     string = "sha256:b43be03ec8a210471d4eaaf044b44e765551730447516c26f15c0f4b27628d45"
)

type Kroki struct{}

func New() *Kroki {
	kroki := &Kroki{}

	return kroki
}

func (kroki *Kroki) Container() *Container {
	container := dag.Container().
		From(ImageRegistry + "/" + ImageRepository + ":" + ImageTag + "@" + ImageDigest).
		WithExposedPort(8000)

	return container
}

func (kroki *Kroki) Server() *Service {
	server := kroki.Container().
		AsService()

	return server
}
