-- =============================================================
-- DML: Seed / default data for user_management schema
-- Run AFTER ddl.sql has been executed
-- =============================================================

INSERT INTO user_management.users (username, email, full_name, phone, is_active)
VALUES
    ('admin',    'admin@example.com',    'Administrator',    '081200000001', TRUE),
    ('john_doe', 'john@example.com',     'John Doe',         '081200000002', TRUE),
    ('jane_doe', 'jane@example.com',     'Jane Doe',         '081200000003', TRUE),
    ('budi_s',   'budi@example.com',     'Budi Santoso',     '081200000004', TRUE),
    ('siti_r',   'siti@example.com',     'Siti Rahayu',      '',             FALSE)
ON CONFLICT DO NOTHING;
