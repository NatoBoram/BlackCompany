package filter

import "github.com/aiseeq/s2l/lib/scl"

// NotIn returns units that are not in the list.
func NotIn(units scl.Units) scl.Filter {
	return func(u *scl.Unit) bool {
		return units.ByTag(u.Tag) == nil
	}
}

// InTurretRange returns a filter that checks if a unit is in range of any
// turret
func InTurretRange(turrets scl.Units) scl.Filter {
	turret := turrets.First()
	if turret == nil {
		// No turrets, so all units are not in turret range
		return func(unit *scl.Unit) bool { return false }
	}

	return func(unit *scl.Unit) bool {
		return turrets.CloserThan(turret.AirRange(), unit).Exists()
	}
}

// NotInTurretRange returns a filter that checks if a unit is not in range of
// any turret
func NotInTurretRange(turrets scl.Units) scl.Filter {
	turret := turrets.First()
	if turret == nil {
		// No turrets, so all units are not in turret range
		return func(unit *scl.Unit) bool { return true }
	}

	return func(unit *scl.Unit) bool {
		return turrets.CloserThan(turret.AirRange(), unit).Empty()
	}
}
