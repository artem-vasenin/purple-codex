import { backendRequest, cookieOptions, forwardBackendError, accessCookie, refreshCookie } from "@/lib/auth";
import { NextResponse } from "next/server";

export async function POST(request: Request) {
  const response = await backendRequest("/auth/register", { method: "POST", body: JSON.stringify(await request.json()) });
  if (!response.ok) return forwardBackendError(response);
  const data = await response.json();
  const result = NextResponse.json({ authenticated: true });
  const refresh = response.headers.get("set-cookie")?.match(/refresh_token=([^;]+)/)?.[1];
  if (data.accessToken) result.cookies.set(accessCookie, data.accessToken, cookieOptions(data.expiresIn));
  if (refresh) result.cookies.set(refreshCookie, refresh, cookieOptions(30 * 24 * 60 * 60));
  return result;
}
