import { Inject, Injectable } from '@nestjs/common';
import { Pool } from 'pg';
import { CreateUserPayload } from './interfaces/user.interface';
import { AuthQueries } from './sql/auth.queries';
import { uuidv7 } from 'uuidv7';
import { Cache } from 'cache-manager';
import { CACHE_KEYS } from 'src/constant/cache-keys.constants';

@Injectable()
export class AuthRepository {
  constructor(
    @Inject('PG_POOL') private readonly db: Pool,
    @Inject('CACHE_MANAGER_AUTH_CACHE') private readonly cacheManager: Cache,
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
      const formattedUser = {
        id: user.id,
        email: user.email,
        password: user.password,
        first_name: user.first_name,
        last_name: user.last_name,
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
        created_at: user.created_at,
      };
      await this.cacheManager.set(cacheKey, formattedUser);
      return formattedUser;
    }
    return user;
  }

  async saveRefreshToken(token: string, userId: string) {
    await this.db.query(AuthQueries.saveRefreshToken, [token, userId]);
    const cacheKey = CACHE_KEYS.REFRESH_TOKEN(token);
    await this.cacheManager.set(cacheKey, { token, userId }, 604800);
  }

  async findRefreshToken(token: string) {
    const cacheKey = CACHE_KEYS.REFRESH_TOKEN(token);
    const cached = await this.cacheManager.get(cacheKey);

    if (cached) {
      return cached;
    }

    const result = await this.db.query(AuthQueries.findRefreshToken, [token]);
    const refreshToken = result.rows[0];
    if (refreshToken) {
      await this.cacheManager.set(cacheKey, refreshToken, 604800);
    }
    return refreshToken;
  }

  async deleteRefreshToken(token: string) {
    await this.db.query(AuthQueries.deleteRefreshToken, [token]);
    const cacheKey = CACHE_KEYS.REFRESH_TOKEN(token);
    await this.cacheManager.del(cacheKey);
  }

  private async clearUserCache(id: string, email: string) {
    await this.cacheManager.del(CACHE_KEYS.USER_ID(id));
    await this.cacheManager.del(CACHE_KEYS.USER_EMAIL(email));
  }
}
