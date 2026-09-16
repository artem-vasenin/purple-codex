"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";

export function AuthForm({ mode }: { mode: "login" | "register" }) {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [pending, setPending] = useState(false);
  const register = mode === "register";

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault(); setPending(true); setError("");
    try {
      const response = await fetch(`/api/auth/${mode}`, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ email, password }) });
      if (!response.ok) { const body = await response.json().catch(() => ({})); throw new Error(body.error ?? "Unable to authenticate"); }
      router.push("/dashboard"); router.refresh();
    } catch (cause) { setError(cause instanceof Error ? cause.message : "Unable to authenticate"); } finally { setPending(false); }
  }

  return <form className="auth-form" onSubmit={submit}>
    <label>Email<input type="email" value={email} onChange={(event) => setEmail(event.target.value)} required autoComplete="email" /></label>
    <label>Password<input type="password" value={password} onChange={(event) => setPassword(event.target.value)} required autoComplete={register ? "new-password" : "current-password"} /></label>
    {error && <p className="form-error" role="alert">{error}</p>}
    <button disabled={pending}>{pending ? "Please wait..." : register ? "Create account" : "Sign in"}</button>
  </form>;
}
