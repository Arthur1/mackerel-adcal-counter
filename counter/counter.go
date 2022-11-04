package counter

type CountResult struct {
	Entries      uint64
	Participants uint64
	Subscribers  uint64
}

type Counter interface {
	Count() (*CountResult, error)
}
