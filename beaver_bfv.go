package main

import (
	"encoding/binary"
	"math"
	"net"

	"github.com/ldsec/lattigo/bfv"
	"github.com/ldsec/lattigo/ring"
)

//BeaverProtocol stores all data that is reused between 2 runs
type BeaverProtocol struct {
	*LocalParty
	Chan  chan BeaverMessage
	Peers map[PartyID]*BeaverRemoteParty

	n uint64
	//elements in ring
	a      []uint64
	b      []uint64
	c      []uint64
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

//NewBeaverProtocol beaver protocol, creates the protocol
func (lp *LocalParty) NewBeaverProtocol() *BeaverProtocol {

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
	bep.n = uint64(1 << bep.params.LogN)
	T := bep.params.T

	bep.a = NewRandomVec(bep.n, T)
	bep.b = NewRandomVec(bep.n, T)
	//fmt.Println("party ", bep.ID, "made a = ", bep.a[:3], " and b = ", bep.b[:3])
	bep.c = MulVec(&bep.a, &bep.b, T)

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

	//convert to ring element
	plainA := bfv.NewPlaintext(bep.params)
	plainB := bfv.NewPlaintext(bep.params)

	encoder.EncodeUint(bep.a, plainA)
	encoder.EncodeUint(bep.b, plainB)

	cipherA := encryptorPk.EncryptNew(plainA)

	msg, err := cipherA.MarshalBinary()
	check(err)

	for _, peer := range bep.Peers {
		if peer.ID != bep.ID {
			peer.Chan <- BeaverMessage{bep.ID, uint64(len(msg)), msg, 0}
		}
	}
	//fmt.Println("sent all cipher")

	received := make(map[PartyID]*bfv.Ciphertext)
	var cipherCPrime bfv.Ciphertext
	counter := 0

	for m := range bep.Chan {
		//fmt.Println("party", bep.ID, " channel length is ", len(bep.Chan))
		//fmt.Println("type of m = ", m.TypeM)
		d := bfv.NewCiphertext(bep.params, 1)
		err := d.UnmarshalBinary(m.Marshal)
		check(err)

		if m.TypeM == 0 {
			received[m.Party] = d
			r := NewRandomVec(bep.n, bep.params.T)
			bep.c = SubVec(&bep.c, &r, bep.params.T)

			plainR := bfv.NewPlaintext(bep.params)

			context, _ := ring.NewContextWithParams(bep.n, bep.params.Qi)
			e0 := context.SampleGaussianNew(bep.params.Sigma, uint64(math.Floor(6.0*bep.params.Sigma)))
			e1 := context.SampleGaussianNew(bep.params.Sigma, uint64(math.Floor(6.0*bep.params.Sigma)))
			gaussian := bfv.NewCiphertext(bep.params, 1)
			gaussian.SetValue([]*ring.Poly{e0, e1})
			encoder.EncodeUint(r, plainR)

			evaluator.Mul(d, plainB, d)
			evaluator.Add(d, plainR, d)
			evaluator.Add(d, gaussian, d)
			//fmt.Println("degrees: (d, b, r, gauss)", d.Degree(), plainB.Degree(), plainR.Degree(), gaussian.Degree())

			//sent back to same party
			msg, err := d.MarshalBinary()
			check(err)

			//fmt.Println("my party is ", bep.Party, "sending response to ", m.Party)
			bep.Peers[m.Party].Chan <- BeaverMessage{bep.ID, uint64(len(msg)), msg, 1}
		} else {
			if counter == 0 {
				cipherCPrime = *d
			} else {
				evaluator.Add(cipherCPrime, d, &cipherCPrime)
			}
			counter++
			if counter == len(bep.Peers)-1 {
				//fmt.Println("close")
				close(bep.Chan)
			}
		}
	}

	decryptCipher := decryptorSk.DecryptNew(&cipherCPrime)
	decodedCypher := encoder.DecodeUint(decryptCipher)
	bep.c = AddVec(&bep.c, &decodedCypher, bep.params.T)

	//fmt.Println("party ", bep.ID, "managed to make c equal to ", bep.c[:3])

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
			//fmt.Println(bep, "writing", m.Party, len(m.Marshal), m.TypeM, "from", brp)
		}(conn, brp)
	}
}
