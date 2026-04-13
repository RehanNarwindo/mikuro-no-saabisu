export const CACHE_KEYS = {
  USER_EMAIL: (email: string) => `auth:user:email:${email}`,
  USER_ID: (id: string) => `auth:user:id:${id}`,

  REFRESH_TOKEN: (token: string) => `auth:refresh:token:${token}`,
  USER_REFRESH_TOKENS: (userId: string) => `auth:refresh:user:${userId}`,

  SESSION: (sessionId: string) => `auth:session:${sessionId}`,
  USER_SESSIONS: (userId: string) => `auth:user:sessions:${userId}`,
} as const;

export type CacheKeyType = keyof typeof CACHE_KEYS;

export const CACHE_PREFIX = {
  AUTH: 'auth:',
  USER: 'auth:user:',
  REFRESH: 'auth:refresh:',
  SESSION: 'auth:session:',
} as const;
