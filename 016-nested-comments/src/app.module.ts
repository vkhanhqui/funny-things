import { Module } from '@nestjs/common';
import { CommentsModule } from './comments/comments.module';
import { PrismaModule } from './prisma/prisma.module';

@Module({
  imports: [PrismaModule, CommentsModule],
  controllers: [],
  providers: [],
})
export class AppModule {}
