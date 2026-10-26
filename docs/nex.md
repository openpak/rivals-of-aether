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

## Outstanding keys (blockers)

- TODO-KEY: NEX key candidate `8f0c5daa` is UNPAIRED — present once in main
  strings, no server-ID pair recovered (`pairs=[]`) [F nex_key]. Next: trace
  the `8f0c5daa` load site in main.bin for the paired server ID
  (movz/movk w2 construction near the NEX config call).
- Base-v0 NCA titlekey unknown: static analysis used the update (v655360)
  main; a base-vs-update NEX surface diff is still open.

## Registry entry sketch (illustrative, not an edit)

```go
{"rivals-of-aether", 23160, "<key once paired>", func(e *protocols.Env) {
    // MatchmakeExtension + Referee + Ranking/Ranking2; secure scheme TBD.
}}
```
