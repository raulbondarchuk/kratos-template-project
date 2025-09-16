INSERT INTO `templates` (`id`, `type_id`, `name`, `created_at`, `updated_at`)
VALUES
  (1, 1, 'Template 1', NOW(), NOW()),
  (2, 2, 'Template 2', NOW(), NOW())
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`);