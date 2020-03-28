package kodenames

import (
	engineio "github.com/googollee/go-engine.io"
	"github.com/googollee/go-socket.io"
	"net/http"
	"time"
)

type customServer struct {
	Server *socketio.Server
}

var Server *socketio.Server

func (s *customServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	origin := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	s.Server.ServeHTTP(w, r)
}

func SetupSocketIO() {
	var err error
	Server, err = socketio.NewServer(&engineio.Options{
		PingTimeout:        4*time.Second,
		PingInterval:       1*time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	// General handlers
	Server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Infof("Session connected: %s", s.ID())
		return nil
	})
	Server.OnError("/", func(s socketio.Conn, e error) {
		log.Error(e)
	})
	Server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Infof("Player disconnected: %s", reason)
		player, ok := s.Context().(*Player)
		if ok {
			game := GetOrSetupGame(player.GameCode)
			game.RemovePlayer(player)
		}
	})

	// Game handlers
	Server.OnEvent("/", "JoinGame", HandleJoinGame)
	Server.OnEvent("/", "WordSelected", HandleWordSelected)
	Server.OnEvent("/", "NextTurn", HandleNextTurn)

	// Start up!
	go Server.Serve()
	defer Server.Close()

	custServ := &customServer{Server: Server}

	http.Handle("/socket.io/", custServ)
	http.Handle("/", http.FileServer(http.Dir("../public")))
	log.Infof("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}