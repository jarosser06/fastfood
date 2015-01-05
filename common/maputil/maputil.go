package maputil

// Merge a higher priority map[string]string with a lower priorit map[string]string
func Merge(hMap map[string]string, lMap map[string]string) map[string]string {
	nMap := hMap
	for k, v := range lMap {
		nMap[k] = v
	}

	return nMap
}
