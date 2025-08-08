const BASE_URL = "http://localhost:3000";
const POST_ID = "01982686-4c00-7abe-826b-bf325c510b7b";
const USER_ID = "01982686-75ad-78ad-8587-d4c9815c5741";

export { POST_ID, USER_ID };

export async function fetchComments() {
  const res = await fetch(`${BASE_URL}/comments/post/${POST_ID}`, { cache: "no-store" });

  if (!res.ok) throw new Error("Failed to fetch comments");

  const data = await res.json();
  return Array.isArray(data) ? data : [];
}

export async function createComment({
  content,
  path,
}: {
  content: string;
  path?: string;
}) {
  const res = await fetch(`${BASE_URL}/comments`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      postId: POST_ID,
      userId: USER_ID,
      content,
      path,
    }),
  });

  if (!res.ok) throw new Error("Failed to create comment");
  return res.json();
}
