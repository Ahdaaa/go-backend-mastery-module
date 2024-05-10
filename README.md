# go-backend-mastery-module

## Database CDM/PDM Setup

1. Im using [dbdiagram.io](https://dbdiagram.io)  as a quick and simple tool for creating conceptual database model, design the database as you wish.

2. Convert the ready design to postgresql format.

## PostgreSQL Docker Setup

1. Make sure docker desktop is running, then pull postgre image

```sh
docker pull postgres
```

2. Create and run container with postgres image,

```sh
docker run --name postgres16 -p 5433:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=insertpw -d postgres:16
```

**notes** -> 5433:5432 means that if we want to connect to port 5433 from host machine, it will map the connection to port 5432 in container machine. Im not using 5432:5432, since i already have postgresql in my host machine that runs on that port.

Alternatively, you can start container with,

```sh
docker start postgres16
```

To locate running container, use ``docker ps``

3. To open postgres terminal use,

```sh
docker exec -it postgres16 psql -U postgres
```

Based on the image documentation, 

> The PostgreSQL image sets up trust authentication locally so you may notice a password is not required when connecting from localhost

4. To see logs for a container, use

```sh
docker logs postgres16
```

## Managing Database with Tableplus (Optional)

1. Download tableplus from [](https://dbdiagram.io) https://tableplus.com/download

2. Create new postgresql connection,

![Alt text](./images/image.png)

3. Open and run the ``bank.sql`` from the previous step.

## Migrating Database

1. With [scoop](https://scoop.sh), install golang-migrate.

```sh
scoop install migrate
```

2. Create new migration

```sh
migrate create -ext sql -dir db/migration -seq init_schema
```

There will be 2 new files in db/migration, up and down migrate. 

``up migrate`` : old db -> run sequentially by the order of prefix version -> new db

``down migrate`` : reverse order of up migration

3. Run postgres16 docker container shell cli by using,

```sh
docker exec -it postgres16 /bin/sh
```

4. Create new database, connect to psql with user postgres, then go to simple_bank db

```sh
createdb --username=postgres --owner=postgres simple_bank
psql -U postgres
\c simple_bank
```

Alternatively, you can createdb without going to shell cli by using,

```sh
docker exec -it postgres16 createdb --username=postgres --owner=postgres simple_bank
```

and access the db console by using,

```sh
docker exec -it postgres16 psql -U postgres simple_bank
```

5. Create Makefile (createdb and dropdb) to help setting up the project easily, to use it you need to install make using scoop.

```sh
scoop install make
```

then, lets say you want to start the container and create database in a new device, then use

```sh
make postgres
make createdb
```

Now, with tableplus, you should be able to connect to simple_bank database with port 5433.

6. Migrate the database from previous schema by using,

```sh
migrate -path db/migration -database "postgresql://postgres:Anaana123@localhost:5433/simple_bank?sslmode=disable" -verbose up
```

If there are no errors, try to refresh your tableplus, you should be able to see the created tables.

7. Add the migrate to Makefile.

Use ``make migrateup`` to create all the tables, and use ``make migratedown`` to drop all of the tables.

## Things to Consider - Golang CRUD

1. Using low level standard library - [database/sql](https://pkg.go.dev/database/sql)

For example:

```go
id := 123
var username string
var created time.Time
err := db.QueryRowContext(ctx, "your sql query", id).Scan(&username, &created)
```

This may look Very fast & straightforward but we have to manually mapping all the sql fields to variables, easy to make mistakes and not caught until runtime.

2. Using Object Relational Mapping Library for Golang - [GORM](https://gorm.io/docs/index.html).

In [GORM](https://gorm.io/docs/index.html), the CRUD Functions are already implemented, this can make the production code shorter. So, we must learn to understand the library to do a specific assignment, such as complex queries. Also, Based on benchmarks on the internet, GORM may run slowly on high loads.

3. Using Middleway Approach - [SQLX Library](https://pkg.go.dev/github.com/jmoiron/sqlx)

SQLX runs nearly as fast as a standard library and very easy to use, the fields mapping are done via query text or struct tags. However, the code that we need to write are relatively longer than our previous method AND the errors of queries wont be occur until runtime.

4. Using [SQLC Library](https://sqlc.dev/)

Just like database/sql, this library runs very fast and easy to use, the most unique thing is that we just need to write sql queries then the golang code will be generated. Also, the library will catch sql query errors before generating the codes. 

This library will be valuable to increase efficiency and maintain code consistency. However, as a developer we need to complement and understand the generated code, we cant be over-reliance on generators without understanding it.

## Getting Started with SQLC

1. Go to the documentation page and install the library with ``go install``.

```sh
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

2. Run ``sqlc init`` inside your project directory, there will be a new file called ``sqlc.yaml``

Using documentation, fill the file with this lines of code,

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "./db/query/"
    schema: "./db/migration"
    gen:
      go:
        package: "db"
        out: "./db/sqlc"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_exact_table_names: true
```

Based on that yaml file, we will use ``postgresql`` engine and we want to store all our sql queries inside the directory ``./db/query`` while migration schemas in ``./db/migration``. Then, we will name our go package as ``db`` then store the generated go codes inside ``./db/sqlc``.

3. Create sql insertion example file inside its directory.

```sql
-- name: CreateAccount :one
INSERT INTO accounts (
  owner, 
  balance,
  currency
) VALUES (
  $1, $2, $3
) RETURNING *;
```

4. Generate sqlc and add its method inside Makefile.

``sqlc generate``

or 

``make sqlc`` if you already put it inside makefile.

After that, you will see 3 new generated go files inside ``./db/sqlc``

``models.go`` : contains struct definition of database model

``db.go`` : contains dbtx interface, the functions inside it will allows us to freely use either db or transaction to execute queries.

``account.sql.go`` : contains all related code with inserting new account row into database.

Then, you might need to install some needed import packages by using ``go get``, make sure you already initialize go module.

Also, when we're working with sqlc, we should not modify the generated file because everytime we run ``make sqlc``, all those files will be regenerated then our changes will be gone. To solve this, make sure to create a new go files to modify its generated code.

5. Add Read, Update, and Delete operation to ``account.sql`` then try to regenerate sqlc.

Notes:

``:one``  : will return one row
``:many`` : will return >1 row
``:exec`` : will not return anything

**read**

```sql
-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;
```

**update**

```sql
-- name: UpdateAccount :one
UPDATE accounts 
SET balance = $2
WHERE id = $1 
RETURNING *;
```

**delete**

```sql
-- name: DeleteAccount :exec
DELETE FROM accounts 
WHERE id = $1;
```

Then run ``make sqlc`` and the generated code will be updated.

## Go CRUD Testing

Please see ``account_test.go`` for further details.

## Go DB Transactions

In a simple bank scenario, here is an example of transfering balance (steps):

1. Create transfer record with amount=x
2. Create entry for account1 with amount -= x
3. Create entry for account2 with amount += x
4. Subtract x from account1 and Add x to account2

DB Transactions is important to isolate programs that accesses the db concurrently. Refer to ACID property for further details.
For Go implementation, refer to ``store.go`` and ``store_test.go``