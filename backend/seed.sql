-- Password for "password123" (bcrypt cost 12)
INSERT INTO users (id, name, email, password, created_at)
VALUES 
('d3b8f0a0-5e8e-4b0e-9f0a-8e8e8e8e8e8e', 
 'Test User', 
 'test@example.com', 
 '$2a$12$6b7b0b3e6e2a5f6d4c5b2a6d7e8f9g0h1i2j3k4l5m6n7o8p9q0r1s', 
 NOW())
ON CONFLICT (email) DO NOTHING;

-- Project
INSERT INTO projects (id, name, description, owner_id, created_at)
VALUES 
('a1b2c3d4-e5f6-7890-abcd-ef1234567890',
 'Website Redesign',
 'Q2 redesign project',
 'd3b8f0a0-5e8e-4b0e-9f0a-8e8e8e8e8e8e',
 NOW())
ON CONFLICT DO NOTHING;

-- Tasks
INSERT INTO tasks (id, title, description, status, priority, project_id, assignee_id, due_date)
VALUES 
('11111111-1111-1111-1111-111111111111', 'Design homepage', 'Create modern homepage', 'in_progress', 'high', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'd3b8f0a0-5e8e-4b0e-9f0a-8e8e8e8e8e8e', '2026-04-20'),
('22222222-2222-2222-2222-222222222222', 'Implement login', '', 'todo', 'medium', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', NULL, NULL),
('33333333-3333-3333-3333-333333333333', 'Write tests', 'Unit tests for auth', 'done', 'low', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'd3b8f0a0-5e8e-4b0e-9f0a-8e8e8e8e8e8e', '2026-04-15')
ON CONFLICT DO NOTHING;