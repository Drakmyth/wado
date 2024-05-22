package cmd

import (
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"regexp"
	"slices"
	"strconv"

	"github.com/Drakmyth/wado/wad"
	"github.com/spf13/cobra"
)

var TEXTURE_REPLACEMENTS = map[string]string{
	"AASTINKY": "DOORSTOP",
	"ASHWALL":  "ASHWALL2",
	"BLODGR1":  "CEMENT9",
	"BLODGR2":  "CEMENT9",
	"BLODGR3":  "CEMENT9",
	"BLODGR4":  "CEMENT9",
	"BRNBIGC":  "MIDGRATE",
	"BRNBIGL":  "MIDGRATE",
	"BRNBIGR":  "MIDGRATE",
	"BRNPOIS2": "BROWN96",
	"BROVINE":  "BROWN1",
	"BROWNWEL": "BROWNHUG",
	"CEMPOIS":  "CEMENT1",
	"COMP2":    "COMPTALL",
	"COMPOHSO": "COMPWERD",
	"COMPTILE": "COMPWERD",
	"COMPUTE1": "COMPSTA1",
	"COMPUTE2": "COMPTALL",
	"COMPUTE3": "COMPTALL",
	"DOORHI":   "TEKBRON2",
	"GRAYDANG": "GRAY5",
	"ICKDOOR1": "DOOR1",
	"ICKWALL6": "ICKWALL5",
	"LITE2":    "BROWN1",
	"LITE4":    "LITE5",
	"LITE96":   "BROWN96",
	"LITEBLU2": "LITEBLU1",
	"LITEBLU3": "LITEBLU1",
	"LITEMET":  "METAL1",
	"LITERED":  "DOORRED",
	"LITESTON": "STONE2",
	"MIDVINE1": "MIDGRATE",
	"MIDVINE2": "MIDGRATE",
	"NUKESLAD": "SLADWALL",
	"PLANET1":  "COMPSTA2",
	"REDWALL1": "REDWALL",
	"SKINBORD": "SKINMET1",
	"SKINTEK1": "SKINMET2",
	"SKINTEK2": "SKSPINE1",
	"SKULWAL3": "SKSPINE1",
	"SKULWALL": "SKSPINE1",
	"SLADRIP1": "SLADWALL",
	"SLADRIP2": "SLADWALL",
	"SLADRIP3": "SLADWALL",
	"SP_DUDE3": "SP_DUDE4",
	"SP_DUDE6": "SP_DUDE4",
	"SP_ROCK2": "SP_ROCK1",
	"STARTAN1": "STARTAN2",
	"STONGARG": "STONE3",
	"STONPOIS": "STONE",
	"TEKWALL2": "TEKWALL4",
	"TEKWALL3": "TEKWALL4",
	"TEKWALL5": "TEKWALL4",
	"WOODSKUL": "WOODGARG",
}

// These changed from 64x128 to 128x128 in Doom 2
var SHIFT_TEXTURES = []string{"BRNPOIS", "NUKEPOIS", "SW1BRN1", "SW1STON2", "SW1STONE", "SW2BRN1", "SW2STON2", "SW2STONE"}

var LUMP_REPLACEMENTS = map[string]string{
	"D_INTER":  "D_DM2INT",
	"D_INTRO":  "D_DM2TTL",
	"D_VICTOR": "D_READ_M",
	"SKY1":     "RSKY1",
	"SKY2":     "RSKY2",
	"SKY3":     "RSKY3",
	"DEMO1":    "DEMO1_D",
	"DEMO2":    "DEMO2_D",
	"DEMO3":    "DEMO3_D",
}

var D2_REPLACEMENT_CANDIDATES = []int16{wad.ENEMY_SHOTGUN, wad.ENEMY_IMP, wad.ENEMY_PINKY, wad.ENEMY_BARON, wad.ENEMY_PISTOL, wad.ENEMY_CACO, wad.ENEMY_SOUL}

var seed uint64
var flagUpdateThings bool
var flagUpdateSidedefs bool

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.PersistentFlags().Uint64VarP(&seed, "seed", "s", rand.Uint64(), "rando me baby")
	convertCmd.PersistentFlags().BoolVarP(&flagUpdateThings, "things", "T", false, "Stick some D2 in my D")
	convertCmd.PersistentFlags().BoolVarP(&flagUpdateSidedefs, "textures", "t", false, "Fill in the gaps")
}

var convertCmd = &cobra.Command{
	Use:   "convert <path>",
	Short: "Convert the WAD",
	Long:  `Converts Doom WADs to Doom 2 WADs`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires at least two args")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := convert(args[0], args[1])
		if err != nil {
			panic(err)
		}
	},
}

func convert(in_filepath string, out_filepath string) error {
	// Copy to output file so we don't have to worry about messing up the format or the source file
	err := copyFile(in_filepath, out_filepath)
	if err != nil {
		return err
	}

	// Open file
	wf, err := wad.OpenFile(out_filepath)
	if err != nil {
		return err
	}
	defer wf.Close()

	fmt.Printf("Seed: %d", seed)
	rng := rand.New(rand.NewPCG(seed, seed))

	// For each lump...
	for i, lump := range wf.Lumps {
		// If lump needs to be renamed...
		newName, rename := LUMP_REPLACEMENTS[lump.Name]
		if rename {
			wf.Lumps[i].Name = newName
		}
	}

	// For each level...
	levelNameRegexp := regexp.MustCompile(`^E(\d)M(\d)$`)
	for i, level := range wf.Levels {
		// Skip non-Doom1 levels
		if !level.IsLevelFromGame(wad.GAME_DOOM) {
			continue
		}

		// Get episode number
		parts := levelNameRegexp.FindStringSubmatch(level.Name)
		episodeNumber, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}

		// Get mission number
		missionNumber, err := strconv.Atoi(parts[2])
		if err != nil {
			return err
		}

		// Convert map name from ExMy to MAPxx
		mapNumber := ((episodeNumber - 1) * 9) + missionNumber
		newName := fmt.Sprintf("MAP%02d", mapNumber)
		wf.Levels[i].Name = newName

		// Replace things
		if flagUpdateThings {
			updateThings(&wf.Levels[i], rng)
		}

		// Fix textures
		if flagUpdateSidedefs {
			updateSidedefs(&wf.Levels[i])
		}
	}

	return wf.Save()
}

func copyFile(srcPath string, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	return err
}

func updateThings(level *wad.Level, rng *rand.Rand) {
	// Replace all shotguns with SSGs
	shotguns := level.FindAllThings(wad.THING_SHOTGUN)
	for _, shotgun := range shotguns {
		shotgun.Type = wad.THING_SSG
	}

	// Generate 1 Megasphere, 1 Archvile, 1 Berserk, and 1 SSG
	d2Things := level.FindAllThings(D2_REPLACEMENT_CANDIDATES...)
	wad.ReplaceThingsCount(d2Things, map[int16]int16{
		wad.THING_MEGASPHERE: 1,
		wad.ENEMY_ARCHVILE:   1,
		wad.THING_BERSERK:    1,
		wad.THING_SSG:        1,
	}, rng)

	// Replace 20% of Imps with Chaingunners
	imps := level.FindAllThings(wad.ENEMY_IMP)
	wad.ReplaceThingsWeighted(imps, map[int16]float64{
		wad.ENEMY_CHAINGUNNER: 0.2,
	}, rng)

	// Replace 10% of Cacodemons with Pain Elementals
	cacos := level.FindAllThings(wad.ENEMY_CACO)
	wad.ReplaceThingsWeighted(cacos, map[int16]float64{
		wad.ENEMY_PAIN: 0.1,
	}, rng)

	// Replace 10% of Barons with Arachnotrons, 10% with Revenants, and 30% with Hell Knights
	barons := level.FindAllThings(wad.ENEMY_BARON)
	wad.ReplaceThingsWeighted(barons, map[int16]float64{
		wad.ENEMY_ARACH:    0.1,
		wad.ENEMY_REVENANT: 0.1,
		wad.ENEMY_KNIGHT:   0.3,
	}, rng)

	// Replace 10% of Pistol Zombies with Chaingunners, 5% with Medikits, 10% with Stimpacks, and 20% with Health Pots
	pistols := level.FindAllThings(wad.ENEMY_PISTOL)
	wad.ReplaceThingsWeighted(pistols, map[int16]float64{
		wad.ENEMY_CHAINGUNNER: 0.1,
		wad.THING_MEDKIT:      0.05,
		wad.THING_STIM:        0.1,
		wad.THING_HEALTH:      0.2,
	}, rng)
}

func updateSidedefs(level *wad.Level) {
	// Update texture names in sidedefs
	for i, sidedef := range level.Sidedefs {
		if shouldShiftTex(sidedef) {
			sidedef.XOffset += 32
		}

		sidedef.UpperTex = getNewTexName(sidedef.UpperTex)
		sidedef.MiddleTex = getNewTexName(sidedef.MiddleTex)
		sidedef.LowerTex = getNewTexName(sidedef.LowerTex)
		level.Sidedefs[i] = sidedef
	}
}

func shouldShiftTex(sidedef wad.Sidedef) bool {
	return slices.ContainsFunc(SHIFT_TEXTURES, func(tex string) bool {
		return tex == sidedef.UpperTex || tex == sidedef.LowerTex || tex == sidedef.MiddleTex
	})
}

func getNewTexName(oldName string) string {
	newName, replaced := TEXTURE_REPLACEMENTS[oldName]
	if !replaced {
		newName = oldName
	}

	return newName
}
