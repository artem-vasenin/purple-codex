import { accessCookie, backendRequest, cookieOptions, forwardBackendError, refreshCookie } from "@/lib/auth";
import { cookies } from "next/headers";
import { NextResponse } from "next/server";

export async function POST() {
  const refresh = (await cookies()).get(refreshCookie)?.value;
  if (!refresh) return Response.json({ error: "Not authenticated" }, { status: 401 });
  const response = await backendRequest("/auth/refresh", { method: "POST", headers: { Cookie: `${refreshCookie}=${refresh}` } });
  if (!response.ok) return forwardBackendError(response);
  const data = await response.json();
  const result = NextResponse.json({ authenticated: true });
  const nextRefresh = response.headers.get("set-cookie")?.match(/refresh_token=([^;]+)/)?.[1];
  if (data.accessToken) result.cookies.set(accessCookie, data.accessToken, cookieOptions(data.expiresIn));
  if (nextRefresh) result.cookies.set(refreshCookie, nextRefresh, cookieOptions(30 * 24 * 60 * 60));
  return result;
}
