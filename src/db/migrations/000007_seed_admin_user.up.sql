INSERT INTO users (username, email, password_hash, role)
VALUES (
    'admin',
    'admin@localhost',
    '$2a$12$Zt.Aw8/vcBbu6Rn15hhOuuJMsXfy/pQFtBqrMNDjxZU/KfUej2bfS',
    'admin'
)
ON CONFLICT (username) DO NOTHING;
