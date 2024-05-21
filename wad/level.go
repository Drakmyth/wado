package wad

type Level struct {
	Name       string
	Things     []Thing
	Linedefs   Lump
	Sidedefs   []Sidedef
	Vertexes   Lump
	Segments   Lump
	Subsectors Lump
	Nodes      Lump
	Sectors    Lump
	Reject     Lump
	Blockmap   Lump
}

func (l Level) IsLevelFromGame(game Game) bool {
	return isLevelFromGame(l.Name, game)
}

func (l Level) toLumps() []Lump {
	levelHeader := Lump{
		Name: l.Name,
		Data: []byte{},
	}

	lumps := make([]Lump, 0, 11)
	lumps = append(lumps, levelHeader)
	lumps = append(lumps, makeThingsLump(l.Things))
	lumps = append(lumps, l.Linedefs)
	lumps = append(lumps, makeSidedefsLump(l.Sidedefs))
	lumps = append(lumps, l.Vertexes)
	lumps = append(lumps, l.Segments)
	lumps = append(lumps, l.Subsectors)
	lumps = append(lumps, l.Nodes)
	lumps = append(lumps, l.Sectors)
	lumps = append(lumps, l.Reject)
	lumps = append(lumps, l.Blockmap)

	return lumps
}
