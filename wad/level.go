package wad

import "slices"

type Level struct {
	Name       string
	Things     []Thing
	Linedefs   []Linedef
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
	lumps = append(lumps, makeLinedefsLump(l.Linedefs))
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

func (l Level) FindAllThings(thingTypes ...int16) []*Thing {
	found := make([]*Thing, 0, 10) // Arbitrarily start with 10 capacity since we don't know how many things we'll find
	for i, thing := range l.Things {
		if slices.Contains(thingTypes, thing.Type) {
			found = append(found, &l.Things[i])
		}
	}

	return found
}
