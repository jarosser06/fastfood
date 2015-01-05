package cmd

import (
	"os"
	"path"
	"strings"
)

const (
	EnvTempPack      = "FASTFOOD_TEMPLATE_PACK"
	EnvCookbooksPath = "COOKBOOKS"
)

// Templatepack is the path to the templatepack
// CookbookPath is the path to cookbooks
// Manifest is loaded using LoadManifest
type Common struct {
	TemplatePack string
	CookbookPath string
	Force        bool
}

// Translates key:value strings into a map
func MapArgs(args []string) map[string]string {
	var argMap map[string]string
	argMap = make(map[string]string)

	for _, arg := range args {
		if strings.Contains(arg, ":") {
			// Split at the first : in an arg
			splitArg := strings.SplitN(arg, ":", 2)

			argMap[splitArg[0]] = splitArg[1]
		}
	}

	return argMap
}

func DefaultTempPack() string {
	packEnv := os.Getenv(EnvTempPack)
	if packEnv == "" {
		return path.Join(os.Getenv("HOME"), "fastfood")
	} else {
		return packEnv
	}
}

func DefaultCookbooksPath() string {
	pathEnv := os.Getenv(EnvCookbooksPath)
	if pathEnv == "" {
		return path.Join(os.Getenv("HOME"), "cookbooks")
	} else {
		return pathEnv
	}
}
