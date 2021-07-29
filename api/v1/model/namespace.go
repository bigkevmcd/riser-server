package model

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v3"
)

/*
bannedNamespacePrefixes is a list of prefixes that we don't allow as they collide with system and possible future namespaces.
This is for usability purposes only and should not be substituted for a robust RBAC policy to prevent the creation or deployment
to certain namespaces.
*/
var bannedNamespacePrefixes = []string{"riser-", "kube-", "knative-", "istio-"}

// Namespace represents a namespace in a deployment cluster.
type Namespace struct {
	Name NamespaceName `json:"name"`
}

// Validate performs some checks and returns an error if the Namespace is not
// valid.
func (v Namespace) Validate() error {
	return validation.ValidateStruct(&v, validation.Field(&v.Name, validation.Required))
}

// NamespaceName is a string with validation rules.
type NamespaceName string

func (v NamespaceName) Validate() error {
	return validation.Validate(string(v), append(RulesNamingIdentifier(), validation.Required, validation.By(bannedNamespaceRule))...)
}

func bannedNamespaceRule(v interface{}) error {
	vStr := v.(string)
	for _, prefix := range bannedNamespacePrefixes {
		if strings.HasPrefix(vStr, prefix) {
			return fmt.Errorf("namespace names may not begin with %q", prefix)
		}
	}

	return nil
}
