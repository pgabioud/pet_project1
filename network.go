package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
)

//CONNECT_ATTEMPTS is the maximum number of connection attempts
const CONNECT_ATTEMPTS = 4

//CONNECT_ATTEMPTS_DELAY is the delay between two connection attempts
const CONNECT_ATTEMPTS_DELAY = 100

//Network simple interface to connect a local party
type Network interface {
	Connect(party *LocalParty) error
}

//TCPNetworkStruct uses mutex locks for connections
type TCPNetworkStruct struct {
	Conns    map[PartyID]net.Conn
	connLock sync.RWMutex

	ready sync.WaitGroup
}

//NewTCPNetwork creates new tcp network struct
func NewTCPNetwork(party *LocalParty) (*TCPNetworkStruct, error) {
	netw := &TCPNetworkStruct{}
	netw.Conns = make(map[PartyID]net.Conn, len(party.Peers))
	return netw, nil
}

//Connect connects local party with tcp network
func (tnw *TCPNetworkStruct) Connect(lp *LocalParty) error {
	//var err error
	waitFor, dialFor := make(map[PartyID]*RemoteParty), make(map[PartyID]*RemoteParty)

	for _, rp := range lp.Peers {
		if lp.ID > rp.ID {
			waitFor[rp.ID] = rp
		}
		if lp.ID < rp.ID {
			dialFor[rp.ID] = rp
		}
	}

	//fmt.Println(lp, "dialFor:", dialFor, "waitFor", waitFor)

	tnw.ready.Add(len(waitFor) + len(dialFor))

	go func() {
		listener, err := net.Listen("tcp", lp.Addr)
		if err != nil {
			panic(fmt.Errorf("cannot create listening socket: %s", err))
		}
		//fmt.Println(lp, "now listening on", listener.Addr())

		for range waitFor {
			conn, err := listener.Accept()
			check(err)
			var partyID PartyID
			check(binary.Read(conn, binary.BigEndian, &partyID))
			if rp, known := waitFor[partyID]; known {
				fmt.Println(lp, "now connected with", rp)

				tnw.connLock.Lock()
				tnw.Conns[partyID] = conn
				tnw.connLock.Unlock()

				tnw.ready.Done()
			} else {
				panic(fmt.Sprintf("%s: unexpected party ID: %d", lp, partyID))
			}
		}
		check(listener.Close())
	}()

	//<- time.After(time.Second)

	for _, p := range dialFor {
		go func(rp *RemoteParty) {
			var conn net.Conn
			var err error
			for attempt := 0; conn == nil && attempt < CONNECT_ATTEMPTS; attempt++ {
				if attempt > 0 {
					//fmt.Println("retrying:", err)
					<-time.After(CONNECT_ATTEMPTS_DELAY * time.Millisecond)
				}
				conn, err = net.Dial("tcp", rp.Addr)
			}
			if conn == nil {
				fmt.Println(lp, "couldn't connect to", rp, ":", err)
			}
			tnw.connLock.Lock()
			tnw.Conns[rp.ID] = conn
			tnw.connLock.Unlock()
			check(binary.Write(conn, binary.BigEndian, lp.ID))
			tnw.ready.Done()
		}(p)
	}

	tnw.ready.Wait()
	return nil
}

//GetTestingTCPNetwork check if can connect
func GetTestingTCPNetwork(P []*LocalParty) []*TCPNetworkStruct {
	var err error
	netws := make([]*TCPNetworkStruct, len(P), len(P))
	for i, lp := range P {
		netws[i], err = NewTCPNetwork(lp)
		check(err)
	}

	wgc := &sync.WaitGroup{}
	for i, lp := range P {
		wgc.Add(1)
		go func(netw *TCPNetworkStruct, lp *LocalParty) {
			err = netw.Connect(lp)
			check(err)
			wgc.Done()
		}(netws[i], lp)
	}
	wgc.Wait()
	return netws
}
