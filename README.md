# dbtesting
### installation
go get github.com/yqyzxd/dbtesting


### for mysql usage
```go
	func TestMain(m *testing.M) {
		os.Exit(dbtesting.RunDBInDocker(m, &dbtesting.Config{
			Image:         "mysql:5.6",
			User:          "root",
			Password:      "123456",
			DBName:        "mu_test_db",
			DB:            dbtesting.Mysql,
			ContainerPort: "3306/tcp",
		}))
	}
	func NewEngine() (*xorm.Engine, error) {
		if dbtesting.ConnURI == "" {
			return nil, fmt.Errorf("conn uri is nil")
		}
		return xorm.NewEngine("mysql", dbtesting.ConnURI)
	}
```

for mongo usage:
```go
	func TestMain(m *testing.M) {
		os.Exit(dbtesting.RunDBInDocker(m, &dbtesting.Config{
				Image:         "mongo",
				User:          "admin",
				Password:      "123456",
				DB:            dbtesting.Mongo,
				ContainerPort: "27017/tcp",
			}))
	}
	func NewClient(c context.Context) (*mongo.Client, error) {
		if dbtesting.ConnURI == "" {
			return nil, fmt.Errorf("conn uri is nil")
		}
		return mongo.Connect(c, options.Client().ApplyURI(dbtesting.ConnURI))
	}
```
