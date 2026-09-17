import { accessCookie, backendRequest, forwardBackendError } from "@/lib/auth";
import { cookies } from "next/headers";

async function authorizedRequest(path: string, init: RequestInit = {}) {
  const access = (await cookies()).get(accessCookie)?.value;
  if (!access) return null;
  return backendRequest(path, { ...init, headers: { Authorization: `Bearer ${access}`, ...init.headers } });
}

export async function GET() {
  const response = await authorizedRequest("/auth/profile");
  if (!response) return Response.json({ error: "Not authenticated" }, { status: 401 });
  if (!response.ok) return forwardBackendError(response);
  return Response.json(await response.json());
}

export async function PATCH(request: Request) {
  const response = await authorizedRequest("/auth/profile", { method: "PATCH", body: JSON.stringify(await request.json()) });
  if (!response) return Response.json({ error: "Not authenticated" }, { status: 401 });
  if (!response.ok) return forwardBackendError(response);
  return Response.json(await response.json());
}
