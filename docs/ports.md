# Ports

| Port | Use |
| --- | --- |
| 23570 | plain HTTP: telemetry (`/api/games/*`) and news CDN (`/files/ingame_assets/*`) stubs behind the TLS terminator |
| 23571 | `/healthz` |
| 23572–23579 | reserved |
| 23160–23162 | requested NEX block for the switch-nex registry entry (auth, secure, healthz); see [nex.md](nex.md) — not served by this repo |

23570–23579 is a free block of the 23500–23999 range (23500 unrailed, 23530/23533
fall-guys, 23540/23541 plants-vs-zombies, 23550/23552 stumble-guys, 23560/23561
grounded).
