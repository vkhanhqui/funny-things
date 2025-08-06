import { Controller, Post, Body, Get, Param } from '@nestjs/common';
import { CommentsService } from './comments.service';
import { CreateCommentDto } from './dto/create-comment.dto';

@Controller('comments')
export class CommentsController {
  constructor(private readonly service: CommentsService) {}

  @Post()
  create(@Body() dto: CreateCommentDto) {
    return this.service.create(dto);
  }

  @Get('post/:postId')
  listByPost(@Param('postId') postId: string) {
    return this.service.listByPost(postId);
  }
}
