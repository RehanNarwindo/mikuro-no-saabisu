import * as jwt from 'jsonwebtoken';
import type { SignOptions } from 'jsonwebtoken';

export function generateAccessToken(payload: object): string {
  const JWT_SECRET = process.env.JWT_SECRET;
  const JWT_EXPIRES_IN = process.env.JWT_EXPIRES_IN || '1h';

  if (!JWT_SECRET) {
    throw new Error('JWT_SECRET is not defined');
  }

  const options: SignOptions = {
    expiresIn: JWT_EXPIRES_IN as SignOptions['expiresIn'],
  };

  return jwt.sign(payload, JWT_SECRET, options);
}

export function verifyAccessToken(token: string): jwt.JwtPayload | null {
  const JWT_SECRET = process.env.JWT_SECRET;

  if (!JWT_SECRET) {
    throw new Error('JWT_SECRET is not defined');
  }

  try {
    return jwt.verify(token, JWT_SECRET) as jwt.JwtPayload;
  } catch {
    return null;
  }
}
