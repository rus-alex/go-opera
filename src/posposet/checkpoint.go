package posposet

import (
	"github.com/Fantom-foundation/go-lachesis/src/hash"
	"github.com/Fantom-foundation/go-lachesis/src/inter"
	"github.com/Fantom-foundation/go-lachesis/src/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/src/logger"
	"github.com/Fantom-foundation/go-lachesis/src/posposet/wire"
)

// checkpoint is for persistent storing.
type checkpoint struct {
	SuperFrameN       idx.SuperFrame
	LastDecidedFrameN idx.Frame
	LastConsensusTime inter.Timestamp
	LastBlockN        idx.Block
	Genesis           hash.Hash
	TotalCap          inter.Stake
}

// ToWire converts to proto.Message.
func (cp *checkpoint) ToWire() *wire.Checkpoint {
	return &wire.Checkpoint{
		SuperFrameN:        uint64(cp.SuperFrameN),
		LastFinishedFrameN: uint32(cp.LastDecidedFrameN),
		LastBlockN:         uint64(cp.LastBlockN),
		Genesis:            cp.Genesis.Bytes(),
		TotalCap:           uint64(cp.TotalCap),
	}
}

// WireToState converts from wire.
func WireToCheckpoint(w *wire.Checkpoint) *checkpoint {
	if w == nil {
		return nil
	}
	return &checkpoint{
		SuperFrameN:       idx.SuperFrame(w.SuperFrameN),
		LastDecidedFrameN: idx.Frame(w.LastFinishedFrameN),
		LastBlockN:        idx.Block(w.LastBlockN),
		Genesis:           hash.FromBytes(w.Genesis),
		TotalCap:          inter.Stake(w.TotalCap),
	}
}

/*
 * Poset's methods:
 */

// State saves checkpoint.
func (p *Poset) saveCheckpoint() {
	p.store.SetCheckpoint(p.checkpoint)
}

// Bootstrap restores checkpoint from store.
func (p *Poset) Bootstrap() {
	if p.checkpoint != nil {
		return
	}
	// restore checkpoint
	p.checkpoint = p.store.GetCheckpoint()
	if p.checkpoint == nil {
		p.Fatal("Apply genesis for store first")
	}
	// restore current super-frame
	p.initSuperFrame()

	// TODO: reload some datas
}

// GetGenesisHash is a genesis getter.
func (p *Poset) GetGenesisHash() hash.Hash {
	return p.Genesis
}

// GenesisHash calcs hash of genesis balances.
func genesisHash(balances map[hash.Peer]inter.Stake) hash.Hash {
	s := NewMemStore()
	defer s.Close()

	if err := s.ApplyGenesis(balances); err != nil {
		logger.Get().Fatal(err)
	}

	return s.GetCheckpoint().Genesis
}
