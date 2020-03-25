package main

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/ldsec/lattigo/bfv"
)

//BeaverProtocol stores all data that is reused between 2 runs
type BeaverProtocol struct {
	*LocalParty
	Chan  chan BeaverMessage
	Peers map[PartyID]*BeaverRemoteParty

	c []uint64
	n uint64
	//elements in ring
	a      *bfv.Plaintext
	b      *bfv.Plaintext
	params *bfv.Parameters
	sk     *bfv.SecretKey
}

//BeaverRemoteParty sends beaver Messages (ciphertexts) over chans
type BeaverRemoteParty struct {
	*RemoteParty
	Chan chan BeaverMessage
}

//BeaverMessage the value in message passed is a ring element
type BeaverMessage struct {
	Party PartyID
	d     bfv.Ciphertext
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
	bep.Chan = make(chan BeaverMessage, 32)
	bep.Peers = make(map[PartyID]*BeaverRemoteParty, len(lp.Peers))

	for i, brp := range lp.Peers {
		bep.Peers[i] = &BeaverRemoteParty{
			RemoteParty: brp,
			Chan:        make(chan BeaverMessage, 32),
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

	//bep.b = bep.params.NewPolyQ()
	//bep.b.SetCoefficients([][]uint64{b})
	encoder := bfv.NewEncoder(bep.params)

	bep.a = bfv.NewPlaintext(bep.params)
	bep.b = bfv.NewPlaintext(bep.params)

	encoder.EncodeUint(a, bep.a)
	encoder.EncodeUint(b, bep.b)
	//prepare encryption

	kgen := bfv.NewKeyGenerator(bep.params)
	bep.sk, _ = kgen.GenKeyPair()

	return bep
}

//BeaverRun runs beaver prot
func (bep *BeaverProtocol) BeaverRun() {

	evaluator := bfv.NewEvaluator(bep.params)
	encryptorSk := bfv.NewEncryptorFromSk(bep.params, bep.sk)
	encoder := bfv.NewEncoder(bep.params)
	plaintext := bfv.NewPlaintext(bep.params)
	//encoder.EncodeUint(bep.a.Coeffs[0], plaintext)

	ciphertext := encryptorSk.EncryptNew(bep.a)
	fmt.Println("made cipher")
	for _, peer := range bep.Peers {
		if peer.ID != bep.ID {
			peer.Chan <- BeaverMessage{bep.ID, *ciphertext}
		}
	}
	fmt.Println("sent all cipher")
	received := make(map[PartyID]bfv.Ciphertext)
	/*   foreach other party j do
	receive dj from j
	r_ij <- Z^n _t
	c_i = c_i - r_ij
	encode r_ ij = to ring R
	(e^0_ij, e^ij) <- xi _err in R^2
	d_ij = Add(Mul(d_j, bi), r_ij) + (e^o _ij, e^ _ij)
	send d_ij to Pj back

	*/
	for m := range bep.Chan {
		received[m.Party] = m.d
		if len(received) == len(bep.Peers)-1 {
			fmt.Println(bep, "received is ", received)
		}
		r := NewRandomVec(bep.n, bep.params.T)
		bep.c = SubVec(&bep.c, &r, bep.params.T)
		//encR := bep.params.NewPolyQ()
		//encR.SetCoefficients([][]uint64{r})
		encoder.EncodeUint(r, plaintext)
		evaluator.Mul(&m.d, bep.b, &m.d)
		evaluator.Add(&m.d, plaintext, &m.d) //how do we add noise?

		//have to send back to same guy
		bep.Peers[m.Party].Chan <- BeaverMessage{bep.ID, m.d}
	}

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

//BeaverBindNetwork need to send BeaverMessage, not Message
func (bep *BeaverProtocol) BeaverBindNetwork(nw *TCPNetworkStruct) {
	for partyID, conn := range nw.Conns {

		if partyID == bep.ID {
			continue
		}

		brp := bep.Peers[partyID]

		// Receiving loop from remote
		go func(conn net.Conn, brp *BeaverRemoteParty) {
			for {
				var id uint64
				var cipher bfv.Ciphertext
				var err error
				err = binary.Read(conn, binary.BigEndian, &id)
				check(err)
				err = binary.Read(conn, binary.BigEndian, &cipher)
				check(err)
				msg := BeaverMessage{
					Party: PartyID(id),
					d:     cipher,
				}
				//fmt.Println(cep, "receiving", msg, "from", rp)
				bep.Chan <- msg
			}
		}(conn, brp)

		// Sending loop of remote
		go func(conn net.Conn, brp *BeaverRemoteParty) {
			var m BeaverMessage
			var open = true
			for open {
				m, open = <-brp.Chan
				//fmt.Println(cep, "sending", m, "to", rp)
				check(binary.Write(conn, binary.BigEndian, m.Party))
				check(binary.Write(conn, binary.BigEndian, m.d))
			}
		}(conn, brp)
	}
}
