import { Inject, Injectable } from '@nestjs/common';
import { Pool } from 'pg';
import { CreateUserPayload } from './interfaces/user.interface';
import { AuthQueries } from './sql/auth.queries';
import { uuidv7 } from 'uuidv7';
import { CACHE_MANAGER } from '@nestjs/cache-manager';
import { Cache } from 'cache-manager';
import { CACHE_KEYS } from 'src/constant/cache-keys.constants';

@Injectable()
export class AuthRepository {
  constructor(
    @Inject('PG_POOL') private readonly db: Pool,
    @Inject(CACHE_MANAGER) private readonly cacheManager: Cache,
  ) {}
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

    const newUser = result.rows[0];
    if (newUser) {
      await this.clearUserCache(newUser.id, user.email);
    }
    return newUser;
  }
  async findByEmail(email: string) {
    const cacheKey = CACHE_KEYS.USER_EMAIL(email);
    const cached = await this.cacheManager.get(cacheKey);

    if (cached) {
      return cached;
    }

    const result = await this.db.query(AuthQueries.findByEmail, [email]);
    const user = result.rows[0];
    if (user) {
      await this.cacheManager.set(cacheKey, user);
    }
    return user;
  }

  async findById(id: string) {
    const cacheKey = CACHE_KEYS.USER_ID(id);

    const cached = await this.cacheManager.get(cacheKey);
    if (cached) {
      return cached;
    }

    const result = await this.db.query(AuthQueries.findById, [id]);
    const user = result.rows[0];

    if (user) {
      await this.cacheManager.set(cacheKey, user);
    }

    return user;
  }

  async saveRefreshToken(token: string, userId: string) {
    const query = `
      INSERT INTO refresh_tokens (token, user_id, expires_at)
      VALUES ($1, $2, NOW() + INTERVAL '7 days')
    `;
    await this.db.query(query, [token, userId]);

    const cacheKey = CACHE_KEYS.REFRESH_TOKEN(token);
    await this.cacheManager.set(cacheKey, { token, userId }, 604800);
  }

  async findRefreshToken(token: string) {
    const cacheKey = CACHE_KEYS.REFRESH_TOKEN(token);
    const cached = await this.cacheManager.get(cacheKey);

    if (cached) {
      return cached;
    }

    const query = `
      SELECT * FROM refresh_tokens
      WHERE token = $1
      AND expires_at > NOW()
    `;
    const result = await this.db.query(query, [token]);
    const refreshToken = result.rows[0];
    if (refreshToken) {
      await this.cacheManager.set(cacheKey, refreshToken, 604800);
    }
    return refreshToken;
  }

  async deleteRefreshToken(token: string) {
    const query = `DELETE FROM refresh_tokens WHERE token = $1`;
    await this.db.query(query, [token]);
    const cacheKey = CACHE_KEYS.REFRESH_TOKEN(token);
    await this.cacheManager.del(cacheKey);
  }

  private async clearUserCache(id: string, email: string) {
    await this.cacheManager.del(CACHE_KEYS.USER_ID(id));
    await this.cacheManager.del(CACHE_KEYS.USER_EMAIL(email));
  }
}
