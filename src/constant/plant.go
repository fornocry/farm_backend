package constant

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
