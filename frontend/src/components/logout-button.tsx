"use client";

import { useRouter } from "next/navigation";

export function LogoutButton() { const router = useRouter(); return <button className="logout" role="menuitem" onClick={async () => { await fetch("/api/auth/logout", { method: "POST" }); router.push("/login"); router.refresh(); }}>Logout</button>; }
