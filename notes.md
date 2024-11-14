# Notes

> notes taken during the course

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
