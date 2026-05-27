# MySQL Initialisation

## 1. Purpose

This directory contains local MySQL initialisation scripts for bfstore.

These scripts are mounted into the MySQL container by Docker Compose:

```text
./db/mysql-init:/docker-entrypoint-initdb.d:ro
```

They run when the MySQL data volume is first created.

---

## 2. Files

```text
001-create-service-databases.sql
002-create-service-users.sql
```

---

## 3. Service-Owned Databases

bfstore uses one database per service.

Initial local databases:

```text
bfstore_catalog
bfstore_basket
bfstore_inventory
bfstore_order
bfstore_payment
bfstore_shipping
bfstore_notification
```

Each service owns its own database and must not directly access another service's database.

Services communicate through:

```text
gRPC for commands and queries
Kafka for events
```

---

## 4. Service Users

Each service has a dedicated MySQL user.

Examples:

```text
bfstore_catalog_user
bfstore_order_user
bfstore_payment_user
```

Each user receives privileges only on its own database.

This supports least-privilege access even in local development.

---

## 5. Local Development Note

The passwords in these scripts are for local development only.

They must not be used in production, staging, shared cloud environments, or client infrastructure.

---

## 6. Resetting Local MySQL

To force the init scripts to run again, remove the local MySQL volume:

```sh
make down-volumes
make up
```

Warning: this deletes local database state.

---

## 7. Client-Facing Rationale

These scripts demonstrate:

```text
microservice data ownership
database isolation
local repeatability
least-privilege access
clear operational setup
```

This makes the project easier to review and easier to explain in client conversations.
