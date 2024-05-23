module github.com/bizflycloud/bizflyctl

go 1.16

require (
	github.com/bizflycloud/gobizfly v1.1.2
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.1
	gopkg.in/yaml.v2 v2.2.8
)

require github.com/go-openapi/strfmt v0.19.5 // indirect

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab
