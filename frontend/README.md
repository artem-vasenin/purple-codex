# Frontend

Next.js 16+ frontend scaffold.

## Authentication

The frontend uses same-origin Next.js route handlers as a BFF for the Go API. Set the backend URL before starting the app:

```sh
export NEXT_PUBLIC_API_URL=http://localhost:8080
npm run dev
```

Registration and login are available at `/register` and `/login`. Both access and refresh tokens are stored by the BFF in `HttpOnly` cookies and are not exposed to browser JavaScript. The authenticated page is `/dashboard`.

Install dependencies when network access is available:

```sh
npm install
```
