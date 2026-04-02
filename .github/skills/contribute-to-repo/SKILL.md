---
name: contribute-to-repo
description: >
  Full workflow for committing changes and opening a pull request against the
  upstream repository. Use this when asked to commit work, push, or open a PR
  on this repo.
---

## Repo layout

Discover the remotes dynamically:

```bash
git remote -v
```

Identify which remote points to the **upstream** (canonical) repo and which is the **fork**.
PRs must always target the **upstream** repo, opened from a fork branch.

## Step-by-step workflow

### 1. Identify co-authors

Ask the user if anyone should be co-authored on the commit (e.g. a pair-programming partner).
If so, look up their email from git history:

```bash
git log --format="%an <%ae>" | grep -i "<name>" | head -1
```

If no match is found, ask the user for the email directly.

### 2. Lint, format, and test

Before committing, run the linter (with auto-fix) and the test suite:

```bash
golangci-lint-v2 run ./... --fix
go test ./...
```

Resolve any lint errors or test failures before proceeding.

### 3. Stage and commit

Stage all modified files and commit using **Conventional Commits** format:

```
<type>(<scope>): <short summary>

<body — what changed and why>

Co-authored-by: <name> <email>
Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>
```

Common types: `feat`, `fix`, `refactor`, `chore`, `docs`, `test`.

```bash
git add -A
git commit
```

### 4. Push to the fork

```bash
git push <fork-remote> <current-branch>
```

### 5. Open a PR against upstream

Use the GitHub MCP tool `create_pull_request` with:
- `owner`: upstream org (from `git remote -v`)
- `repo`: repo name
- `head`: `<fork-owner>:<current-branch>`
- `base`: `main`
- Include a clear PR body summarising **what** changed and **why**.
