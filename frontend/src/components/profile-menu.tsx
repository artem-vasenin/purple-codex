"use client";

import { useState } from "react";
import { LogoutButton } from "@/components/logout-button";

export function ProfileMenu({ email }: { email: string }) {
  const [open, setOpen] = useState(false);
  const [avatarFailed, setAvatarFailed] = useState(false);

  return (
    <div className="profile-menu">
      <button className="profile-trigger" type="button" aria-expanded={open} aria-haspopup="menu" onClick={() => setOpen((current) => !current)}>
        <span className="profile-icon" aria-hidden="true">
          {avatarFailed ? <svg viewBox="0 0 24 24" role="img"><circle cx="12" cy="8" r="3.5" /><path d="M5 20c.7-3.5 3.1-5.5 7-5.5s6.3 2 7 5.5" /></svg> : <img src="/api/profile/avatar" alt="" onError={() => setAvatarFailed(true)} />}
        </span>
        <span className="profile-name">{email}</span>
        <span className={`profile-chevron${open ? " is-open" : ""}`} aria-hidden="true">⌄</span>
      </button>
      {open && <div className="profile-dropdown" role="menu"><a href="/profile" role="menuitem">Profile</a><LogoutButton /></div>}
    </div>
  );
}
