package minecraft

type Root struct {
	DataVersion int32
	Level       Chunk
}

type Chunk struct {
	X                int32        `nbt:"xPos"` // X position of the chunk.
	Z                int32        `nbt:"zPos"` // Z position of the chunk.
	LastUpdate       int64        // Tick when the chunk was last saved.
	LightPopulated   bool         // false = recalculates light.
	TerrainPopulated bool         // false = NMC resets world.
	V                bool         // Likely a version tag, always 1 (ATM)
	InhabitedTime    int64        // How many ticks compounded by players that have occupied this chunk.
	Biomes           [256]byte    // -1 = NMC reset biome.
	HeightMap        [256]int32   // Lowest Y light is at full strength. ZX.
	Sections         [16]*Section // `Sections`         // 16x16x16 blocks.
	Entities         []Entity     // `Entities`         // List of NBT Compound.
	TileEntities     []TileEntity // `TileEntities`     // List of NBT Compound.
	TileTicks        []TileTick   // `TileTicks`        // List of NBT Compound.
}

type Section struct {
	Y          byte       // Y section index. 0~15 bottom to top.
	Blocks     [4096]byte // 8b/block. YZX.
	Add        [2048]byte // 4b/block. YZX. Add << 8 | Blocks
	Data       [2048]byte // 4b/block. YZX.
	BlockLight [2048]byte // 4b/block. YZX.
	SkyLight   [2048]byte // 4b/block. YZX.
}
