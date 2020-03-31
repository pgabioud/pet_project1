package main

import (
	"fmt"
	"math"
	"sync"
	"testing"
)

func TestProtocol_simpleBeaver(t *testing.T) {
	peers := map[PartyID]string{
		0: "localhost:6660",
		1: "localhost:6661",
		2: "localhost:6662",
	}

	//Circuit Id to test
	var circuitID CircuitID = 1

	N := uint64(len(peers))
	P := make([]*LocalParty, N, N)
	dummyProtocol := make([]*Protocol, N, N)

	beaverProtocol := make([]*BeaverProtocol, N, N)
	wgBeaver := new(sync.WaitGroup)

	for i := range peers {
		P[i], _ = NewLocalParty(i, peers)
		P[i].WaitGroup = wgBeaver

		beaverProtocol[i] = P[i].NewBeaverProtocol()

	}

	networkBeaver := GetTestingTCPNetwork(P)
	fmt.Println("parties connected")

	for i, Pi := range beaverProtocol {
		Pi.BeaverBindNetwork(networkBeaver[i])
	}

	for _, p := range beaverProtocol {
		p.Add(1)
		go p.BeaverRun()

	}

	wgBeaver.Wait()

	var err error
	wg := new(sync.WaitGroup)
	for i := range peers {
		P[i], err = NewLocalParty(i, peers)
		P[i].WaitGroup = wg
		check(err)

		dummyProtocol[i] = P[i].NewProtocol(TestCircuits[circuitID-1].Inputs[i][GateID(i)], circuitID, beaverProtocol[i])
	}

	network := GetTestingTCPNetwork(P)
	fmt.Println("parties connected")

	for i, Pi := range dummyProtocol {
		Pi.BindNetwork(network[i])
	}

	for _, p := range dummyProtocol {
		p.Add(1)
		go p.Run()
	}
	wg.Wait()

	for _, p := range dummyProtocol {
		fmt.Println(p, "completed with output", p.Output)
	}

	fmt.Println("test completed")
}

/*
func TestProtocol_simpleBeaver(t *testing.T) {
	peers := map[PartyID]string{
		0: "localhost:6660",
		1: "localhost:6661",
		2: "localhost:6662",
	}

	//Circuit Id to test
	var circuitID CircuitID = 1

	genSharedBeavers := GenAllBeaverTriplets(CircuitID(circuitID))

	N := uint64(len(peers))
	P := make([]*LocalParty, N, N)
	dummyProtocol := make([]*Protocol, N, N)

	var err error
	wg := new(sync.WaitGroup)
	for i := range peers {
		P[i], err = NewLocalParty(i, peers)
		P[i].WaitGroup = wg
		check(err)

		if len(genSharedBeavers) > 0 {
			dummyProtocol[i] = P[i].NewProtocol(TestCircuits[circuitID-1].Inputs[i][GateID(i)], circuitID, &genSharedBeavers[i])
		} else {
			nullBeavers := [][3]uint64{{0, 0, 0}}
			dummyProtocol[i] = P[i].NewProtocol(TestCircuits[circuitID-1].Inputs[i][GateID(i)], circuitID, &nullBeavers)
		}

	}

	network := GetTestingTCPNetwork(P)
	fmt.Println("parties connected")

	for i, Pi := range dummyProtocol {
		Pi.BindNetwork(network[i])
	}

	for _, p := range dummyProtocol {
		p.Add(1)
		go p.Run()
	}
	wg.Wait()

	for _, p := range dummyProtocol {
		fmt.Println(p, "completed with output", p.Output)
	}

	fmt.Println("test completed")
}
*/

func TestBeaver(t *testing.T) {

	t.Run("countMultGate", func(t *testing.T) {
		n := CountMultGate(7)
		if n != 3 {
			t.Errorf("counting failed")
		}
	})
	t.Run("genBeavers", func(t *testing.T) {
		var s = int64(math.Pow(2, 16)) + 1
		sharedBeavers := GenBeavers(3, s)
		for i := range sharedBeavers {
			if mod(int64(sharedBeavers[i][0])*int64(sharedBeavers[i][1]), s) != sharedBeavers[i][2] {
				t.Errorf("bad beaver triplet")
			}
		}
	})
	t.Run("genSharedBeavers", func(t *testing.T) {
		var s = int64(math.Pow(2, 16)) + 1
		beavers := GenBeavers(3, s)
		sharedBeavers := GenSharedBeavers(&beavers, 3, s)
		fmt.Println(sharedBeavers)
	})

}

func TestEval(t *testing.T) {
	t.Run("circuit1", func(t *testing.T) {
		test(1, t)
	})

	t.Run("circuit2", func(t *testing.T) {
		test(2, t)
	})

	t.Run("circuit3", func(t *testing.T) {
		test(3, t)
	})

	t.Run("circuit4", func(t *testing.T) {
		test(4, t)
	})

	t.Run("circuit5", func(t *testing.T) {
		test(5, t)
	})

	t.Run("circuit6", func(t *testing.T) {
		test(6, t)
	})

	t.Run("circuit7", func(t *testing.T) {
		test(7, t)
	})

	t.Run("circuit8", func(t *testing.T) {
		test(8, t)
	})

	t.Run("circuit9", func(t *testing.T) {
		test(9, t)
	})
}

func test(circuitID CircuitID, t *testing.T) {

	peers := TestCircuits[circuitID-1].Peers

	N := uint64(len(peers))
	P := make([]*LocalParty, N, N)
	dummyProtocol := make([]*Protocol, N, N)

	beaverProtocol := make([]*BeaverProtocol, N, N)
	wgBeaver := new(sync.WaitGroup)

	for i := range peers {
		P[i], _ = NewLocalParty(i, peers)
		P[i].WaitGroup = wgBeaver

		beaverProtocol[i] = P[i].NewBeaverProtocol()

	}

	networkBeaver := GetTestingTCPNetwork(P)
	fmt.Println("parties connected")

	for i, Pi := range beaverProtocol {
		Pi.BeaverBindNetwork(networkBeaver[i])
	}

	for _, p := range beaverProtocol {
		p.Add(1)
		go p.BeaverRun()

	}

	wgBeaver.Wait()

	var err error
	wg := new(sync.WaitGroup)
	for i := range peers {
		P[i], err = NewLocalParty(i, peers)
		P[i].WaitGroup = wg
		check(err)

		dummyProtocol[i] = P[i].NewProtocol(TestCircuits[circuitID-1].Inputs[i][GateID(i)], circuitID, beaverProtocol[i])
	}

	network := GetTestingTCPNetwork(P)
	fmt.Println("parties connected")

	for i, Pi := range dummyProtocol {
		Pi.BindNetwork(network[i])
	}

	for _, p := range dummyProtocol {
		p.Add(1)
		go p.Run()
	}
	wg.Wait()

	for _, p := range dummyProtocol {
		fmt.Println(p, "completed with output", p.Output)
		if p.Output != TestCircuits[circuitID-1].ExpOutput {
			t.Errorf("Party %d got wrong answer: %d when expected output was: %d", p.ID, p.Output, TestCircuits[circuitID-1].ExpOutput)
		}
	}

	fmt.Println("test completed")
}

/*
func test_simpleBeaver(circuitID CircuitID, t *testing.T) {

	peers := TestCircuits[circuitID-1].Peers
	genSharedBeavers := GenAllBeaverTriplets(CircuitID(circuitID))

	N := uint64(len(peers))
	P := make([]*LocalParty, N, N)
	dummyProtocol := make([]*Protocol, N, N)

	var err error
	wg := new(sync.WaitGroup)
	for i := range peers {
		P[i], err = NewLocalParty(i, peers)
		P[i].WaitGroup = wg
		check(err)

		if len(genSharedBeavers) > 0 {
			dummyProtocol[i] = P[i].NewProtocol(TestCircuits[circuitID-1].Inputs[i][GateID(i)], circuitID, &genSharedBeavers[i])
		} else {
			nullBeavers := [][3]uint64{{0, 0, 0}}
			dummyProtocol[i] = P[i].NewProtocol(TestCircuits[circuitID-1].Inputs[i][GateID(i)], circuitID, &nullBeavers)
		}
	}

	network := GetTestingTCPNetwork(P)
	fmt.Println("parties connected")

	for i, Pi := range dummyProtocol {
		Pi.BindNetwork(network[i])
	}

	for _, p := range dummyProtocol {
		p.Add(1)
		go p.Run()
	}
	wg.Wait()

	for _, p := range dummyProtocol {
		fmt.Println(p, "completed with output", p.Output)
		if p.Output != TestCircuits[circuitID-1].ExpOutput {
			t.Errorf("Party %d got wrong answer: %d when expected output was: %d", p.ID, p.Output, TestCircuits[circuitID-1].ExpOutput)
		}
	}

	fmt.Println("test completed")
}
*/

func TestVectorOperations(t *testing.T) {

	t.Run("NewRandomVec", func(t *testing.T) {
		v := NewRandomVec(20, 13)
		if len(v) != 20 {
			t.Errorf("wrong size")
		}
		for _, x := range v {
			if x < 0 || x >= 13 {
				t.Errorf("wrong random value")
			}
		}
		fmt.Println(v)
	})

	t.Run("AddVec", func(t *testing.T) {
		v1 := []uint64{1, 2, 3, 4, 5}
		v2 := []uint64{6, 7, 8, 9, 10}
		v1addv2 := AddVec(&v1, &v2, 11)
		expectedRes := []uint64{7, 9, 0, 2, 4}
		for i := 0; i < len(v1addv2); i++ {
			if v1addv2[i] != expectedRes[i] {
				t.Errorf("wrong addition: %d, should have been %d", v1addv2[i], expectedRes[i])
			}
		}
		//fmt.Println(v1addv2)
	})

	t.Run("SubVec", func(t *testing.T) {
		v1 := []uint64{1, 2, 3, 4, 5}
		v2 := []uint64{6, 7, 8, 9, 10}
		v1Subv2 := SubVec(&v1, &v2, 11)
		expectedRes := []uint64{6, 6, 6, 6, 6}
		for i := 0; i < len(v1Subv2); i++ {
			if v1Subv2[i] != expectedRes[i] {
				t.Errorf("wrong substraction: %d, should have been %d", v1Subv2[i], expectedRes[i])
			}
		}
		//fmt.Println(v1Subv2)
	})

	t.Run("MulVec", func(t *testing.T) {
		v1 := []uint64{1, 2, 3, 4, 5}
		v2 := []uint64{6, 7, 8, 9, 10}
		v1Mulv2 := MulVec(&v1, &v2, 11)
		expectedRes := []uint64{6, 3, 2, 3, 6}
		for i := 0; i < len(v1Mulv2); i++ {
			if v1Mulv2[i] != expectedRes[i] {
				t.Errorf("wrong multiplication: %d, should have been %d", v1Mulv2[i], expectedRes[i])
			}
		}
		//fmt.Println(v1Mulv2)
	})

	t.Run("NegVec", func(t *testing.T) {
		v1 := []uint64{1, 2, 3, 4, 5}
		v1Neg := NegVec(&v1, 11)
		expectedRes := []uint64{10, 9, 8, 7, 6}
		for i := 0; i < len(v1Neg); i++ {
			if v1Neg[i] != expectedRes[i] {
				t.Errorf("wrong negation: %d, should have been %d", v1Neg[i], expectedRes[i])
			}
		}
		//fmt.Println(v1Neg)
	})

}

func TestBVF(t *testing.T) {
	t.Run("New", func(t *testing.T) {

		var modulus uint64 = 65537
		peers := TestCircuits[7].Peers

		// Create the network for the circuit

		N := uint64(len(peers))
		P := make([]*LocalParty, N, N)
		beaverProtocol := make([]*BeaverProtocol, N, N)
		wg := new(sync.WaitGroup)

		for i := range peers {
			P[i], _ = NewLocalParty(i, peers)
			P[i].WaitGroup = wg

			beaverProtocol[i] = P[i].NewBeaverProtocol()

		}

		network := GetTestingTCPNetwork(P)
		fmt.Println("parties connected")

		for i, Pi := range beaverProtocol {
			Pi.BeaverBindNetwork(network[i])
		}

		for _, p := range beaverProtocol {
			p.Add(1)
			go p.BeaverRun()

		}

		wg.Wait()
		for i := 0; i < 20; i++ {
			var resulta uint64 = 0
			var resultb uint64 = 0
			var resultc uint64 = 0
			for _, p := range beaverProtocol {
				resultc += p.c[i]
				resulta += p.a[i]
				resultb += p.b[i]
			}

			res1 := mod(int64(resulta*resultb), int64(modulus))
			res2 := mod(int64(resultc), int64(modulus))
			if res1 != res2 {
				t.Error("wrong, youre fake news")
			}
		}

		wg.Wait()
		fmt.Println("test completed")
	})

}
