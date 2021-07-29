package core

import (
	"fmt"
	"strings"
)

const DefaultNamespace = "apps"

// Namespace represents a namespace.
type Namespace struct {
	Name string
}

// NamespacedName is a typed name with a namespace.
//
// https://pkg.go.dev/k8s.io/apimachinery/pkg/types#NamespacedName
type NamespacedName struct {
	Name      string
	Namespace string
}

func (v *NamespacedName) String() string {
	return strings.ToLower(fmt.Sprintf("%s.%s", v.Name, v.Namespace))
}

// ParseNamespacedName parses a string into a NamespacedName.
func ParseNamespacedName(namespacedName string) *NamespacedName {
	parts := strings.Split(strings.ToLower(namespacedName), ".")
	if len(parts) == 2 {
		return &NamespacedName{parts[0], parts[1]}
	}

	return &NamespacedName{Name: strings.ToLower(namespacedName)}
}

// NewNamespacedName creates and returns a NamespacedName.
//
// The provided name and namespace are converted to lower-case strings.
//
// If the namespace is empty, then it defaults to DefaultNamespace.
func NewNamespacedName(name, namespace string) *NamespacedName {
	if namespace == "" {
		namespace = DefaultNamespace
	}

	return &NamespacedName{Name: strings.ToLower(name), Namespace: strings.ToLower(namespace)}
}
