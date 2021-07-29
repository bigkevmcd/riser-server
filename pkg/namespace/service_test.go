package namespace

import (
	"errors"
	"testing"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Create(t *testing.T) {
	namespaces := &core.FakeNamespaceRepository{
		CreateFn: func(namespace *core.Namespace) error {
			assert.Equal(t, "myns", namespace.Name)
			return nil
		},
	}
	environments := &core.FakeEnvironmentRepository{
		ListFn: func() ([]core.Environment, error) {
			return []core.Environment{
				{Name: "myenv1"},
				{Name: "myenv2"},
			}, nil
		},
	}
	svc := &service{namespaces, environments}

	err := svc.Create("myns")

	assert.NoError(t, err)
	assert.Equal(t, 1, namespaces.CreateCallCount)
}

func TestService_Create_when_namespace_create_err(t *testing.T) {
	namespaces := &core.FakeNamespaceRepository{
		CreateFn: func(namespace *core.Namespace) error {
			return errors.New("test")
		},
	}
	svc := &service{namespaces: namespaces}

	err := svc.Create("myns")

	assert.EqualError(t, err, `error creating namespace "myns": test`)
}

func TestService_EnsureDefaultNamespace_returns_err(t *testing.T) {
	namespaces := &core.FakeNamespaceRepository{
		GetFn: func(string) (*core.Namespace, error) {
			return nil, errors.New("test")
		},
	}
	svc := &service{namespaces: namespaces}

	err := svc.EnsureDefaultNamespace()

	assert.EqualError(t, err, "test")
}

func TestService_EnsureDefaultNamespace_when_exists_noop(t *testing.T) {
	namespaces := &core.FakeNamespaceRepository{
		GetFn: func(string) (*core.Namespace, error) {
			return &core.Namespace{}, nil
		},
	}
	svc := &service{namespaces: namespaces}

	err := svc.EnsureDefaultNamespace()

	assert.NoError(t, err)
}

func TestService_ValidateDeployable_NamespaceExists(t *testing.T) {
	namespaces := &core.FakeNamespaceRepository{
		GetFn: func(namespaceArg string) (*core.Namespace, error) {
			assert.Equal(t, "myns", namespaceArg)
			return &core.Namespace{Name: namespaceArg}, nil
		},
	}
	svc := &service{namespaces: namespaces}

	err := svc.ValidateDeployable("myns")

	assert.NoError(t, err)
	assert.Equal(t, 1, namespaces.GetCallCount)
}

func TestService_ValidateDeployable_namespace_missing(t *testing.T) {
	namespaces := &core.FakeNamespaceRepository{
		GetFn: func(namespaceArg string) (*core.Namespace, error) {
			return nil, core.ErrNotFound
		},
		ListFn: func() ([]core.Namespace, error) {
			return []core.Namespace{
				{Name: "ns1"},
				{Name: "ns2"},
			}, nil
		},
	}
	svc := &service{namespaces: namespaces}

	err := svc.ValidateDeployable("myns")

	require.IsType(t, &core.ValidationError{}, err, err.Error())
	assert.EqualError(t, err, `invalid namespace "myns". Must be one of: ns1, ns2`)
}

func TestService_ValidateDeployable_namespace_missing_list_error(t *testing.T) {
	namespaces := &core.FakeNamespaceRepository{
		GetFn: func(namespaceArg string) (*core.Namespace, error) {
			return nil, core.ErrNotFound
		},
		ListFn: func() ([]core.Namespace, error) {
			return nil, errors.New("test")
		},
	}
	svc := &service{namespaces: namespaces}

	err := svc.ValidateDeployable("myns")

	assert.EqualError(t, err, "test")
}

func TestService_ValidateDeployable_get_error(t *testing.T) {
	namespaces := &core.FakeNamespaceRepository{
		GetFn: func(namespaceArg string) (*core.Namespace, error) {
			return nil, errors.New("test")
		},
	}
	svc := &service{namespaces: namespaces}

	err := svc.ValidateDeployable("myns")

	assert.EqualError(t, err, "test")
}
