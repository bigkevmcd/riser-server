package postgres

import (
	"database/sql"

	"github.com/riser-platform/riser-server/pkg/core"
)

// NewNamespaceRepository creates and returns a new repository for storing
// namespaces.
func NewNamespaceRepository(db *sql.DB) core.NamespaceRepository {
	return &namespaceRepository{db: db}
}

type namespaceRepository struct {
	db *sql.DB
}

func (r *namespaceRepository) Create(namespace *core.Namespace) error {
	_, err := r.db.Exec("INSERT INTO namespace (name) VALUES ($1)", namespace.Name)
	return err
}

func (r *namespaceRepository) Get(namespaceName string) (*core.Namespace, error) {
	ns := &core.Namespace{}
	// Effectively used just to make sure that the namespace exists. This will do in the future as we add more fields.
	err := r.db.QueryRow("SELECT name FROM namespace WHERE name = $1", namespaceName).Scan(&ns.Name)
	if err != nil {
		return nil, noRowsErrorHandler(err)
	}
	return ns, nil
}

func (r *namespaceRepository) List() ([]core.Namespace, error) {
	namespaces := []core.Namespace{}
	rows, err := r.db.Query("SELECT name FROM namespace")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		ns := core.Namespace{}
		err := rows.Scan(&ns.Name)
		if err != nil {
			return nil, err
		}
		namespaces = append(namespaces, ns)
	}

	return namespaces, nil
}
