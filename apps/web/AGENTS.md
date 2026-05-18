<!-- BEGIN:nextjs-agent-rules -->
# This is NOT the Next.js you know

This version has breaking changes — APIs, conventions, and file structure may all differ from your training data. Read the relevant guide in `node_modules/next/dist/docs/` before writing any code. Heed deprecation notices.
<!-- END:nextjs-agent-rules -->

## Project Safety Rules

- Backend-only request: do not edit anything under `apps/web/` unless the user explicitly asks for frontend changes.
- If a change might impact UI, stop and ask for confirmation before editing `apps/web/app/page.tsx` or `apps/web/app/globals.css`.
- Put new backend integration UI on dedicated routes/components first; never overwrite the homepage as a first step.
- Before large edits, create a checkpoint commit when a git repository is available.
