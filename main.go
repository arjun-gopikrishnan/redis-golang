package main

import (
	"log"
	"log/slog"
	"net"
)

const defaultListenAddress = ":6379"

type Config struct {
	ListenAddress string
}

type Server struct {
	Config    Config
	ln        net.Listener
	peers     map[*Peer]bool
	addPeerCh chan *Peer
	quitChan  chan struct{}
}

func NewServer(cfg Config) *Server {

	if len(cfg.ListenAddress) == 0 {
		cfg.ListenAddress = defaultListenAddress

	}

	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitChan:  make(chan struct{}),
	}
}

func (s *Server) Start() error {

	ln, err := net.Listen("tcp", s.Config.ListenAddress)
	if err != nil {
		return err
	}

	s.ln = ln

	// go s.acceptLoop()
	go s.loop()

	slog.Info("server running", "listenAddr", s.Config.ListenAddress)

	return s.acceptLoop()
}

func (s *Server) loop() {
	for {
		select {
		case <-s.quitChan:
			return
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		}
	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept error", "err", err)
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn)
	s.addPeerCh <- peer

	slog.Info("New peer connected", "remoteAddr", conn.RemoteAddr())

	go peer.readLoop()
}

func main() {
	server := NewServer(Config{})

	log.Fatal(server.Start())
}
