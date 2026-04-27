CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,  
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  first_name TEXT,
  last_name TEXT,
  role TEXT NOT NULL DEFAULT 'user',
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  CONSTRAINT valid_role CHECK (role IN ('admin', 'user'))
);
