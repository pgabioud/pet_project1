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
	//elements in ring
	a      *ring.Poly
	b      *ring.Poly
	params *bfv.Parameters
	sk     uint64
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
	d     ring.Poly
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
	n := uint64(1 << bep.params.LogN)
	T := bep.params.T
	a := NewRandomVec(n, T)
	b := NewRandomVec(n, T)
	bep.c = MulVec(&a, &b, T)

	//convert to ring element
	//JE CROIS
	Q := uint64(bep.params.LogQP())
	bep.a = ring.NewPoly(n, Q)
	bep.a.SetCoefficients([][]uint64{a})
	bep.b = ring.NewPoly(n, Q)
	fmt.Println(bep.b.GetDegree())
	fmt.Printf("initialization done with n being %d, T being %d and Q is %d \n", n, T, Q)
	//sk <- xi_(1/3)
	context := ring.NewContext()
	fmt.Println(1.0 / 3.0)
	sk := context.SampleTernaryNew(1.0 / 3.0)
	fmt.Println("sk is sampled ", sk)
	return bep
}

//BeaverRun runs beaver prot
func BeaverRun() {

	/*
		di = Enc_sk_i (a_i) in R^2
		foreach other party j
		send di to j
	*/

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
