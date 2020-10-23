## Gira

[![Coverage Status](https://coveralls.io/repos/github/asankov/gira/badge.svg?branch=main&service=github)](https://coveralls.io/github/asankov/gira?branch=main)

Gira is like Jira, but for tracking your video games progress

### How to run
Via Docker compose:
```
 $ docker-compose -f docker/docker-compose.yml up
```
Then, initialize the database:
```
$ go get -u github.com/pressly/goose/cmd/goose
$ goose -dir sql/ postgres 'host=localhost port=5432 user=gira dbname=gira password=password sslmode=disable' up
```
Now you should be able to open the browser on [localhost:4000](localhost:4000) and see the UI.

### License
This work is licensed under MIT license. For more info see [LICENSE.md](LICENSE.md)
