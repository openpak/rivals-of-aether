# Rivals of Aether

OpenPak's Rivals of Aether Server.

| | |
| --- | --- |
| Game | Rivals of Aether |
| Network | Nintendo Switch |
| Title ID | `0100D2700DE7C000` |
| Licence | AGPL-3.0-only |

## What it serves

| Host | What | Served by |
| --- | --- | --- |
| `n.n.srv.nintendo.net`, `s.n.srv.nintendo.net` + secure PRUDP | NEX 4.6.4 auth, matchmaking, referee, ranking | [`switch-nex`](https://github.com/openpak/switch-nex) (registry entry pending, see [docs/nex.md](docs/nex.md)) |
| `stats.rivalsofaether.com` | match start/finish telemetry (plain HTTP) | this repo (`rivals-of-aetherd`) |
| `www.rivalsofaether.com` | in-game news CDN | this repo (`rivals-of-aetherd`) |

## Run

```sh
go run ./cmd/rivals-of-aetherd
```

| Variable | Default | Meaning |
| --- | --- | --- |
| `RIVALS_ADDR` | `:23570` | telemetry + news stub listen address (behind the TLS terminator) |
| `RIVALS_HEALTH_ADDR` | `:23571` | `/healthz` listen address |
| `RIVALS_NEX_SECURE_URL` | — | switch-nex secure upstream, once the title profile lands |
| `RIVALS_NEWS_DIR` | — | directory serving the news CDN files; empty answers `204` |
| `NX_INTERNAL_URL` | `http://127.0.0.1:20070` | nx-baas internal API |
| `NX_INTERNAL_KEY` | — | its `X-Internal-Key` (or `OPENPAK_INTERNAL_KEY`) |

## Docs

- [Protocol](docs/protocol.md)
- [NEX registry spec](docs/nex.md)
- [Ports](docs/ports.md)

## Licence

AGPL-3.0-only. See [LICENSE](LICENSE).
