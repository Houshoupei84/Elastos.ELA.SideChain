package types

import (
	"bytes"
	"io"

	"github.com/elastos/Elastos.ELA.SideChain/auxpow"

	"github.com/elastos/Elastos.ELA/common"
)

type Header struct {
	Version    uint32
	Previous   common.Uint256
	MerkleRoot common.Uint256
	Timestamp  uint32
	Bits       uint32
	Nonce      uint32
	Height     uint32
	SideAuxPow auxpow.SideAuxPow
}

func (header *Header) Serialize(w io.Writer) error {
	err := header.serializeNoAux(w)
	if err != nil {
		return err
	}

	err = header.SideAuxPow.Serialize(w)
	if err != nil {
		return err
	}

	w.Write([]byte{byte(1)})
	return nil
}

func (header *Header) Deserialize(r io.Reader) error {
	err := common.ReadElements(r,
		&header.Version,
		&header.Previous,
		&header.MerkleRoot,
		&header.Timestamp,
		&header.Bits,
		&header.Nonce,
		&header.Height,
	)
	if err != nil {
		return err
	}

	// SideAuxPow
	err = header.SideAuxPow.Deserialize(r)
	if err != nil {
		return err
	}

	r.Read(make([]byte, 1))

	return nil
}

func (header *Header) serializeNoAux(w io.Writer) error {
	return common.WriteElements(w,
		header.Version,
		&header.Previous,
		&header.MerkleRoot,
		header.Timestamp,
		header.Bits,
		header.Nonce,
		header.Height,
	)
}

func (header *Header) Hash() common.Uint256 {
	buf := new(bytes.Buffer)
	header.serializeNoAux(buf)
	return common.Sha256D(buf.Bytes())
}
