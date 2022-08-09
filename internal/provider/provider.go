package provider

type Variables map[string]string

func (vars Variables) Merge(other Variables) Variables {
	for key, val := range other {
		vars[key] = val
	}
	return vars
}
