package chef

import (
	"encoding/json"
	"fmt"

	"github.com/jarosser06/fastfood/common/maputil"
)

type Options struct {
	BerksDeps    map[string]Berks  `json:"berks_dependencies"`
	Dependencies []string          `json:"dependencies"`
	Directories  []string          `json:"directories"`
	Files        map[string]string `json:"files"`
	Partials     []string
}

// Merge Two stencils together giving local a higher priority
func Merge(global Options, local Options) Options {
	o := Options{}

	o.Dependencies = append(global.Dependencies, local.Dependencies...)
	// Can't use maputils Merge since its special
	if local.BerksDeps != nil {
		if global.BerksDeps != nil {
			o.BerksDeps = global.BerksDeps
			for k, v := range local.BerksDeps {
				o.BerksDeps[k] = v
			}
		} else {
			o.BerksDeps = local.BerksDeps
		}
	}

	o.Files = maputil.Merge(local.Files, global.Files)
	o.Partials = append(global.Partials, local.Partials...)
	o.Directories = append(global.Directories, local.Directories...)

	return o
}

// Given the framework information create the framework stencil
func NewOptions(conf string) (Options, error) {
	newOptions := Options{}
	err := json.Unmarshal([]byte(conf), &newOptions)
	if err != nil {
		return newOptions, fmt.Errorf("error parsing json %v", err)
	}

	if newOptions.Files == nil {
		newOptions.Files = make(map[string]string)
	}

	return newOptions, nil
}
