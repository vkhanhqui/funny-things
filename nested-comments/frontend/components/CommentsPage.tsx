"use client";

import { fetchComments, createComment } from "@/lib/api";
import CommentItem from "@/components/CommentItem";
import { useEffect, useState } from "react";

export default function CommentsPage() {
  const [comments, setComments] = useState<any[]>([]);

  useEffect(() => {
    fetchComments().then(setComments);
  }, []);

  return (
    <main className="max-w-xl mx-auto p-6">
      <h1 className="text-2xl mb-4">Nested Comments</h1>

      <NewCommentForm />

      <div className="mt-6">
        {comments.map((comment) => (
          <CommentItem key={comment.id} comment={comment} />
        ))}
      </div>
    </main>
  );
}

function NewCommentForm() {
  const [content, setContent] = useState("");

  const handleSubmit = async () => {
    await createComment({ content });
    window.location.reload();
  };

  return (
    <div>
      <textarea
        placeholder="Write a root comment..."
        value={content}
        onChange={(e) => setContent(e.target.value)}
        className="border w-full p-4"
      />
      <button
        onClick={handleSubmit}
        className="mt-2 p-2 bg-blue-500 text-white"
      >
        Submit
      </button>
    </div>
  );
}
