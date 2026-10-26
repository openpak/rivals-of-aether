# Rivals of Aether wire protocol (Switch 0100D2700DE7C000)

Read out of the Switch client's own executable strings (GameMaker Studio 2
build, display version 2.0.0). Every fact below comes from the workspace recon
at `scratch/decompile-rivals-of-aether/RECON.md` and its
`scratch/decompile-rivals-of-aether/findings.json` — cited per section as
`[R §x]` / `[F key]`. Keep this file in step with `cmd/rivals-of-aetherd`.

## Netplay: NEX 4.6.4 over secure PRUDP

- NEX banners `NEX-4_6_4`, `NEX_MM-4_6_4`, `NEX_R2-4_6_4`, `NEX_VS-4_6_4`:
  modules MM/R2/VS [R §Engine] [F middleware.nex].
- Protocol clients present: MatchmakeExtension (+ DDL), MatchmakeReferee
  (round arbitration), MatchmakeSession / RendezVous session passwords,
  Ranking and Ranking2, secure PRUDP transport (`prudp:/`, `prudps:` URL
  schemes, `PRUDPEndPoint.cpp`) [R §Engine] [F middleware.nex_apis].
- Auth/discovery via nncs: `n.n.srv.nintendo.net`, `s.n.srv.nintendo.net`
  [R §Hosts] [F hosts].
- No Pia banner, no Photon / PlayFab / EOS / Demonware / Unity / Unreal
  markers [R §Engine] [F middleware.absent].
- NN SDK surface: `nn::socket`, `nn::ssl`, `nn::nsd::ResolveEx`,
  `nn::nifm` proxy; libcurl `NintendoSDK_libcurl-9_5_0-Release`
  [R §Engine] [F libcurl, F middleware.nn_sdk].

Netplay is served by `switch-nex`, not this repo: see [nex.md](nex.md) for the
registry spec. Outstanding: the NEX key candidate `8f0c5daa` is UNPAIRED (no
server ID recovered) [R §NEX key] [F nex_key].

## Telemetry: stats.rivalsofaether.com

First-party match telemetry over **plain HTTP** [R §Hosts] [F hosts]:

| Request | rivals-of-aetherd answers |
| --- | --- |
| `/api/games/start` | `200` with `{}`, body discarded |
| `/api/games/finish` | `200` with `{}`, body discarded |

Static strings name only the paths, so the request method and body shape are
unconfirmed: the stub accepts any method and logs the body size only. A
runtime capture of an online match pins the method and fields down.

## News CDN: www.rivalsofaether.com

In-game news files [R §Hosts] [F hosts]:

| Request | rivals-of-aetherd answers |
| --- | --- |
| `/files/ingame_assets/news` | file from `RIVALS_NEWS_DIR`, else `204` |
| `/files/ingame_assets/news_data.txt` | file from `RIVALS_NEWS_DIR`, else `204` |

## What is not here

- Crash reporting goes to sentry.io (project DSN embedded in the client);
  OpenPak does not answer it [R §Hosts] [F hosts].
- smash.gg shop links, steampowered cross-promo links and symbiotestudios
  merch are content/ads, not netplay [R §Hosts] [F hosts].
- Matchmaking, lobbies and gameplay are NEX (see above), not REST.
