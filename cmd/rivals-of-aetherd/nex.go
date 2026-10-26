package main

// Middleware client stub for the game's NEX stack.
//
// Netplay is NEX 4.6.4 over secure PRUDP, served by switch-nex once the title
// profile in docs/nex.md is registered there. This stub only tracks whether
// that upstream is configured; it performs no NEX I/O itself.

// TODO-KEY: the NEX Kerberos key candidate 8f0c5daa is UNPAIRED (no server ID
// recovered), so no registry entry or secure URL can be configured yet.
const nexKeyPaired = false

type nexStub struct {
	secureURL string
}

func newNexStub(cfg Config) *nexStub {
	return &nexStub{secureURL: cfg.NexSecureURL}
}

// status reports whether the NEX upstream is usable. ok=false names what is
// missing: the unpaired key first, then the unconfigured upstream URL.
func (n *nexStub) status() (ok bool, missing string) {
	if !nexKeyPaired {
		return false, "nex-key"
	}
	if n.secureURL == "" {
		return false, "nex-secure-url"
	}
	return true, ""
}
