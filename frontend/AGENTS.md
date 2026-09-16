# Frontend Guidelines

## Project Structure & Module Organization

The frontend is a Next.js 16+ application. Routes and shared layout live in `src/app/`; static assets belong in `public/` when needed. Keep components and feature code organized under `src/`.

## Build, Test, and Development Commands

Run commands from the `frontend/` directory:

```sh
npm install          # Install frontend dependencies
npm run dev          # Start Next.js development server
npm run build        # Create a production build
npm run lint         # Run ESLint
```

Use the committed `package-lock.json`; update it through npm when dependencies change.

Authentication uses Next.js Route Handlers as a same-origin BFF. Access and refresh tokens must remain in `HttpOnly` cookies and must never be exposed to client components, `localStorage`, or `sessionStorage`. The Go API URL is configured through `NEXT_PUBLIC_API_URL`.

## Coding Style & Naming Conventions

TypeScript and TSX use two-space indentation and semicolons as enforced by the existing ESLint setup. Use `PascalCase` for React components, `camelCase` for variables and functions, and kebab-case for route segments. Prefer explicit types at public boundaries.

## Testing Guidelines

Frontend tests are not configured yet; add them with the chosen framework before relying on UI coverage.
