package utils

func IntValueOrDefault(v, def int) int {
	if v != 0 {
		return v
	}

	return def
}

func StringTrunc(s string, max int) string {
	if len(s) <= max {
		return s
	}

	return s[:max]
}

type NamedResource struct {
	Name string
}

func MapNamedResource(resource []*NamedResource) map[string]*NamedResource {
	var resources map[string]*NamedResource
	for _, r := range resource {
		resources[r.Name] = r
	}
	return resources
}