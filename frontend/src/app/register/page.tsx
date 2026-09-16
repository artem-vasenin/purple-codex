import Link from "next/link";
import { AuthForm } from "@/components/auth-form";

export default function RegisterPage() { return <main className="auth-shell"><section className="auth-card"><p className="eyebrow">START MONITORING</p><h1>Make it observable.</h1><p className="muted">Create your workspace and get a clear view of uptime from day one.</p><AuthForm mode="register" /><p className="switch">Already have an account? <Link href="/login">Sign in</Link></p></section></main>; }
