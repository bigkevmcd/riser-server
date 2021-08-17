package core

// RuntimeConfig provides config for the server.
type RuntimeConfig struct {
	BootstrapApikey string `split_words:"true"`
	BindAddress     string `split_words:"true" default:":8000"`
	DeveloperMode   bool   `split_words:"true"`
	GitUrl          string `split_words:"true" required:"true"`
	// GitDir is the temp directory to store the contents of the state repo. Warning: this directory is deleted on Riser server startup
	GitDir                   string `split_words:"true" default:"/tmp/riser/git/"`
	GitBranch                string `split_words:"true" default:"main"`
	PostgresUrl              string `split_words:"true" default:"postgres://postgres.riser-system.svc.cluster.local/riserdb?sslmode=disable&connect_timeout=3"`
	PostgresUsername         string `split_words:"true"`
	PostgresPassword         string `split_words:"true"`
	PostgresMigrateOnStartup bool   `split_words:"true" default:"true"`
}
