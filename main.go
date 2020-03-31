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
	if (errCircuitID != nil) || (circuitID > 8) || (circuitID == 0) {
		fmt.Println("Circuit input should be an integer between 1 and 8")
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
	networkBeaver, err := NewTCPNetwork(lp)
	check(err)

	// Connect the circuit network
	err = networkBeaver.Connect(lp)
	check(err)
	fmt.Println(lp, "connected")
	<-time.After(time.Second) // Leave time for others to connect

	fmt.Println("start generating beavers")
	partyBeaverProtocol := lp.NewBeaverProtocol()
	fmt.Println("protocol created, now bind network")
	partyBeaverProtocol.BeaverBindNetwork(networkBeaver)
	fmt.Println("binding on the network done, now run")
	partyBeaverProtocol.BeaverRun()
	fmt.Println("gen beaver finished")

	// Create the network for the circuit
	network, err := NewTCPNetwork(lp)
	check(err)

	// Connect the circuit network
	err = network.Connect(lp)
	check(err)
	fmt.Println(lp, "connected")
	<-time.After(time.Second) // Leave time for others to connect

	// Create a new circuit evaluation protocol
	dummyProtocol := lp.NewProtocol(partyInput, circuitID, partyBeaverProtocol)

	// Bind evaluation protocol to the network
	dummyProtocol.BindNetwork(network)

	// Evaluate the circuit
	dummyProtocol.Run()

	fmt.Println(lp, "completed with output", dummyProtocol.Output)
}
