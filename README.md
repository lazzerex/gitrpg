<div align="center">

# GitRPG

<img width="469" height="155" alt="gitrpg-logo" src="https://github.com/user-attachments/assets/9401cf11-6203-49bd-82f0-bf8e005a9cbe" />

**Transform your GitHub activity into an RPG character.**

Commits earn XP. Repositories unlock classes. Every push levels you up.

> Still in active development. Expect breaking changes.

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
<img src="https://img.shields.io/badge/GitHub_OAuth-181717?style=flat&logo=github&logoColor=white"/>
<img src="https://img.shields.io/badge/Render-46E3B7?style=flat&logo=render&logoColor=white"/>
<img src="https://img.shields.io/badge/Neon-00E599?style=flat"/>
<img src="https://img.shields.io/badge/Upstash-00E9A3?style=flat"/>

</div>

<br>

<img width="1339" height="862" alt="image" src="https://github.com/user-attachments/assets/3a77861c-d3b3-491a-8fec-34ca62563050" />

## What is this

GitRPG reads your public GitHub activity and assigns you a character. Your primary language determines your class. Your commits, pull requests, reviews, and stars determine your level. Everything updates automatically every 6 hours.

You get a README card you can embed anywhere.

![GitRPG](https://gitrpg.onrender.com/card/lazzerex.svg)

## Features

- **XP from real activity** - commits, merged PRs, reviews, closed issues, qualified repositories - across your entire account history
- **9 classes** - assigned by primary language (Go, Rust, TypeScript, JavaScript, Python, C#, Java, C++, plus Wanderer)
- **Level curve** - `XP(N) = 100 × N^1.8`, farming-resistant by design
- **5 stats** - STR, INT, WIS, DEX, CHA derived from different activity signals
- **SVG README card** - pixel-art character card, Redis-cached, embeddable in any README
- **Card evolution** - the card frame upgrades at levels 10, 25, and 50: bronze, silver, gold, animated gold
- **Pixel characters** - a generated 16x16 sprite per class, shown on profiles and playable in the arena
- **Weekly quests** - rotating GitHub-activity challenges (merge PRs, review, close issues, commit) that award bonus XP
- **Equipment** - weapon, shield, and accessory earned from your top languages
- **Achievements** - 9 milestone badges, evaluated on every sync
- **Leaderboard** - global XP ranking
- **Arena mode** - a built-in top-down minigame: pick any class, fight 5 waves and a boss with your real stats
- **Anti-farming rules** - repos need 5 commits, 1 star, or 1 external contributor to count
- **Efficient sync** - cheap activity checks skip idle users; full syncs cover multi-year history in windowed GraphQL queries

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.25, Chi router |
| Frontend | HTMX, TailwindCSS, Go Templates |
| Database | PostgreSQL 16 (Neon) |
| Cache | Redis 7 (Upstash) - SVG cards cached 1 hr |
| Auth | GitHub OAuth + HMAC-SHA256 signed cookies |
| Migrations | Goose |
| Hosting | Render |

## Classes

| Class | Language |
|---|---|
| Guardian | Go |
| Berserker | Rust |
| Paladin | TypeScript |
| Rogue | JavaScript |
| Sage | Python |
| Knight | C# |
| Battlemage | Java |
| Warlord | C++ |
| Wanderer | Other / no dominant language |

## Achievements

| Rarity | Achievement | Condition |
|---|---|---|
| Common | First Commit | 1+ commits |
| Common | First Pull Request | 1+ merged PRs |
| Common | First Repository | 1+ qualified repos |
| Rare | Thousand Commits | 1,000+ commits |
| Rare | Century of PRs | 100+ merged PRs |
| Rare | Code Reviewer | 10+ reviews submitted |
| Rare | Year-Long Streak | 365-day contribution streak |
| Legendary | Open Source Hero | 5+ external repos contributed to |
| Legendary | Star Collector | 10,000+ stars received |

## Getting Started

### Prerequisites

- Go 1.25+
- Docker (for Postgres and Redis)
- A GitHub OAuth App ([create one](https://github.com/settings/developers))

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
make dev             # run server
make build           # compile to bin/server
make test            # run tests
make lint            # run golangci-lint
make migrate-up      # apply pending migrations
make migrate-down    # roll back last migration
make migrate-status  # show migration state
make docker-up       # start postgres + redis
make docker-down     # stop infrastructure
make tidy            # go mod tidy
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
![GitRPG](https://gitrpg.onrender.com/card/YOUR_USERNAME.svg)
```

The card updates automatically every time your character syncs, and its frame evolves with your level: bronze below 10, silver at 10+, gold at 25+, and an animated gold frame at 50+.

## Arena Mode

Signed-in users can enter the arena from their profile - a zero-dependency canvas minigame where your character fights five waves of enemies plus a final boss, The Warden.

- Your real stats matter: AGILITY drives movement speed, POWER drives damage, level drives HP
- Play any of the 9 classes, each with its own attack style - blade arcs, whirlwinds, exploding orbs, knife fans, piercing lances
- Three enemy types (chasers, shooters, teleporters), buff pickups, and a boss with charge, radial burst, and summon skills
- Sprites are generated from text grids by `tools/gen-sprites` - the pixel art is code

## Credits

| Asset | Author | License |
|---|---|---|
| [KayKit Adventurers Character Pack 2.0](https://kaylousberg.itch.io/kaykit-adventurers) | Kay Lousberg | CC0 1.0 |
| [Kyrise's Free 16x16 RPG Icon Pack v1.3](https://kyrise.itch.io/kyrises-free-16x16-rpg-icon-pack) | Kyrise | CC BY 4.0 |
| Free Cute Tileset | itch.io | - |
| [Tiny RPG Character Asset Pack v1.03](https://pixel-boy.itch.io/tiny-rpg-character-asset-pack) | Pixel-Boy | Free |

## Project Structure

```
cmd/
  server/         entry point
  migrate/        migration runner
internal/
  achievements/   badge evaluation and storage
  auth/           GitHub OAuth flow and session management
  characters/     character persistence
  config/         env-based config loader
  crypto/         AES-256 token encryption
  equipment/      language-based loadout
  events/         domain event bus and persisted event log
  github/         GitHub GraphQL sync
  leaderboards/   global XP ranking
  quests/         weekly quest assignment and progress
  server/         HTTP handlers and routing
  stats/          XP, level, and stat formulas
  svg/            SVG card generation
  users/          user store
  worker/         background sync worker
tools/
  gen-sprites/    pixel sprite generator (text grids to PNG)
web/
  static/         assets (sprites, icons, tiles), css, js
  templates/      Go HTML templates
migrations/       SQL migration files (goose)
```

<div align="center">

Built with Go. No JavaScript frameworks. No build pipeline.

</div>
