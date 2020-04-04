package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	prog := os.Args[0]
	args := os.Args[1:]

	if len(args) < 3 {
		fmt.Println("Usage:", prog, "[Party ID] [Input] [CircuitID]")
		os.Exit(1)
	}

	partyID, errPartyID := strconv.ParseUint(args[0], 10, 64)
	if errPartyID != nil {
		fmt.Println("Party ID should be an unsigned integer")
		os.Exit(1)
	}

	partyInput, errPartyInput := strconv.ParseUint(args[1], 10, 64)
	if errPartyInput != nil {
		fmt.Println("Party input should be an unsigned integer")
		os.Exit(1)
	}

	circuitID, errCircuitID := strconv.ParseUint(args[2], 10, 64)
	if (errCircuitID != nil) || (circuitID > uint64(len(TestCircuits))) || (circuitID == 0) {
		fmt.Println("Circuit input should be an integer between 1 and ", len(TestCircuits))
		os.Exit(1)
	}

	if int(partyID) < len(TestCircuits[circuitID-1].Peers) {
		Client(PartyID(partyID), partyInput, CircuitID(circuitID))
	}
}

//Client function
func Client(partyID PartyID, partyInput uint64, circuitID CircuitID) {
	//N := uint64(len(peers))
	peers := TestCircuits[circuitID-1].Peers

	// Create a local party
	lp, err := NewLocalParty(partyID, peers)
	check(err)

	// Create the network for the circuit
	network, err := NewTCPNetwork(lp)
	check(err)

	// Connect the circuit network
	err = network.Connect(lp)
	check(err)
	fmt.Println(lp, "connected")
	<-time.After(time.Second) // Leave time for others to connect

	partyBeaverProtocol := lp.NewBeaverProtocol()

	nbBeaver := CountMultGate(circuitID)
	Beavers := [][3]uint64{{0, 0, 0}}
	if nbBeaver > 0 {
		fmt.Println(lp, " start bever protocol")
		fmt.Println(lp, " beaver protocol binding on the network")
		partyBeaverProtocol.BeaverBindNetwork(network)
		fmt.Println(lp, " beaver protocol running")
		partyBeaverProtocol.BeaverRun()
		fmt.Println(lp, " beaver successfully generated")
	}

	Beavers = partyBeaverProtocol.ReshapeBeaver(circuitID)
	fmt.Println(lp, "did this ", Beavers)
	<-time.After(time.Second)
	// Create a new circuit evaluation protocol
	fmt.Println(lp, " create new dummy protocol")
	dummyProtocol := lp.NewProtocol(partyInput, circuitID, &Beavers)

	// Bind evaluation protocol to the network
	fmt.Println(lp, " dummy protocol binding to the network")
	dummyProtocol.BindNetwork(network)

	// Evaluate the circuit
	fmt.Println(lp, " dummy protocol running")
	dummyProtocol.Run()

	fmt.Println("\n", lp, "completed with output", dummyProtocol.Output)
}
