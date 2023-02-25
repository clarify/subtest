module github.com/clarify/subtest/examples

go 1.13

require (
	github.com/clarify/subtest v0.0.0-00010101000000-000000000000
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/nsf/jsondiff v0.0.0-20190712045011-8443391ee9b6
	golang.org/x/crypto v0.1.0
	gopkg.in/thedevsaddam/gojsonq.v2 v2.3.0
)

replace github.com/clarify/subtest => ./..
