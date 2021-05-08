package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type HandleFunc func(conn net.Conn)

type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HandleFunc
}

func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HandleFunc)}
}

func (s *Server) Register(path string, handler HandleFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := listener.Close(); cerr != nil {
			if err != nil {
				err = cerr
				return
			}
			log.Println(err)
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		err = s.handleConnection(conn)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) (err error) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Println(err)
		}
	}()
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != io.EOF {
		log.Printf("%s", buf[:n])
	}
	if err != nil {
		return err
	}
	log.Printf("%s", buf[:n])

	data := buf[:n]
	requestLineDelim := []byte{'\n', '\r'}
	requestLineEnd := bytes.Index(data, requestLineDelim)
	if requestLineEnd == -1 {

	}
	requestLine := string(data[:requestLineEnd])
	parts := strings.Split(requestLine, " ")
	for path, handler := range s.handlers {
		if parts[1] == path {
			handler(conn)
		}
	}

	return nil
}
