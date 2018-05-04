package main

import (
	"fmt"
	"net"
	"bufio"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"errors"
	"time"
)


var Q bool
var L bool
var OT string
var response []byte
var err error
var check string
var ENQ = byte(0x05)
var EOT = byte(0x04)
var ACK = byte(0x06)
var CR = byte(0x0D)
var LF = byte(0x0A)
var ETXs = byte(0x03)


const (
	ETX = 0x03
	ETB = 23
	STX = 0x02
)

func ASTMCheckSum(frame string) string {

	var sumOfChars uint8

	//take each byte in the string and add the values
	for i := 0; i < len(frame) ; i++ {
		byteVal := frame[i]
		sumOfChars += byteVal

		if byteVal == STX {
			sumOfChars = 0
		}

		if byteVal == ETX || byteVal == ETB {
			break
		}
	}

	// return as hex value in upper case
	return fmt.Sprintf("%02X", sumOfChars)
}


// Estructura a revisar
type InquiryRecord struct {
	RecordType	string
	SequenceNumber string
	StartingRangeIDNumber string
	EndingRangeIDNumber string
	UniversalTestID string
	RangeofRequestTimeLimits string
	StartingDateTimeofResultsRequest string
	EndingDateTimeofResultsRequest string
	RequestingPhysicianName string
	RequestingPhysicianTelephoneNumber string
	UserFieldNo1 string
	UserFieldNo2 string
	RequestedInformationStatusCodes string
}

func verifyQuery(message string) (OT string, Q bool, L bool,  response []byte, err error){
	Q = false
	L = false
	response = []byte{}
	err = nil
	QueryResult := InquiryRecord{}
	verify := message[2:3]
	if verify == "Q" {
		Q = true
		data := message[2:]
		data = strings.TrimSuffix(data, "\r")
		parsed := strings.Split(data, "|")
		QueryResult.RecordType = parsed[0]
		QueryResult.SequenceNumber = parsed[1]
		QueryResult.StartingRangeIDNumber = parsed[2]
		QueryResult.EndingRangeIDNumber = parsed[3]
		QueryResult.UniversalTestID = parsed[4]
		QueryResult.RangeofRequestTimeLimits = parsed[5]
		QueryResult.StartingDateTimeofResultsRequest = parsed[6]
		QueryResult.EndingDateTimeofResultsRequest = parsed[7]
		QueryResult.RequestingPhysicianName = parsed[8]
		QueryResult.RequestingPhysicianTelephoneNumber = parsed[9]
		QueryResult.UserFieldNo1 = parsed[10]
		QueryResult.UserFieldNo2 = parsed[11]
		QueryResult.RequestedInformationStatusCodes = parsed[12]
		OT = QueryResult.StartingRangeIDNumber
	} else if verify == "L"{
		L = true
	}

	if OT != "" && Q == true{
		response = []byte{0x06}
	} else if OT == "" && Q == true {
		response = []byte{0x15}
		err = errors.New("no trae OT la orden")
	} else {
		response = []byte{0x06}
		err = errors.New(verify)
	}

	return OT, Q, L, response, err
}

func main(){
	fmt.Println("Cargando Servidor...")
	// Escuchando en las interfaces
	ln, _ := net.Listen("tcp", ":9999")

Restart:

// Acepta condiciento en puerto indicado
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("error tcp")
	}

	fmt.Println(err)
	// Loop infinito
	for {
	NewMessage:
		OT = ""
		check = ""
		fmt.Println("Inicio mensaje")
		// ENQ
		message, err := bufio.NewReader(conn).ReadString(ENQ)
		if err != nil {
			fmt.Println("desconectado") // Manejo de errores
			break                       // Sale del loop si se desconecta el cliente
		} else {
			fmt.Print("ENQ:\n")
			conn.Write([]byte{0x06})
			for {
				// H Q L
				message, err = bufio.NewReader(conn).ReadString('\r')
				if err != nil {
					fmt.Println("desconectado") // Manejo de errores
					break                       // Sale del loop si se desconecta el cliente
				} else {
					OT, Q, L ,response, err = verifyQuery(message)
				}
				if err != nil {
					conn.Write(response)
					fmt.Println(err)
				} else {
					conn.Write(response)
					fmt.Println(OT)
					check = OT
				}

				if L == true {
					// EOT
					message, err = bufio.NewReader(conn).ReadString(EOT)
					if err != nil {
						fmt.Println("desconectado") // Manejo de errores
						break                       // Sale del loop si se desconecta el cliente
					} else {
						fmt.Println("Fin mensaje")
						conn.Write([]byte{0x06})
						break
					}
				}
			}
			// enviar ENQ
			fmt.Println(check)
			fmt.Println("Envio orden")
			conn.Write([]byte{0x05})
			//respuesta
			_, err = bufio.NewReader(conn).ReadString(ACK)
			if err != nil {
				fmt.Print(err)
			}
			//OT vacia
			if check == ""{
				//crear examen sin OT
			} else { // OT existe
				//query

				//******HEADER**********
				data := "1H|\\^&|||LIS||||||||LIS2-A2|20170615152716"+string(CR)+string(ETXs)
				//fmt.Println(data)
				CheckSum := ASTMCheckSum(data)
				fullData := string(STX)+data+CheckSum+string(CR)+string(LF)
				conn.Write([]byte(fullData))

				time.Sleep(1 * time.Second)

				_, err := bufio.NewReader(conn).ReadString(ACK)
				if err != nil {
					fmt.Print(err)
				}


				//******PERSON**********
				data = "2P|1|UAA-01|2129346||UAA-01^AUTOMATED^URINALYSIS||19890919|F||||||OPOS|||||||||||||||||||||"+string(CR)+string(ETXs)
				//fmt.Println(data)
				CheckSum = ASTMCheckSum(data)
				fullData = string(STX)+data+CheckSum+string(CR)+string(LF)
				conn.Write([]byte(fullData))

				time.Sleep(1 * time.Second)

				_, err = bufio.NewReader(conn).ReadString(ACK)
				if err != nil {
					fmt.Print(err)
				}

				//******ORDER**********
				data = "3O|1|2129346||^^^GLU\\^^^RBC|R||20150319151541||||N||||||||||||||O|||||"+string(CR)+string(ETXs)
				//fmt.Println(data)
				CheckSum = ASTMCheckSum(data)
				fullData = string(STX)+data+CheckSum+string(CR)+string(LF)
				conn.Write([]byte(fullData))

				time.Sleep(1 * time.Second)

				_, err = bufio.NewReader(conn).ReadString(ACK)
				if err != nil {
					fmt.Print(err)
				}

				//******LINE END**********
				data = "4L|1|N"+string(CR)+string(ETXs)
				//fmt.Println(data)
				CheckSum = ASTMCheckSum(data)
				fullData = string(STX)+data+CheckSum+string(CR)+string(LF)
				conn.Write([]byte(fullData))

				time.Sleep(1 * time.Second)

				_, err = bufio.NewReader(conn).ReadString(ACK)
				if err != nil {
					fmt.Print(err)
				}


				//******EOT**********
				conn.Write([]byte{0x04})

				//crear mensaje
			}
			goto NewMessage

		}

	}
	conn.Close() // Cierra conexion TCP
	goto Restart // Reinicia la conexion TCP
}
