import { backendRequest, cookieOptions, accessCookie, refreshCookie, profileCookie, profileCookieOptions } from "@/lib/auth";
import { cookies } from "next/headers";
import { NextResponse } from "next/server";

export async function POST() {
  const refresh = (await cookies()).get(refreshCookie)?.value;
  if (refresh) await backendRequest("/auth/logout", { method: "POST", headers: { Cookie: `${refreshCookie}=${refresh}` } });
  const result = NextResponse.json({ authenticated: false });
  result.cookies.set(accessCookie, "", cookieOptions(0));
  result.cookies.set(refreshCookie, "", cookieOptions(0));
  result.cookies.set(profileCookie, "", profileCookieOptions(0));
  return result;
}
