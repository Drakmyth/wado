package wad

import "slices"

type Level struct {
	Name       string
	Things     []Thing
	Linedefs   []Linedef
	Sidedefs   []Sidedef
	Vertexes   []byte
	Segments   []byte
	Subsectors []byte
	Nodes      []byte
	Sectors    []byte
	Reject     []byte
	Blockmap   []byte
}

func (l Level) IsLevelFromGame(game Game) bool {
	return isLevelFromGame(l.Name, game)
}

func (l Level) HasSecretExit() bool {
	for _, linedef := range l.Linedefs {
		if slices.Contains(SECRET_EXIT_LINETYPES, linedef.SpecialType) {
			return true
		}
	}

	return false
}

func (l Level) toLumps() []Lump {
	levelHeader := Lump{
		Name: l.Name,
		Data: []byte{},
	}

	lumps := make([]Lump, 0, 11)
	lumps = append(lumps, levelHeader)
	lumps = append(lumps, Things(l.Things).toLump())
	lumps = append(lumps, Linedefs(l.Linedefs).toLump())
	lumps = append(lumps, Sidedefs(l.Sidedefs).toLump())
	lumps = append(lumps, Lump{Name: LUMP_VERTEXES, Data: l.Vertexes})
	lumps = append(lumps, Lump{Name: LUMP_SEGMENTS, Data: l.Segments})
	lumps = append(lumps, Lump{Name: LUMP_SUBSECTORS, Data: l.Subsectors})
	lumps = append(lumps, Lump{Name: LUMP_NODES, Data: l.Nodes})
	lumps = append(lumps, Lump{Name: LUMP_SECTORS, Data: l.Sectors})
	lumps = append(lumps, Lump{Name: LUMP_REJECT, Data: l.Reject})
	lumps = append(lumps, Lump{Name: LUMP_BLOCKMAP, Data: l.Blockmap})

	return lumps
}

func (l Level) FindAllThings(thingTypes ...int16) []*Thing {
	found := make([]*Thing, 0, 10) // Arbitrarily start with 10 capacity since we don't know how many things we'll find
	for i, thing := range l.Things {
		if slices.Contains(thingTypes, thing.Type) {
			found = append(found, &l.Things[i])
		}
	}

	return found
}
