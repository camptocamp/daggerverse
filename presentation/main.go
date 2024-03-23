package main

type Presentation struct{}

func New() *Presentation {
	presentation := &Presentation{}

	return presentation
}

func (presentation *Presentation) Init() *Directory {
	template := dag.CurrentModule().Source().Directory("template").
		WithoutFile(".gitignore")

	return template
}

type PresentationBuilder struct {
	// +private
	Directory *Directory
	// +private
	Npmrc *Secret
}

func (presentation *Presentation) Builder(
	directory *Directory,
	npmrc *Secret,
) *PresentationBuilder {
	builder := &PresentationBuilder{
		Directory: directory,
		Npmrc:     npmrc,
	}

	return builder
}

func (builder *PresentationBuilder) Container() *Container {
	kroki := dag.Kroki()

	container := dag.Redhat().Container().
		With(dag.Nodejs(NodejsOpts{
			Npmrc: builder.Npmrc,
		}).Configuration).
		WithMountedDirectory(".", builder.Directory).
		WithExec([]string{"npm clean-install"}).
		WithEntrypoint([]string{"npm", "run", "all"}).
		WithoutDefaultArgs()

	// FIXME
	// Doesn’t work at the time.
	//WithServiceBinding("kroki", kroki.Server())
	_ = kroki

	return container
}

type PresentationBuild struct {
	// +private
	Builder *Container
}

func (builder *PresentationBuilder) Build() *PresentationBuild {
	build := &PresentationBuild{
		Builder: builder.Container().WithExec(nil),
	}

	return build
}

func (build *PresentationBuild) Directory() *Directory {
	directory := build.Builder.Directory("dist")

	return directory
}

func (build *PresentationBuild) Container() *Container {
	directory := build.Directory()
	container := dag.Caddy(directory).Container()

	return container
}

func (build *PresentationBuild) Server() *Service {
	server := build.Container().AsService()

	return server
}
