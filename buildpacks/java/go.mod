module kn-fn/java-function-buildpack

go 1.19

require (
	github.com/buildpacks/libcnb v1.27.0
	github.com/onsi/gomega v1.21.1
	github.com/paketo-buildpacks/libpak v1.63.0
	github.com/sclevine/spec v1.4.0
	kn-fn/buildpacks v0.0.0
	knative.dev/kn-plugin-func v0.19.0
)

require (
	github.com/BurntSushi/toml v1.2.0 // indirect
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/Microsoft/go-winio v0.6.0 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20220930113650-c6815a8c17ad // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/cloudflare/circl v1.2.0 // indirect
	github.com/creack/pty v1.1.18 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.3.1 // indirect
	github.com/go-git/go-git/v5 v5.4.2 // indirect
	github.com/gobuffalo/here v0.6.7 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/h2non/filetype v1.1.3 // indirect
	github.com/heroku/color v0.0.6 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/markbates/pkger v0.17.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/xanzy/ssh-agent v0.3.2 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	golang.org/x/crypto v0.0.0-20221010152910-d6f0a8c073c2 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/net v0.0.0-20221004154528-8021a29435af // indirect
	golang.org/x/sys v0.0.0-20221010170243-090e33056c14 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.12 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apimachinery v0.25.2 // indirect
	k8s.io/utils v0.0.0-20220922133306-665eaaec4324 // indirect
)

replace kn-fn/buildpacks => ../common
