package minecraft

type Entity struct {
	Id                string       // Entity ID. Doesn't exist for players.
	Pos               [3]float64   // current X,Y,Z position of the entity.
	Motion            [3]float64   // current dX,dY,dZ velocity of the entity in meters per tick.
	Rotation          [2]float32   // rotation in degrees, 0 < yaw < 390 degrees, pitch: +-90 degrees
	FallDistance      float32      // Distance the entity has fallen.
	Fire              int16        // Fire ticks left or inmune ticks iff Fire < 0.
	Air               int16        // Air ticks left. Max: 200 (10s). Decreases under water.
	OnGround          bool         // Captain Obvious!
	NoGravity         bool         // If entity will not fall
	Dimension         int32        // -1 Nether, 0 Overworld, 1 The End
	Invulnerable      bool         // Can entity take damage
	PortalCooldown    int32        // Starts at 900 ticks (45s) and decrements.
	UUIDMost          int64        // The most significant bits of this entity's UUID
	UUIDLeast         int64        // The least significant bits of this entity's UUID
	UUID              string       // (removed 1.9)
	CustomName        string       // Appears in player death messages and villager trading interfaces and above entity
	CustomNameVisible bool         //  if true, and this entity has a custom name, it will always appear above them
	Silent            bool         //  if true, this entity will not make sound. May not exist.
	Riding            *Entity      // (deprecated 1.9) Entity being ridden. Recursive.
	Passengers        []*Entity    // Passengers riding. Recursive.
	Glowing           bool         // true if the entity has a glowing outline.
	Tags              []string     // List of custom string data.
	CommandStats      CommandStats //  Information identifying scoreboard parameters to modify relative to the last command run
}

type CommandStats struct {
	SuccessCountObjective     string // Objective's name about the number of successes of the last command (will be an int)
	SuccessCountName          string // Fake player name about the number of successes of the last command
	AffectedBlocksObjective   string // Objective's name about how many blocks were modified in the last command (will be an int)
	AffectedBlocksName        string // Fake player name about how many blocks were modified in the last command
	AffectedEntitiesObjective string // Objective's name about how many entities were altered in the last command (will be an int)
	AffectedEntitiesName      string // Fake player name about how many entities were altered in the last command
	AffectedItemsObjective    string // Objective's name about how many items were altered in the last command (will be an int)
	AffectedItemsName         string // Fake player name about how many items were altered in the last command
	QueryResultObjective      string // Objective's name about the query result of the last command
	QueryResultName           string // Fake player name about the query result of the last command
}
