package values

type PepSource = string

const (
	PepSourceCOAF         PepSource = "COAF"
	PepSourceSelfDeclared PepSource = "SELF_DECLARED"
	PepSourceWatchlist    PepSource = "WATCHLIST"
)
