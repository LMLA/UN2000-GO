package main

import (
	"fmt"
	"log"
	"os"
	"net"
	"github.com/joho/godotenv"
	"bufio"
	"time"
)

func main()  {
	fmt.Println("Cargando Servidor...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error cargando archivo .env")
	}
	// archivo .env con la informacion de conexión
	tcpPortServer := os.Getenv("TCP_PORT_SERVER")

	// String conexion MySQL

	ln, _ := net.Listen("tcp", ":"+ tcpPortServer)

Retry:

// Acepta condiciento en puerto indicado
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("error tcp")
	}
	for {
		message, err := bufio.NewReader(conn).ReadString(ENQ)
		if err != nil {
			fmt.Println("desconectado") // Manejo de errores
			break                       // Sale del loop si se desconecta el cliente
		} else {
			fmt.Print("ENQ:\n")
			time.Sleep(1 * time.Second)
			conn.Write([]byte{0x06})
			for {
				message, err = bufio.NewReader(conn).ReadString('\r')
				if err != nil {
					fmt.Println("desconectado") // Manejo de errores
					break // Sale del loop si se desconecta el cliente
				} else {
					// verificar si es L
				}
				goto Retry
			}
		}
	}
}