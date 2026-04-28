import { Module, Global } from '@nestjs/common';
import { RedisProvider } from './redis.providers';

@Global()
@Module({
  providers: [RedisProvider],
  exports: [RedisProvider],
})
export class RedisModule {}
