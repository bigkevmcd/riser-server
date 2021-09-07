package postgres

import (
	"database/sql"
	"net/url"

	"github.com/pkg/errors"
	// Required for pq lib dynamic driver loading
	_ "github.com/lib/pq"

	"github.com/riser-platform/riser-server/pkg/core"
)

// NewDB a database using the provided connection string.
//
// The newly opened database is pinged to check that it is available and an
// error is returned if this fails.
func NewDB(postgresConn string) (*sql.DB, error) {
	var err error
	db, err := sql.Open("postgres", postgresConn)
	if err != nil {
		return nil, errors.Wrap(err, "error opening connection to posgres")
	}

	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "error pinging postgres")
	}

	return db, nil
}

// AddAuthToConnString adds authentication credentials to the provided data
// source name.
//
// If the URL already has credentials then these are not replaced.
func AddAuthToConnString(postgresConn string, username string, password string) (string, error) {
	postgresUrl, err := url.Parse(postgresConn)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse connection string")
	}
	if postgresUrl.User == nil {
		postgresUrl.User = url.UserPassword(username, password)
	}
	return postgresUrl.String(), nil
}

func resultHasRows(r sql.Result) bool {
	if r == nil {
		return false
	}
	rows, err := r.RowsAffected()
	return err == nil && rows > 0
}

func noRowsErrorHandler(err error) error {
	if err == sql.ErrNoRows {
		return core.ErrNotFound
	}

	return errors.WithStack(err)
}
