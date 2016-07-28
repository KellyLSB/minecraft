package minecraft

type TileTick struct {
	Id      int32 // Block Id.
	Ticks   int32 // Ticks until processing. Iff Ticks < 0: overdue.
	X, Y, Z int32 // Pos.
}
