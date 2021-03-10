package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Player struct
type Player struct {
	no       string
	terrains []string
}

//Game struct
type Game struct {
	totalPlayers  int
	activePlayers int
	leftTokens    []string
}

var p Player
var g Game
var opponents []Player
var tokenMap = [24]string{"1F", "2F", "3F", "4F", "5F", "6F", "7F", "8F", "1B", "2B", "3B", "4B", "5B", "6B", "7B", "8B", "1M", "2M", "3M", "4M", "5M", "6M", "7M", "8M"}

type fn func(string) string

func selectedFunction(f fn, val string) string { // selectedFunction provides functionality to call specific function by its id [:2] of args string
	return f(val)
}

var functions = map[string]fn{
	"01": g.playerNO,
	"02": p.readMyTerrain,
	"03": g.leftoverTokens,
	"04": playerTurn,
	//"05": chooseDice,
	"06": SendInterrogation,
	//"07": guessTokens,
	"08": guessCorrect,
	"09": guessIncorrect,
	"10": tokenInfoSwap,
	"11": remainingWinner,
	"99": errorMsg,
}

func errorMsg(arg string) string {
	fmt.Println(arg)
	return ""
}

func (gm *Game) playerNO(args string) string {
	gm.totalPlayers, _ = strconv.Atoi(args[len(args)-1:])
	gm.activePlayers = gm.totalPlayers
	fmt.Printf("There are %d players, and you are player %s.\n", gm.totalPlayers, p.no)
	initOpponents(gm.totalPlayers)
	return ""
}

func initOpponents(totalPlayers int) {
	for i := 1; i <= totalPlayers; i++ {
		if strconv.Itoa(i) != p.no {
			player := Player{no: strconv.Itoa(i)}
			opponents = append(opponents, player)
		}
	}
}

func (pl *Player) readMyTerrain(args string) string {
	pl.terrains = strings.Split(args[6:], ",")
	fmt.Println("Terrains are: " + strings.Join(pl.terrains, ", "))
	return ""
}

func (gm *Game) leftoverTokens(args string) string {
	gm.leftTokens = strings.Split(args[3:], ",")
	fmt.Printf("Leftover tokens: ")
	for _, x := range gm.leftTokens {
		fmt.Printf(x + " ")
	}
	fmt.Printf("\n")
	return ""
}

func tokenInfoSwap(args string) string {
	message := strings.Split(args[3:], ",")

	if string(message[0][1]) == p.no {
		fmt.Printf("You let %s know you got a token %s\n", message[1], message[2])
	} else {
		fmt.Printf("You acknowledge %s got a token %s\n", message[0], message[2])
	}
	return ""
}

func remainingWinner(args string) string {
	message := strings.Split(args[3:], ",")

	fmt.Printf("%s wins as the only remaining player. All others have guessed incorrectly and been disqualified. The treasures are located at %s and %s\n",
		message[0], message[1], message[2])
	return ""
}

func playerTurn(args string) string {
	message := strings.Split(args[3:], ",")

	fmt.Println("Player " + message[0][1:] + " has rolled " + message[1] + "," + message[2] + "," + message[3])
	if "P"+p.no != message[0] {

		return ""
	}
	var response string
	fmt.Println("Would you like to guess? Y/N")
	fmt.Scanln(&response)
	response = strings.ToUpper(response)
	for response != "Y" && response != "N" {
		fmt.Println("Would you like to guess? Y/N")
		fmt.Scanln(&response)
		response = strings.ToUpper(response)
	}
	if response == "Y" {
		return guessTokens()
	}
	fmt.Println("Choose any two dice options from the following or choose A")
	for j := 1; j < len(message); j++ {
		fmt.Printf("%d. %s\n", j, message[j])
	}
	return chooseDice(args)
}

func terrainParser(t1 string, t2 string) string {
	terrainMap := map[string]string{"B": "Beach", "F": "Forest", "M": "Mountain", "A": "All terrians"}
	var t string
	if t1 == "W" && t2 == "W" {
		fmt.Println("Choose Terrian:")
		for k, v := range terrainMap {
			fmt.Printf("%s: %s\n", k, v)
		}
		fmt.Scanf("%s", &t)
		t = strings.ToUpper(t)
		_, found := terrainMap[t]
		for !found {
			fmt.Println("Invalid terrain, Choose Terrian:")
			fmt.Scanf("%s", &t)
			t = strings.ToUpper(t)
			_, found = terrainMap[t]
		}
	} else if t1 == "W" && t2 != "W" {
		fmt.Println("Choose Terrian:")
		fmt.Printf("%s: %s\n", t2, terrainMap[t2])
		fmt.Printf("A: %s\n", terrainMap["A"])
		fmt.Scanf("%s", &t)
		t = strings.ToUpper(t)
		for t != "A" && t != t2 {
			fmt.Println("Invalid terrain, Choose Terrian:")
			fmt.Scanf("%s", &t)
			t = strings.ToUpper(t)
		}
	} else if t1 != "W" && t2 == "W" {
		fmt.Println("Choose Terrian:")
		fmt.Printf("%s: %s\n", t1, terrainMap[t1])
		fmt.Printf("A: %s\n", terrainMap["A"])
		fmt.Scanf("%s", &t)
		t = strings.ToUpper(t)
		for t != "A" && t != t1 {
			fmt.Println("Invalid terrain, Choose Terrian:")
			fmt.Scanf("%s", &t)
			t = strings.ToUpper(t)
		}
	} else if t1 == t2 {
		return t1
	}
	return "A"
}

func chooseDice(args string) string {
	rolledDice := strings.Split(args[3:], ",")
	var n int
	var die1, die2, terrain, player string

	fmt.Println("Choose first die by number")
	_, err := fmt.Scan(&n)
	for err != nil || n > 3 || n < 1 {
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Out of range, choose first die by number")
		}
		_, err = fmt.Scan(&n)
	}
	die1 = rolledDice[n]

	fmt.Println("Choose second die by number")
	_, err = fmt.Scan(&n)
	for err != nil || n > 3 || n < 1 || rolledDice[n] == die1 {
		if err != nil {
			fmt.Println(err)
		} else if n > 3 || n < 1 {
			fmt.Println("Out of range, choose second die by number")
		} else if rolledDice[n] == die1 {
			fmt.Println("Die has chosen, enter another number for second die")
		}
		_, err = fmt.Scan(&n)
	}
	die2 = rolledDice[n]

	terrain = terrainParser(string(die1[2]), string(die2[2]))

	fmt.Println("Choose Player that you want to interrogate by number")
	for i, opponent := range opponents {
		fmt.Printf("%d. Player%s\n", i+1, opponent.no)
	}
	_, err = fmt.Scan(&n)
	for err != nil || n > len(opponents) || n < 1 {
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Out of range, choose Player that you want to interrogate by number")
		}
		_, err = fmt.Scan(&n)
	}
	player = opponents[n-1].no

	var temp string = "05:" + die1 + "," + die2 + "," + terrain + ",P" + player
	return temp
}

func SendInterrogation(args string) string {
	stringSlice := strings.Split(args, ":")
	stringSlice2 := strings.Split(stringSlice[1], ",")

	fmt.Printf("%s asks %s how many locations they've searched between %s and %s in %s terrain.\n",
		stringSlice2[5], stringSlice2[4], stringSlice2[0], stringSlice2[1], stringSlice2[2])

	fmt.Printf("%s responds %s.\n",
		stringSlice2[4], stringSlice2[3])
	return ""
}

func isValidToken(token string) bool {
	for _, t := range tokenMap {
		if strings.ToUpper(token) == t {
			return true
		}
	}
	return false
}

func guessTokens() string {
	var first_token string
	var second_token string

	fmt.Println("Choose first token: ")
	fmt.Scanf("%s", &first_token)
	for !isValidToken(first_token) {
		fmt.Println("Invalid token, please choose the first token: ")
		fmt.Scanf("%s", &first_token)
	}

	fmt.Println("Choose second token: ")
	fmt.Scanf("%s", &second_token)
	for !isValidToken(second_token) || second_token == first_token {
		if second_token == first_token {
			fmt.Println("Same guess, please choose the second token: ")
		} else {
			fmt.Println("Invalid token, please choose the second token: ")
		}
		fmt.Scanf("%s", &second_token)
	}
	var temp string = "07:P" + p.no + "," + strings.ToUpper(first_token) + "," + strings.ToUpper(second_token)
	return temp
}

func guessCorrect(args string) string {
	stringSlice := strings.Split(args, ":")
	stringSlice2 := strings.Split(stringSlice[1], ",")
	fmt.Printf("Player %s is correct! They have won the game.\n",
		stringSlice2[0])
	fmt.Printf("The treasures were located at %s and %s.\n",
		stringSlice2[1], stringSlice2[2])
	return ""
}

func guessIncorrect(args string) string {
	message := strings.Split(args, ":")
	fmt.Printf("Player %s is submitting a guess at the treasure locations!. Player %s was wrong. They are now disqualified from winning.\n",
		message[1], message[1])
	return ""
}

// Reads from "toPN" named pipe
func readFromPipe(fd *os.File, rd *bufio.Reader) string {
	buff, err := rd.ReadString('\n')
	if err == io.EOF {
		return "exit"
	}
	if err != nil {
		fmt.Println(err)
	}
	if len(buff) > 0 {
		return buff
	}
	return ""
}

// Writes to "fromPN" named pipe
func writeToPipe(fd1 *os.File, args string) {
	fd1.Write([]byte(args))
}

func main() {
	var pipeName, toPN, fromPN string

	fmt.Println("Enter your Player Number and pipe Prefixed Name: (separated by space)")
	fmt.Scanf("%s%s", &p.no, &pipeName)
	toPN = "/tmp/" + pipeName + "toP" + p.no
	fromPN = "/tmp/" + pipeName + "fromP" + p.no
	fmt.Println(toPN, fromPN)

	fd, err := os.OpenFile(toPN, os.O_RDONLY, os.ModeNamedPipe) // opens toPN named pipe
	if err != nil {
		fmt.Println(err)
	}
	rd := bufio.NewReader(fd)

	fd1, err1 := os.OpenFile(fromPN, os.O_WRONLY, 0) // opens fromPN named pipe
	if err1 != nil {
		fmt.Println(err1)
	}

	for {
		serverSaid := readFromPipe(fd, rd)
		if serverSaid == "exit" {
			break
		}
		playerReply := selectedFunction(functions[serverSaid[:2]], strings.TrimSpace(serverSaid))
		if playerReply != "" {
			writeToPipe(fd1, playerReply)
		}
	}
}
