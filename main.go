package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"crypto/tls"
	"encoding/hex"
	"os"
)

// Update if needed
var (
	name      = "Prahaladh HN"
	email     = "prahaladhhn96@gmail.com"
	country   = "India"
	skype     = "N/A"
	birthdate = "15.07.1996"
	addr1     = "Agananooru Street"
	addr2     = "Chennai 600062"
)

func write(conn net.Conn, msg string) {
	conn.Write([]byte(msg + "\n"))
	fmt.Println("Client:", msg)
}

func send(conn *tls.Conn, authdata, challenge, value string) {
	hash := sha1.Sum([]byte(authdata + challenge))
	hashStr := hex.EncodeToString(hash[:])

	msg := fmt.Sprintf("%s %s\n", hashStr, value)
	conn.Write([]byte(msg))
}

func solvePOW(authdata string, difficulty int) string {
	prefix := strings.Repeat("0", difficulty)
	base := []byte(authdata)

	resultChan := make(chan string)
	workers := runtime.NumCPU()

	for w := 0; w < workers; w++ {
		go func(start int) {

			buf := make([]byte, 0, len(base)+20)

			for i := start; ; i += workers {
				suffix := strconv.Itoa(i)

				buf = buf[:0]
				buf = append(buf, base...)
				buf = append(buf, suffix...)

				hash := sha1.Sum(buf)
				hashStr := hex.EncodeToString(hash[:])

				if strings.HasPrefix(hashStr, prefix) {
					resultChan <- suffix
					return
				}
			}
		}(w)
	}

	return <-resultChan
}

func main() {

	// Load cert + key from env
	certFile := os.Getenv("TLS_CERT")
	keyFile := os.Getenv("TLS_KEY")

	if certFile == "" || keyFile == "" {
		fmt.Println("Set TLS_CERT and TLS_KEY environment variables")
		return
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}

	config := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	// Try multiple ports
	ports := []string{"3336", "8083", "8446", "49155", "3481", "65532"}

	var conn *tls.Conn

	for _, port := range ports {
		address := "18.202.148.130:" + port
		fmt.Println("Trying port:", port)

		conn, err = tls.Dial("tcp", address, config)
		if err == nil {
			fmt.Println("Connected on port:", port)
			break
		}

		fmt.Println("Failed:", err)
	}

	if conn == nil {
		panic("Could not connect to any port")
	}

	defer conn.Close()

	reader := bufio.NewReader(conn)

	var authdata string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		line = strings.TrimSpace(line)
		fmt.Println("Server:", line)

		args := strings.Split(line, " ")

		switch args[0] {

		case "HELO":
			conn.Write([]byte("TOAKUEI\n"))

		case "POW":
			authdata = args[1]
			difficulty, _ := strconv.Atoi(args[2])

			fmt.Println("Solving POW...")
			suffix := solvePOW(authdata, difficulty)

			conn.Write([]byte(suffix + "\n"))

		case "NAME":
			send(conn, authdata, args[1], name)

		case "MAILNUM":
			send(conn, authdata, args[1], "1")

		case "MAIL1":
			send(conn, authdata, args[1], email)

		case "SKYPE":
			send(conn, authdata, args[1], skype)

		case "BIRTHDATE":
			send(conn, authdata, args[1], birthdate)

		case "COUNTRY":
			send(conn, authdata, args[1], country)

		case "ADDRNUM":
			send(conn, authdata, args[1], "2")

		case "ADDRLINE1":
			send(conn, authdata, args[1], addr1)

		case "ADDRLINE2":
			send(conn, authdata, args[1], addr2)

		case "END":
			conn.Write([]byte("OK\n"))
			fmt.Println("Completed successfully")
			return

		case "ERROR":
			fmt.Println("Server error:", strings.Join(args[1:], " "))
			return
		}
	}
}