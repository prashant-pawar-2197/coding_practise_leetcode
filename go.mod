module standByCB

go 1.19

require com.Go/checkpack v0.0.0

require (
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	golang.org/x/exp v0.0.0-20230515195305-f3d0a9c9a5cc // indirect
)

require (
	github.com/antlr/antlr4/runtime/Go/antlr v1.4.10 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/couchbase/gocb/v2 v2.6.0 // indirect
	github.com/couchbase/gocbcore/v10 v10.2.0 // indirect
	github.com/etherlabsio/healthcheck/v2 v2.0.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	golang.org/x/tools v0.4.1-0.20221208213631-3f74d914ae6d // indirect
	honnef.co/go/tools v0.4.2 // indirect
	ruleEngine v0.0.0
)

replace (
	com.Go/checkpack v0.0.0 => ./checkPack
	ruleEngine v0.0.0 => ./ruleEngine
)
