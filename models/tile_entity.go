package minecraft

type TileEntity struct {
	Id      string // Tile entity Id.
	X, Y, Z int32  // Pos (not sure why the Entity object uses Pos)
}
