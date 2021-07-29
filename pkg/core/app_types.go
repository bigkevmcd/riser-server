package core

import (
	"github.com/google/uuid"
)

// App represents a named application.
type App struct {
	Id        uuid.UUID
	Name      string
	Namespace string
}

// AppStatus represents the state of an application including its deployments.
type AppStatus struct {
	AppId             uuid.UUID
	EnvironmentStatus []EnvironmentStatus
	// Deployments returns the whole deployment. We should probably use a different type here with less data, but we can't just pass
	// Deployment.Doc.Status as we also need the DeploymentName and the environment.
	Deployments []Deployment
}
