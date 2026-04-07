import { ConfigModule } from '@nestjs/config';
import { Module } from '@nestjs/common';
import database from './database.config';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      load: [database],
    }),
  ],
})
export class AppModule {}
