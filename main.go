package main

import (
	"database/sql"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	echolog "github.com/onrik/logrus/echo"
	"github.com/sirupsen/logrus"

	"github.com/riser-platform/riser-server/api"
	apiv1 "github.com/riser-platform/riser-server/api/v1"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/environment"
	"github.com/riser-platform/riser-server/pkg/login"
	"github.com/riser-platform/riser-server/pkg/namespace"
	"github.com/riser-platform/riser-server/pkg/postgres"
	"github.com/riser-platform/riser-server/pkg/util"
)

// All env vars are prefixed with RISER_
const envPrefix = "RISER"

var logger = logrus.StandardLogger()

func main() {
	logger.SetFormatter(&logrus.JSONFormatter{})

	var rc core.RuntimeConfig
	exitIfError(envconfig.Process(envPrefix, &rc), "Error loading environment variables")

	if rc.DeveloperMode {
		logger.SetFormatter(&logrus.TextFormatter{})
		logger.Info("Developer mode active")
	}

	logger.Infof("Server version %s", util.VersionString)
	logger.Info("Initializing postgres")
	postgresConn, err := postgres.AddAuthToConnString(rc.PostgresUrl, rc.PostgresUsername, rc.PostgresPassword)
	exitIfError(err, "Error creating postgres connection url")

	postgresDb, err := postgres.NewDB(postgresConn)
	exitIfError(err, "Error initializing postgres")

	if rc.PostgresMigrateOnStartup {
		logger.Info("Applying Postgres migrations")
		err = postgres.Migrate(postgresConn, logger)
		exitIfError(err, "Error performing Postgres migrations")
	}

	repoSettings := environment.RepoSettings{
		URL:        rc.GitUrl,
		BaseGitDir: rc.GitDir,
	}
	repoCache := environment.NewBranchPerEnvRepoCache(repoSettings)

	bootstrapApiKey(postgresDb, &rc)
	bootstrapDefaultNamespace(postgresDb)

	e := echo.New()
	e.HideBanner = true

	e.Logger = echolog.NewLogger(logger, "")
	e.Use(echolog.Middleware(echolog.DefaultConfig))
	e.HTTPErrorHandler = api.ErrorHandler
	e.Binder = &api.DataBinder{}

	apiv1.RegisterRoutes(e, repoCache, postgresDb)
	err = e.Start(rc.BindAddress)
	exitIfError(err, "Error starting server")
}

func bootstrapDefaultNamespace(db *sql.DB) {
	namespaceService := namespace.NewService(postgres.NewNamespaceRepository(db), postgres.NewEnvironmentRepository(db))
	err := namespaceService.EnsureDefaultNamespace()
	exitIfError(err, "Error ensuring default namespace")
}

func bootstrapApiKey(db *sql.DB, rc *core.RuntimeConfig) {
	loginService := login.NewService(postgres.NewUserRepository(db), postgres.NewApiKeyRepository(db))
	err := loginService.BootstrapRootUser(rc.BootstrapApikey)
	if err != nil {
		if err == login.ErrRootUserExists {
			logger.Info("Ignoring environment variable RISER_BOOTSTRAP_APIKEY: root user already exists.")
		} else {
			exitIfError(err, "Unable to bootstrap API KEY")
		}
	}
}

func exitIfError(err error, message string) {
	if err != nil {
		logger.Fatalf("%s: %s", message, err)
	}
}
