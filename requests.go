package gdp

type Request struct {
	H0   []byte
	L    int
	S    int
	ISM1 int // WARNING: this is NOT in the paper!
	TSM1 int
	MSM1 []byte
	I1   int
	IS   int
}

type Reply struct {
	HS   []byte
	MS   []byte
	ISP1 int
	TS   int
}

type FinalRequest struct {
	H0   []byte
	HL   []byte
	L    int
	MLM1 int
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
