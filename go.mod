module github.com/containerssh/authintegration

go 1.14

require (
	github.com/containerssh/auth v0.9.6
	github.com/containerssh/geoip v0.9.4
	github.com/containerssh/http v0.9.9
	github.com/containerssh/log v0.9.13
	github.com/containerssh/metrics v0.9.8
	github.com/containerssh/service v0.9.3
	github.com/containerssh/sshserver v0.9.19
	github.com/containerssh/structutils v0.9.0
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/mattn/go-shellwords v1.0.11 // indirect
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	golang.org/x/sys v0.0.0-20210309074719-68d13333faf2 // indirect
)

// Fixes CVE-2020-9283
exclude (
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
)

// Fixes CVE-2020-14040
exclude (
	golang.org/x/text v0.3.0
	golang.org/x/text v0.3.1
	golang.org/x/text v0.3.2
)

// Fixes CVE-2019-11254
exclude (
	gopkg.in/yaml.v2 v2.2.0
	gopkg.in/yaml.v2 v2.2.1
	gopkg.in/yaml.v2 v2.2.2
	gopkg.in/yaml.v2 v2.2.3
	gopkg.in/yaml.v2 v2.2.4
	gopkg.in/yaml.v2 v2.2.5
	gopkg.in/yaml.v2 v2.2.6
	gopkg.in/yaml.v2 v2.2.7
)
