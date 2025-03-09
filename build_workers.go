package main

import (
	"log"

	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/enums/ability"
)

const (
	// MaxWorkers is the maximum number of workers that can be trained.
	MaxWorkers = 80
)

// BuildWorker trains SCVs from command centers.
//
//   - When SCVs can be afforded and there's less than 80 of them
//   - Find mineral fields that aren't depleted and count the missing SCVs to saturate them
//   - Find vespene geysers that aren't exhausted and count the missing SCVs to saturate them
//
// For each town halls:
//
//   - Get the closest resource that's not saturated
//   - Set the rally point to that resource
//   - Train a SCV
func (b *Bot) BuildWorker() {
	if !b.CanBuy(ability.Train_SCV) || b.findMiners().Len() >= MaxWorkers {
		return
	}

	townHalls := b.findTownHalls()
	if townHalls.Empty() {
		return
	}

	resources := b.findUnsaturatedResourcesNearTownHalls(townHalls)
	if resources.Empty() {
		return
	}

	idleTownHalls := townHalls.Filter(scl.Ready, scl.Idle, func(u *scl.Unit) bool {
		_, ok := b.state.ccForExp[u.Tag]
		return !ok
	})
	if idleTownHalls.Empty() {
		return
	}

	for _, cc := range idleTownHalls {
		if !b.CanBuy(ability.Train_SCV) {
			break
		}

		resource := resources.ClosestTo(cc)
		if resource == nil {
			break
		}

		log.Printf("Training SCV for resource %v", resource.Point())
		cc.CommandTag(ability.Rally_CommandCenter, resource.Tag)
		cc.CommandQueue(ability.Train_SCV)
		b.DeductResources(ability.Train_SCV)
	}
}
