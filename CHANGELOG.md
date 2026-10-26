# Changelog — rivals-of-aether

## v0.1.0 — unreleased

- `rivals-of-aetherd`: skeleton serving the first-party HTTP side. Match
  start/finish telemetry is accepted and not stored, the news CDN serves
  `RIVALS_NEWS_DIR` or `204`, and `/v1/nex/status` reports `501` until the NEX
  key is paired. NEX registry spec delivered in `docs/nex.md` (port base
  request 23160); key candidate `8f0c5daa` still unpaired (TODO-KEY).
