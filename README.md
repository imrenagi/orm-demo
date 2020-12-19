# ORM DEMO

Just wanted to make self-documentation about tradeoff between using ORM like (gorm) and raw SQL query using `database/sql` package in golang.

Summary:
* Use a simple model which has multiple relationships (has one, belongs to, has many and many-to-many).
* Create two repositories: one that is using `*gorm.DB` and another one that is using `*sql.DB`. Those two data repositories are demonstrating the API differences between `gorm` and `sql`. You can see how big differences `gorm` makes in simplifying the database queries and struct construction.
* You can refer to the unit test to see that those two repositories are fetching the same number of data and fields.

Below is the benchmark of the experiments. I dont want to explain, look at these data below by yourself and try to make sense of it by looking at the source code for each package.


## Benchmark with SQLMock

### Gorm
```
pkg: github.com/imrenagi/orm-demo/db/orm/
====
BenchmarkFindAll_WithMock-16                       10000            172983 ns/op
BenchmarkFindByIDWithJoin_WithMock-16              10000            220720 ns/op
BenchmarkFindCompletedByID_WithMock-16              3914           1142499 ns/op
```

### Raw SQL
```
pkg: github.com/imrenagi/orm-demo/db/sql/ 
===
BenchmarkFindAll_WithMock-16                       10000            157948 ns/op
BenchmarkFindByIDWithJoin_WithMock-16              10000            183222 ns/op
BenchmarkFindCompletedByID_WithMock-16              7498           1062009 ns/op
```

## Benchmark with Postgres Running Locally

### Gorm
```
pkg: github.com/imrenagi/orm-demo/db/orm
====
BenchmarkFindAll_WithDB-16                           538           1959670 ns/op
BenchmarkFindByIDWithJoin_WithDB-16                  618           1755951 ns/op
BenchmarkFindCompletedByID_WithDB-16                 160           6792193 ns/op
```

### Raw SQL
```
pkg: github.com/imrenagi/orm-demo/db/sql
===
BenchmarkFindAll_WithDB-16                           699           1672481 ns/op
BenchmarkFindByIDWithJoin_WithDB-16                  615           1756467 ns/op
BenchmarkFindCompletedByID_WithDB-16                 618           1743564 ns/op
```