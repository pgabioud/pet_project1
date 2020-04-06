package main

import "fmt"

//Evaluate evaluates different gates for cep using circuitID and shares in modulus ring s
func Evaluate(cep *Protocol, secrets *map[PartyID]uint64, s int64) {
	wire := make([]uint64, len(TestCircuits[cep.circuitID-1].Circuit))
	for _, op := range TestCircuits[cep.circuitID-1].Circuit {
		switch op.(type) {

		case *Input:
			wire[op.Output()] = (*secrets)[op.(*Input).Party]

		case *Add:
			wire[op.Output()] = mod(int64(wire[op.(*Add).In1])+int64(wire[op.(*Add).In2]), s)

		case *Sub:
			wire[op.Output()] = mod(int64(wire[op.(*Sub).In1])-int64(wire[op.(*Sub).In2]), s)

		case *MultCst:
			wire[op.Output()] = mod(int64(wire[op.(*MultCst).In])*int64(op.(*MultCst).CstValue), s)

		case *AddCst:
			if cep.ID == 0 {
				wire[op.Output()] = mod(int64(wire[op.(*AddCst).In])+int64(op.(*AddCst).CstValue), s)
			} else {
				wire[op.Output()] = wire[op.(*AddCst).In]
			}

		case *Reveal:
			cep.Output = wire[op.(*Reveal).In]
			Revealgate(cep, s)

		case *Mult:
			//1. compute (x-a)
			cep.Output = mod(int64(wire[op.(*Mult).In1])-int64(cep.Beavers[0][0]), s)
			Revealgate(cep, s)
			xa := cep.Output
			//2. compute (y-b)
			cep.Output = mod(int64(wire[op.(*Mult).In2])-int64(cep.Beavers[0][1]), s)
			Revealgate(cep, s)
			yb := cep.Output
			//3. comp z = c + x*(y-b) + y*(x-a) - (x-a)*(y-b)
			z := mod(int64(cep.Beavers[0][2])+int64(wire[op.(*Mult).In1]*yb)+int64(wire[op.(*Mult).In2]*xa), s)
			if cep.ID == 0 {
				z = mod(int64(z)-int64(xa*yb), s)
			}
			wire[op.Output()] = z
			//remove beaver used beaver triplet
			cep.Beavers = cep.Beavers[:1]

		default:
			fmt.Println("op not implemented or does not exist")
		}
	}

}

//Revealgate gate
func Revealgate(cep *Protocol, s int64) {
	for _, peer := range cep.Peers {
		if peer.ID != cep.ID {
			peer.Chan <- Message{cep.ID, cep.Output}
		}
	}
	received := 0
	for m := range cep.Chan {
		cep.Output = uint64(mod(int64(m.Value)+int64(cep.Output), s))

		received++
		if received == len(cep.Peers)-1 {
			break
		}

	}
}

/*
//AddCnstGate gate
func AddCnstGate(cep *DummyProtocol, wire []uint64, op Operation, s int64) {
	if cep.ID == 0 {
		wire[op.Output()] = uint64(mod(int64(wire[op.(*AddCst).In])+int64(op.(*AddCst).CstValue), s))
	} else {
		wire[op.Output()] = wire[op.(*AddCst).In]
	}
}

*/
