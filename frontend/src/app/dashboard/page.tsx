import { redirect } from "next/navigation";
import { hasAccessCookie } from "@/lib/auth";
import { LogoutButton } from "@/components/logout-button";

export default async function DashboardPage() { if (!(await hasAccessCookie())) redirect("/login"); return <main className="dashboard"><div className="dashboard-top"><p className="eyebrow">CONTROL ROOM</p><LogoutButton /></div><h1>Your uptime, in focus.</h1><p className="muted">Your authenticated workspace is ready. Service monitors can be connected here next.</p><div className="empty-state"><span>00</span><p>No monitors yet</p><small>Add your first endpoint to begin tracking availability.</small></div></main>; }
