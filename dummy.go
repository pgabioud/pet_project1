package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"net"
	"time"
)

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

type DummyProtocol struct {
	*LocalParty
	Chan  chan DummyMessage
	Peers map[PartyID]*DummyRemote

	circuitID CircuitID

	Input  uint64
	Output uint64
}

type DummyRemote struct {
	*RemoteParty
	Chan chan DummyMessage
}

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

//func (cep *DummyProtocol) Run(circuit *TestCircuit) {
func (cep *DummyProtocol) Run() {

	fmt.Println(cep, "is running")
	rand.Seed(time.Now().UTC().UnixNano() + int64(cep.ID))
	var s = uint64(math.Pow(2, 16)) + 1 //prime number used for modulus ring
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

	fmt.Println("we're making secrets! ", secretshares, "at party ", cep.ID, " total was ", tot, " and input was ", cep.Input)
	for i, peer := range cep.Peers {
		if peer.ID != cep.ID {
			peer.Chan <- DummyMessage{cep.ID, secretshares[i]}
		}
	}

	received := make(map[PartyID]uint64)
	received[cep.ID] = secretshares[cep.ID]
	for m := range cep.Chan {
		fmt.Println(cep, "received message from", m.Party, ":", m.Value)

		received[m.Party] = m.Value
		fmt.Println(cep, "received is ", received)
		if len(received) == len(cep.Peers) {
			wire := make([]uint64, len(TestCircuits[cep.circuitID-1].Circuit))
			fmt.Println("circuitID = ", cep.circuitID)
			for _, op := range TestCircuits[cep.circuitID-1].Circuit {
				fmt.Println(op)
				switch op.(type) {
				case *Input:
					wire[op.Output()] = secretshares[op.(*Input).Party]
				case *Add:
					wire[op.Output()] = uint64(mod(int64(wire[op.(*Add).In1])+int64(wire[op.(*Add).In2]), int64(s)))
				case *Reveal:
					cep.Output = wire[op.(*Reveal).In]
					for _, peer := range cep.Peers {
						if peer.ID != cep.ID {
							peer.Chan <- DummyMessage{cep.ID, cep.Output}
						}
					}
					received := 0
					for m := range cep.Chan {
						cep.Output = uint64(mod(int64(m.Value)+int64(cep.Output), int64(s)))
						received++
						if received == len(cep.Peers)-1 {
							close(cep.Chan)
						}
					}

					fmt.Println("HUZZAH")
				default:
					fmt.Println("op not implemented or does not exist")
				}
			}
		}
	}
	if cep.WaitGroup != nil {
		cep.WaitGroup.Done()
	}
}
