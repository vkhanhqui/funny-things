export class CommentNode {
  id: string;
  postId: string;
  userId: string;
  content: string;
  createdAt: Date;
  updatedAt: Date;
  path: string | null;
  replies: CommentNode[];

  constructor(comment: {
    id: string;
    postId: string;
    userId: string;
    content: string;
    createdAt: Date;
    updatedAt: Date;
    path: string | null;
  }) {
    this.id = comment.id;
    this.postId = comment.postId;
    this.userId = comment.userId;
    this.content = comment.content;
    this.createdAt = comment.createdAt;
    this.updatedAt = comment.updatedAt;
    this.replies = [];
    this.path = comment.path;
  }
}
