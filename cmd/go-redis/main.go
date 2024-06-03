package main

import (
	"log"
	"net"

	"github.com/arjun/redis-go/internal/keystore"
	. "github.com/arjun/redis-go/internal/peer"
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
	Cache *keystore.Store
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddress) == 0 {
		cfg.ListenAddress = defaultListenAddress
	}
	cache,_ := keystore.NewStore()
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitChan:  make(chan struct{}),
		Cache: cache,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.Config.ListenAddress)
	if err != nil {
		return err
	}

	s.ln = ln

	go s.loop()

	log.Println("server running", "listenAddr", s.Config.ListenAddress)

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
			log.Println("accept error", "err", err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn)
	s.addPeerCh <- peer

	log.Println("New peer connected", "remoteAddr", conn.RemoteAddr())

	go peer.ReadLoop(s.Cache)
}

func main() {
	server := NewServer(Config{})
	log.Fatal(server.Start())
}
