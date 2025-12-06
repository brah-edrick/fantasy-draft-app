CREATE TABLE draft_rooms (
    id UUID PRIMARY KEY,
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Insert a dummy room to test
INSERT INTO draft_rooms (id, status) 
VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'WAITING');
