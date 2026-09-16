import { cookies } from "next/headers";

export const accessCookie = "access_token";
export const refreshCookie = "refresh_token";

export function backendUrl() {
  return process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";
}

export async function backendRequest(path: string, init: RequestInit = {}) {
  return fetch(`${backendUrl()}${path}`, {
    ...init,
    headers: { "Content-Type": "application/json", ...init.headers },
    cache: "no-store",
  });
}

export async function hasAccessCookie() {
  return Boolean((await cookies()).get(accessCookie)?.value);
}

export function cookieOptions(maxAge: number) {
  return {
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    sameSite: "lax" as const,
    path: "/",
    maxAge,
  };
}

export async function forwardBackendError(response: Response) {
  const body = await response.json().catch(() => ({ error: "Request failed" }));
  return Response.json({ error: body.error ?? "Request failed" }, { status: response.status });
}
