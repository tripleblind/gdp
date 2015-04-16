package gdp

// TODO: client state MUST be checked with invariants (L == len())
// TODO: limit memory consumption!

type Tour struct {
	H    [][]byte
	M    [][]byte
	T    []int
	I    []int
	L    int
	S    int
	Link string
}

type Request struct {
	H     [][]byte // WARNING: replace h0 with this!
	M     [][]byte // WARNING: replace MSMS1 with this!
	I     []int    // WARNING: replace I with this / client state!
	H0    []byte
	L     int
	S     int
	ISM1  int // WARNING: this is NOT in the pape
	TSM1  int
	MSM1  []byte
	I1    int
	IS    int
	Link  string
	Reply *Reply
}

type Reply struct {
	H    [][]byte // WARNING: client state!
	HS   []byte
	MS   []byte
	ISP1 int
	TS   int
}

type FinalRequest struct {
	H0   []byte
	HL   []byte
	L    int
	MLM1 []byte
	I    []int
}

type FinalReply struct {
	HSOL []byte
}

type Solution struct {
	H0   []byte
	HSOL []byte
	T0   int
	TL   int
	I1   int
}
