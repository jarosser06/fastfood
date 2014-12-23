package cmd

type Builder struct {
	CookbooksPath string
}

func (b *Builder) Run(args []string) int {

	return 0
}

func (b *Builder) Synopsis() string {
	return "Creates a cookbook w/ providers from a config file"
}

func (b *Builder) Help() string {
	return `
Usage:
  fastfood build <flags> [config_file]

  This command will build a new cookbook or update an
  existing one using a configuration file.
`
}
