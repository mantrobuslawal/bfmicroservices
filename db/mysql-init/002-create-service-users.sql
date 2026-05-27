-- bfstore local MySQL service users
--
-- Purpose:
-- Create one database user per service for local development.
--
-- Security rule:
-- Each service user receives privileges only for its own database.
--
-- Local development note:
-- Passwords here are intentionally simple and must not be reused outside
-- local development.

CREATE USER IF NOT EXISTS 'bfstore_catalog_user'@'%' IDENTIFIED BY 'bfstore_catalog_password';
CREATE USER IF NOT EXISTS 'bfstore_basket_user'@'%' IDENTIFIED BY 'bfstore_basket_password';
CREATE USER IF NOT EXISTS 'bfstore_inventory_user'@'%' IDENTIFIED BY 'bfstore_inventory_password';
CREATE USER IF NOT EXISTS 'bfstore_order_user'@'%' IDENTIFIED BY 'bfstore_order_password';
CREATE USER IF NOT EXISTS 'bfstore_payment_user'@'%' IDENTIFIED BY 'bfstore_payment_password';
CREATE USER IF NOT EXISTS 'bfstore_shipping_user'@'%' IDENTIFIED BY 'bfstore_shipping_password';
CREATE USER IF NOT EXISTS 'bfstore_notification_user'@'%' IDENTIFIED BY 'bfstore_notification_password';

GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, DROP
ON bfstore_catalog.* TO 'bfstore_catalog_user'@'%';

GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, DROP
ON bfstore_basket.* TO 'bfstore_basket_user'@'%';

GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, DROP
ON bfstore_inventory.* TO 'bfstore_inventory_user'@'%';

GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, DROP
ON bfstore_order.* TO 'bfstore_order_user'@'%';

GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, DROP
ON bfstore_payment.* TO 'bfstore_payment_user'@'%';

GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, DROP
ON bfstore_shipping.* TO 'bfstore_shipping_user'@'%';

GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, DROP
ON bfstore_notification.* TO 'bfstore_notification_user'@'%';

FLUSH PRIVILEGES;
