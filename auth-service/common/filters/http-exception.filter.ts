import {
  ExceptionFilter,
  Catch,
  ArgumentsHost,
  HttpException,
  HttpStatus,
  Logger,
} from '@nestjs/common';

@Catch()
export class GlobalExceptionsFilter implements ExceptionFilter {
  private readonly logger = new Logger(GlobalExceptionsFilter.name);

  catch(exception: any, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const response = ctx.getResponse();
    const request = ctx.getRequest();

    let status = HttpStatus.INTERNAL_SERVER_ERROR;
    let message = 'Internal server error';

    if (exception instanceof HttpException) {
      status = exception.getStatus();
      const errorResponse = exception.getResponse();
      message =
        typeof errorResponse === 'string'
          ? errorResponse
          : (errorResponse as any).message;
    } else if (exception?.code === '23505') {
      status = HttpStatus.CONFLICT;
      message = 'Email already registered';
    } else if (exception?.code === 'ECONNREFUSED') {
      status = HttpStatus.SERVICE_UNAVAILABLE;
      message = 'Database connection failed. Please try again later.';
    }

    this.logger.error(
      `[${request.method}] ${request.url} - Status: ${status}`,
      exception.stack || exception.message,
    );
    response.status(status).json({
      message: message,
      timestamp: new Date().toISOString(),
      path: request.url,
    });
  }
}
