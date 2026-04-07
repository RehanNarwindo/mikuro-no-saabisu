import { Pool } from 'pg';

export const databaseProviders = [
  {
    provide: 'PG_POOL',
    useFactory: async () => {
      const pool = new Pool({
        host: process.env.DB_HOST,
        port: Number.parseInt(process.env.DB_PORT || '5432'),
        user: process.env.DB_USER,
        password: process.env.DB_PASS,
        database: process.env.DB_NAME,
      });

      await pool.connect();
      return pool;
    },
  },
];
