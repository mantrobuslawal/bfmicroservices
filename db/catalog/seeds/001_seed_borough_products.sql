-- bfstore Borough Furniture Store seed data
--
-- Purpose:
-- Provide memorable local development product data for Catalogue Service.
--
-- The data is intentionally playful, but the schema usage is serious:
-- - categories
-- - products
-- - variants
-- - category-scoped attributes
-- - product attribute values
-- - product images

USE bfstore_catalog;

-- -------------------------------------------------------------------
-- Categories
-- -------------------------------------------------------------------

INSERT INTO categories (
  category_id,
  parent_category_id,
  name,
  slug,
  description,
  status,
  display_order
) VALUES
  ('11111111-1111-1111-1111-111111111111', NULL, 'Developer Homeware', 'developer-homeware', 'Homeware for engineers who have strong opinions about tabs, spaces, and distributed systems.', 'active', 1),
  ('22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', 'Lighting', 'lighting', 'Desk lamps and ambient lighting for late-night debugging sessions.', 'active', 10),
  ('33333333-3333-3333-3333-333333333333', '11111111-1111-1111-1111-111111111111', 'Soft Furnishings', 'soft-furnishings', 'Cushions, blankets, and textile goods with excellent runtime comfort.', 'active', 20),
  ('44444444-4444-4444-4444-444444444444', '11111111-1111-1111-1111-111111111111', 'Wall Art', 'wall-art', 'Decorative pieces for tasteful homes and over-engineered office corners.', 'active', 30),
  ('55555555-5555-5555-5555-555555555555', '11111111-1111-1111-1111-111111111111', 'Secure Storage', 'secure-storage', 'Lockboxes and storage for people who read the security notes.', 'active', 40),
  ('66666666-6666-6666-6666-666666666666', '11111111-1111-1111-1111-111111111111', 'Rugs', 'rugs', 'Floor coverings for pathfinding, graph traversal, and general comfort.', 'active', 50)
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  description = VALUES(description),
  status = VALUES(status),
  display_order = VALUES(display_order);

-- -------------------------------------------------------------------
-- Attribute definitions
-- -------------------------------------------------------------------

INSERT INTO product_attribute_definitions (
  attribute_id,
  category_id,
  code,
  display_name,
  description,
  data_type,
  unit,
  is_required,
  is_filterable,
  is_variant_defining,
  status
) VALUES
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', '22222222-2222-2222-2222-222222222222', 'bulb_type', 'Bulb Type', 'Supported bulb type for the lamp.', 'option', NULL, TRUE, TRUE, FALSE, 'active'),
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2', '22222222-2222-2222-2222-222222222222', 'max_wattage', 'Maximum Wattage', 'Maximum supported bulb wattage.', 'number', 'W', TRUE, TRUE, FALSE, 'active'),
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa3', '33333333-3333-3333-3333-333333333333', 'fabric_type', 'Fabric Type', 'Primary fabric used for the item.', 'option', NULL, TRUE, TRUE, FALSE, 'active'),
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa4', '33333333-3333-3333-3333-333333333333', 'care_instructions', 'Care Instructions', 'Recommended cleaning or care instructions.', 'string', NULL, FALSE, FALSE, FALSE, 'active'),
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa5', '44444444-4444-4444-4444-444444444444', 'wall_art_size', 'Wall Art Size', 'Display size for wall art.', 'option', NULL, TRUE, TRUE, TRUE, 'active'),
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa6', '55555555-5555-5555-5555-555555555555', 'security_rating', 'Security Rating', 'Borough internal fun security rating.', 'option', NULL, TRUE, TRUE, FALSE, 'active'),
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa7', '55555555-5555-5555-5555-555555555555', 'lock_type', 'Lock Type', 'Type of lock mechanism.', 'option', NULL, TRUE, TRUE, FALSE, 'active'),
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa8', '66666666-6666-6666-6666-666666666666', 'rug_shape', 'Rug Shape', 'Shape of the rug.', 'option', NULL, TRUE, TRUE, TRUE, 'active'),
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa9', '66666666-6666-6666-6666-666666666666', 'pile_height', 'Pile Height', 'Approximate rug pile height.', 'number', 'mm', FALSE, TRUE, FALSE, 'active')
ON DUPLICATE KEY UPDATE
  display_name = VALUES(display_name),
  description = VALUES(description),
  data_type = VALUES(data_type),
  unit = VALUES(unit),
  is_required = VALUES(is_required),
  is_filterable = VALUES(is_filterable),
  is_variant_defining = VALUES(is_variant_defining),
  status = VALUES(status);

-- -------------------------------------------------------------------
-- Attribute options
-- -------------------------------------------------------------------

INSERT INTO product_attribute_options (
  option_id,
  attribute_id,
  value,
  display_name,
  display_order,
  status
) VALUES
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb001', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 'e27', 'E27', 1, 'active'),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb002', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 'led_integrated', 'Integrated LED', 2, 'active'),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb003', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa3', 'cotton', 'Cotton', 1, 'active'),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb004', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa3', 'woven_polyester', 'Woven Polyester', 2, 'active'),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb005', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa5', 'medium', 'Medium', 1, 'active'),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb006', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa5', 'large', 'Large', 2, 'active'),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb007', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa6', 'secure', 'Secure', 1, 'active'),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb008', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa6', 'super_secure', 'Super Secure', 2, 'active'),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb009', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa7', 'keypad', 'Keypad', 1, 'active'),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb010', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa7', 'biometric', 'Biometric', 2, 'active'),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb011', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa8', 'rectangle', 'Rectangle', 1, 'active'),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbb012', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa8', 'runner', 'Runner', 2, 'active')
ON DUPLICATE KEY UPDATE
  display_name = VALUES(display_name),
  display_order = VALUES(display_order),
  status = VALUES(status);

-- -------------------------------------------------------------------
-- Products
-- -------------------------------------------------------------------

INSERT INTO products (
  product_id,
  category_id,
  name,
  slug,
  description,
  brand,
  status,
  base_price_minor,
  currency_code
) VALUES
  ('cccccccc-cccc-cccc-cccc-cccccccc0001', '22222222-2222-2222-2222-222222222222', 'Gopher Desk Lamp', 'gopher-desk-lamp', 'A cheerful blue Gopher desk lamp for debugging after sunset.', 'Borough Originals', 'active', 4599, 'GBP'),
  ('cccccccc-cccc-cccc-cccc-cccccccc0002', '33333333-3333-3333-3333-333333333333', 'Gopher Cushion Set', 'gopher-cushion-set', 'Soft Gopher-themed cushions for distributed lounging.', 'Borough Originals', 'active', 3299, 'GBP'),
  ('cccccccc-cccc-cccc-cccc-cccccccc0003', '44444444-4444-4444-4444-444444444444', 'Rob Pike Wall Tapestry', 'rob-pike-wall-tapestry', 'A tasteful wall tapestry inspired by Go culture and engineering simplicity.', 'Borough Gallery', 'active', 6999, 'GBP'),
  ('cccccccc-cccc-cccc-cccc-cccccccc0004', '55555555-5555-5555-5555-555555555555', 'Rivest Super-Secure Lockbox', 'rivest-super-secure-lockbox', 'A cryptography-themed lockbox for secrets, snacks, and serious-looking stationery.', 'Borough Secure', 'active', 8999, 'GBP'),
  ('cccccccc-cccc-cccc-cccc-cccccccc0005', '66666666-6666-6666-6666-666666666666', 'Dijkstra Pathfinding Rug', 'dijkstra-pathfinding-rug', 'A graph-inspired rug for finding the shortest path from sofa to snacks.', 'Borough Floors', 'active', 11999, 'GBP'),
  ('cccccccc-cccc-cccc-cccc-cccccccc0006', '33333333-3333-3333-3333-333333333333', 'Grace Hopper Debugging Blanket', 'grace-hopper-debugging-blanket', 'A cosy debugging blanket for chilly incident reviews and warm retrospectives.', 'Borough Originals', 'active', 5499, 'GBP')
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  description = VALUES(description),
  brand = VALUES(brand),
  status = VALUES(status),
  base_price_minor = VALUES(base_price_minor),
  currency_code = VALUES(currency_code);

-- -------------------------------------------------------------------
-- Variants
-- -------------------------------------------------------------------

INSERT INTO product_variants (
  variant_id,
  product_id,
  sku,
  variant_name,
  status,
  price_minor,
  currency_code
) VALUES
  ('dddddddd-dddd-dddd-dddd-dddddddd0001', 'cccccccc-cccc-cccc-cccc-cccccccc0001', 'BFS-LAMP-GOPHER-BLUE', 'Blue Gopher Lamp', 'active', 4599, 'GBP'),
  ('dddddddd-dddd-dddd-dddd-dddddddd0002', 'cccccccc-cccc-cccc-cccc-cccccccc0002', 'BFS-CUSH-GOPHER-SET2', 'Set of 2 Cushions', 'active', 3299, 'GBP'),
  ('dddddddd-dddd-dddd-dddd-dddddddd0003', 'cccccccc-cccc-cccc-cccc-cccccccc0003', 'BFS-WALL-PIKE-MED', 'Medium Tapestry', 'active', 6999, 'GBP'),
  ('dddddddd-dddd-dddd-dddd-dddddddd0004', 'cccccccc-cccc-cccc-cccc-cccccccc0003', 'BFS-WALL-PIKE-LRG', 'Large Tapestry', 'active', 8999, 'GBP'),
  ('dddddddd-dddd-dddd-dddd-dddddddd0005', 'cccccccc-cccc-cccc-cccc-cccccccc0004', 'BFS-LOCK-RIVEST-KEYPAD', 'Keypad Lockbox', 'active', 8999, 'GBP'),
  ('dddddddd-dddd-dddd-dddd-dddddddd0006', 'cccccccc-cccc-cccc-cccc-cccccccc0005', 'BFS-RUG-DIJKSTRA-RECT', 'Rectangle Rug', 'active', 11999, 'GBP'),
  ('dddddddd-dddd-dddd-dddd-dddddddd0007', 'cccccccc-cccc-cccc-cccc-cccccccc0006', 'BFS-BLANKET-HOPPER-STD', 'Standard Blanket', 'active', 5499, 'GBP')
ON DUPLICATE KEY UPDATE
  variant_name = VALUES(variant_name),
  status = VALUES(status),
  price_minor = VALUES(price_minor),
  currency_code = VALUES(currency_code);

-- -------------------------------------------------------------------
-- Product attribute values
-- -------------------------------------------------------------------

INSERT INTO product_attribute_values (
  product_attribute_value_id,
  product_id,
  variant_id,
  attribute_id,
  value_string,
  value_number,
  value_boolean,
  value_json,
  unit
) VALUES
  ('eeeeeeee-eeee-eeee-eeee-eeeeeeee0001', 'cccccccc-cccc-cccc-cccc-cccccccc0001', NULL, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 'led_integrated', NULL, NULL, NULL, NULL),
  ('eeeeeeee-eeee-eeee-eeee-eeeeeeee0002', 'cccccccc-cccc-cccc-cccc-cccccccc0001', NULL, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2', NULL, 12, NULL, NULL, 'W'),

  ('eeeeeeee-eeee-eeee-eeee-eeeeeeee0003', 'cccccccc-cccc-cccc-cccc-cccccccc0002', NULL, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa3', 'cotton', NULL, NULL, NULL, NULL),
  ('eeeeeeee-eeee-eeee-eeee-eeeeeeee0004', 'cccccccc-cccc-cccc-cccc-cccccccc0002', NULL, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa4', 'Machine washable at 30°C. Do not deploy to tumble dryer.', NULL, NULL, NULL, NULL),

  ('eeeeeeee-eeee-eeee-eeee-eeeeeeee0005', 'cccccccc-cccc-cccc-cccc-cccccccc0003', 'dddddddd-dddd-dddd-dddd-dddddddd0003', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa5', 'medium', NULL, NULL, NULL, NULL),
  ('eeeeeeee-eeee-eeee-eeee-eeeeeeee0006', 'cccccccc-cccc-cccc-cccc-cccccccc0003', 'dddddddd-dddd-dddd-dddd-dddddddd0004', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa5', 'large', NULL, NULL, NULL, NULL),

  ('eeeeeeee-eeee-eeee-eeee-eeeeeeee0007', 'cccccccc-cccc-cccc-cccc-cccccccc0004', NULL, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa6', 'super_secure', NULL, NULL, NULL, NULL),
  ('eeeeeeee-eeee-eeee-eeee-eeeeeeee0008', 'cccccccc-cccc-cccc-cccc-cccccccc0004', NULL, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa7', 'keypad', NULL, NULL, NULL, NULL),

  ('eeeeeeee-eeee-eeee-eeee-eeeeeeee0009', 'cccccccc-cccc-cccc-cccc-cccccccc0005', NULL, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa8', 'rectangle', NULL, NULL, NULL, NULL),
  ('eeeeeeee-eeee-eeee-eeee-eeeeeeee0010', 'cccccccc-cccc-cccc-cccc-cccccccc0005', NULL, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa9', NULL, 11, NULL, NULL, 'mm'),

  ('eeeeeeee-eeee-eeee-eeee-eeeeeeee0011', 'cccccccc-cccc-cccc-cccc-cccccccc0006', NULL, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa3', 'woven_polyester', NULL, NULL, NULL, NULL),
  ('eeeeeeee-eeee-eeee-eeee-eeeeeeee0012', 'cccccccc-cccc-cccc-cccc-cccccccc0006', NULL, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa4', 'Gentle wash. Best used after finding the bug.', NULL, NULL, NULL, NULL)
ON DUPLICATE KEY UPDATE
  value_string = VALUES(value_string),
  value_number = VALUES(value_number),
  value_boolean = VALUES(value_boolean),
  value_json = VALUES(value_json),
  unit = VALUES(unit);

-- -------------------------------------------------------------------
-- Product images
-- -------------------------------------------------------------------

INSERT INTO product_images (
  image_id,
  product_id,
  url,
  alt_text,
  display_order,
  is_primary
) VALUES
  ('ffffffff-ffff-ffff-ffff-ffffffff0001', 'cccccccc-cccc-cccc-cccc-cccccccc0001', 'https://example.local/images/gopher-desk-lamp.png', 'Blue Gopher desk lamp on a wooden desk.', 1, TRUE),
  ('ffffffff-ffff-ffff-ffff-ffffffff0002', 'cccccccc-cccc-cccc-cccc-cccccccc0002', 'https://example.local/images/gopher-cushion-set.png', 'Set of two Gopher-themed cushions on a sofa.', 1, TRUE),
  ('ffffffff-ffff-ffff-ffff-ffffffff0003', 'cccccccc-cccc-cccc-cccc-cccccccc0003', 'https://example.local/images/rob-pike-wall-tapestry.png', 'Rob Pike inspired wall tapestry in a modern office.', 1, TRUE),
  ('ffffffff-ffff-ffff-ffff-ffffffff0004', 'cccccccc-cccc-cccc-cccc-cccccccc0004', 'https://example.local/images/rivest-lockbox.png', 'Rivest super-secure lockbox on a hallway table.', 1, TRUE),
  ('ffffffff-ffff-ffff-ffff-ffffffff0005', 'cccccccc-cccc-cccc-cccc-cccccccc0005', 'https://example.local/images/dijkstra-pathfinding-rug.png', 'Dijkstra pathfinding rug on a living room floor.', 1, TRUE),
  ('ffffffff-ffff-ffff-ffff-ffffffff0006', 'cccccccc-cccc-cccc-cccc-cccccccc0006', 'https://example.local/images/grace-hopper-debugging-blanket.png', 'Grace Hopper debugging blanket folded on a chair.', 1, TRUE)
ON DUPLICATE KEY UPDATE
  url = VALUES(url),
  alt_text = VALUES(alt_text),
  display_order = VALUES(display_order),
  is_primary = VALUES(is_primary);
