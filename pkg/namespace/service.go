package namespace

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
)

// Service provides functionality for creating and validating namespaces.
type Service interface {
	// Create creates a new namespace with the provided name in the namespaces
	// repository.
	Create(namespaceName string) error

	// EnsureDefaultNamespace ensures that the default namespace has been
	// provisioned. Designed to be used only at server startup.
	EnsureDefaultNamespace() error

	// ValidateDeployable validates that a namespace is deployable. Returns a
	// ValidationError if it is not.
	ValidateDeployable(namespaceName string) error
}

type service struct {
	namespaces   core.NamespaceRepository
	environments core.EnvironmentRepository
}

// NewService creates and returns a new namespace service.
func NewService(namespaces core.NamespaceRepository, environments core.EnvironmentRepository) Service {
	return &service{namespaces, environments}
}

// EnsureDefaultNamespace checks for the default namespace being in the
// namespaces repository, and if it doesn't exist, it automatically creates it.
func (s *service) EnsureDefaultNamespace() error {
	_, err := s.namespaces.Get(core.DefaultNamespace)
	if err == nil {
		return nil
	}
	if err == core.ErrNotFound {
		return s.Create(core.DefaultNamespace)
	}
	return err
}

// Create creates a namespace in the namespaces repository with the provided
// name.
func (s *service) Create(namespaceName string) error {
	err := s.namespaces.Create(&core.Namespace{Name: namespaceName})
	if err != nil {
		return errors.Wrapf(err, "error creating namespace %q", namespaceName)
	}
	return nil
}

// ValidateDeployable returns an error if the provided name does not exist in
// the namespaces repository.
func (s *service) ValidateDeployable(namespaceName string) error {
	_, err := s.namespaces.Get(namespaceName)
	if err != core.ErrNotFound {
		return err
	}
	namespaces, nsListErr := s.namespaces.List()
	if nsListErr != nil {
		return nsListErr
	}
	validNamespaceNames := toNameList(namespaces)
	return core.NewValidationErrorMessage(fmt.Sprintf("invalid namespace %q. Must be one of: %s", namespaceName, strings.Join(validNamespaceNames, ", ")))
}

func toNameList(namespaces []core.Namespace) []string {
	names := []string{}
	for _, namespace := range namespaces {
		names = append(names, namespace.Name)
	}
	return names
}
