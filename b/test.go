package main

import (
	"github.com/isaacmaddox/mc-project-bot/db"

	_ "github.com/ncruces/go-sqlite3/driver"

	_ "github.com/ncruces/go-sqlite3/embed"
)

func main() {
	db.Init_database()

	var p db.Project
	p.Create("Example Project", "A testing project created by the seed script")

	p.AddResource("Acacia Button", 0, 64)
	p.AddResource("Birch Log", 64, 128)
	p.AddResource("Fence posts", 0, 24)
	p.AddResource("Oak Log", 0, 192)
	p.AddResource("Spruce Plank", 10, 50)
	p.AddResource("Stone Brick", 50, 200)
	p.AddResource("Iron Ingot", 5, 30)
	p.AddResource("Gold Nugget", 3, 15)
	p.AddResource("Diamond Block", 1, 5)
	p.AddResource("Redstone Dust", 20, 100)
	p.AddResource("Lapis Lazuli", 10, 40)
	p.AddResource("Nether Quartz", 15, 60)
	p.AddResource("Cobblestone", 100, 300)
	p.AddResource("Sandstone", 30, 120)
	p.AddResource("Glass Pane", 25, 80)
	p.AddResource("Glowstone", 5, 20)
	p.AddResource("Ender Pearl", 2, 10)
	p.AddResource("Blaze Rod", 4, 15)
	p.AddResource("Obsidian", 10, 50)
	p.AddResource("Emerald", 5, 25)
	p.AddResource("Clay Block", 20, 90)
	p.AddResource("Wool", 15, 70)
	p.AddResource("Leather", 10, 40)
	p.AddResource("String", 20, 60)
	p.AddResource("Feather", 5, 30)
	p.AddResource("Gunpowder", 10, 50)
	p.AddResource("Bone", 8, 40)
	p.AddResource("Slime Ball", 5, 20)
	p.AddResource("Magma Cream", 3, 15)
	p.AddResource("Nether Wart", 7, 25)
	p.AddResource("Sugar Cane", 20, 100)
	p.AddResource("Cactus", 15, 80)
	p.AddResource("Pumpkin", 10, 50)
	p.AddResource("Melon", 8, 40)
	p.AddResource("Carrot", 12, 60)
	p.AddResource("Potato", 14, 70)
	p.AddResource("Beetroot", 6, 30)
}
