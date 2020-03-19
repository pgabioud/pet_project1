package main

import (
	"fmt"
	"math"
	"sync"
	"testing"
)

func TestDummyProtocol(t *testing.T) {
	peers := map[PartyID]string{
		0: "localhost:6660",
		1: "localhost:6661",
		2: "localhost:6662",
	}

	//Circuit Id to test
	//var circuitID CircuitID = 1

	N := uint64(len(peers))
	P := make([]*LocalParty, N, N)
	dummyProtocol := make([]*DummyProtocol, N, N)

	var err error
	wg := new(sync.WaitGroup)
	for i := range peers {
		P[i], err = NewLocalParty(i, peers)
		P[i].WaitGroup = wg
		check(err)

		//dummyProtocol[i] = P[i].NewDummyProtocol(uint64(i+10), circuitID)
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

func TestBeaver(t *testing.T) {

	t.Run("countMultGate", func(t *testing.T) {
		n := countMultGate(7)
		if n != 3 {
			t.Errorf("counting failed")
		}
	})
	t.Run("genBeavers", func(t *testing.T) {
		sharedBeavers := genBeavers(3)
		var s = int64(math.Pow(2, 16)) + 1
		for i := range sharedBeavers {
			if mod(int64(sharedBeavers[i][0])*int64(sharedBeavers[i][1]), s) != sharedBeavers[i][2] {
				t.Errorf("bad beaver triplet")
			}
		}
	})
	t.Run("genSharedBeavers", func(t *testing.T) {
		beavers := genBeavers(3)
		//var s = int64(math.Pow(2, 16)) + 1
		sharedBeavers := genSharedBeavers(&beavers, 3)
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
}

func test(circuitID CircuitID, t *testing.T) {

	peers := TestCircuits[circuitID-1].Peers

	N := uint64(len(peers))
	P := make([]*LocalParty, N, N)
	dummyProtocol := make([]*DummyProtocol, N, N)

	var err error
	wg := new(sync.WaitGroup)
	for i := range peers {
		P[i], err = NewLocalParty(i, peers)
		P[i].WaitGroup = wg
		check(err)

		//dummyProtocol[i] = P[i].NewDummyProtocol(TestCircuits[circuitID-1].Inputs[i][GateID(i)], circuitID)
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
