package environment

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/riser-platform/riser-server/pkg/git"
)

func Test_getRepo(t *testing.T) {
	settings := git.RepoSettings{}

	cache := NewFakeRepoCache()

	repo1, err := cache.getRepo("env1", settings)
	require.NoError(t, err)
	repo2, err := cache.getRepo("env1", settings)
	require.NoError(t, err)
	repo3, err := cache.getRepo("env2", settings)
	require.NoError(t, err)

	assert.IsType(t, &git.FakeRepo{}, repo1)
	assert.Same(t, repo1, repo2)
	assert.NotSame(t, repo1, repo3)
}

func Test_getRepo_emptyName(t *testing.T) {
	cache := NewFakeRepoCache()

	repo, err := cache.getRepo("", git.RepoSettings{})

	assert.Nil(t, repo)
	assert.Equal(t, "Environment name cannot be empty", err.Error())
}

func Test_getRepo_NewErr(t *testing.T) {
	expectedErr := errors.New("test")
	cache := RepoCache{
		newFunc: func(git.RepoSettings) (git.Repo, error) {
			return nil, expectedErr
		},
		cache: map[string]git.Repo{},
	}

	repo, err := cache.getRepo("env", git.RepoSettings{})

	assert.Nil(t, repo)
	assert.Equal(t, expectedErr, err)
}

func Test_newGitSettingsForEnv(t *testing.T) {
	settings := RepoSettings{
		URL:        "git@my.org/state",
		BaseGitDir: "/tmp/riserstate",
	}
	result := newGitSettingsForEnv("env1", settings)

	assert.Equal(t, settings.URL, result.URL)
	assert.Equal(t, "env1", result.Branch)
	assert.Equal(t, "/tmp/riserstate/env/env1", result.BaseWorkspaceDir)
}
