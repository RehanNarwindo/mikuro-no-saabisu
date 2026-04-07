import { Inject, Injectable } from '@nestjs/common';
import { Pool } from 'pg';
import { CreateUserPayload } from './interfaces/user.interface';
import { AuthQueries } from './sql/auth.queries';
import { uuidv7 } from 'uuidv7';

@Injectable()
export class AuthRepository {
  constructor(@Inject('PG_POOL') private readonly db: Pool) {}

  async createUser(user: CreateUserPayload) {
    const id = uuidv7();

    const values = [
      id,
      user.email,
      user.password,
      user.firstName,
      user.lastName,
    ];
    const result = await this.db.query(AuthQueries.create, values);
    return result.rows[0];
  }
  async findByEmail(email: string) {
    const result = await this.db.query(AuthQueries.findByEmail, [email]);
    return result.rows[0];
  }

  async findById(id: string) {
    const result = await this.db.query(AuthQueries.findById, [id]);
    return result.rows[0];
  }

  async saveRefreshToken(token: string, userId: string) {
    const query = `
      INSERT INTO refresh_tokens (token, user_id, expires_at)
      VALUES ($1, $2, NOW() + INTERVAL '7 days')
    `;

    await this.db.query(query, [token, userId]);
  }

  async findRefreshToken(token: string) {
    const query = `
      SELECT * FROM refresh_tokens
      WHERE token = $1
      AND expires_at > NOW()
    `;

    const result = await this.db.query(query, [token]);
    return result.rows[0];
  }

  async deleteRefreshToken(token: string) {
    const query = `DELETE FROM refresh_tokens WHERE token = $1`;
    await this.db.query(query, [token]);
  }
}
