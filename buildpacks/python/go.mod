module kn-fn/python-function-buildpack

go 1.19

require (
	github.com/buildpacks/libcnb v1.27.0
	github.com/golang/mock v1.6.0
	github.com/onsi/gomega v1.21.1
	github.com/paketo-buildpacks/libpak v1.63.0
	github.com/sclevine/spec v1.4.0
	k8s.io/utils v0.0.0-20220922133306-665eaaec4324
	kn-fn/buildpacks v0.0.0
	knative.dev/kn-plugin-func v0.34.1
)

require (
	github.com/BurntSushi/toml v1.2.0 // indirect
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20211112122917-428f8eabeeb3 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/cloudevents/sdk-go/v2 v2.10.1 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/creack/pty v1.1.18 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.3.1 // indirect
	github.com/go-git/go-git/v5 v5.4.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/h2non/filetype v1.1.3 // indirect
	github.com/heroku/color v0.0.6 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/whilp/git-urls v1.0.0 // indirect
	github.com/xanzy/ssh-agent v0.3.1 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.23.0 // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/net v0.0.0-20220826154423-83b083e8dc8b // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apimachinery v0.25.3 // indirect
	knative.dev/pkg v0.0.0-20220914154704-5f66ecf267fe // indirect
)

replace (
	cloud.google.com/go => cloud.google.com/go v0.100.2
	golang.org/x/net => golang.org/x/net v0.1.1-0.20221104162952-702349b0e862
	k8s.io/client-go => k8s.io/client-go v0.25.2
	kn-fn/buildpacks => ../common
)
