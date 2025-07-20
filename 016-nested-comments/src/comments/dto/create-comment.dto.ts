import { IsUUID, IsOptional, IsString, Matches } from 'class-validator';
import { MaxPathDepth } from '../validators/max-path-depth.validator';

export class CreateCommentDto {
  @IsUUID()
  postId: string;

  @IsUUID()
  userId: string;

  @IsString()
  content: string;

  @IsOptional()
  @IsString()
  @Matches(/^([0-9a-fA-F-]{36})(\/[0-9a-fA-F-]{36})*$/, {
    message: 'path must be a slash-separated list of UUIDs',
  })
  @MaxPathDepth(10)
  path?: string;
}
