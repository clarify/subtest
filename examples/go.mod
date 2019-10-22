module github.com/searis/subtest/examples

go 1.13

require (
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/searis/subtest v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	gopkg.in/thedevsaddam/gojsonq.v2 v2.3.0
)

replace github.com/searis/subtest => ./..
