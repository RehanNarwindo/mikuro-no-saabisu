// custom-cache.module.ts
import { Module } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import {
  CACHE_MANAGER,
  CacheModule,
  CacheModuleOptions,
} from '@nestjs/cache-manager';

@Module({})
export class CustomCacheModule {
  static register(options: { name: string; config: CacheModuleOptions }) {
    return {
      module: CustomCacheModule,
      imports: [
        CacheModule.registerAsync({
          useFactory: () => ({
            ttl: options.config.ttl,
            max: options.config.max,
          }),
          inject: [ConfigService],
        }),
      ],
      providers: [
        {
          provide: `CACHE_MANAGER_${options.name.toUpperCase()}`,
          useExisting: CACHE_MANAGER,
        },
      ],
      exports: [`CACHE_MANAGER_${options.name.toUpperCase()}`],
    };
  }
}
