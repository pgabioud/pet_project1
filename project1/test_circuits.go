package main

// TestCircuit is the structure representing a circuit
type TestCircuit struct {
	Peers     map[PartyID]string            // Mapping from PartyID to network addresses
	Inputs    map[PartyID]map[GateID]uint64 // The partys' input for each gate
	Circuit   []Operation                   // Circuit definition
	ExpOutput uint64                        // Expected output
}

// TestCircuits is the slice of all circuits
var TestCircuits = []*TestCircuit{&circuit1, &circuit2, &circuit3, &circuit4, &circuit5, &circuit6, &circuit7, &circuit8, &circuit9}

var circuit1 = TestCircuit{
	// f(a,b,c) = a + b + c
	Peers: map[PartyID]string{
		0: "localhost:6660",
		1: "localhost:6661",
		2: "localhost:6662",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 18},
		1: {1: 7},
		2: {2: 42},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Input{
			Party: 2,
			Out:   2,
		},
		&Add{
			In1: 0,
			In2: 1,
			Out: 3,
		},
		&Add{
			In1: 2,
			In2: 3,
			Out: 4,
		},
		&Reveal{
			In:  4,
			Out: 5,
		},
	},
	ExpOutput: 67,
}

var circuit2 = TestCircuit{ // TODO check the ordering of the wires
	// f(a,b) = a - b
	Peers: map[PartyID]string{
		0: "localhost:6660",
		1: "localhost:6661",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 17},
		1: {1: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Sub{
			In1: 0,
			In2: 1,
			Out: 2,
		},
		&Reveal{
			In:  2,
			Out: 3,
		},
	},
	ExpOutput: 10,
}

var circuit3 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:6660",
		1: "localhost:6661",
		2: "localhost:6662",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
		2: {2: 11},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Input{
			Party: 2,
			Out:   2,
		},
		&Add{
			In1: 0,
			In2: 1,
			Out: 3,
		},
		&Add{
			In1: 2,
			In2: 3,
			Out: 4,
		},
		&MultCst{
			In:       4,
			CstValue: 5,
			Out:      5,
		},
		&Reveal{
			In:  5,
			Out: 6,
		},
	},
	ExpOutput: 115,
}

var circuit4 = TestCircuit{
	// f(a,b,c) = (a + b + c) + K
	Peers: map[PartyID]string{
		0: "localhost:6660",
		1: "localhost:6661",
		2: "localhost:6662",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
		2: {2: 11},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Input{
			Party: 2,
			Out:   2,
		},
		&Add{
			In1: 0,
			In2: 1,
			Out: 3,
		},
		&Add{
			In1: 2,
			In2: 3,
			Out: 4,
		},
		&AddCst{
			In:       4,
			CstValue: 7,
			Out:      5,
		},
		&Reveal{
			In:  5,
			Out: 6,
		},
	},
	ExpOutput: 30,
}

var circuit5 = TestCircuit{
	// f(a,b,c) = (a*K0 + b - c) + K1
	Peers: map[PartyID]string{
		0: "localhost:6660",
		1: "localhost:6661",
		2: "localhost:6662",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 4},
		1: {1: 2},
		2: {2: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Input{
			Party: 2,
			Out:   2,
		},
		&MultCst{
			In:       0,
			CstValue: 8,
			Out:      3,
		},
		&Add{
			In1: 3,
			In2: 1,
			Out: 4,
		},
		&Sub{
			In1: 4,
			In2: 2,
			Out: 5,
		},
		&AddCst{
			In:       5,
			CstValue: 8,
			Out:      6,
		},
		&Reveal{
			In:  6,
			Out: 7,
		},
	},
	ExpOutput: 35,
}

var circuit6 = TestCircuit{
	// f(a,b,c,d) = a+b+c+d
	Peers: map[PartyID]string{
		0: "localhost:6660",
		1: "localhost:6661",
		2: "localhost:6662",
		3: "localhost:6663",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 18},
		1: {1: 7},
		2: {2: 42},
		3: {3: 73},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Input{
			Party: 2,
			Out:   2,
		},
		&Input{
			Party: 3,
			Out:   3,
		},
		&Add{
			In1: 0,
			In2: 1,
			Out: 4,
		},
		&Add{
			In1: 2,
			In2: 3,
			Out: 5,
		},
		&Add{
			In1: 4,
			In2: 5,
			Out: 6,
		},
		&Reveal{
			In:  6,
			Out: 7,
		},
	},
	ExpOutput: 140,
}

var circuit7 = TestCircuit{
	// f(a,b,c) = (a*b) + (b*c) + (c*a)
	Peers: map[PartyID]string{
		0: "localhost:6660",
		1: "localhost:6661",
		2: "localhost:6662",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 7},
		1: {1: 3},
		2: {2: 14},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Input{
			Party: 2,
			Out:   2,
		},
		&Mult{
			In1: 0,
			In2: 1,
			Out: 3,
		},
		&Mult{
			In1: 1,
			In2: 2,
			Out: 4,
		},
		&Mult{
			In1: 0,
			In2: 2,
			Out: 5,
		},
		&Add{
			In1: 3,
			In2: 4,
			Out: 6,
		},
		&Add{
			In1: 5,
			In2: 6,
			Out: 7,
		},
		&Reveal{
			In:  7,
			Out: 8,
		},
	},
	ExpOutput: 161,
}

var circuit8 = TestCircuit{
	// f(a,b,c,d,e) = ((a+K0) + b*K1 - c)*(d+e)
	Peers: map[PartyID]string{
		0: "localhost:6660",
		1: "localhost:6661",
		2: "localhost:6662",
		3: "localhost:6663",
		4: "localhost:6664",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 11},
		2: {2: 17},
		3: {3: 2},
		4: {4: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Input{
			Party: 2,
			Out:   2,
		},
		&Input{
			Party: 3,
			Out:   3,
		},
		&Input{
			Party: 4,
			Out:   4,
		},
		&AddCst{
			In:       0,
			CstValue: 42,
			Out:      5,
		},
		&MultCst{
			In:       1,
			CstValue: 4,
			Out:      6,
		},
		&Add{
			In1: 3,
			In2: 4,
			Out: 7,
		},
		&Add{
			In1: 5,
			In2: 6,
			Out: 8,
		},
		&Sub{
			In1: 8,
			In2: 2,
			Out: 9,
		},
		&Mult{
			In1: 7,
			In2: 9,
			Out: 10,
		},
		&Reveal{
			In:  10,
			Out: 11,
		},
	},
	ExpOutput: 666,
}

var circuit9 = TestCircuit{
	// f(a,b,c,d) = ((a*K0) + (b*c)) - ((d+K1)*a)
	Peers: map[PartyID]string{
		0: "localhost:6660",
		1: "localhost:6661",
		2: "localhost:6662",
		3: "localhost:6663",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 10},
		1: {1: 11},
		2: {2: 12},
		3: {3: 13},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Input{
			Party: 2,
			Out:   2,
		},
		&Input{
			Party: 3,
			Out:   3,
		},
		&MultCst{
			In:       0,
			CstValue: 42,
			Out:      4,
		},
		&Mult{
			In1: 1,
			In2: 2,
			Out: 5,
		},
		&Add{
			In1: 4,
			In2: 5,
			Out: 6,
		},
		&AddCst{
			In:       3,
			CstValue: 24,
			Out:      7,
		},
		&Mult{
			In1: 0,
			In2: 7,
			Out: 8,
		},
		&Sub{
			In1: 6,
			In2: 8,
			Out: 9,
		},
		&Reveal{
			In:  9,
			Out: 10,
		},
	},
	ExpOutput: 182,
}
