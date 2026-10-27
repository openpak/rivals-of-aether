# NEX registry spec — rivals-of-aether (for switch-nex)

Delivery spec for the serial NEX registry integration. This repo does NOT edit
`switch-nex`; the registry owner copies this spec into
`cmd/switch-nexd/main.go` once TODO-KEY closes.

Wire facts from `scratch/decompile-rivals-of-aether/RECON.md` and
`scratch/decompile-rivals-of-aether/findings.json` (`F` keys).

| Field | Value |
| --- | --- |
| Title | Rivals of Aether (Dan Fornace, display version 2.0.0) [F control] |
| Title ID | `0100D2700DE7C000` [F title_id] |
| NEX version | 4.6.4, modules MM/R2/VS [F middleware.nex] |
| Engine | GameMaker Studio 2 [F engine] |
| Auth/discovery | nncs: `n.n.srv.nintendo.net`, `s.n.srv.nintendo.net` [F hosts] |
| Secure transport | PRUDP; client knows `prudp:/` and `prudps:` schemes [F middleware.nex_apis] |
| Requested port base | **23160** (auth 23160, secure 23161, healthz 23162) — next free in the 230xx third-party batch after tetris-99 (23150); verified free of the registry on 2026-10-26 |

## Requested protocol set

From the protocol clients present in the executable [F middleware.nex_apis]:

- MatchmakeExtension (+ DDL)
- MatchmakeReferee (round arbitration)
- MatchmakeSession / RendezVous session passwords
- Ranking + Ranking2
- SecureConnection / Utility (standard session envelope)

Unconfirmed from static strings: NAT traversal use (no `natcheck`-style
marker recovered) and the secure station scheme choice (`prudp:` vs
`prudps:`). Needs a runtime online-menu DNS + traffic trace before the entry
is finalized.

## Recovered keys (TODO-KEY closed 2026-10-27)

- Game server ID `2A28FB00` (706199552), access key `8f0c5daa`.
  Same GameMaker literal pool in the update main: key at file offset
  `0x417ce95` between the NEX login format strings, ID double
  `706199552.0` at `0x417cee8`; both literals occur exactly once.
  Same-pool evidence; runtime/DNS confirmation still open (derived
  SNI `g2a28fb00-lp1.s.n.srv.nintendo.net`).
- Base-v0 NCA titlekey unknown: static analysis used the update (v655360)
  main; a base-vs-update NEX surface diff is still open.

## Registry entry (landed in switch-nex)

```go
{"rivals-of-aether", 23160, "8f0c5daa", func(e *protocols.Env) {
    e.StationScheme = "prudps"
    e.Pia = protocols.LegacyPia()
    e.Match = protocols.NewMatchmakingStore()
    e.MatchmakeReferee = true
    e.FinalizeCreateSession = true
    e.JoinEchoSession = true
    e.JoinRecapParticipate = true
}}
```

Secure scheme `prudps` follows the NEX 4.6.x peers (client knows both
`prudp:/` and `prudps:`); Ranking/Ranking2 stay unregistered until the
first boot's NotImplemented log. Provenance:
`switch-nex/docs/rivals-of-aether.md`.
