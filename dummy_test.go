package main

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"testing"
	"time"
)

func Test_genSimpleBeaver(t *testing.T) {

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

	fmt.Println("Test_genSimpleBeaver COMPLETED")
}

func TestProtocol_simpleBeaver(t *testing.T) {
	//Circuit Id to test
	var circuitID CircuitID = 1

	peers := TestCircuits[circuitID-1].Peers

	count := CountMultGate(circuitID)

	N := uint64(len(peers))
	P := make([]*LocalParty, N, N)
	dummyProtocol := make([]*Protocol, N, N)

	var err error
	wg := new(sync.WaitGroup)
	for i := range peers {
		P[i], err = NewLocalParty(i, peers)
		P[i].WaitGroup = wg
		check(err)

		if count > 0 {
			genSharedBeavers := GenAllBeaverTriplets(CircuitID(circuitID))
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
		if TestCircuits[circuitID-1].ExpOutput != p.Output {
			t.Errorf("Party %d got wrong answer: %d when expected output was: %d", p.ID, p.Output, TestCircuits[circuitID-1].ExpOutput)
		}
	}

	fmt.Println("TestProtocol_simpleBeaver COMPLETED")
}

func Test_genBvfBeaver(t *testing.T) {
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
		fmt.Println("New Test_genBvfBeaver COMPLETED")
	})

	t.Run("Reshape", func(t *testing.T) {

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

		for i := 0; i < 3; i++ {
			var resulta uint64 = 0
			var resultb uint64 = 0
			var resultc uint64 = 0
			for _, p := range beaverProtocol {
				Beaver := p.ReshapeBeaver(7)
				resultc += Beaver[i][2]
				resulta += Beaver[i][0]
				resultb += Beaver[i][1]
			}

			res1 := mod(int64(resulta*resultb), int64(modulus))
			res2 := mod(int64(resultc), int64(modulus))
			if res1 != res2 {
				t.Error("wrong, youre fake news")
			}
		}

		wg.Wait()
		fmt.Println("Reshape Test_genBvfBeaver COMPLETED")
	})

}

func TestProtocol_bvfBeaver(t *testing.T) {
	//Circuit Id to test
	var circuitID CircuitID = 9

	peers := TestCircuits[circuitID-1].Peers

	N := uint64(len(peers))
	P := make([]*LocalParty, N, N)
	dummyProtocol := make([]*Protocol, N, N)

	count := CountMultGate(circuitID)
	beaverProtocol := make([]*BeaverProtocol, N, N)

	if count > 0 {
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
	}

	var err error
	wg := new(sync.WaitGroup)
	for i := range peers {
		P[i], err = NewLocalParty(i, peers)
		P[i].WaitGroup = wg
		check(err)
		Beaver := [][3]uint64{{0, 0, 0}}
		if count > 0 {
			Beaver = beaverProtocol[i].ReshapeBeaver(circuitID)
		}
		dummyProtocol[i] = P[i].NewProtocol(TestCircuits[circuitID-1].Inputs[i][GateID(i)], circuitID, &Beaver)
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
		if TestCircuits[circuitID-1].ExpOutput != p.Output {
			t.Errorf("Party %d got wrong answer: %d when expected output was: %d", p.ID, p.Output, TestCircuits[circuitID-1].ExpOutput)
		}
	}

	fmt.Println("TestProtocol_bvfBeaver COMPLETED")
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

	count := CountMultGate(circuitID)
	beaverProtocol := make([]*BeaverProtocol, N, N)

	if count > 0 {
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
	}

	var err error
	wg := new(sync.WaitGroup)
	for i := range peers {
		P[i], err = NewLocalParty(i, peers)
		P[i].WaitGroup = wg
		check(err)
		Beaver := [][3]uint64{{0, 0, 0}}
		if count > 0 {
			Beaver = beaverProtocol[i].ReshapeBeaver(circuitID)
		}
		dummyProtocol[i] = P[i].NewProtocol(TestCircuits[circuitID-1].Inputs[i][GateID(i)], circuitID, &Beaver)
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
		if TestCircuits[circuitID-1].ExpOutput != p.Output {
			t.Errorf("Party %d got wrong answer: %d when expected output was: %d", p.ID, p.Output, TestCircuits[circuitID-1].ExpOutput)
		}
	}

	fmt.Println("TestEval_bfvBeaver on circuit ", circuitID, " COMPLETED")
}

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

	fmt.Println("TestVectorOperations COMPLETED")
}

func TestPerformance_bvfBeaver(t *testing.T) {
	times := 7
	for i := 0; i < len(TestCircuits); i++ {
		circuitStr := "circuit" + strconv.Itoa(i+1)
		fmt.Println(circuitStr)
		t.Run(circuitStr, func(t *testing.T) {
			avg := 0
			for j := 0; j < times; j++ {
				start := time.Now()
				test(CircuitID(i+1), t)
				elapsed := time.Since(start)
				avg += int(elapsed)
			}
			avg = avg / times
			fmt.Println(circuitStr, " performance averaged over ", times, " is ", time.Duration(avg))
		})
	}
}

func testSimple(circuitID CircuitID, t *testing.T) {

	peers := TestCircuits[circuitID-1].Peers

	count := CountMultGate(circuitID)

	N := uint64(len(peers))
	P := make([]*LocalParty, N, N)
	dummyProtocol := make([]*Protocol, N, N)

	var err error
	wg := new(sync.WaitGroup)
	for i := range peers {
		P[i], err = NewLocalParty(i, peers)
		P[i].WaitGroup = wg
		check(err)

		if count > 0 {
			genSharedBeavers := GenAllBeaverTriplets(CircuitID(circuitID))
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
		if TestCircuits[circuitID-1].ExpOutput != p.Output {
			t.Errorf("Party %d got wrong answer: %d when expected output was: %d", p.ID, p.Output, TestCircuits[circuitID-1].ExpOutput)
		}
	}

	fmt.Println("TestProtocol_simpleBeaver COMPLETED")
}

func TestPerformance_simpleBeaver(t *testing.T) {
	times := 7
	for i := 0; i < len(TestCircuits); i++ {
		circuitStr := "circuit" + strconv.Itoa(i+1)
		fmt.Println(circuitStr)
		t.Run(circuitStr, func(t *testing.T) {
			avg := 0
			for j := 0; j < times; j++ {
				start := time.Now()
				testSimple(CircuitID(i+1), t)
				elapsed := time.Since(start)
				avg += int(elapsed)
			}
			avg = avg / times
			fmt.Println(circuitStr, " performance averaged over ", times, " is ", time.Duration(avg))
		})
	}
}
