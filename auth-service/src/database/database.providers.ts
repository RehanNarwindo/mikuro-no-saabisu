import { Pool } from 'pg';

export const databaseProviders = [
  {
    provide: 'PG_POOL',
    useFactory: async () => {
      const pool = new Pool({
        host: process.env.DB_HOST,
        port: Number.parseInt(process.env.DB_PORT || '5432'),
        user: process.env.DB_USERNAME,
        password: process.env.DB_PASSWORD,
        database: process.env.DB_NAME,
        max: 20,
        idleTimeoutMillis: 30000,
        connectionTimeoutMillis: 2000,
        maxUses: 7500,
      });

      await pool.connect();
      return pool;
    },
  },
];
