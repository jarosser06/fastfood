package fastfood

type FrameworkOptions struct {
	Destination string
	BaseFiles   []string
	BaseDirs    []string
	Force       bool
	Name        string
}

type Framework interface {
	// Initializes the framework
	// gives the framework the opportunity to set things up for later
	Init(FrameworkOptions) error
	// Generate Empty generates a new base directory and files for a framework
	// it returns a slice of files that were modified and an error
	GenerateBase() ([]string, error)
	// Generate Stencil generates a stencil and returns a slice of
	// files that were modified and an error
	// Accepts a stencil name, a stencilset and options
	GenerateStencil(string, StencilSet, map[string]string) ([]string, error)
}
