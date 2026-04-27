export const AuthQueries = {
  findById: `
    SELECT id, email, first_name, last_name, role, created_at 
    FROM users 
    WHERE id = $1
  `,
  findByEmail: `
    SELECT id, email, password, first_name, last_name, role, created_at 
    FROM users 
    WHERE email = $1
  `,
  create: `
    INSERT INTO users (id, email, password, first_name, last_name, role)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id, email, first_name, last_name, role, created_at
  `,

  saveRefreshToken: `
    INSERT INTO refresh_tokens (token, user_id, expires_at)
    VALUES ($1, $2, NOW() + INTERVAL '7 days')
  `,

  findRefreshToken: `
    SELECT * FROM refresh_tokens
    WHERE token = $1
    AND expires_at > NOW()
  `,

  deleteRefreshToken: `
    DELETE FROM refresh_tokens WHERE token = $1
  `,
};
