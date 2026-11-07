module github.com/kamalyes/go-sqlbuilder

go 1.25.0

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/kamalyes/go-toolbox v0.19.3
	github.com/kamalyes/go-toolbox/pkg/protox v0.19.3
	github.com/lib/pq v1.10.9
	github.com/stretchr/testify v1.12.1
	google.golang.org/protobuf v1.36.12
	gorm.io/driver/mysql v1.6.0
	gorm.io/driver/sqlite v1.6.0
	gorm.io/gorm v1.30.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/text v0.34.0 // indirect
)

// 本地开发替换
// replace github.com/kamalyes/go-argus => ../go-argus

// replace github.com/kamalyes/go-toolbox => ../go-toolbox

// replace github.com/kamalyes/go-logger => ../go-logger
