import { Inject, Injectable } from '@nestjs/common';
import { Pool } from 'pg';
import { CreateUserPayload } from './interfaces/user.interface';
import { AuthQueries } from './sql/auth.queries';
import { uuidv7 } from 'uuidv7';
import { Cache } from 'cache-manager';
import { CACHE_KEYS } from 'src/constant/cache-keys.constants';
import { RedisClientType } from 'redis';

@Injectable()
export class AuthRepository {
  constructor(
    @Inject('PG_POOL') private readonly db: Pool,
    @Inject('CACHE_MANAGER_AUTH_CACHE') private readonly cacheManager: Cache,
    @Inject('REDIS_CLIENT') private readonly redisClient: RedisClientType,
  ) {}
  async createUser(user: CreateUserPayload) {
    const id = uuidv7();

    const values = [
      id,
      user.email,
      user.password,
      user.firstName,
      user.lastName,
      user.role || 'user',
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
      const formattedUser = {
        id: user.id,
        email: user.email,
        password: user.password,
        first_name: user.first_name,
        last_name: user.last_name,
        role: user.role,
        created_at: user.created_at,
      };
      await this.cacheManager.set(cacheKey, formattedUser);
      return formattedUser;
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
      const formattedUser = {
        id: user.id,
        email: user.email,
        password: user.password,
        first_name: user.first_name,
        last_name: user.last_name,
        role: user.role,
        created_at: user.created_at,
      };
      await this.cacheManager.set(cacheKey, formattedUser);
      return formattedUser;
    }
    return user;
  }

  private getRefreshTokenKey(token: string): string {
    return `auth:refresh_token:${token}`;
  }

  async saveRefreshToken(token: string, userId: string) {
    const key = this.getRefreshTokenKey(token);
    await this.redisClient.setEx(key, 604800, userId);
  }

  async findRefreshToken(token: string) {
    const key = this.getRefreshTokenKey(token);
    const userId = await this.redisClient.get(key);

    if (!userId) {
      return null;
    }

    return { token, user_id: userId };
  }

  async deleteRefreshToken(token: string) {
    const key = this.getRefreshTokenKey(token);
    await this.redisClient.del(key);
  }

  async clearUserCache(id: string, email: string) {
    await this.cacheManager.del(CACHE_KEYS.USER_ID(id));
    await this.cacheManager.del(CACHE_KEYS.USER_EMAIL(email));
  }
}
