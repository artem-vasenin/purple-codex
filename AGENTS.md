# Repository Guidelines

## Project Structure & Module Organization

This repository contains two independently runnable applications: `backend/` and `frontend/`.

Keep code, configuration, and tests within the application they belong to. Tests should live beside the code they exercise or in the framework's established test location.

## Testing Guidelines

Every behavior change should include or update focused tests. Do not rely on frontend UI coverage until a frontend test framework is configured.

## Commit & Pull Request Guidelines

There is no commit history yet, so no repository-specific convention exists. Use imperative, concise subjects such as `Add health check endpoint`, keeping commits focused. Pull requests should explain the change, identify backend/frontend impact, list validation commands, link related issues, and include screenshots for visible frontend changes.

## Security & Configuration Tips

Never commit secrets or local environment files. Document required variables with safe example values, and keep production configuration separate from source code.
