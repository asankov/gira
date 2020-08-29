## Gira

[![Coverage Status](https://coveralls.io/repos/github/asankov/gira/badge.svg?branch=main&service=github)](https://coveralls.io/github/asankov/gira?branch=main)

Gira is like Jira, but for tracking your video games progress

### How to run
First, run a PostgreSQL database:
```
$ docker run --name postgres -e POSTGRES_PASSWORD=password -e POSTGRES_USER=antonsankov -p 5432:5432 -d postgres
```
Initialize the database:
```
$ docker exec postgres
# psql -U antonsankov
## CREATE DATABASE gira;
## \c gira;
## <content of init.sql>
```

Now, run the front-end and api services:
```
$ go run cmd/front-end/main.go
$ go run cmd/api/main.go -db_pass password
```

TODO: simplify this via Docker compose

### License
This work is licensed under MIT license. For more info see [LICENSE.md](LICENSE.md)
