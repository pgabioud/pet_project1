package main

import "fmt"

//evaluates different gates for cep using circuitID and shares in modulus ring s
func evaluate(cep *DummyProtocol, secrets map[PartyID]uint64, s uint64) {
	wire := make([]uint64, len(TestCircuits[cep.circuitID-1].Circuit))
	for _, op := range TestCircuits[cep.circuitID-1].Circuit {
		switch op.(type) {

		case *Input:
			wire[op.Output()] = secrets[op.(*Input).Party]

		case *Add:
			wire[op.Output()] = uint64(mod(int64(wire[op.(*Add).In1])+int64(wire[op.(*Add).In2]), int64(s)))

		case *Sub:
			wire[op.Output()] = uint64(mod(int64(wire[op.(*Sub).In1])-int64(wire[op.(*Sub).In2]), int64(s)))

		case *MultCst:
			wire[op.Output()] = uint64(mod(int64(wire[op.(*MultCst).In])*int64(op.(*MultCst).CstValue), int64(s)))

		case *AddCst:
			if cep.ID == 0 {
				wire[op.Output()] = uint64(mod(int64(wire[op.(*AddCst).In])+int64(op.(*AddCst).CstValue), int64(s)))
			} else {
				wire[op.Output()] = wire[op.(*AddCst).In]
			}

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

		default:
			fmt.Println("op not implemented or does not exist")
		}
	}
}
