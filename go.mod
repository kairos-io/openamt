module provider-amt

go 1.19

require (
	github.com/kairos-io/kairos-sdk v0.0.1
	github.com/mudler/go-pluggable v0.0.0-20230126220627-7710299a0ae5
	github.com/sirupsen/logrus v1.9.0
	github.com/stretchr/testify v1.8.2
	gopkg.in/yaml.v2 v2.4.0
	rpc v0.0.0-00010101000000-000000000000
)

require (
	github.com/chuckpreslar/emission v0.0.0-20170206194824-a7ddd980baf9 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace rpc => github.com/open-amt-cloud-toolkit/rpc-go v0.0.0-20230306150617-96cc235bf732
