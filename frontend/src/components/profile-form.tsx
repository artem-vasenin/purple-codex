"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";

type Profile = { email: string; firstName: string; lastName: string; phone: string; hasAvatar: boolean };

export function ProfileForm({ initialProfile }: { initialProfile: Profile }) {
  const router = useRouter();
  const [profile, setProfile] = useState(initialProfile);
  const [error, setError] = useState("");
  const [saved, setSaved] = useState(false);
  const [pending, setPending] = useState(false);
  const [avatarPending, setAvatarPending] = useState(false);
  const [avatarError, setAvatarError] = useState("");

  function update(field: keyof Profile, value: string) {
    setProfile((current) => ({ ...current, [field]: value }));
    setSaved(false);
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setPending(true); setError(""); setSaved(false);
    try {
      const response = await fetch("/api/profile", { method: "PATCH", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ firstName: profile.firstName, lastName: profile.lastName, phone: profile.phone }) });
      const data = await response.json();
      if (!response.ok) throw new Error(data.error ?? "Unable to save profile");
      setProfile(data); setSaved(true); router.refresh();
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "Unable to save profile");
    } finally { setPending(false); }
  }

  async function uploadAvatar(file: File) {
    setAvatarPending(true); setAvatarError("");
    try {
      const body = new FormData(); body.append("avatar", file);
      const response = await fetch("/api/profile/avatar", { method: "POST", body });
      const data = await response.json();
      if (!response.ok) throw new Error(data.error ?? "Unable to upload avatar");
      setProfile((current) => ({ ...current, hasAvatar: true })); router.refresh();
    } catch (cause) { setAvatarError(cause instanceof Error ? cause.message : "Unable to upload avatar"); }
    finally { setAvatarPending(false); }
  }

  return <form className="profile-form" onSubmit={submit}>
    <div className="avatar-editor"><label className="avatar-picker" htmlFor="avatar-upload"><span className="avatar-circle">{profile.hasAvatar ? <img src="/api/profile/avatar" alt="Current avatar" /> : <span className="avatar-fallback" aria-hidden="true">♙</span>}</span><span>{avatarPending ? "Uploading..." : "Change avatar"}</span></label><input id="avatar-upload" className="avatar-input" type="file" accept="image/png,image/jpeg,image/gif" disabled={avatarPending} onChange={(event) => { const file = event.target.files?.[0]; if (file) void uploadAvatar(file); event.currentTarget.value = ""; }} />{avatarError && <p className="form-error">{avatarError}</p>}<small>PNG, JPEG or GIF. Maximum 1024x1024.</small></div>
    <label>Email<input type="email" value={profile.email} readOnly /></label>
    <label>First name<input value={profile.firstName} onChange={(event) => update("firstName", event.target.value)} autoComplete="given-name" /></label>
    <label>Last name<input value={profile.lastName} onChange={(event) => update("lastName", event.target.value)} autoComplete="family-name" /></label>
    <label>Phone<input type="tel" value={profile.phone} onChange={(event) => update("phone", event.target.value)} autoComplete="tel" /></label>
    {error && <p className="form-error">{error}</p>}
    {saved && <p className="form-success">Profile saved.</p>}
    <button disabled={pending}>{pending ? "Saving..." : "Save changes"}</button>
  </form>;
}
