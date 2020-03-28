package kodenames

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Game struct {
	Code string
	Players map[string][]Player
	Words []Word
	CurrentTeam string
	Winner string
}

type Player struct {
	GameCode string
	Name string
	Team string
	Spy bool
}

type Word struct {
	Word string
	Color string
	Found bool
	HintedBy []string
}

var Games = map[string]*Game{}

func GetOrSetupGame(code string) *Game {
	if game, ok := Games[code]; ok {
		// Return existing game
		return game
	} else {
		// Setup game
		game = &Game{
			Code: code,
			Players: map[string][]Player{
				"red":  {},
				"blue": {},
			},
		}

		game.NewGame()

		Games[code] = game
		return game
	}
}

func (g *Game) ScheduleNewGame() {
	go func() {
		time.Sleep(10 * time.Second)
		g.NewGame()

		// Turn spies
		log.Infof("Rotating spies!")
		for ti, team := range g.Players {
			removedCurrentSpy := false
			assignedNewSpy := false

			for pi, p := range team {
				if p.Spy {
					log.Infof("%s was a spy, not anymore", p.Name)
					p.Spy = false
					removedCurrentSpy = true
				} else if !assignedNewSpy {
					log.Infof("%s is the new %s spy", p.Name, p.Team)
					p.Spy = true
					assignedNewSpy = true
				}

				g.Players[ti][pi] = p

				if removedCurrentSpy && assignedNewSpy {
					break
				}
			}
		}

		go g.Broadcast()
	}()
}

func (g *Game) NewGame() {
	g.Winner = ""

	// Setup 25 random words
	rand.Seed(time.Now().UnixNano())

	words := GetWords(25)
	var colors []string
	if rand.Intn(2) == 0 {
		// Blue starts, 9 blue cards, 8 red cards, 7 lime cards, 1 black card
		g.CurrentTeam = "blue"
		colors = []string{
			"blue", "blue", "blue", "blue", "blue", "blue", "blue", "blue", "blue",
			"red", "red", "red", "red", "red", "red", "red", "red",
			"lime", "lime", "lime", "lime", "lime", "lime", "lime",
			"black"}
	} else {
		// Red starts
		g.CurrentTeam = "red"
		colors = []string{
			"blue", "blue", "blue", "blue", "blue", "blue", "blue", "blue",
			"red", "red", "red", "red", "red", "red", "red", "red", "red",
			"lime", "lime", "lime", "lime", "lime", "lime", "lime",
			"black"}
	}

	// shuffle colors and match them with words
	rand.Shuffle(len(colors), func (i, j int) { colors[i], colors[j] = colors[j], colors[i] })

	g.Words = nil
	for i, word := range words {
		g.Words = append(g.Words, Word{
			Word:     word,
			Color:    colors[i],
			Found:    false,
			HintedBy: []string{},
		})
	}
}

func (g *Game) Broadcast() {
	Server.BroadcastToRoom("/", g.Code, "GameUpdate", g)
}

func (g *Game) Notice(mode string, msg string, args ...interface{}) {
	Server.BroadcastToRoom("/", g.Code, "Notice", Notify{
		Type:    mode,
		Message: fmt.Sprintf(msg, args...),
	})
}

func (g *Game) JoinPlayer(request JoinGameRequest) Player {
	spy := false
	if len(g.Players[request.Team]) == 0 {
		spy = true
	}

	player := Player{
		GameCode: request.GameCode,
		Name: request.Name,
		Team: request.Team,
		Spy:  spy,
	}
	g.Players[request.Team] = append(g.Players[request.Team], player)

	go g.Broadcast()
	go g.Notice("info", "%s a rejoint l'équipe %s !", request.Name, strings.ToUpper(request.Team))

	return player
}

func (g *Game) SelectWord(player *Player, word string) {
	// Lookup the word on the board
	for wi, w := range g.Words {
		if w.Word == word {
			// Add player as hinter on this word. If all the players of the team (minus the spy master) hinted the word,
			// go for it.
			w.HintedBy = append(w.HintedBy, player.Name)
			g.Words[wi] = w
			go g.Notice( "info", "%s suggère le mot \"%s\"", player.Name, word)
			go g.Broadcast()

			if len(w.HintedBy) == len(g.Players[player.Team]) - 1 {
				log.Infof("Every players of team %s suggested %s, validate it", player.Team, word)
				go g.ValidateWord(player.Team, word)
			}

			break
		}
	}
}

func (g *Game) ValidateWord(team string, word string) {
	// Lookup the word on the board
	for wi, w := range g.Words {
		// Reset hinted status in any way
		w.HintedBy = []string{}
		g.Words[wi] = w

		// Then check if it's the validated word
		if w.Word == word {
			w.Found = true

			// If the word is of the proper team
			if w.Color == team {
				go g.Notice("success", "Le mot \"%s\" appartenait à l'équipe %s !", word, team)
			} else if w.Color == "lime" {
				go g.Notice("info", "Le mot \"%s\" n'appartenait à aucune équipe.", word)
				g.NextTurn()
			} else if w.Color == "black" {
				// Black word, this team has lost
				go g.Notice("error", "Le mot \"%s\" était le mot noir !", word)
				g.CurrentTeam = ""

				if team == "red" {
					g.Winner = "blue"
					go g.ScheduleNewGame()
				} else {
					g.Winner = "red"
					go g.ScheduleNewGame()
				}
			} else {
				go g.Notice("warning", "Le mot \"%s\" appartenait à l'équipe adverse !", word)
				g.NextTurn()
			}

			g.Words[wi] = w
			go g.Broadcast()
			g.CheckWinners()
		}
	}
}

func (g *Game) NextTurn() {
	log.Infof("Next Turn! Current team is %s", g.CurrentTeam)

	if g.CurrentTeam == "red" {
		g.CurrentTeam = "blue"
	} else {
		g.CurrentTeam = "red"
	}

	log.Infof("Current team is now %s", g.CurrentTeam)

	g.Broadcast()
	g.Notice("info", "Au tour de l'équipe %s !", strings.ToUpper(g.CurrentTeam))
}

func (g *Game) CheckWinners() {
	if g.Winner != "" {
		// There's already a winner, do nothing
		return
	}

	redWins := true
	blueWins := true

	for _, w := range g.Words {
		if w.Color == "red" && !w.Found {
			redWins = false
		}
		if w.Color == "blue" && !w.Found {
			blueWins = false
		}

		if !redWins && !blueWins {
			break
		}
	}

	if redWins {
		log.Infof("Red wins")
		g.Winner = "red"
		g.CurrentTeam = ""
		go g.Broadcast()
		go g.Notice("success", "L'équipe RED a trouvé tous ses mots et remporte la partie !")
		go g.ScheduleNewGame()
	} else if blueWins {
		log.Infof("Blue wins")
		g.Winner = "blue"
		g.CurrentTeam = ""
		go g.Broadcast()
		go g.Notice("success", "L'équipe BLUE a trouvé tous ses mots et remporte la partie !")
		go g.ScheduleNewGame()
	} else {
		log.Infof("No winner for now")
	}
}

func (g *Game) RemovePlayer(player *Player) {
	for pi, p := range g.Players[player.Team] {
		if p.Name == player.Name {
			g.Players[player.Team] = append(g.Players[player.Team][:pi], g.Players[player.Team][pi+1:]...)
			break
		}
	}

	go g.Broadcast()
}