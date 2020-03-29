package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"time"

	"github.com/ldsec/lattigo/bfv"
	"github.com/ldsec/lattigo/ring"
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
	pk     *bfv.PublicKey
}

//BeaverRemoteParty sends beaver Messages (ciphertexts) over chans
type BeaverRemoteParty struct {
	*RemoteParty
	Chan chan BeaverMessage
}

//BeaverMessage the value in message passed is a ring element
type BeaverMessage struct {
	Party   PartyID
	Size    uint64
	Marshal []byte
	TypeM   uint8
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
	fmt.Println("party ", bep.ID, "made a = ", a[:3], " and b = ", b[:3])
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
	bep.sk, bep.pk = kgen.GenKeyPair()

	return bep
}

//BeaverRun runs beaver prot
func (bep *BeaverProtocol) BeaverRun() {

	evaluator := bfv.NewEvaluator(bep.params)
	encryptorPk := bfv.NewEncryptorFromPk(bep.params, bep.pk)
	encoder := bfv.NewEncoder(bep.params)
	decryptorSk := bfv.NewDecryptor(bep.params, bep.sk)

	//encoder.EncodeUint(bep.a.Coeffs[0], plaintext)

	ciphertext := encryptorPk.EncryptNew(bep.a)

	msg, err := ciphertext.MarshalBinary()
	check(err)
	for _, peer := range bep.Peers {
		if peer.ID != bep.ID {
			time.Sleep(100 * time.Millisecond)
			peer.Chan <- BeaverMessage{bep.ID, uint64(len(msg)), msg, 0}
		}
	}

	//fmt.Println("sent all cipher")
	received := make(map[PartyID]*bfv.Ciphertext)
	/*   foreach other party j do
	receive dj from j
	r_ij <- Z^n _t
	c_i = c_i - r_ij
	encode r_ ij = to ring R
	(e^0_ij, e^ij) <- xi _err in R^2
	d_ij = Add(Mul(d_j, bi), r_ij) + (e^o _ij, e^ _ij)
	send d_ij to Pj back

	*/
	var cPrime []uint64
	for i := uint64(0); i < bep.n; i++ {
		cPrime = append(cPrime, uint64(0))
	}

	plainCPrime := bfv.NewPlaintext(bep.params)
	encoder.EncodeUint(cPrime, plainCPrime)
	cipherCPrime := encryptorPk.EncryptNew(plainCPrime)
	counter := 0

	for m := range bep.Chan {
		//fmt.Println("party", bep.ID, " channel length is ", len(bep.Chan))
		//fmt.Println("type of m = ", m.TypeM)
		d := bfv.NewCiphertext(bep.params, 1)
		//fmt.Println(d.Degree())
		err := d.UnmarshalBinary(m.Marshal)
		check(err)
		if m.TypeM == 0 {
			received[m.Party] = d
			//fmt.Println("len of received = ", len(received))
			if len(received) == len(bep.Peers)-1 {
				//fmt.Println(bep, "received is ", received)
			}
			r := NewRandomVec(bep.n, bep.params.T)
			bep.c = SubVec(&bep.c, &r, bep.params.T)
			//encR := bep.params.NewPolyQ()
			//encR.SetCoefficients([][]uint64{r})
			plaintext := bfv.NewPlaintext(bep.params)

			context, _ := ring.NewContextWithParams(bep.n, bep.params.Qi)
			e0 := context.SampleGaussianNew(bep.params.Sigma, uint64(math.Floor(6.0*bep.params.Sigma)))
			e1 := context.SampleGaussianNew(bep.params.Sigma, uint64(math.Floor(6.0*bep.params.Sigma)))
			gaussian := bfv.NewCiphertext(bep.params, 1)
			gaussian.SetValue([]*ring.Poly{e0, e1})
			//encoder.EncodeUint(e0, e1, plaintext)
			encoder.EncodeUint(r, plaintext)

			evaluator.Mul(d, bep.b, d)

			evaluator.Add(d, plaintext, d) //how do we add noise?
			evaluator.Add(d, gaussian, d)
			fmt.Println("degrees: (d, b, r, gauss", d.Degree(), bep.b.Degree(), plaintext.Degree(), gaussian.Degree())
			//have to send back to same party
			msg, err := d.MarshalBinary()
			check(err)
			//time.Sleep(100 * time.Millisecond)
			//fmt.Println("my party is ", bep.Party, "sending response to ", m.Party)
			bep.Peers[m.Party].Chan <- BeaverMessage{bep.ID, uint64(len(msg)), msg, 1}
		} else {
			counter++
			evaluator.Add(cipherCPrime, d, cipherCPrime)
			//fmt.Println(bep.ID, " counter is ", counter)
			if counter == len(bep.Peers)-1 {
				//fmt.Println("close")
				close(bep.Chan)
			}
		}
	}

	decryptCipher := decryptorSk.DecryptNew(cipherCPrime)
	decodedCypher := encoder.DecodeUint(decryptCipher)
	bep.c = AddVec(&bep.c, &decodedCypher, bep.params.T)

	fmt.Println("party ", bep.ID, "managed to make c equal to ", bep.c[:3])
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
	if bep.WaitGroup != nil {
		bep.WaitGroup.Done()
	}

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
				var typeMsg uint8
				var sizeCypher uint64
				var err error

				err = binary.Read(conn, binary.BigEndian, &id)
				check(err)
				err = binary.Read(conn, binary.BigEndian, &sizeCypher)
				check(err)
				cipher := make([]byte, sizeCypher)
				err = binary.Read(conn, binary.BigEndian, &cipher)
				check(err)
				err = binary.Read(conn, binary.BigEndian, &typeMsg)
				check(err)

				msg := BeaverMessage{
					Party:   PartyID(id),
					Size:    sizeCypher,
					Marshal: cipher,
					TypeM:   typeMsg,
				}
				//fmt.Println(bep, "receiving", msg.Party, len(msg.Marshal), msg.TypeM, "from", brp)
				bep.Chan <- msg
			}
		}(conn, brp)

		// Sending loop of remote
		go func(conn net.Conn, brp *BeaverRemoteParty) {
			var m BeaverMessage
			var open = true
			for open {
				m, open = <-brp.Chan

				check(binary.Write(conn, binary.BigEndian, m.Party))
				check(binary.Write(conn, binary.BigEndian, m.Size))
				check(binary.Write(conn, binary.BigEndian, m.Marshal))
				check(binary.Write(conn, binary.BigEndian, m.TypeM))
			}
		}(conn, brp)
	}
}
