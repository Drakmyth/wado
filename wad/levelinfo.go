package wad

import "fmt"

type LevelInfo struct {
	Name        string
	Label       string
	Next        string
	NextSecret  string
	EndGame     bool
	BossActions []BossAction
}

type BossAction struct {
	Boss        Boss
	SpecialType int16
	Tag         int16
}

func (ba BossAction) String() string {
	return fmt.Sprintf("%s, %d, %d", ba.Boss, ba.SpecialType, ba.Tag)
}

type Boss string

const (
	BOSS_CYBERDEMON  Boss = "Cyberdemon"
	BOSS_SPIDERDEMON Boss = "SpiderMastermind"
	BOSS_BARON       Boss = "BaronOfHell"
	BOSS_MANCUBUS    Boss = "Fatso"
	BOSS_ARACHNOTRON Boss = "Arachnotron"
)

var DEFAULT_LEVELINFOS = map[string]LevelInfo{
	"E1M1": {
		Name:       "Hangar",
		Label:      "E1M1",
		Next:       "E1M2",
		NextSecret: "E1M1",
	},
	"E1M2": {
		Name:       "Nuclear Plant",
		Label:      "E1M2",
		Next:       "E1M3",
		NextSecret: "E1M2",
	},
	"E1M3": {
		Name:       "Toxin Refinery",
		Label:      "E1M3",
		Next:       "E1M4",
		NextSecret: "E1M9",
	},
	"E1M4": {
		Name:       "Command Control",
		Label:      "E1M4",
		Next:       "E1M5",
		NextSecret: "E1M4",
	},
	"E1M5": {
		Name:       "Phobos Lab",
		Label:      "E1M5",
		Next:       "E1M6",
		NextSecret: "E1M5",
	},
	"E1M6": {
		Name:       "Central Processing",
		Label:      "E1M6",
		Next:       "E1M7",
		NextSecret: "E1M6",
	},
	"E1M7": {
		Name:       "Computer Station",
		Label:      "E1M7",
		Next:       "E1M8",
		NextSecret: "E1M7",
	},
	"E1M8": {
		Name:        "Phobos Anomaly",
		Label:       "E1M8",
		Next:        "E1M9",
		NextSecret:  "E1M8",
		EndGame:     true,
		BossActions: []BossAction{{Boss: BOSS_BARON, SpecialType: 23, Tag: 666}}, // S1 Floor Lower to Lowest Floor
	},
	"E1M9": {
		Name:       "Military Base",
		Label:      "E1M9",
		Next:       "E1M4",
		NextSecret: "E1M9",
	},
	"E2M1": {
		Name:       "Deimos Anomaly",
		Label:      "E2M1",
		Next:       "E2M2",
		NextSecret: "E2M1",
	},
	"E2M2": {
		Name:       "Containment Area",
		Label:      "E2M2",
		Next:       "E2M3",
		NextSecret: "E2M2",
	},
	"E2M3": {
		Name:       "Refinery",
		Label:      "E2M3",
		Next:       "E2M4",
		NextSecret: "E2M3",
	},
	"E2M4": {
		Name:       "Deimos Lab",
		Label:      "E2M4",
		Next:       "E2M5",
		NextSecret: "E2M4",
	},
	"E2M5": {
		Name:       "Command Center",
		Label:      "E2M5",
		Next:       "E2M6",
		NextSecret: "E2M9",
	},
	"E2M6": {
		Name:       "Halls of the Damned",
		Label:      "E2M6",
		Next:       "E2M7",
		NextSecret: "E2M6",
	},
	"E2M7": {
		Name:       "Spawning Vats",
		Label:      "E2M7",
		Next:       "E2M8",
		NextSecret: "E2M7",
	},
	"E2M8": {
		Name:        "Tower of Babel",
		Label:       "E2M8",
		Next:        "E2M9",
		NextSecret:  "E2M8",
		EndGame:     true,
		BossActions: []BossAction{{Boss: BOSS_CYBERDEMON, SpecialType: 11, Tag: 0}}, // S1 Exit Level
	},
	"E2M9": {
		Name:       "Fortress of Mystery",
		Label:      "E2M9",
		Next:       "E2M6",
		NextSecret: "E2M9",
	},
	"E3M1": {
		Name:       "Hell Keep",
		Label:      "E3M1",
		Next:       "E3M2",
		NextSecret: "E3M1",
	},
	"E3M2": {
		Name:       "Slough of Despair",
		Label:      "E3M2",
		Next:       "E3M3",
		NextSecret: "E3M2",
	},
	"E3M3": {
		Name:       "Pandemonium",
		Label:      "E3M3",
		Next:       "E3M4",
		NextSecret: "E3M3",
	},
	"E3M4": {
		Name:       "House of Pain",
		Label:      "E3M4",
		Next:       "E3M5",
		NextSecret: "E3M4",
	},
	"E3M5": {
		Name:       "Unholy Cathedral",
		Label:      "E3M5",
		Next:       "E3M6",
		NextSecret: "E3M5",
	},
	"E3M6": {
		Name:       "Mt. Erebus",
		Label:      "E3M6",
		Next:       "E3M9",
		NextSecret: "E3M6",
	},
	"E3M7": {
		Name:       "Limbo",
		Label:      "E3M7",
		Next:       "E3M8",
		NextSecret: "E3M7",
	},
	"E3M8": {
		Name:        "Dis",
		Label:       "E3M8",
		Next:        "E3M9",
		NextSecret:  "E3M8",
		EndGame:     true,
		BossActions: []BossAction{{Boss: BOSS_SPIDERDEMON, SpecialType: 11, Tag: 0}}, // S1 Exit Level
	},
	"E3M9": {
		Name:       "Warrens",
		Label:      "E3M9",
		Next:       "E3M7",
		NextSecret: "E3M9",
	},
	"E4M1": {
		Name:       "Hell Beneath",
		Label:      "E4M1",
		Next:       "E4M2",
		NextSecret: "E4M1",
	},
	"E4M2": {
		Name:       "Perfect Hatred",
		Label:      "E4M2",
		Next:       "E4M9",
		NextSecret: "E4M2",
	},
	"E4M3": {
		Name:       "Sever the Wicked",
		Label:      "E4M3",
		Next:       "E4M4",
		NextSecret: "E4M3",
	},
	"E4M4": {
		Name:       "Unruly Evil",
		Label:      "E4M4",
		Next:       "E4M5",
		NextSecret: "E4M4",
	},
	"E4M5": {
		Name:       "They Will Repent",
		Label:      "E4M5",
		Next:       "E4M6",
		NextSecret: "E4M5",
	},
	"E4M6": {
		Name:        "Against Thee Wickedly",
		Label:       "E4M6",
		Next:        "E4M7",
		NextSecret:  "E4M6",
		BossActions: []BossAction{{Boss: BOSS_CYBERDEMON, SpecialType: 112, Tag: 666}}, // S1 Door Open Stay (fast)
	},
	"E4M7": {
		Name:       "And Hell Followed",
		Label:      "E4M7",
		Next:       "E4M8",
		NextSecret: "E4M7",
	},
	"E4M8": {
		Name:        "Unto the Cruel",
		Label:       "E4M8",
		Next:        "E4M9",
		NextSecret:  "E4M8",
		EndGame:     true,
		BossActions: []BossAction{{Boss: BOSS_SPIDERDEMON, SpecialType: 23, Tag: 666}}, // S1 Floor Lower to Lowest Floor
	},
	"E4M9": {
		Name:       "Fear",
		Label:      "E4M9",
		Next:       "E4M3",
		NextSecret: "E4M9",
	},
	"MAP01": {
		Name:       "Entryway",
		Label:      "Level 1",
		Next:       "MAP02",
		NextSecret: "MAP01",
	},
	"MAP02": {
		Name:       "Underhalls",
		Label:      "Level 2",
		Next:       "MAP03",
		NextSecret: "MAP02",
	},
	"MAP03": {
		Name:       "The Gantlet",
		Label:      "Level 3",
		Next:       "MAP04",
		NextSecret: "MAP03",
	},
	"MAP04": {
		Name:       "The Focus",
		Label:      "Level 4",
		Next:       "MAP05",
		NextSecret: "MAP04",
	},
	"MAP05": {
		Name:       "The Waste Tunnels",
		Label:      "Level 5",
		Next:       "MAP06",
		NextSecret: "MAP05",
	},
	"MAP06": {
		Name:       "The Crusher",
		Label:      "Level 6",
		Next:       "MAP07",
		NextSecret: "MAP06",
	},
	"MAP07": {
		Name:       "Dead Simple",
		Label:      "Level 7",
		Next:       "MAP08",
		NextSecret: "MAP07",
		BossActions: []BossAction{
			{Boss: BOSS_MANCUBUS, SpecialType: 23, Tag: 666},    // S1 Floor Lower to Lowest Floor
			{Boss: BOSS_ARACHNOTRON, SpecialType: 30, Tag: 667}, // W1 Floor Raise by Shortest Lower Texture - NOTE: Bugged in PrBoom-based ports
		},
	},
	"MAP08": {
		Name:       "Tricks and Traps",
		Label:      "Level 8",
		Next:       "MAP09",
		NextSecret: "MAP08",
	},
	"MAP09": {
		Name:       "The Pit",
		Label:      "Level 9",
		Next:       "MAP10",
		NextSecret: "MAP09",
	},
	"MAP10": {
		Name:       "Refueling Base",
		Label:      "Level 10",
		Next:       "MAP11",
		NextSecret: "MAP10",
	},
	"MAP11": {
		Name:       "'O' of Destruction!",
		Label:      "Level 11",
		Next:       "MAP12",
		NextSecret: "MAP11",
	},
	"MAP12": {
		Name:       "The Factory",
		Label:      "Level 12",
		Next:       "MAP13",
		NextSecret: "MAP12",
	},
	"MAP13": {
		Name:       "Downtown",
		Label:      "Level 13",
		Next:       "MAP14",
		NextSecret: "MAP13",
	},
	"MAP14": {
		Name:       "The Inmost Dens",
		Label:      "Level 14",
		Next:       "MAP15",
		NextSecret: "MAP14",
	},
	"MAP15": {
		Name:       "Industrial Zone",
		Label:      "Level 15",
		Next:       "MAP16",
		NextSecret: "MAP31",
	},
	"MAP16": {
		Name:       "Suburbs",
		Label:      "Level 16",
		Next:       "MAP17",
		NextSecret: "MAP16",
	},
	"MAP17": {
		Name:       "Tenements",
		Label:      "Level 17",
		Next:       "MAP18",
		NextSecret: "MAP17",
	},
	"MAP18": {
		Name:       "The Courtyard",
		Label:      "Level 18",
		Next:       "MAP19",
		NextSecret: "MAP18",
	},
	"MAP19": {
		Name:       "The Citadel",
		Label:      "Level 19",
		Next:       "MAP20",
		NextSecret: "MAP19",
	},
	"MAP20": {
		Name:       "Gotcha!",
		Label:      "Level 20",
		Next:       "MAP21",
		NextSecret: "MAP20",
	},
	"MAP21": {
		Name:       "Nirvana",
		Label:      "Level 21",
		Next:       "MAP22",
		NextSecret: "MAP21",
	},
	"MAP22": {
		Name:       "The Catacombs",
		Label:      "Level 22",
		Next:       "MAP23",
		NextSecret: "MAP22",
	},
	"MAP23": {
		Name:       "Barrels o' Fun",
		Label:      "Level 23",
		Next:       "MAP24",
		NextSecret: "MAP23",
	},
	"MAP24": {
		Name:       "The Chasm",
		Label:      "Level 24",
		Next:       "MAP25",
		NextSecret: "MAP24",
	},
	"MAP25": {
		Name:       "Bloodfalls",
		Label:      "Level 25",
		Next:       "MAP26",
		NextSecret: "MAP25",
	},
	"MAP26": {
		Name:       "The Abandoned Mines",
		Label:      "Level 26",
		Next:       "MAP27",
		NextSecret: "MAP26",
	},
	"MAP27": {
		Name:       "Monster Condo",
		Label:      "Level 27",
		Next:       "MAP28",
		NextSecret: "MAP27",
	},
	"MAP28": {
		Name:       "The Spirit World",
		Label:      "Level 28",
		Next:       "MAP29",
		NextSecret: "MAP28",
	},
	"MAP29": {
		Name:       "The Living End",
		Label:      "Level 29",
		Next:       "MAP30",
		NextSecret: "MAP29",
	},
	"MAP30": {
		Name:       "Icon of Sin",
		Label:      "Level 30",
		Next:       "MAP31",
		NextSecret: "MAP30",
		EndGame:    true,
	},
	"MAP31": {
		Name:       "Wolfenstein",
		Label:      "Level 31",
		Next:       "MAP16",
		NextSecret: "MAP32",
	},
	"MAP32": {
		Name:       "Grosse",
		Label:      "Level 32",
		Next:       "MAP16",
		NextSecret: "MAP32",
	},
}
