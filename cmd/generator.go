package cmd

type Generator struct {
}

func (g *Generator) Help() string {
	return "TODO: Add help text for gen command"
}

func (g *Generator) Run(args []string) int {
	return 0
}

func (g *Generator) Synopsis() string {
	return "Generates a new recipe for an existing cookbook"
}
