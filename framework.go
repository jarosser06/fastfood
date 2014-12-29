package fastfood

type FrameworkOptions struct {
	Destination string
	Force       bool
	Name        string
}

type Framework interface {
	// Generate Empty generates a new base directory and files for a framework
	// it returns a slice of files that were modified and an error
	GenerateEmpty(FrameworkOptions) ([]string, error)
	// Generate Stencil generates a stencil and returns a slice of
	// files that were modified and an error
	GenerateStencil(StencilSet, FrameworkOptions) ([]string, error)
}
