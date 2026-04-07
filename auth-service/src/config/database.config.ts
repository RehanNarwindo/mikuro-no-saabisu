const database = () => ({
  port: Number.parseInt(process.env.PORT, 10) || 3000,

  database: {
    host: process.env.DB_HOST,
    port: Number.parseInt(process.env.DB_PORT, 10),
    username: process.env.DB_USERNAME,
    password: process.env.DB_PASSWORD,
    name: process.env.DB_NAME,
  },

  jwt: {
    secret: process.env.JWT_SECRET,
    expiresIn: process.env.JWT_EXPIRES_IN,
  },
});

export default database;
