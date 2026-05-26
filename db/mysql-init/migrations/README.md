# `db/mysql-init`

## 1. Purpose

This directory contains MySQL initialisation scripts for the local bfstore development environment.

These scripts create the local MySQL schemas, database users, and permissions required by the service-owned database model.

This directory is intended to support local development, CI integration tests, and client review of bfstore’s database boundary approach.

---

## 2. Scope

This directory is responsible for local initialisation only.

It may contain:

```text
schema creation scripts
service database user creation scripts
local permission grants
safe local seed prerequisites
```

It should not contain:

```text
service table migrations
production credentials
real customer data
cloud-specific database provisioning
manual production database changes
```

Service table migrations belong under:

```text
db/<service>/migrations/
```

Cloud database provisioning belongs in the platform infrastructure repository.

---

## 3. Expected Files

Recommended layout:

```text
db/mysql-init/
├── README.md
├── 001-create-databases.sql
├── 002-create-users.sql
├── 003-grant-permissions.sql
└── 004-create-local-seed-prerequisites.sql
```

The fourth file is optional and should only contain safe local setup data.

---

## 4. Schemas Created

Initial local schemas should include:

```text
bfstore_catalog
bfstore_inventory
bfstore_basket
bfstore_order
bfstore_payment
bfstore_shipping
bfstore_notification
```

Deferred schemas may include:

```text
bfstore_auth
bfstore_customer
bfstore_review
bfstore_search
bfstore_recommendation
```

The initial implementation should focus on the checkout vertical slice before creating every possible schema.

---

## 5. Service-Owned Database Principle

Each service owns its own logical schema.

Rules:

```text
catalog-service owns bfstore_catalog
inventory-service owns bfstore_inventory
basket-service owns bfstore_basket
order-service owns bfstore_order
payment-service owns bfstore_payment
shipping-service owns bfstore_shipping
notification-service owns bfstore_notification
```

A service must not read or write another service’s schema.

---

## 6. Database Users

Each service should have its own database user.

Recommended local users:

```text
bfstore_catalog_user
bfstore_inventory_user
bfstore_basket_user
bfstore_order_user
bfstore_payment_user
bfstore_shipping_user
bfstore_notification_user
```

Optional deferred users:

```text
bfstore_auth_user
bfstore_customer_user
bfstore_review_user
bfstore_search_user
bfstore_recommendation_user
```

---

## 7. Permission Model

Runtime users should receive only the permissions needed by their owning service.

Typical local runtime permissions:

```sql
SELECT, INSERT, UPDATE, DELETE
```

Migration users may require additional permissions such as:

```sql
CREATE, ALTER, DROP, INDEX
```

For local development, a simplified permission model may be acceptable, but the intended production-style model should remain clear.

---

## 8. Example Permission Intent

Example:

```sql
GRANT SELECT, INSERT, UPDATE, DELETE
ON bfstore_order.*
TO 'bfstore_order_user'@'%';
```

Forbidden:

```sql
GRANT ALL PRIVILEGES
ON *.*
TO 'bfstore_order_user'@'%';
```

The Order Service user should not have access to Catalogue, Inventory, Payment, Shipping, or Notification schemas.

---

## 9. Local Development Usage

Expected local flow:

```sh
make dev-up
make migrate-up
make seed-local
make smoke-test
```

The MySQL container should run these init scripts when the local database is first created.

If the database volume already exists, init scripts may not automatically rerun. Developers may need to reset local volumes deliberately.

---

## 10. Safety Rules

These scripts must not include:

```text
real passwords
production credentials
real customer data
raw payment data
personal addresses
tokens
cloud secrets
```

Use safe dummy values for local development only.

---

## 11. Relationship to Migrations

This directory creates schemas and users.

It does not create service tables.

Service tables are created by service-owned migrations:

```text
db/catalog/migrations/
db/inventory/migrations/
db/basket/migrations/
db/order/migrations/
db/payment/migrations/
db/shipping/migrations/
db/notification/migrations/
```

This separation keeps environment bootstrap separate from service schema evolution.

---

## 12. Review Checklist

Before approving changes to this directory, check:

```text
Are schemas service-owned?
Are database users service-specific?
Are permissions least-privilege where practical?
Are production secrets excluded?
Are local dummy credentials clearly local-only?
Are service table changes placed in migrations instead?
Does the setup support the checkout vertical slice?
```

---

## 13. Related Documents

```text
docs/data/data-ownership.md
docs/data/service-database-design.md
docs/data/mysql-standards.md
docs/data/migrations.md
docs/architecture/service-boundaries.md
```

---

## 14. Summary

`db/mysql-init` bootstraps the local MySQL environment for bfstore.

It should make service-owned database boundaries visible from the start by creating separate schemas and service-specific users.
