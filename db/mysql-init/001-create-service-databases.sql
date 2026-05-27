-- bfstore local MySQL database initialisation
--
-- Purpose:
-- Create one database per service for local development.
--
-- Design rule:
-- Each microservice owns its own database/schema.
-- Services must not directly read or write another service's database.

CREATE DATABASE IF NOT EXISTS bfstore_catalog
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

CREATE DATABASE IF NOT EXISTS bfstore_basket
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

CREATE DATABASE IF NOT EXISTS bfstore_inventory
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

CREATE DATABASE IF NOT EXISTS bfstore_order
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

CREATE DATABASE IF NOT EXISTS bfstore_payment
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

CREATE DATABASE IF NOT EXISTS bfstore_shipping
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

CREATE DATABASE IF NOT EXISTS bfstore_notification
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;
