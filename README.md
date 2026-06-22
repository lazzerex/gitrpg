<div align="center">

# GitRPG

**Transform your GitHub activity into an RPG character.**

Commits earn XP. Repositories unlock classes. Every push levels you up.

<br>

<a href="https://github.com/lazzerex/gitrpg/stargazers"><img src="https://img.shields.io/github/stars/lazzerex/gitrpg?style=flat&color=FFD700&label=Stars"/></a>
<img src="https://img.shields.io/badge/License-MIT-6633CC?style=flat"/>


<img src="https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white"/>
<img src="https://img.shields.io/badge/Chi-Router-00ADD8?style=flat"/>
<img src="https://img.shields.io/badge/PostgreSQL-4169E1?style=flat&logo=postgresql&logoColor=white"/>
<img src="https://img.shields.io/badge/Redis-DC382D?style=flat&logo=redis&logoColor=white"/>
<img src="https://img.shields.io/badge/HTMX-3D72D7?style=flat"/>
<img src="https://img.shields.io/badge/Tailwind_CSS-06B6D4?style=flat&logo=tailwindcss&logoColor=white"/>
<img src="https://img.shields.io/badge/Three.js-000000?style=flat&logo=threedotjs&logoColor=white"/>
<img src="https://img.shields.io/badge/Docker-2496ED?style=flat&logo=docker&logoColor=white"/>
<img src="https://img.shields.io/badge/Goose-Migrations-4169E1?style=flat"/>
<img src="https://img.shields.io/badge/GitHub_OAuth-181717?style=flat&logo=github&logoColor=white"/>


</div>


## What is this

GitHub RPG reads your public GitHub activity and assigns you a character. Your primary language determines your class. Your commits, pull requests, reviews, and stars determine your level. Everything updates automatically.

You get a README card you can embed anywhere.


## Features

- **XP from real activity** - commits, merged PRs, reviews, issues, repository creation
- **9 classes** - assigned by your primary language (Go, Rust, TypeScript, JavaScript, Python, C#, Java, C++, and a catch-all)
- **Level curve** - `XP(N) = 100 * N^1.8`, scales to discourage farming
- **SVG README card** - classic and chart styles, auto-updating, cacheable
- **Achievements** - milestone badges synced on each login
- **3D character viewer** - interactive portrait on the landing page
- **Anti-farming rules** - repos need 5 commits, 1 star, or 1 external contributor to qualify


## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.25, Chi router |
| Frontend | HTMX, TailwindCSS, Go Templates |
| Database | PostgreSQL 16 |
| Cache / SVG cache | Redis 7 |
| Auth | GitHub OAuth |
| Migrations | Goose |
| 3D viewer | Three.js |


## Classes

| Class | Language | Description |
|---|---|---|
| Guardian | Go | Reliable defender. Builds scalable systems. |
| Berserker | Rust | Unstoppable force. Conquers memory fearlessly. |
| Paladin | TypeScript | Code with honor. Brings order to dynamic systems. |
| Rogue | JavaScript | Moves fast. Ships features at the speed of thought. |
| Sage | Python | Data sorcerer. Wields algorithms as spells. |
| Knight | C# | Enterprise warrior. Structured and disciplined. |
| Battlemage | Java | Hybrid power. Combines structure with versatility. |
| Warlord | C++ | Low-level conqueror. Masters performance and control. |
| Wanderer | Other | Adapts to any terrain. No language holds them back. |


## Getting Started

### Prerequisites

- Go 1.25+
- Docker (for Postgres and Redis)
- A GitHub OAuth App

### 1. Clone and configure

```bash
git clone https://github.com/lazzerex/gitrpg.git
cd gitrpg
cp .env.example .env
```

Edit `.env`:

```env
GITHUB_CLIENT_ID=your_client_id
GITHUB_CLIENT_SECRET=your_client_secret
GITHUB_CALLBACK_URL=http://localhost:8080/auth/github/callback

DATABASE_URL=postgres://gitrpg:gitrpg@localhost:5433/gitrpg?sslmode=disable
REDIS_URL=redis://localhost:6380

SESSION_SECRET=a-long-random-string
TOKEN_ENCRYPTION_KEY=   # 64 hex chars (32 bytes) - required in production
```

### 2. Start infrastructure

```bash
make docker-up
```

### 3. Run migrations

```bash
make migrate-up
```

### 4. Start the server

```bash
make dev
```

Open `http://localhost:8080`.


## Development

```bash
make dev          # run server with live reload
make build        # compile to bin/server
make test         # run tests
make lint         # run golangci-lint
make migrate-up   # apply pending migrations
make migrate-down # roll back last migration
make docker-up    # start postgres + redis
make docker-down  # stop infrastructure
make tidy         # go mod tidy
```


## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `PORT` | no | `8080` | HTTP server port |
| `ENV` | no | `development` | Set to `production` to enforce required vars |
| `DATABASE_URL` | yes | local default | PostgreSQL connection string |
| `REDIS_URL` | yes | local default | Redis connection string |
| `GITHUB_CLIENT_ID` | prod only | - | GitHub OAuth App client ID |
| `GITHUB_CLIENT_SECRET` | prod only | - | GitHub OAuth App client secret |
| `GITHUB_CALLBACK_URL` | no | localhost | OAuth callback URL |
| `SESSION_SECRET` | prod only | - | Cookie signing secret |
| `TOKEN_ENCRYPTION_KEY` | prod only | - | 64 hex chars (AES-256) for token encryption |


## README Card

After signing in and syncing, embed your card in any GitHub README:

```markdown
![GitHub RPG](https://gitrpg.app/card/YOUR_USERNAME.svg)
```

Chart style:

```markdown
![GitHub RPG](https://gitrpg.app/card/YOUR_USERNAME.svg?style=chart)
```

Compact badge:

```markdown
![GitHub RPG](https://gitrpg.app/card/compact/YOUR_USERNAME.svg)
```


## Project Structure

```
cmd/
  server/       entry point
  migrate/      migration runner
internal/
  achievements/ badge evaluation and storage
  auth/         GitHub OAuth flow
  characters/   class and stat computation
  config/       env-based config loader
  crypto/       AES-256 token encryption
  github/       GitHub API sync
  server/       HTTP handlers and routing
  stats/        XP and level formulas
  svg/          card SVG generation
  users/        user store
  worker/       background sync worker
web/
  static/       assets (sprites, icons, 3D models)
  templates/    Go HTML templates
migrations/     SQL migration files
```


<div align="center">

Built with Go. No JavaScript frameworks. No build pipeline.

</div>
