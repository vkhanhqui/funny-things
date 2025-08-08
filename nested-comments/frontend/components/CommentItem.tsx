"use client";

import { useState } from "react";
import { createComment } from "@/lib/api";

type Comment = {
  id: string;
  content: string;
  path: string;
  replies: Comment[];
};

export default function CommentItem({ comment }: { comment: Comment }) {
  const [replying, setReplying] = useState(false);
  const [text, setText] = useState("");

  const handleReply = async () => {
    await createComment({
      content: text,
      path: comment.path ? `${comment.path}/${comment.id}` : comment.id,
    });
    window.location.reload();
    return;
  };

  return (
    <div className="ml-4 mt-4 border-l pl-4">
      <p>{comment.content}</p>
      <button
        onClick={() => setReplying((v) => !v)}
        className="text-sm text-blue-500"
      >
        Reply
      </button>

      {replying && (
        <div className="mt-2">
          <textarea
            value={text}
            onChange={(e) => setText(e.target.value)}
            className="border w-full p-2"
            rows={2}
          />
          <button
            onClick={handleReply}
            className="mt-1 p-2 bg-blue-500 text-white"
          >
            Submit
          </button>
        </div>
      )}

      {comment.replies?.map((c) => (
        <CommentItem key={c.id} comment={c} />
      ))}
    </div>
  );
}
