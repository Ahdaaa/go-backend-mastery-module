# Rest API : Bank Transactions

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
docker run postgres16
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