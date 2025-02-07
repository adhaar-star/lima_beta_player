package player

import (
	"fmt"
	"strconv"
)

// Status flag -1: no mark, 0: unobtained, 1: obtained, 2: treasure
type terrains struct {
	Beach    int
	Forest   int
	Mountain int
}

// Player struct
// Table's indices corresponse to island map's regions, ex: index 0 = region 1, index 1 = region 2
type Player struct {
	No    string
	Table [8]terrains
	// Format ex: {1:{{1B, 2B, 3B}, {1F, 2F}}, 2:{{2B, 3B, 4B}}...}
	PotentialObtainedTknsList map[int][][]string
	//Special abilities 0: no usage, 1: ready to use
	Pistol int
	Shovel int
	Barrel int
}

var TableIndexMap = map[string]int{"NN": 0, "NE": 1, "EE": 2, "SE": 3, "SS": 4, "SW": 5, "WW": 6, "NW": 7}
var Terrains = [4]string{"A", "F", "M", "B"}

func NewPlayer(playerNo string) Player {
	plr := Player{No: playerNo}
	for i := range plr.Table {
		plr.Table[i].Beach = -1
		plr.Table[i].Forest = -1
		plr.Table[i].Mountain = -1
	}
	plr.Pistol = 1
	plr.Shovel = 1
	plr.Barrel = 1

	return plr
}

func (plr *Player) InitPotentialObtainedTknsList(maxTkns int) {
	plr.PotentialObtainedTknsList = make(map[int][][]string)
	for i := 1; i <= maxTkns; i++ {
		plr.PotentialObtainedTknsList[i] = [][]string{}
	}
	plr.RecordPotentialCandidates(maxTkns, plr.TokensInRegionByStatus("NN", "NN", "A", -1))
}

// Parses each token and makes a record
// Status flag -1: no mark, 0: unobtained, 1: obtained, 2: treasure
// token format ex: 1B, 2F, 3M
func (plr *Player) MakeRecord(token string, tokenStatus int) {
	region, _ := strconv.Atoi(token[0:1])
	terrain := token[1:]
	if terrain == "B" {
		if plr.Table[region-1].Beach == -1 {
			plr.Table[region-1].Beach = tokenStatus
		}
	} else if terrain == "F" {
		if plr.Table[region-1].Forest == -1 {
			plr.Table[region-1].Forest = tokenStatus
		}
	} else if terrain == "M" {
		if plr.Table[region-1].Mountain == -1 {
			plr.Table[region-1].Mountain = tokenStatus
		}
	}
}

func (plr *Player) UseAbility(abilityCode string) bool {
	switch abilityCode {
	case "P":
		if plr.Pistol != 0 {
			plr.Pistol--
			return true
		}
	case "S":
		if plr.Shovel != 0 {
			plr.Shovel--
			return true
		}
	case "B":
		if plr.Barrel != 0 {
			plr.Barrel--
			return true
		}
	}

	return false
}

func (plr *Player) RecordPotentialCandidates(nTokens int, candidates []string) {
	// Checks there's no the same set for same number winners
	for _, set := range plr.PotentialObtainedTknsList[nTokens] {
		if len(set) == len(candidates) {
			var equalityStat = make([]bool, len(set))
			for i, va := range set {
				var found bool = false
				for _, vb := range candidates {
					if va == vb {
						found = true
					}
				}
				equalityStat[i] = found
			}
			for _, equality := range equalityStat {
				// Not a same set for a same number winners in record, storing it
				if equality == false {
					break
				}
				// A same set for a same number winners in record, not storing it
				return
			}
		}
	}
	plr.PotentialObtainedTknsList[nTokens] = append(plr.PotentialObtainedTknsList[nTokens], candidates)
}

// Checks token status
func (plr *Player) StatusByToken(token string) int {
	var status int = -1

	region, _ := strconv.Atoi(token[0:1])
	terrain := token[1:]
	if terrain == "B" {
		status = plr.Table[region-1].Beach
	} else if terrain == "F" {
		status = plr.Table[region-1].Forest
	} else if terrain == "M" {
		status = plr.Table[region-1].Mountain
	}

	return status
}

// Reports tokens in a block according to its status
// The checking order of tokens are fixed in order to match the order in token map: terrain order B F M
func (plr *Player) TokensInRegionByStatus(start string, end string, terrain string, checkedStatus int) []string {
	var itStart, itEnd int
	var tokens []string

	itStart = TableIndexMap[start]
	itEnd = TableIndexMap[end]
	if itEnd <= itStart {
		itEnd += 8
	}
	for itStart != itEnd {
		it := itStart % 8
		if terrain == "B" {
			if plr.Table[it].Beach == checkedStatus {
				tokens = append(tokens, strconv.Itoa(it+1)+"B")
			}
		} else if terrain == "F" {
			if plr.Table[it].Forest == checkedStatus {
				tokens = append(tokens, strconv.Itoa(it+1)+"F")
			}
		} else if terrain == "M" {
			if plr.Table[it].Mountain == checkedStatus {
				tokens = append(tokens, strconv.Itoa(it+1)+"M")
			}
		} else if terrain == "A" {
			if plr.Table[it].Beach == checkedStatus {
				tokens = append(tokens, strconv.Itoa(it+1)+"B")
			}
			if plr.Table[it].Forest == checkedStatus {
				tokens = append(tokens, strconv.Itoa(it+1)+"F")
			}
			if plr.Table[it].Mountain == checkedStatus {
				tokens = append(tokens, strconv.Itoa(it+1)+"M")
			}
		}
		itStart++
	}

	return tokens
}

// Checks whether ready to guess, if true, return answer tokens
func (plr *Player) IsGuessingAndGetAnswer() (bool, []string) {
	var treasures, potentialTreasures []string

	treasures = plr.TokensInRegionByStatus("NN", "NN", "A", 2)
	if len(treasures) == 2 {
		return true, treasures
	}
	potentialTreasures = plr.TokensInRegionByStatus("NN", "NN", "A", -1)
	if len(treasures)+len(potentialTreasures) == 2 {
		return true, append(treasures, potentialTreasures...)
	}

	return false, nil
}

// Checks tokens in a block of status -1
func (plr *Player) UnfirmedOneTokensInRegion(start string, end string, terrain string) []string {
	return plr.TokensInRegionByStatus(start, end, terrain, -1)
}

// Checks tokens in a block of status 2
func (plr *Player) UnfirmedTwoTokensInRegion(start string, end string, terrain string) []string {
	return plr.TokensInRegionByStatus(start, end, terrain, 2)
}


// Prints the matrix of current table
func (plr *Player) DisplayTable() {
	fmt.Printf("--------------Player%s--------------\n", plr.No)
	fmt.Printf("%10sNN NE EE SE SS SW WW NW NN\n", "Direction:")
	fmt.Printf("%11s", "Region: ")
	for i := 1; i < 9; i++ {
		fmt.Printf("R%d ", i)
	}
	fmt.Print("\n")
	fmt.Printf("%11s", "Beach: ")
	for _, trs := range plr.Table {
		fmt.Printf("%2d ", trs.Beach)
	}
	fmt.Print("\n")
	fmt.Printf("%11s", "Forest: ")
	for _, trs := range plr.Table {
		fmt.Printf("%2d ", trs.Forest)
	}
	fmt.Print("\n")
	fmt.Printf("%11s", "Mountain: ")
	for _, trs := range plr.Table {
		fmt.Printf("%2d ", trs.Mountain)
	}
	fmt.Print("\n")
}

// Prints the potential tokens report
func (plr *Player) DisplayPotentialTokensReport() {
	fmt.Printf("-------------------Player %s's potential token set report-------------------\n", plr.No)
	for j := 1; j <= len(plr.PotentialObtainedTknsList); j++ {
		fmt.Printf("-Winners: %-64d-\n", j)
		fmt.Printf("-%-73s-\n", "List of Potential candidates set:")
		for _, set := range plr.PotentialObtainedTknsList[j] {
			setString := fmt.Sprintf("%s", set)
			fmt.Printf("-%-73s-\n", setString)
		}
		fmt.Printf("-%73s-\n", "")
	}
	fmt.Println("---------------------------------------------------------------------------")
}

func (plr *Player) DisplayUsageSpecialAbilities() {
	fmt.Printf("--------------------\n")
	fmt.Println("Pistol usage:", plr.Pistol)
	fmt.Println("Shovel usage:", plr.Shovel)
	fmt.Println("Barrel usage:", plr.Barrel)
	fmt.Printf("--------------------\n")
}
