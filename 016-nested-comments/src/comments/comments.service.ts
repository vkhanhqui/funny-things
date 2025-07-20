import { Injectable } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';
import { Comment } from '@prisma/client';
import { v7 as uuidv7 } from 'uuid';
import { CommentNode } from './dto/comment-node.dto';

@Injectable()
export class CommentsService {
  constructor(private readonly prisma: PrismaService) {}

  async create(input: {
    postId: string;
    userId: string;
    content: string;
    path?: string;
  }): Promise<Comment> {
    return this.prisma.comment.create({
      data: {
        id: uuidv7(),
        postId: input.postId,
        userId: input.userId,
        content: input.content,
        path: input.path ?? null,
        createdAt: new Date(),
        updatedAt: new Date(),
      },
    });
  }

  async listByPost(postId: string): Promise<CommentNode[]> {
    const comments = await this.prisma.comment.findMany({
      where: { postId },
      orderBy: { createdAt: 'desc' },
    });

    const idNodeMap = new Map<string, CommentNode>();
    const roots: CommentNode[] = [];

    for (const comment of comments) {
      const node = new CommentNode(comment);
      idNodeMap.set(comment.id, node);
    }

    for (const comment of comments) {
      const node = idNodeMap.get(comment.id)!;

      if (!comment.path) {
        roots.push(node);
      } else {
        const parentId = comment.path.split('/').slice(-1)[0];
        const parent = idNodeMap.get(parentId);
        parent?.replies.push(node);
      }
    }
    return roots;
  }
}
