package bot

import (
	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
)

func (b *Bot) MostThreatenedBase(enemyInBases map[api.UnitTag]scl.Units) (*scl.Unit, scl.Units) {
	var mostEnemyBase *scl.Unit
	var mostEnemyUnits scl.Units

	for base, units := range enemyInBases {
		if len(units) > mostEnemyUnits.Len() {
			mostEnemyBase = b.Units.MyAll.ByTag(base)
			mostEnemyUnits = units
		}
	}

	return mostEnemyBase, mostEnemyUnits
}

func FindClusterAtBase(base *scl.Unit, enemies scl.Units) scl.Units {
	closest := enemies.ClosestTo(base)
	cluster := ClusterBySight(closest, enemies)
	return cluster
}

func ClusterBySight(target *scl.Unit, units scl.Units) scl.Units {
	cluster := make(scl.Units, 0, len(units))
	cluster.Add(target)

	// For each unit in the cluster, add all units in sight to the cluster.
	for _, unit := range cluster {
		inSight := units.Filter(filter.NotIn(cluster), filter.InSightOf(unit))
		cluster = append(cluster, inSight...)
	}

	return cluster
}

func (b *Bot) FindEnemyClusterAtHome() scl.Units {
	enemyInBases := b.FindEnemiesInBases()
	if len(enemyInBases) == 0 {
		return scl.Units{}
	}

	base, enemies := b.MostThreatenedBase(enemyInBases)
	if base == nil || enemies.Empty() {
		return scl.Units{}
	}

	// Where is the enemy cluster at that base?
	cluster := FindClusterAtBase(base, enemies)
	if cluster.Empty() {
		log.Error("No cluster found at base %v", base.Point())
		return scl.Units{}
	}

	return cluster
}
