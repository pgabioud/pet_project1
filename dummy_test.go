package main

import (
	"fmt"
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
	var circuitID CircuitID = 1

	N := uint64(len(peers))
	P := make([]*LocalParty, N, N)
	dummyProtocol := make([]*DummyProtocol, N, N)

	var err error
	wg := new(sync.WaitGroup)
	for i := range peers {
		P[i], err = NewLocalParty(i, peers)
		P[i].WaitGroup = wg
		check(err)

		dummyProtocol[i] = P[i].NewDummyProtocol(uint64(i+10), circuitID)
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

func TestTesting(t *testing.T) {
	got := 1
	if got != 1 {
		t.Errorf("Abs(-1 = %d; want 1", got)
	}
}

func TestEval(t *testing.T) {
	a := 2
	t.Run("circuit1", func(t *testing.T) {
		test(1, t)
		fmt.Println("we managed it!", a)
		a++

	})
	t.Run("circuit2", func(t *testing.T) {
		fmt.Println("circuit2", a)
		test(2, t)
	})

	t.Run("circuit3", func(t *testing.T) {
		fmt.Println("circuit3", a)
		test(3, t)
	})
}

func test(ciruitID CircuitID, t *testing.T) {

	peers := TestCircuits[ciruitID-1].Peers

	N := uint64(len(peers))
	P := make([]*LocalParty, N, N)
	dummyProtocol := make([]*DummyProtocol, N, N)

	var err error
	wg := new(sync.WaitGroup)
	for i := range peers {
		P[i], err = NewLocalParty(i, peers)
		P[i].WaitGroup = wg
		check(err)

		dummyProtocol[i] = P[i].NewDummyProtocol(TestCircuits[ciruitID-1].Inputs[i][GateID(i)], ciruitID)
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
