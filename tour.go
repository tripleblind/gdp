package tour

// TODO: client state MUST be checked with invariants (L == len())
// TODO: limit memory consumption!

type Tour struct { // TODO: add unmarshalling routines to tour
	H    [][]byte
	M    [][]byte
	T    []int
	I    []int
	L    int
	S    int
	Sol  []byte
	Link string
}
