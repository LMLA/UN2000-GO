package main

import (
	"fmt"
	"log"
	"os"
	"net"
	"github.com/joho/godotenv"
	"bufio"
	"time"
	"io"
	"errors"
	"github.com/secsy/goftp"
)

var Lr bool
var responser []byte
var errr error
var ENQr = byte(0x05)
var EOTr = byte(0x04)
var ACKr = byte(0x06)

func UploadFTP(client *goftp.Client, filename string, location string) error {
	bigFile, err := os.Open(filename) //"./test/valid/" + filename
	if err != nil {
		return err
	}

	err = client.Store(location, bigFile) //"iib/test/pmla/"+filename
	if err != nil {
		return err
	}

	return nil
}

func verifyQueryReceive(message string) (L bool, response []byte, err error) {
	L = false
	response = []byte{0x06}
	err = nil
	verify := message[2:3]
	if verify == "L" {
		L = true
	}
		response = []byte{0x06}
		err = errors.New(verify)

	return L, response, err
}


func main() {
	fmt.Println("Cargando Servidor...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error cargando archivo .env")
	}
	// archivo .env con la informacion de conexión
	tcpPort := os.Getenv("TCP_PORT_SERVER")
	ftpAddress := os.Getenv("FTP_ADDRESS")



	// String conexion MySQL

	ln, _ := net.Listen("tcp", ":"+ tcpPort)

Retry:

// Acepta condiciento en puerto indicado
	conn, err := ln.Accept()
	fmt.Println(conn.RemoteAddr().String())
	if err != nil {
		fmt.Println("error tcp", err)
	}

	config := goftp.Config{
		User:               "conlab97",
		Password:           "lab3000",
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
		Logger:             os.Stderr,
	}

	client, ftpconnerr := goftp.DialConfig(config, ftpAddress)
	if ftpconnerr != nil {
		fmt.Println(ftpconnerr)
		c := time.Tick(10 * time.Second) // Reconexion TCP
		for now := range c {
			fmt.Println(now)
			goto Retry
		}
	}

	fmt.Println(err)
	if err != nil {
		conn.Close()
		goto Retry
	}
	// Loop infinito
	for {
	NewMessage:
		fmt.Println("Inicio mensaje")
		t := time.Now()
		timestamp := t.Format("20060102150405")
		filename := timestamp+".txt"
		// ENQ
		message, err := bufio.NewReader(conn).ReadString(ENQr)
		if err != nil {
			fmt.Println("timeout") // Manejo de errores
			if io.EOF == err {
				fmt.Println("connection dropped message", err)
				goto Retry
			}
			goto Retry // Sale del loop si se desconecta el cliente
		} else {
			fmt.Print("ENQ:\n")
			time.Sleep(100 * time.Millisecond)
			_ , err = conn.Write([]byte{0x06})
			fmt.Print("ACK sent: ", err)
			for {
				// H Q L
				message, err = bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					fmt.Println("desconectado") // Manejo de errores
					break // Sale del loop si se desconecta el cliente
				} else {
					// If the file doesn't exist, create it, or append to the file
					f, err := os.OpenFile("results/"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						log.Fatal(err)
					}
					if _, err := f.Write([]byte(message)); err != nil {
						log.Fatal(err)
					}
					if err := f.Close(); err != nil {
						log.Fatal(err)
					}
					Lr, responser, err = verifyQueryReceive(message)
				}
				if err != nil {
					time.Sleep(100 * time.Millisecond)
					conn.Write(responser)
					fmt.Println(err)
				} else {
					time.Sleep(100 * time.Millisecond)
					conn.Write(responser)
				}

				if Lr == true {
					// EOT
					message, err = bufio.NewReader(conn).ReadString(EOTr)
					if err != nil {
						fmt.Println("desconectado") // Manejo de errores
						break // Sal*e del loop si se desconecta el cliente
					} else {
						fmt.Println("Fin mensaje")
						time.Sleep(100 * time.Millisecond)
						break
					}
				}
			}
			ftpUpErr := UploadFTP(client, "results/" + filename, "iib/071/"+filename)
			if ftpUpErr != nil {
				fmt.Println(ftpUpErr)
			} else {
				// se proceso correctamente se puede eliminar
				fmt.Println("procesado")
				os.Remove("results/" + filename)
			}
			goto NewMessage
		}

	}
	conn.Close() // Cierra conexion TCP
	goto Retry // Reinicia la conexion TCP
}