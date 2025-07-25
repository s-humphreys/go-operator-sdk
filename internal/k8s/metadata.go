package k8s

type Labels map[string]string

// CreateLabels creates a set of common labels for a Kubernetes resource.
func CreateLabels(name string) Labels {
	return Labels{
		"app": name,
	}
}
