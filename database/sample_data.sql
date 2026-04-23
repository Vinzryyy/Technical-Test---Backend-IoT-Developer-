-- Optional extra demo data. Run manually (Supabase SQL editor or psql)
-- AFTER the core schema.sql has been applied. Idempotent — safe to re-run.
--
-- Adds:
--   - 2 extra locations (Yogyakarta id=4, Medan id=5)
--   - 16 extra devices across all 5 locations
--   - 2 extra staff users (Yogyakarta-only and Medan-only)
--     passwords: both 'user123'
--   - Grants the existing 'supervisor' access to Yogyakarta too

-- ─── Locations ──────────────────────────────────────────────────────────────
INSERT INTO locations (id, name) VALUES
    (4, 'Yogyakarta'),
    (5, 'Medan')
ON CONFLICT (name) DO NOTHING;

-- Keep sequence ahead of manual inserts so auto-create works
SELECT setval('locations_id_seq', GREATEST((SELECT COALESCE(MAX(id), 0) FROM locations), 1));

-- ─── Devices ────────────────────────────────────────────────────────────────
INSERT INTO devices (id, name, location_id, status) VALUES
    -- Jakarta extras
    (10, 'JKT Motion Sensor 01',       1, 'online'),
    (11, 'JKT Smoke Detector 01',      1, 'offline'),
    (12, 'JKT Door Lock 02',           1, 'online'),

    -- Surabaya extras
    (13, 'SBY Humidity Sensor 01',     2, 'online'),
    (14, 'SBY CCTV Camera 01',         2, 'online'),
    (15, 'SBY Smart Meter 02',         2, 'offline'),

    -- Bandung extras
    (16, 'BDG Smart Meter 01',         3, 'online'),
    (17, 'BDG CCTV Camera 01',         3, 'offline'),
    (18, 'BDG Motion Sensor 01',       3, 'online'),

    -- Yogyakarta
    (19, 'YGY Temperature Sensor 01',  4, 'online'),
    (20, 'YGY Humidity Sensor 01',     4, 'online'),
    (21, 'YGY Smart Meter 01',         4, 'offline'),
    (22, 'YGY Door Lock 01',           4, 'online'),

    -- Medan
    (23, 'MDN Temperature Sensor 01',  5, 'online'),
    (24, 'MDN CCTV Camera 01',         5, 'online'),
    (25, 'MDN Smoke Detector 01',      5, 'offline')
ON CONFLICT (id) DO NOTHING;

SELECT setval('devices_id_seq', GREATEST((SELECT COALESCE(MAX(id), 0) FROM devices), 1));

-- ─── Users ──────────────────────────────────────────────────────────────────
-- Bcrypt hash of 'user123' (cost=10).
INSERT INTO users (id, name, email, password, role) VALUES
    ('ffffffff-0000-0000-0000-000000000001', 'Yogyakarta Staff', 'yogya@example.com',
     '$2a$10$13ZgxRrgAuZ35A7JXPPcTeZ7WvYkNZ085XAnp8byk5SUAriivplNm', 'user'),
    ('ffffffff-0000-0000-0000-000000000002', 'Medan Staff', 'medan@example.com',
     '$2a$10$13ZgxRrgAuZ35A7JXPPcTeZ7WvYkNZ085XAnp8byk5SUAriivplNm', 'user')
ON CONFLICT (email) DO NOTHING;

-- Grant location access for extra users
INSERT INTO user_locations (user_id, location_id) VALUES
    ('ffffffff-0000-0000-0000-000000000001', 4),
    ('ffffffff-0000-0000-0000-000000000002', 5)
ON CONFLICT DO NOTHING;

-- Expand supervisor coverage to include Yogyakarta — only if the supervisor
-- user already exists (it is created by seed.go at app startup, not by
-- schema.sql). Safe to run before OR after the app seeds; re-run after the
-- first successful app boot to activate this grant.
INSERT INTO user_locations (user_id, location_id)
SELECT 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 4
WHERE EXISTS (
    SELECT 1 FROM users WHERE id = 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee'
)
ON CONFLICT DO NOTHING;
