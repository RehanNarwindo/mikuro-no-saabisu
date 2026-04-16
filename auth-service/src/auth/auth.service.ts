import {
  Injectable,
  ConflictException,
  UnauthorizedException,
  InternalServerErrorException,
} from '@nestjs/common';
import { RegisterDto } from './dto/register.dto';
import { LoginDto } from './dto/login.dto';
import { hashPass, comparePass } from './helpers/bcrypt.helper';
import { generateAccessToken } from './helpers/jwt.helper';
import { AuthRepository } from './auth.repository';

@Injectable()
export class AuthService {
  constructor(private readonly authRepository: AuthRepository) {}

  async register(registerDto: RegisterDto) {
    const existing = await this.authRepository.findByEmail(registerDto.email);
    if (existing) {
      throw new ConflictException('Email already registered');
    }

    try {
      const hashedPassword = await hashPass(registerDto.password);

      const user = await this.authRepository.createUser({
        email: registerDto.email,
        password: hashedPassword,
        firstName: registerDto.firstName,
        lastName: registerDto.lastName,
      });

      if (!user) {
        throw new InternalServerErrorException('Failed to create user');
      }

      const tokens = generateAccessToken({
        sub: user.id,
        email: user.email,
      });

      return { user, tokens };
    } catch (error) {
      if (
        error instanceof ConflictException ||
        error instanceof InternalServerErrorException
      ) {
        throw error;
      }
      throw new InternalServerErrorException(
        'Registration failed. Please try again.',
      );
    }
  }

  async login(loginDto: LoginDto) {
    const user = await this.authRepository.findByEmail(loginDto.email);

    if (!user) {
      throw new UnauthorizedException('Invalid credentials');
    }

    const isValid = await comparePass(loginDto.password, user.password);

    if (!isValid) {
      throw new UnauthorizedException('Invalid credentials');
    }

    const tokens = generateAccessToken({
      sub: user.id,
      email: user.email,
    });

    return {
      user: {
        id: user.id,
        email: user.email,
        firstName: user.first_name,
        lastName: user.last_name,
      },
      tokens,
    };
  }

  async me(userId: string) {
    const user = await this.authRepository.findById(userId);
    if (!user) {
      throw new UnauthorizedException('User not found');
    }
    return user;
  }
}
