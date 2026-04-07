export const AuthQueries = {
  findById: `
        SELECT * FROM users 
        WHERE id = $1
    `,
  findByEmail: `
        SELECT * FROM users 
        WHERE email = $1
    `,
  create: `
        INSERT INTO users (id, email, password, first_name, last_name)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING *
    `,
};
