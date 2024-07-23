module github.com/kubeedge/keink

go 1.20

require (
	github.com/spf13/cobra v1.6.1
	github.com/spf13/pflag v1.0.5
	sigs.k8s.io/kind v0.0.0
)

require (
	github.com/google/safetext v0.0.0-20220905092116-b49f7bc46da2 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/alessio/shellescape v1.4.1 // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/tools v0.20.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace sigs.k8s.io/kind => github.com/kubeedge/kind v0.21.0-kubeedge1
