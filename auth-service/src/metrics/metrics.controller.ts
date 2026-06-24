import { Controller, Get, Res } from '@nestjs/common';
import { Response } from 'express';
import { register } from './metrics.provider';

@Controller('metrics')
export class MetricsController {
  @Get()
  async getMetrics(@Res() res: Response) {
    res.set('Content-Type', register.contentType);
    res.send(await register.metrics());
  }
}
