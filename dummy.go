package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"net"
	"time"
)

//CircuitID uint to select the circuit to run
type CircuitID uint64

func mod(a, b int64) int64 {
	m := a % b
	if a < 0 && b < 0 {
		m -= b
	}
	if a < 0 && b > 0 {
		m += b
	}
	return m
}

//DummyMessage is simple value message for party
type DummyMessage struct {
	Party PartyID
	Value uint64
}

//DummyProtocol basic structure for a protocol
type DummyProtocol struct {
	*LocalParty
	Chan  chan DummyMessage
	Peers map[PartyID]*DummyRemote

	circuitID CircuitID
	Beavers   []uint64

	Input  uint64
	Output uint64
}

//DummyRemote template for remote party
type DummyRemote struct {
	*RemoteParty
	Chan chan DummyMessage
}

//NewDummyProtocol initializes SMC for defined circuit with ID = circuitID
func (lp *LocalParty) NewDummyProtocol(input uint64, circuitID CircuitID) *DummyProtocol {
	cep := new(DummyProtocol)
	cep.LocalParty = lp
	cep.circuitID = circuitID
	cep.Chan = make(chan DummyMessage, 32)
	cep.Peers = make(map[PartyID]*DummyRemote, len(lp.Peers))
	for i, rp := range lp.Peers {
		cep.Peers[i] = &DummyRemote{
			RemoteParty: rp,
			Chan:        make(chan DummyMessage, 32),
		}
	}

	cep.Input = input
	cep.Output = input
	return cep
}

//BindNetwork connects parties
func (cep *DummyProtocol) BindNetwork(nw *TCPNetworkStruct) {
	for partyID, conn := range nw.Conns {

		if partyID == cep.ID {
			continue
		}

		rp := cep.Peers[partyID]

		// Receiving loop from remote
		go func(conn net.Conn, rp *DummyRemote) {
			for {
				var id, val uint64
				var err error
				err = binary.Read(conn, binary.BigEndian, &id)
				check(err)
				err = binary.Read(conn, binary.BigEndian, &val)
				check(err)
				msg := DummyMessage{
					Party: PartyID(id),
					Value: val,
				}
				//fmt.Println(cep, "receiving", msg, "from", rp)
				cep.Chan <- msg
			}
		}(conn, rp)

		// Sending loop of remote
		go func(conn net.Conn, rp *DummyRemote) {
			var m DummyMessage
			var open = true
			for open {
				m, open = <-rp.Chan
				//fmt.Println(cep, "sending", m, "to", rp)
				check(binary.Write(conn, binary.BigEndian, m.Party))
				check(binary.Write(conn, binary.BigEndian, m.Value))
			}
		}(conn, rp)
	}
}

//Run runs SMC protocol
func (cep *DummyProtocol) Run() {

	fmt.Println(cep, "is running")
	rand.Seed(time.Now().UTC().UnixNano() + int64(cep.ID))
	var s = int64(math.Pow(2, 16)) + 1 //prime number used for modulus ring
	fmt.Println(s)
	var secretshares = make([]uint64, len(cep.Peers))
	//get N-1 random values
	var tot uint64 = 0
	for i := range cep.Peers {
		if i != (cep.ID) {
			secretshares[i] = uint64(rand.Int63n(int64(s)))

			tot = tot + secretshares[i]

		}
	}
	secretshares[cep.ID] = uint64(mod(int64(cep.Input)-int64(tot), int64(s)))

	//fmt.Println("we're making secrets! ", secretshares, "at party ", cep.ID, " total was ", tot, " and input was ", cep.Input)
	for i, peer := range cep.Peers {
		if peer.ID != cep.ID {
			peer.Chan <- DummyMessage{cep.ID, secretshares[i]}
		}
	}

	received := make(map[PartyID]uint64)
	received[cep.ID] = secretshares[cep.ID]
	for m := range cep.Chan {
		received[m.Party] = m.Value
		if len(received) == len(cep.Peers) {
			//fmt.Println(cep, "received is ", received)

			//evaluate circuit in gates.go
			evaluate(cep, received, s)
			close(cep.Chan)
		}
	}
	if cep.WaitGroup != nil {
		cep.WaitGroup.Done()

	}

}
