package state

import "github.com/riser-platform/riser-server/pkg/core"

var _ Committer = (*DryRunCommitter)(nil)

// DryRunCommit stores a commit in memory.
type DryRunCommit struct {
	Message string
	Files   []core.ResourceFile
}

// DryRunCommitter is an implementation of the Committer interface that does not
// write the files or commit them, changes are retained in memory.
type DryRunCommitter struct {
	Commits []DryRunCommit
}

// NewDryRunCommitter creates and returns an initialised DryRunCommitter.
func NewDryRunCommitter() *DryRunCommitter {
	return &DryRunCommitter{
		Commits: []DryRunCommit{},
	}
}

// Commit stores the provided commit details in memory.
func (committer *DryRunCommitter) Commit(message string, files []core.ResourceFile) error {
	committer.Commits = append(committer.Commits, DryRunCommit{Message: message, Files: files})
	return nil
}
