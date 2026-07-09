package player

// CloneVariables returns a shallow copy of variables so parallel scenarios
// cannot mutate a shared map.
func CloneVariables(variables map[string]string) map[string]string {
	if len(variables) == 0 {
		return map[string]string{}
	}
	cloned := make(map[string]string, len(variables))
	for key, value := range variables {
		cloned[key] = value
	}
	return cloned
}
