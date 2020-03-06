package main

type WireID uint64

type GateID uint64

type Operation interface {
	Output() WireID
	Inputs() []WireID
}

type Input struct {
	Party PartyID
	Out   WireID
}

func (io Input) Output() WireID {
	return io.Out
}

func (io Input) Inputs() []WireID {
	return []WireID{WireID(io.Party)}
}

type Add struct {
	In1 WireID
	In2 WireID
	Out WireID
}

func (ao Add) Output() WireID {
	return ao.Out
}
func (ao Add) Inputs() []WireID {
	return []WireID{ao.In1, ao.In2}
}

type AddCst struct {
	In       WireID
	CstValue uint64
	Out      WireID
}

func (aco AddCst) Output() WireID {
	return aco.Out
}
func (aco AddCst) Inputs() []WireID {
	return []WireID{aco.In}
}

type Sub struct {
	In1 WireID
	In2 WireID
	Out WireID
}

func (so Sub) Output() WireID {
	return so.Out
}
func (so Sub) Inputs() []WireID {
	return []WireID{so.In1, so.In2}
}

type Mult struct {
	In1 WireID
	In2 WireID
	Out WireID
}

func (mo Mult) Output() WireID {
	return mo.Out
}
func (mo Mult) Inputs() []WireID {
	return []WireID{mo.In1, mo.In2}
}

type MultCst struct {
	In       WireID
	CstValue uint64
	Out      WireID
}

func (mco MultCst) Output() WireID {
	return mco.Out
}
func (mco MultCst) Inputs() []WireID {
	return []WireID{mco.In}
}

type Reveal struct {
	In  WireID
	Out WireID
}

func (ro Reveal) Output() WireID {
	return ro.Out
}
func (ro Reveal) Inputs() []WireID {
	return []WireID{ro.In}
}
