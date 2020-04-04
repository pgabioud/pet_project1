package main

/*
import (
	"encoding/binary"
	"math"
	"math/rand"
	"net"
	"time"
)

//CircuitID uint to select the circuit to run
type CircuitID uint64

func mod(a, b int64) uint64 {
	m := a % b
	if a < 0 && b < 0 {
		m -= b
	}
	if a < 0 && b > 0 {
		m += b
	}
	return uint64(m)
}

//Message is simple value message for party
type Message struct {
	Party PartyID
	Value uint64
}

//Protocol structure for a protocol
type Protocol struct {
	*LocalParty
	Chan  chan Message
	Peers map[PartyID]*Remote

	circuitID      CircuitID
	beaverProtocol *BeaverProtocol

	Input  uint64
	Output uint64
}

//Remote template for remote party
type Remote struct {
	*RemoteParty
	Chan chan Message
}

//NewProtocol initializes SMC for defined circuit with ID = circuitID
func (lp *LocalParty) NewProtocol(input uint64, circuitID CircuitID, beaverProtocol *BeaverProtocol) *Protocol {
	cep := new(Protocol)
	cep.LocalParty = lp
	cep.circuitID = circuitID
	cep.Chan = make(chan Message, 32)
	cep.Peers = make(map[PartyID]*Remote, len(lp.Peers))

	cep.beaverProtocol = beaverProtocol

	for i, rp := range lp.Peers {
		cep.Peers[i] = &Remote{
			RemoteParty: rp,
			Chan:        make(chan Message, 32),
		}
	}

	cep.Input = input
	cep.Output = input
	return cep
}

//BindNetwork connects parties
func (cep *Protocol) BindNetwork(nw *TCPNetworkStruct) {
	for partyID, conn := range nw.Conns {

		if partyID == cep.ID {
			continue
		}

		rp := cep.Peers[partyID]

		// Receiving loop from remote
		go func(conn net.Conn, rp *Remote) {
			for {
				var id, val uint64
				var err error
				err = binary.Read(conn, binary.BigEndian, &id)
				check(err)
				err = binary.Read(conn, binary.BigEndian, &val)
				check(err)
				msg := Message{
					Party: PartyID(id),
					Value: val,
				}
				//fmt.Println(cep, "receiving", msg, "from", rp)
				cep.Chan <- msg
			}
		}(conn, rp)

		//Sending loop of remote
		go func(conn net.Conn, rp *Remote) {
			var m Message
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
func (cep *Protocol) Run() {

	rand.Seed(time.Now().UTC().UnixNano() + int64(cep.ID))
	var s = int64(math.Pow(2, 16)) + 1 //prime number used for modulus ring
	var secretshares = make([]uint64, len(cep.Peers))
	//get N-1 random values
	var tot uint64 = 0
	for i := range cep.Peers {
		if i != (cep.ID) {
			secretshares[i] = uint64(rand.Int63n(s))

			tot = tot + secretshares[i]

		}
	}
	secretshares[cep.ID] = uint64(mod(int64(cep.Input)-int64(tot), s))

	//fmt.Println("we're making secrets! ", secretshares, "at party ", cep.ID, " total was ", tot, " and input was ", cep.Input)
	for i, peer := range cep.Peers {
		if peer.ID != cep.ID {
			peer.Chan <- Message{cep.ID, secretshares[i]}
		}
	}

	received := make(map[PartyID]uint64)
	received[cep.ID] = secretshares[cep.ID]
	for m := range cep.Chan {
		received[m.Party] = m.Value
		if len(received) == len(cep.Peers) {
			//fmt.Println(cep, "received is ", received)

			//evaluate circuit in gates.go
			Evaluate(cep, &received, s)
			close(cep.Chan)
		}
	}
	if cep.WaitGroup != nil {
		cep.WaitGroup.Done()
	}
}
*/
