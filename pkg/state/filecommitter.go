package state

import (
	"io/ioutil"
	"path/filepath"

	"github.com/riser-platform/riser-server/pkg/util"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
)

var _ Committer = (*FileCommitter)(nil)

// FileCommiter writes the committed files to the filesystem relative to the
// basePath.
//
// No commit is performed.
type FileCommitter struct {
	basePath string
}

// NewFileCommitter creates and returns an initialised FileCommitter rooted at
// the provided base path.
func NewFileCommitter(basePath string) *FileCommitter {
	return &FileCommitter{basePath}
}

// Commit writes the provided files tot he filesystem.
func (committer *FileCommitter) Commit(_ string, files []core.ResourceFile) error {
	for _, file := range files {
		fullpath := filepath.Join(committer.basePath, file.Name)
		err := util.EnsureDir(fullpath, 0755)
		if err != nil {
			return errors.Wrapf(err, "error creating directory for file %q", fullpath)
		}
		err = ioutil.WriteFile(fullpath, file.Contents, 0644)
		if err != nil {
			return errors.Wrapf(err, "error writing file %q", fullpath)
		}
	}
	return nil
}
