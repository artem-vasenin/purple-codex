import { accessCookie, backendRequest, forwardBackendError } from "@/lib/auth";
import { cookies } from "next/headers";

async function accessToken() { return (await cookies()).get(accessCookie)?.value; }

export async function GET() {
  const access = await accessToken();
  if (!access) return Response.json({ error: "Not authenticated" }, { status: 401 });
  const response = await backendRequest("/auth/profile/avatar", { headers: { Authorization: `Bearer ${access}` } });
  if (!response.ok) return forwardBackendError(response);
  return new Response(await response.arrayBuffer(), { status: response.status, headers: { "Content-Type": response.headers.get("Content-Type") ?? "application/octet-stream", "Cache-Control": "private, no-store" } });
}

export async function POST(request: Request) {
  const access = await accessToken();
  if (!access) return Response.json({ error: "Not authenticated" }, { status: 401 });
  const response = await backendRequest("/auth/profile/avatar", { method: "POST", headers: { Authorization: `Bearer ${access}`, "Content-Type": request.headers.get("Content-Type") ?? "" }, body: await request.arrayBuffer() });
  if (!response.ok) return forwardBackendError(response);
  return Response.json(await response.json());
}
