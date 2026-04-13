import { Module } from '@nestjs/common';
import { CacheModule } from '@nestjs/cache-manager';
import { AuthController } from './auth.controller';
import { AuthService } from './auth.service';
import { AuthRepository } from './auth.repository';
import { DatabaseModule } from 'src/database/database.module';

@Module({
  imports: [
    CacheModule.register({
      ttl: 300000,
      max: 100,
      name: 'auth_cache',
    }),
    DatabaseModule,
  ],
  controllers: [AuthController],
  providers: [AuthService, AuthRepository],
})
export class AuthModule {}
