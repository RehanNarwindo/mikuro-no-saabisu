import { Injectable, NestMiddleware } from '@nestjs/common';
import { Request, Response, NextFunction } from 'express';
import { httpRequestsTotal, httpRequestDuration } from './metrics.provider';

@Injectable()
export class MetricsMiddleware implements NestMiddleware {
  use(req: Request, res: Response, next: NextFunction) {
    const start = Date.now();
    const { method, originalUrl } = req;

    if (originalUrl === '/metrics') {
      return next();
    }

    res.on('finish', () => {
      const duration = (Date.now() - start) / 1000;
      const status = res.statusCode.toString();

      httpRequestsTotal.inc({ method, endpoint: originalUrl, status });
      httpRequestDuration.observe({ method, endpoint: originalUrl }, duration);
    });

    next();
  }
}
