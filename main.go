package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Server struct{}

func (s *Server) start() {
	ln, err := net.Listen("tcp", ":4000")

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go s.reader(conn)
	}

}

func (s *Server) reader(conn net.Conn) {
	buff := new(bytes.Buffer)
	for {
		var size int64
		binary.Read(conn, binary.LittleEndian, &size)
		n, err := io.CopyN(buff, conn, size)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(buff.Bytes())
		fmt.Printf("recieved %d  bytes over the network\n", n)
	}
}

func sendFile(size int) error {
	file := make([]byte, size)

	_, err := io.ReadFull(rand.Reader, file)

	if err != nil {
		log.Fatal(err)
		return err
	}

	conn, err := net.Dial("tcp", ":4000")

	if err != nil {
		log.Fatal(err)
		return err
	}

	binary.Write(conn, binary.LittleEndian, int64(size))
	n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("written %d bytes over the network\n", n)

	return nil
}
func main() {

	go func() {
		time.Sleep(3 * time.Second)
		sendFile(200000)
	}()
	server := &Server{}

	server.start()
}
