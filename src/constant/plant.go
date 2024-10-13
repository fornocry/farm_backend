package constant

import "time"

type Plant string

const (
	MONEY          Plant = "MONEY"
	STRAWBERRY     Plant = "STRAWBERRY"
	ROSE           Plant = "ROSE"
	SUNFLOWER      Plant = "SUNFLOWER"
	CHRISTMAS_TREE Plant = "CHRISTMAS_TREE"
)

var Plants = []Plant{
	MONEY,
	STRAWBERRY,
	ROSE,
	SUNFLOWER,
	CHRISTMAS_TREE,
}

type PlantInfo struct {
	Name     Plant
	GrowTime time.Duration // Grow time in days
	Reward   int           // Reward for harvesting
}

var plantData = []PlantInfo{
	{Name: MONEY, GrowTime: time.Hour * 5, Reward: 1},
	{Name: STRAWBERRY, GrowTime: time.Hour * 5, Reward: 1},
	{Name: ROSE, GrowTime: time.Hour * 5, Reward: 1},
	{Name: SUNFLOWER, GrowTime: time.Hour * 5, Reward: 2},
	{Name: CHRISTMAS_TREE, GrowTime: time.Hour * 5, Reward: 1},
}

func IsValidPlant(plant Plant) bool {
	for _, validPlant := range Plants {
		if plant == validPlant {
			return true
		}
	}
	return false
}
