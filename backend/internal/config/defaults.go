package config

var DefaultExcludes = []string{
	".git",
	".git*",
	"node_modules",
	"vendor",
	"tests",
	".env",
	".env.*",
	"*.log",
	"storage/logs/*",
	"storage/framework/cache/*",
	"*.bak",
	"*.tmp",
	".DS_Store",
}

var DefaultBackupPrefix = ".env.bak"
var DefaultEndpoint = "_deploy_%s.php"

func ApplyDefaults(p *Project) {
	if p.Excludes == nil {
		p.Excludes = DefaultExcludes
	}
	if p.Deploy.BackupPrefix == "" {
		p.Deploy.BackupPrefix = DefaultBackupPrefix
	}
	if p.Deploy.Endpoint == "" {
		p.Deploy.Endpoint = DefaultEndpoint
	}
	if p.FTPS.Port == 0 {
		p.FTPS.Port = 21
	}
}