package config

type Config struct {
	Global   GlobalConfig          `toml:"global"`
	Projects map[string]*Project  `toml:"projects"`
}

type GlobalConfig struct {
	SecretEnv  string `toml:"secret_env"`
	HistoryDB  string `toml:"history_db"`
	Fallback   bool   `toml:"fallback"`
}

type Project struct {
	Name         string   `toml:"name"`
	WorkingDir   string   `toml:"working_dir"`
	Excludes     []string `toml:"excludes"`
	RunMigrations bool    `toml:"run_migrations"`
	FTPS         FTPSConfig `toml:"ftps"`
	Deploy      DeployConfig `toml:"deploy"`
}

type FTPSConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
	User string `toml:"user"`
	Pass string `toml:"pass"`
}

type DeployConfig struct {
	TargetPath   string `toml:"target_path"`
	BackupPrefix string `toml:"backup_prefix"`
	VerifyURL    string `toml:"verify_url"`
	Endpoint    string `toml:"endpoint"`
}

func DefaultConfig() *Config {
	return &Config{
		Global: GlobalConfig{
			SecretEnv: "ENIAC_DEPLOY_SECRET",
			HistoryDB: "~/.eniac-deploy/history.db",
			Fallback:  false,
		},
		Projects: make(map[string]*Project),
	}
}