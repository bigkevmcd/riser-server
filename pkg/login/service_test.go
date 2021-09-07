package login

import (
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/riser-platform/riser-server/pkg/core"
)

const testValidKey = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func TestService_LoginWithApiKey(t *testing.T) {
	plainText := "aabbccdd"
	var hash []byte
	user := &core.User{Username: "test"}
	userRepository := &core.FakeUserRepository{
		GetByApiKeyFn: func(hashArg []byte) (*core.User, error) {
			hash = hashArg
			return user, nil
		},
	}
	service := service{users: userRepository}

	result, err := service.LoginWithApiKey(plainText)

	assert.Equal(t, user, result)
	assert.NoError(t, err)
	assert.Equal(t, hashApiKey([]byte(plainText)), hash)
}

func TestService_LoginWithApiKey_trims_api_key(t *testing.T) {
	plainText := " aabbccdd "
	var hash []byte
	user := &core.User{Username: "test"}
	userRepository := &core.FakeUserRepository{
		GetByApiKeyFn: func(hashArg []byte) (*core.User, error) {
			hash = hashArg
			return user, nil
		},
	}
	service := service{users: userRepository}

	result, err := service.LoginWithApiKey(plainText)

	assert.Equal(t, user, result)
	assert.NoError(t, err)
	assert.Equal(t, hashApiKey([]byte("aabbccdd")), hash)
}

func TestService_LoginWithApiKey_not_found_returns_error(t *testing.T) {
	userRepository := &core.FakeUserRepository{
		GetByApiKeyFn: func([]byte) (*core.User, error) {
			return nil, core.ErrNotFound
		},
	}
	service := service{users: userRepository}

	result, err := service.LoginWithApiKey("nope")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidLogin)
}

func TestService_LoginWithApiKey_error_returns_error(t *testing.T) {
	userRepository := &core.FakeUserRepository{
		GetByApiKeyFn: func([]byte) (*core.User, error) {
			return nil, errors.New("test")
		},
	}
	service := service{users: userRepository}

	result, err := service.LoginWithApiKey("nope")

	assert.Nil(t, result)
	assert.EqualError(t, err, "test")
}

func TestService_BootstrapRootUser(t *testing.T) {
	var rootUserId uuid.UUID
	userRepository := &core.FakeUserRepository{
		GetByUsernameFn: func(username string) (*core.User, error) {
			assert.Equal(t, RootUsername, username)
			return nil, core.ErrNotFound
		},
		CreateFn: func(newUser *core.NewUser) error {
			assert.NotEqual(t, uuid.Nil, newUser.Id)
			assert.Equal(t, RootUsername, newUser.Username)
			rootUserId = newUser.Id
			return nil
		},
	}
	apikeyRepository := &core.FakeApiKeyRepository{
		CreateFn: func(userId uuid.UUID, keyHash []byte) error {
			assert.Equal(t, rootUserId, userId)
			assert.Equal(t, hashApiKey([]byte(testValidKey)), keyHash)
			return nil
		},
	}

	service := service{userRepository, apikeyRepository}

	err := service.BootstrapRootUser(testValidKey)

	assert.NoError(t, err)
	assert.Equal(t, 1, userRepository.CreateCallCount)
	assert.Equal(t, 1, apikeyRepository.CreateCallCount)
}

func TestService_BootstrapRootUser_UnableToCreateApiKey_ReturnsError(t *testing.T) {
	userRepository := &core.FakeUserRepository{
		GetByUsernameFn: func(username string) (*core.User, error) {
			return nil, core.ErrNotFound
		},
		CreateFn: func(newUser *core.NewUser) error {
			return nil
		},
	}
	apikeyRepository := &core.FakeApiKeyRepository{
		CreateFn: func(uuid.UUID, []byte) error {
			return errors.New("test")
		},
	}

	service := service{userRepository, apikeyRepository}

	err := service.BootstrapRootUser(testValidKey)

	assert.EqualError(t, err, "error creating root API key: test")
}

func TestService_BootstrapRootUser_UserWithKeyExists_ReturnsError(t *testing.T) {
	userRepository := &core.FakeUserRepository{
		GetByUsernameFn: func(string) (*core.User, error) {
			return &core.User{Id: uuid.New()}, nil
		},
	}
	apikeyRepository := &core.FakeApiKeyRepository{
		GetByUserIdFn: func(uuid.UUID) ([]core.ApiKey, error) {
			return []core.ApiKey{core.ApiKey{}}, nil
		},
	}

	service := service{userRepository, apikeyRepository}

	err := service.BootstrapRootUser(testValidKey)

	assert.Equal(t, ErrRootUserExists, err)
}

func TestService_BootstrapRootUser_UnableToCreateUser_ReturnsError(t *testing.T) {
	userRepository := &core.FakeUserRepository{
		GetByUsernameFn: func(string) (*core.User, error) {
			return nil, core.ErrNotFound
		},
		CreateFn: func(newUser *core.NewUser) error {
			return errors.New("test")
		},
	}
	apikeyRepository := &core.FakeApiKeyRepository{}

	service := service{userRepository, apikeyRepository}

	err := service.BootstrapRootUser(testValidKey)

	assert.EqualError(t, err, "unable to create root user: test")
}

func TestService_BootstrapRootUser_UnableToQueryKeys_ReturnsError(t *testing.T) {
	userRepository := &core.FakeUserRepository{
		GetByUsernameFn: func(string) (*core.User, error) {
			return &core.User{Id: uuid.New()}, nil
		},
	}
	apikeyRepository := &core.FakeApiKeyRepository{
		GetByUserIdFn: func(uuid.UUID) ([]core.ApiKey, error) {
			return nil, errors.New("test")
		},
	}

	service := service{userRepository, apikeyRepository}

	err := service.BootstrapRootUser(testValidKey)

	assert.EqualError(t, err, "unable to retrieve root API keys: test")
}

func TestService_BootstrapRootUser_UnableToQueryUser_ReturnsError(t *testing.T) {
	userRepository := &core.FakeUserRepository{
		GetByUsernameFn: func(string) (*core.User, error) {
			return nil, errors.New("test")
		},
	}

	service := service{users: userRepository}

	err := service.BootstrapRootUser(testValidKey)

	assert.EqualError(t, err, "unable to retrieve root user: test")
}

func TestService_BootstrapRootUser_ShortKey_ReturnsError(t *testing.T) {
	service := service{}

	err := service.BootstrapRootUser("oops")

	assert.EqualError(t, err, "API Key must be a minimum of 32 characters. It is highly recommended to use `riser ops generate-apikey` to generate the key")
}

func TestService_BootstrapRootUser_EmptyKeyWithUsers_NOOP(t *testing.T) {
	userRepository := &core.FakeUserRepository{
		GetActiveCountFn: func() (int, error) {
			return 1, nil
		},
	}
	service := service{users: userRepository}

	err := service.BootstrapRootUser("")

	assert.NoError(t, err)
}

func TestService_BootstrapRootUser_EmptyKeyNoUsers_NOOP(t *testing.T) {
	userRepository := &core.FakeUserRepository{
		GetActiveCountFn: func() (int, error) {
			return 0, nil
		},
	}
	service := service{users: userRepository}

	err := service.BootstrapRootUser("")

	assert.EqualError(t, err, "You must specify RISER_BOOTSTRAP_APIKEY is required when there are no users. Use \"riser ops generate-apikey\" to generate the key.")
}
