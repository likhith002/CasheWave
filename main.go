package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
)

const defaultPort = ":5001"

type Config struct {
	Address string
}

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerch chan *Peer
	quitChan  chan struct{}
	msgChan   chan []byte
}

func NewServer(cfg Config) *Server {
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerch: make(chan *Peer),
		quitChan:  make(chan struct{}),
		msgChan:   make(chan []byte),
	}
}

func (s *Server) StartServer() error {
	ln, err := net.Listen("tcp", s.Address)

	if err != nil {
		log.Fatal("Unable to Start server")
		return err
	}
	s.ln = ln
	go s.loop()
	slog.Info("Server Started", "listenAddress", s.Address)
	return s.acceptConns()

}

func (s *Server) handleRawMsg(rawmsg []byte) error {
	// fmt.Println(string(rawmsg))
	cmd, err := parseMessage(string(rawmsg))

	if err != nil {
		return err
	}
	switch v := cmd.(type) {
	case SetCommand:
		slog.Info("trying to set a key in table", "key", v.key, "value", v.value)

	}
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeerch:
			s.peers[peer] = true

		case rawmsg := <-s.msgChan:
			// fmt.Println(rawmsg)
			if err := s.handleRawMsg(rawmsg); err != nil {
				slog.Error("Unable to read message by server", err)
			}

		case <-s.quitChan:
			fmt.Println("Quit chan triggered")
			return

		}

	}
}

func (s *Server) acceptConns() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Unable to accept connection")
			continue
		}
		go s.handleConn(conn)
	}

}

func (s *Server) handleConn(conn net.Conn) {

	peer := NewPeer(conn, s.msgChan)

	s.addPeerch <- peer
	slog.Info("new peer connected", "Addrs", conn.RemoteAddr())
	if err := peer.readLoop(); err != nil {
		slog.Error("Unable to read data", "err", err)
	}
}

func main() {
	server := NewServer(Config{Address: defaultPort})
	log.Fatal(server.StartServer())
}
