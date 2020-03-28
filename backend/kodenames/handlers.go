package kodenames

import (
	socketio "github.com/googollee/go-socket.io"
)

type Notify struct {
	Type string
	Message string
}

type JoinGameRequest struct {
	Name     string
	Team     string
	GameCode string
}

type JoinGameResponse struct {
	Player Player
	Game   *Game
}

type WordSelectedRequest struct {
	Word string
}

func HandleJoinGame(s socketio.Conn, req JoinGameRequest) {
	log.Infof("Player %s joining game %s on team %s", req.Name, req.GameCode, req.Team)
	game := GetOrSetupGame(req.GameCode)
	player := game.JoinPlayer(req)

	s.Join(req.GameCode)
	s.SetContext(&player)

	s.Emit("JoinAccepted", JoinGameResponse{
		Player: player,
		Game:   game,
	})
}

func HandleWordSelected(s socketio.Conn, word WordSelectedRequest) {
	player := s.Context().(*Player)

	game := GetOrSetupGame(player.GameCode)
	log.Infof("Player %s selected word %s", player.Name, word.Word)
	game.SelectWord(player, word.Word)
}

func HandleNextTurn(s socketio.Conn) {
	player := s.Context().(*Player)

	game := GetOrSetupGame(player.GameCode)
	game.NextTurn()
}
