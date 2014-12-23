package cmd

import (
	"os"
	"path"
	"strings"
)

const (
	tempPackEnvVar     = "FASTFOOD_TEMPLATE_PACK"
	cookbookPathEnvVar = "COOKBOOKS"
)

func DefaultTempPack() string {
	packEnv := os.Getenv(tempPackEnvVar)
	if packEnv == "" {
		return path.Join(os.Getenv("HOME"), "fastfood")
	} else {
		return packEnv
	}
}

func DefaultCookbookPath() string {
	cookbookPath := os.Getenv(cookbookPathEnvVar)

	return cookbookPath
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
