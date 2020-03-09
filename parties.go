package main

import (
	"fmt"
	"sync"
)

//PartyID party id
type PartyID uint64

//Party structure wwwith id and addr
type Party struct {
	ID   PartyID
	Addr string
}

//LocalParty define a party that is local
type LocalParty struct {
	Party
	*sync.WaitGroup
	Peers map[PartyID]*RemoteParty
}

func check(err error) {
	if err != nil {
		panic(err.Error())
	}
}

//NewLocalParty local party initializer
func NewLocalParty(id PartyID, peers map[PartyID]string) (*LocalParty, error) {
	// Create a new local party from the peers map
	p := &LocalParty{}
	p.ID = id

	p.Peers = make(map[PartyID]*RemoteParty, len(peers))
	p.Addr = peers[id]

	var err error
	for pId, pAddr := range peers {
		p.Peers[pId], err = NewRemoteParty(pId, pAddr)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (lp *LocalParty) String() string {
	// Print the party number
	return fmt.Sprintf("party-%d", lp.ID)
}

//RemoteParty to define a remote party
type RemoteParty struct {
	Party
}

func (rp *RemoteParty) String() string {
	return fmt.Sprintf("party-%d", rp.ID)
}

//NewRemoteParty initialize a remote party
func NewRemoteParty(id PartyID, addr string) (*RemoteParty, error) {
	p := &RemoteParty{}
	p.ID = id
	p.Addr = addr
	return p, nil
}
