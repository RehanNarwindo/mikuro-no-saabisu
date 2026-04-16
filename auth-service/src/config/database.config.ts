const database = () => ({
  port: process.env.PORT ? Number.parseInt(process.env.PORT, 10) : 3000,

  database: {
    host: process.env.DB_HOST,
    port: process.env.PORT ? Number.parseInt(process.env.PORT, 10) : 3000,
    username: process.env.DB_USERNAME,
    password: process.env.DB_PASSWORD,
    name: process.env.DB_NAME,
  },

  jwt: {
    secret: process.env.JWT_SECRET,
    expiresIn: process.env.JWT_EXPIRES_IN || '1h',
  },
});

export default database;
