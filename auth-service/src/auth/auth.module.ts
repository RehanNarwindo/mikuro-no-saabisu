import { Module } from '@nestjs/common';
import { AuthController } from './auth.controller';
import { AuthService } from './auth.service';
import { AuthRepository } from './auth.repository';
import { DatabaseModule } from 'src/database/database.module';
import { CustomCacheModule } from 'src/cache/custom-cache.module';

@Module({
  imports: [
    CustomCacheModule.register({
      name: 'auth_cache',
      config: { ttl: 300000, max: 100 },
    }),
    DatabaseModule,
  ],
  controllers: [AuthController],
  providers: [AuthService, AuthRepository],
})
export class AuthModule {}
