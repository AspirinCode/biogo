// Copyright ©2011-2012 Dan Kortschak <dan.kortschak@adelaide.edu.au>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package pals

import (
	"bytes"
	"code.google.com/p/biogo/seq"
	"code.google.com/p/biogo/util"
	"fmt"
)

var (
	binSize    int = 1 << 10
	minPadding int = 50
)

// A Packer collects a set of sequence into a Packed sequence.
type Packer struct {
	Packed  *seq.Seq
	lastPad int
	length  int
}

// Create a new Packer.
func NewPacker(id string) (p *Packer) {
	return &Packer{
		Packed: &seq.Seq{
			ID:     id,
			Strand: 1,
			Meta:   seqMap{},
		},
	}
}

// Pack a sequence into the Packed sequence. Returns a string giving diagnostic information.
func (self *Packer) Pack(sequence *seq.Seq) string {
	m := self.Packed.Meta.(seqMap)

	c := contig{seq: sequence}

	padding := binSize - sequence.Len()%binSize
	if padding < minPadding {
		padding += binSize
	}

	self.length += self.lastPad
	c.from = self.length
	self.length += sequence.Len()
	self.lastPad = padding

	bins := make([]int, (padding+sequence.Len())/binSize)
	for i := 0; i < len(bins); i++ {
		bins[i] = len(m.contigs)
	}

	m.binMap = append(m.binMap, bins...)
	m.contigs = append(m.contigs, c)
	self.Packed.Meta = m

	return fmt.Sprintf("%20s\t%10d\t%7d-%-d", sequence.ID[:util.Min(20, len(sequence.ID))], sequence.Len(), len(m.binMap)-len(bins), len(m.binMap)-1)
}

// Finalise the sequence packing.
func (self *Packer) FinalisePack() {
	lastPad := 0
	self.Packed.Seq = make([]byte, 0, self.length)
	for _, c := range self.Packed.Meta.(seqMap).contigs {
		padding := binSize - c.seq.Len()%binSize
		if padding < minPadding {
			padding += binSize
		}
		self.Packed.Seq = append(self.Packed.Seq, bytes.Repeat([]byte("N"), lastPad)...)
		self.Packed.Seq = append(self.Packed.Seq, c.seq.Seq...)
		lastPad = padding
	}
}

// TODO: The following types should be rationalised to make a true Packed sequence type - include in exp/seq.

// A Contig holds a sequence within a SeqMap.
type contig struct {
	seq  *seq.Seq
	from int
}

// A SeqMap is a collection of sequences mapped to a Packed sequence.
type seqMap struct {
	contigs []contig
	binMap  []int
}
