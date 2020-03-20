package main

import {

}

//each Party i does
/*
	a_i <- Z^n _t 		sampling of a uniform vector of integers mod p of n elements
	b_i <- Z^n _t
	c_i = a_i x b_i (component wise)
	encode a_i to ring R = Z_Q[X]/X^N + 1
	encode b_i to ring R
	sk_i <- xi_(1/3)    sampling of elem of R with ternary coeff [-1 0 1] with distrib [(2/6), 1/3, 2/6]

	di = Enc_sk_i (a_i) in R^2
	foreach other party j
		send di to j
	

*/
//each Party i does
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

//BeaverProtocol stores all data that is reused between 2 runs
type BeaverProtocol struct{
	*LocalParty
	Chan  chan Message
	Peers map[PartyID]*Remote

	c Vector
	//elements in ring
	a Poly
	b Poly

	sk uint64 


}


type BeaverRemoteParty struct {
	*RemoteParty
	Chan chan Message
}

type BeaverMessage struct {
	Party PartyID
	d Poly
}

type BeaverInputs struct {
	//inputs are the BFV scheme parameters
	n uint64 //degree of R
	s int64 //plaintext modulus
}

func  (lp *LocalParty) New() *BeaverProtocol {

}

func Run() {

}

func BindNetwork() {

}