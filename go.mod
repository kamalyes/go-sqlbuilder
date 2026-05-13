module github.com/kamalyes/go-sqlbuilder

go 1.25.0

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/kamalyes/go-logger v0.4.6
	github.com/kamalyes/go-toolbox v0.12.1-0.20260513145936-8b6f54d3138b
	github.com/lib/pq v1.10.9
	github.com/stretchr/testify v1.11.1
	google.golang.org/protobuf v1.36.11
	gorm.io/driver/mysql v1.6.0
	gorm.io/driver/sqlite v1.6.0
	gorm.io/gorm v1.30.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/kamalyes/go-jsonpath v0.0.0-20260129163507-0b67ed48bb28 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/grpc v1.79.3 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// 本地开发替换
// replace github.com/kamalyes/go-toolbox => ../go-toolbox

// replace github.com/kamalyes/go-logger => ../go-logger
