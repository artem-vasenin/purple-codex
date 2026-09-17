import { redirect } from "next/navigation";
import { accessCookie, backendRequest, hasAccessCookie } from "@/lib/auth";
import { cookies } from "next/headers";
import { ProfileForm } from "@/components/profile-form";

export default async function ProfilePage() {
  if (!(await hasAccessCookie())) redirect("/login");
  const access = (await cookies()).get(accessCookie)?.value;
  const response = await backendRequest("/auth/profile", { headers: { Authorization: `Bearer ${access}` } });
  if (response.status === 401) redirect("/login");
  if (!response.ok) throw new Error("Unable to load profile");
  const profile = await response.json();
  return <main className="profile-page"><div className="profile-page-header"><a href="/dashboard" className="back-link">← Dashboard</a><p className="eyebrow">PROFILE</p></div><section className="profile-card"><h1>Edit profile.</h1><p className="muted">Keep your account details up to date.</p><ProfileForm initialProfile={profile} /></section></main>;
}
