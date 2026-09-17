import { redirect } from "next/navigation";
import { cookies } from "next/headers";
import { hasAccessCookie, profileCookie } from "@/lib/auth";
import { ProfileMenu } from "@/components/profile-menu";

export default async function DashboardPage() {
  if (!(await hasAccessCookie())) redirect("/login");
  const email = (await cookies()).get(profileCookie)?.value ?? "Your profile";
  return <main className="dashboard"><header className="dashboard-header"><p className="eyebrow">CONTROL ROOM</p><ProfileMenu email={email} /></header><section className="dashboard-content"><h1>Your uptime, in focus.</h1><p className="muted">Your authenticated workspace is ready. Service monitors can be connected here next.</p><div className="empty-state"><span>00</span><p>No monitors yet</p><small>Add your first endpoint to begin tracking availability.</small></div></section></main>;
}
