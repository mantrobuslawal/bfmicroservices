-- 003-grant-service-permissions.sql
--
-- Grants least-privilege access for bfstore service database users.
--
-- Local development note:
-- This file is intended to be run during first-time MySQL container
-- initialisation via /docker-entrypoint-initdb.d.
--
-- It assumes:
--   001-create-service-databases.sql has already created the databases.
--   002-create-service-users.sql has already created the service users.
--
-- Keep it boring where production matters:
-- each service user should only access its own database.

GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, DROP
ON bfstore_catalog.*
TO 'bfstore_catalog_user'@'%';

GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, DROP
ON bfstore_inventory.*
TO 'bfstore_inventory_user'@'%';

GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, DROP
ON bfstore_basket.*
TO 'bfstore_basket_user'@'%';

GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, DROP
ON bfstore_order.*
TO 'bfstore_order_user'@'%';

GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, DROP
ON bfstore_payment.*
TO 'bfstore_payment_user'@'%';

GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, DROP
ON bfstore_shipping.*
TO 'bfstore_shipping_user'@'%';

GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, DROP
ON bfstore_notification.*
TO 'bfstore_notification_user'@'%';

FLUSH PRIVILEGES;
