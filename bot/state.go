package bot

import (
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
)

type BotState struct {
	// CcForExp marks command centers that are reserved for new expansions.
	CcForExp CcForExp

	// CcForOrbitalCommand marks command centers that are reserved for upgrading
	// to orbital commands.
	CcForOrbitalCommand api.UnitTag

	// CcForPlanetaryFortress marks command centers that are reserved for
	// upgrading to planetary fortresses.
	CcForPlanetaryFortress api.UnitTag

	// BuildingForAddOn marks a barracks as reserved for building a reactor or
	// tech lab.
	BuildingForAddOn api.UnitTag

	// AttackWaves holds the groups of units that are used for attacking.
	AttackWaves AttackWaves
}

func (b *Bot) InitState() {
	b.initCcForExp()
}

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

func (b *Bot) initCcForExp() {
	if b.State.CcForExp == nil {
		b.State.CcForExp = make(map[api.UnitTag]point.Point)
	}

	townHalls := b.FindTownHalls()
	if townHalls.Empty() {
		log.Warn("Couldn't initialize CcForExp because there are no town halls.")
		return
	}

	expansions := append(b.Locs.MyExps, b.Locs.MyStart)
	for _, expansion := range expansions {
		townHall := townHalls.ClosestTo(expansion)
		if townHall.IsCloserThan(1, expansion) {
			b.State.CcForExp[townHall.Tag] = expansion
		}
	}
}
