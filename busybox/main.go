package main

const (
	ImageRegistry   string = "docker.io"
	ImageRepository string = "busybox"
	ImageTag        string = "1.36.1"
	ImageDigest     string = "sha256:ba76950ac9eaa407512c9d859cea48114eeff8a6f12ebaa5d32ce79d4a017dd8"
)

type Busybox struct{}

func (busybox *Busybox) Container() *Container {
	container := dag.Container().
		From(ImageRegistry + "/" + ImageRepository + ":" + ImageTag + "@" + ImageDigest).
		WithEntrypoint([]string{"sh", "-c"}).
		WithoutDefaultArgs().
		WithWorkdir("/home")

	return container
}
