package main

//WireID id for wire
type WireID uint64

//GateID id for gate
type GateID uint64

//Operation interface to have function on gate
type Operation interface {
	Output() WireID
	Inputs() []WireID
}

//Input gate
type Input struct {
	Party PartyID
	Out   WireID
}

//Output getter for input
func (io Input) Output() WireID {
	return io.Out
}

//Inputs getter for input
func (io Input) Inputs() []WireID {
	return []WireID{WireID(io.Party)}
}

//Add gate
type Add struct {
	In1 WireID
	In2 WireID
	Out WireID
}

//Output getter for add gate
func (ao Add) Output() WireID {
	return ao.Out
}

//Inputs getter for add
func (ao Add) Inputs() []WireID {
	return []WireID{ao.In1, ao.In2}
}

//AddCst gate
type AddCst struct {
	In       WireID
	CstValue uint64
	Out      WireID
}

//Output getter for addcst gate
func (aco AddCst) Output() WireID {
	return aco.Out
}

//Inputs getter for addcst
func (aco AddCst) Inputs() []WireID {
	return []WireID{aco.In}
}

//Sub gate
type Sub struct {
	In1 WireID
	In2 WireID
	Out WireID
}

//Output getter for sub gate
func (so Sub) Output() WireID {
	return so.Out
}

//Inputs getter for sub
func (so Sub) Inputs() []WireID {
	return []WireID{so.In1, so.In2}
}

//Mult gate
type Mult struct {
	In1 WireID
	In2 WireID
	Out WireID
}

//Output getter for Mult gate
func (mo Mult) Output() WireID {
	return mo.Out
}

//Inputs getter for mult
func (mo Mult) Inputs() []WireID {
	return []WireID{mo.In1, mo.In2}
}

//MultCst gate
type MultCst struct {
	In       WireID
	CstValue uint64
	Out      WireID
}

//Output getter for MultCst gate
func (mco MultCst) Output() WireID {
	return mco.Out
}

//Inputs getter for multcst
func (mco MultCst) Inputs() []WireID {
	return []WireID{mco.In}
}

//Reveal gate
type Reveal struct {
	In  WireID
	Out WireID
}

//Output getter for reveal gate
func (ro Reveal) Output() WireID {
	return ro.Out
}

//Inputs getter for reveal
func (ro Reveal) Inputs() []WireID {
	return []WireID{ro.In}
}
