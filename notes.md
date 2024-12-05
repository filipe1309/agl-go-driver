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
> grant usage, select on sequence folders_id_seq to imersao;
> grant usage, select on sequence files_id_seq to imersao;
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

## Commands

### Users

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

List

```bash
go run cmd/cli/main.go users list
# OR
./bin/drive users list
```

Get

```bash
go run cmd/cli/main.go users get --id 1
# OR
./bin/drive users get -id 1
```

Update

```bash
go run cmd/cli/main.go users update --id 1 --name "John Doe New Name"
# OR
./bin/drive users update -id 1 -name "John Doe New Name"
```

Delete

```bash
go run cmd/cli/main.go users delete --id 2
# OR
./bin/drive users delete -id 2
```

### Folders

Insert

```bash
go run cmd/api/main.go
go run cmd/cli/main.go folders create --name "My Folder" --user 1
# OR
./bin/drive folders create -name "My Folder" -user 1
```

null types
```bash
go get gopkg.in/guregu/null.v4
```

Update

```bash
go run cmd/cli/main.go folders update --id 1 --name "My Folder New Name"
# OR
./bin/drive folders update -id 1 -name "My Folder New Name"
```

Delete

```bash
go run cmd/cli/main.go folders delete --id 2
# OR
./bin/drive folders delete -id 2
```

Files

Upload

```bash
go run cmd/api/main.go
go run cmd/cli/main.go files upload --filename ./internal/files/testdata/test-image-1.jpg
# OR
./bin/drive files upload -filename ./internal/files/testdata/test-image-1.jpg
```

## Worker

```bash
go run cmd/worker/main.go
```

RabbitMQ

Management:
http://localhost:15672
> guest:guest


### List

```bash
go run cmd/cli/main.go folders list
# OR
./bin/drive folders list
```

## Files

### Update

```bash
go run cmd/cli/main.go files update --id 1 --name "My File New Name"
# OR
./bin/drive files update -id 1 -name "My File New Name"
```

### Delete

```bash
go run cmd/cli/main.go files delete --id 2
# OR
./bin/drive files delete -id 2
```


# gRPC

```bash
cd proto/v1
protoc --go_out=users --go_opt=paths=source_relative user.proto
```
```

```bash
