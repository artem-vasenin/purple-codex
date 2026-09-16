import Link from "next/link";
import { AuthForm } from "@/components/auth-form";

export default function LoginPage() { return <main className="auth-shell"><section className="auth-card"><p className="eyebrow">UPTIME MONITOR</p><h1>Welcome back.</h1><p className="muted">Sign in to keep an eye on every service that matters.</p><AuthForm mode="login" /><p className="switch">New here? <Link href="/register">Create an account</Link></p></section></main>; }
