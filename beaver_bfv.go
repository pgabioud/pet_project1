package main

import (
	"fmt"

	"github.com/ldsec/lattigo/bfv"
	"github.com/ldsec/lattigo/ring"
)

//BeaverProtocol stores all data that is reused between 2 runs
type BeaverProtocol struct {
	*LocalParty
	Chan  chan Message
	Peers map[PartyID]*Remote

	c []uint64
	n uint64
	//elements in ring
	a      *ring.Poly
	b      *ring.Poly
	params *bfv.Parameters
	sk     *bfv.SecretKey
}

//!!!same as remoteParty!!!
/*
type BeaverRemoteParty struct {
	*RemoteParty
	Chan chan Message
}
*/

//BeaverMessage the value in message passed is a ring element
type BeaverMessage struct {
	Party PartyID
	d     *bfv.Ciphertext
}

//BeaverInputs are the BFV scheme parameters
type BeaverInputs struct {
	n uint64 //degree of R
	s int64  //plaintext modulus
}

//New beaver protocol, creates the protocol
func (lp *LocalParty) New() *BeaverProtocol {

	bep := new(BeaverProtocol)
	bep.LocalParty = lp
	bep.Chan = make(chan Message, 32)
	bep.Peers = make(map[PartyID]*Remote, len(lp.Peers))

	for i, rp := range lp.Peers {
		bep.Peers[i] = &Remote{
			RemoteParty: rp,
			Chan:        make(chan Message, 32),
		}
	}

	bep.params = bfv.DefaultParams[bfv.PN13QP218]
	/*
		a_i <- Z^n _t
		b_i <- Z^n _t
		c_i = a_i x b_i
	*/
	bep.n = uint64(1 << bep.params.LogN)
	T := bep.params.T
	a := NewRandomVec(bep.n, T)
	b := NewRandomVec(bep.n, T)
	bep.c = MulVec(&a, &b, T)

	//convert to ring element

	bep.a = bep.params.NewPolyQ()
	bep.a.SetCoefficients([][]uint64{a})
	bep.b = bep.params.NewPolyQ()
	bep.b.SetCoefficients([][]uint64{b})

	//prepare encryption

	kgen := bfv.NewKeyGenerator(bep.params)
	bep.sk, _ = kgen.GenKeyPair()

	return bep
}

//BeaverRun runs beaver prot
func (bep *BeaverProtocol) BeaverRun() {

	encryptorSk := bfv.NewEncryptorFromSk(bep.params, bep.sk)
	encoder := bfv.NewEncoder(bep.params)
	plaintext := bfv.NewPlaintext(bep.params)
	encoder.EncodeUint(bep.a.Coeffs[0], plaintext)

	cyphertext := encryptorSk.EncryptNew(plaintext)

	for _, peer := range bep.Peers {
		if peer.ID != bep.ID {
			peer.Chan <- Message{bep.ID, cyphertext}
		}
	}

	received := make(map[PartyID]uint64)
	for m := range bep.Chan {
		received[m.Party] = m.Value
		if len(received) == len(bep.Peers)-1 {
			fmt.Println(bep, "received is ", received)
		}
		r := NewRandomVec(bep.n, bep.params.T)
		cyphertext = SubVec(cyphertext, &r)
		encR := bep.params.NewPolyQ()
		encR.SetCoefficients([][]uint64{r})

	}

	/*   foreach other party j do
	receive dj from j
	r_ij <- Z^n _t
	c_i = c_i - r_ij
	encode r_ ij = to ring R
	(e^0_ij, e^ij) <- xi _err in R^2
	d_ij = Add(Mul(d_j, bi), r_ij) + (e^o _ij, e^ _ij)
	send d_ij to Pj

	*/
	//each Party i does
	/*
		c' = (0,0) in R^2
		foreach other party j
			d _ji receive from j
			c' = add(c', d_ij)

		c' = Dec_sk(c')
		decode from R c'
		c_i = c_i + c'
	*/
	return
}

//!!!same as given!!!
/*
func BindNetwork() {
	return
}
*/
