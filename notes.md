# Notes

> notes taken during the course

```sh
go mod init github.com/filipe1309/agl-go-driver
```

Executing SQL scripts:

```bash
$ docker run --network=host -it --rm -v $(pwd):/tmp postgres bash
$ cd /tmp/scripts/database
```

```bash
$ psql -h imersao-postgres -U postgres
# On Mac: psql --host docker.for.mac.host.internal -U postgres
> create database imersao;
> \l
> exit
$ psql -h docker.for.mac.host.internal -U postgres -d imersao < users.sql
$ psql -h docker.for.mac.host.internal -U postgres -d imersao < folders.sql
$ psql -h docker.for.mac.host.internal -U postgres -d imersao < files.sql

$ psql -h docker.for.mac.host.internal -U postgres
> \c imersao
> \dt
```

# Worker

```bash
go mod tidy
```

# API

## Testing

```bash
go get github.com/DATA-DOG/go-sqlmock
```

```bash
go test ./internal/users/... -v
```

sql.NullInt64

## Tests

testify

```bash
go get github.com/stretchr/testify
```

```bash
go get github.com/stretchr/testify/suite
```


## Auth

```bash
go get github.com/golang-jwt/jwt/v5
```

## Debug

```bash
make up-all # db and queue
```

```bash
$ docker run --network=host -it --rm -v $(pwd):/tmp postgres bash
$ cd /tmp/scripts/database
$ psql -h docker.for.mac.host.internal -U postgres

> create database imersao;
> \l
> create user imersao with encrypted password '1234';
> grant all privileges on database imersao to imersao;
> grant all privileges on table users to imersao;
> grant all privileges on table folders to imersao;
> grant all privileges on table files to imersao;
> grant usage, select on sequence users_id_seq to imersao;
> exit
$ psql -h docker.for.mac.host.internal -U postgres -d imersao < users.sql
$ psql -h docker.for.mac.host.internal -U postgres -d imersao < folders.sql
$ psql -h docker.for.mac.host.internal -U postgres -d imersao < files.sql

$ psql -h docker.for.mac.host.internal -U postgres
> \c imersao
> \dt
```

```bash
chmod +x scripts/shell/env.sh
```

```bash
. ./scripts/shell/env.sh
echo $DB_HOST
```

Insert

```bash
go run cmd/api/main.go
go run cmd/cli/main.go users create --name John --login johndoe --pass 123456
# OR
./bin/drive users create -name John -login johndoe -pass 123456
```

```bash
# $ docker run --network=host -it --rm -v $(pwd):/tmp postgres bash
$ docker exec -it imersao-postgres bash
$ psql -h docker.for.mac.host.internal -U postgres
> \c imersao
> \dt;
> select * from users;
````


Auth

```bash
go run cmd/api/main.go
go run cmd/cli/main.go auth --user johndoe --pass 123456
# OR
./bin/drive auth -user johndoe -pass 123456
```

```bash
