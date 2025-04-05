package bot

import (
	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
)

type CcForExp map[api.UnitTag]point.Point

func (m *CcForExp) Misplaced(b *Bot) CcForExp {
	misplaced := make(CcForExp)
	for tag, point := range *m {
		unit := b.Units.ByTag[tag]
		if unit == nil {
			delete(*m, tag)
			continue
		}

		if unit.IsFurtherThan(1, point) {
			misplaced[tag] = point
		}
	}

	return misplaced
}

func (m *CcForExp) DeleteTag(tag api.UnitTag) {
	delete(*m, tag)
}

func (m *CcForExp) IsReserved(b *Bot, expansion point.Point) bool {
	for tag, point := range *m {
		unit := b.Units.ByTag[tag]
		if unit == nil {
			delete(*m, tag)
			continue
		}

		if point.IsCloserThan(1, expansion) {
			return true
		}
	}

	return false
}

// Reserved returns the list of expansions that are reserved.
func (m *CcForExp) Reserved(b *Bot) point.Points {
	reserved := make(point.Points, 0, len(*m))
	for tag, point := range *m {
		unit := b.Units.ByTag[tag]
		if unit == nil {
			delete(*m, tag)
			continue
		}

		reserved = append(reserved, point)
	}

	return reserved
}

// ByExpansion returns the town halls that are reserved for a specific
// expansion.
func (m *CcForExp) ByExpansion(b *Bot, expansion point.Point) scl.Units {
	cc := make(scl.Units, 0, len(*m))
	for tag, exp := range *m {
		if exp.IsFurtherThan(1, expansion) {
			continue
		}

		// Dead command centers can be removed from the state
		unit := b.Units.ByTag[tag]
		if unit == nil {
			delete(*m, tag)
			continue
		}

		cc = append(cc, unit)
	}
	return cc
}

// Reserve marks a command center as reserved for a specific expansion. If there
// are any command centers that are already reserved for that expansion, they
// are removed from the state.
func (m *CcForExp) Reserve(b *Bot, unit *scl.Unit, expansion point.Point) {
	cucked := m.ByExpansion(b, expansion)
	if cucked.Exists() {
		for _, cc := range cucked {
			m.DeleteTag(cc.Tag)
		}
	}

	(*m)[unit.Tag] = expansion
}
