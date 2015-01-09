package chef

import (
	"fmt"

	"github.com/jarosser06/fastfood/common/json"
	"github.com/jarosser06/fastfood/common/maputil"
)

type Options struct {
	BerksDeps    map[string]BerksCookbook `json:"berks_dependencies"`
	Dependencies []string                 `json:"dependencies"`
	Directories  []string                 `json:"directories"`
	Files        map[string]string        `json:"files"`
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
func NewOptions(conf []byte) (Options, error) {
	newOptions := Options{}
	err := json.Unmarshal(conf, &newOptions)
	if err != nil {
		return newOptions, fmt.Errorf("parsing json %v", err)
	}

	// Make sure the name is set
	for n, d := range newOptions.BerksDeps {
		d.Name = n
		newOptions.BerksDeps[n] = d
	}

	if newOptions.Files == nil {
		newOptions.Files = make(map[string]string)
	}

	return newOptions, nil
}
