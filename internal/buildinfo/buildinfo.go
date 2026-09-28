package buildinfo

import "runtime"

// Values are injected by the reproducible build. Defaults deliberately make
// an unversioned local binary obvious without preventing local development.
var (
	Service         = "dinapay-connector-template"
	Repository      = "github.com/Germatic/dinapay-connector-template"
	ContractVersion = "connector-v1"
	Version         = "dev"
	Commit          = "unknown"
	BuiltAt         = "unknown"
	Environment     = "local"
)

type Info struct {
	Service         string `json:"service"`
	Repository      string `json:"repository"`
	ContractVersion string `json:"contractVersion"`
	Version         string `json:"version"`
	Commit          string `json:"commit"`
	BuiltAt         string `json:"builtAt"`
	GoVersion       string `json:"goVersion"`
	Environment     string `json:"environment"`
}

func Current() Info {
	return Info{
		Service: Service, Repository: Repository,
		ContractVersion: ContractVersion, Version: Version,
		Commit: Commit, BuiltAt: BuiltAt, GoVersion: runtime.Version(),
		Environment: Environment,
	}
}
