package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	hello "./soap"
	"github.com/fiorix/wsdl2go/soap"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var Q bool
var L bool
var OT string
var response []byte

// var check string
var check = "NC219C2"
var caseDate string
var ENQ = byte(0x05)
var EOT = byte(0x04)
var ACK = byte(0x06)
var CR = byte(0x0D)
var LF = byte(0x0A)
var ETXs = byte(0x03)
var data []*hostQueryData
var caseData []*caseQueryData
var genero string
var verHoraOT string
var soapURL string

const (
	ETX = 0x03
	ETB = 23
	STX = 0x02
)

type hostQueryData struct {
	NumOT           string
	CedulaActual    string
	Nombres         string
	Apellido1       string
	Apellido2       string
	Sexo            string
	FechaNacimiento string
	GrupoSanguineo  string
	RH              string
	CodigoExamen    string
	CODUNIVERSAL    string
	URGENTE         string
	FechaOT         string
	HoraOT          string
}

type Orden struct {
	NumOT          string `json:"NumOT"`
	PendienteNumOT string `json:"PendienteNumOT"`
}

type caseQueryData struct {
	hora string
}

func hostQueryDB(db *sql.DB, check string) (p *hostQueryData, err error) {
	rows, err := db.Query("SELECT EAP.NumOT,EAP.CedulaActual, PAC.Nombres, PAC.Apellido1, PAC.Apellido2, PAC.Sexo, PAC.FechaNacimiento, PAC.GrupoSanguineo, PAC.RH, EAP.CodigoExamen, EXA.COD_UNIVERSAL, EAP.URGENTE, EAP.FechaOT, EAP.HoraOT FROM Pacientes PAC, ExamAPracticar EAP, Examenes EXA WHERE EAP.NumOT = ? AND EAP.CedulaActual = PAC.CedulaActual AND EAP.CodigoExamen = EXA.Codigo AND EXA.Instrumento = '071'", check)
	if err != nil {
		log.Println(err) // Manejo de errores
		return nil, err
	}
	for rows.Next() {
		p = new(hostQueryData)
		err = rows.Scan(
			&p.NumOT,
			&p.CedulaActual,
			&p.Nombres,
			&p.Apellido1,
			&p.Apellido2,
			&p.Sexo,
			&p.FechaNacimiento,
			&p.GrupoSanguineo,
			&p.RH,
			&p.CodigoExamen,
			&p.CODUNIVERSAL,
			&p.URGENTE,
			&p.FechaOT,
			&p.HoraOT)

		if err != nil {
			log.Println(err) // Manejo de errores
			return nil, err
		}

	}
	return p, nil
}

func ASTMCheckSum(frame string) string {

	var sumOfChars uint8

	//take each byte in the string and add the values
	for i := 0; i < len(frame); i++ {
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
	RecordType                         string
	SequenceNumber                     string
	StartingRangeIDNumber              string
	EndingRangeIDNumber                string
	UniversalTestID                    string
	RangeofRequestTimeLimits           string
	StartingDateTimeofResultsRequest   string
	EndingDateTimeofResultsRequest     string
	RequestingPhysicianName            string
	RequestingPhysicianTelephoneNumber string
	UserFieldNo1                       string
	UserFieldNo2                       string
	RequestedInformationStatusCodes    string
}

func verifyQuery(message string) (OT string, Q bool, L bool, response []byte, err error) {
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
	} else if verify == "L" {
		L = true
	}

	// Convert string to a rune and grab the num of chars we need
	if len(OT) > 7 {
		OT = string([]rune(OT)[0:7])
	}

	if OT != "" && Q == true {
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

func activeSample(conn net.Conn, p *hostQueryData, nacimiento string, fechaOT string, horaOT string) {
	t := time.Now()
	//******HEADER**********
	data := "1H|\\^&|||LIS||||||||LIS2-A2|" + t.Format("20060102150405") + string(CR) + string(ETXs)
	//fmt.Println(data)
	CheckSum := ASTMCheckSum(data)
	fullData := string(STX) + data + CheckSum + string(CR) + string(LF)
	time.Sleep(300 * time.Millisecond)
	conn.Write([]byte(fullData))

	_, err := bufio.NewReader(conn).ReadString(ACK)
	if err != nil {
		log.Print(err)
	}

	//******PERSON**********
	data = "2P|1|" + p.Nombres + "|" + p.CedulaActual + "||" + p.Nombres + "^" + p.Apellido1 + "||" + nacimiento + "|" + genero + "||||||OPOS|||||||||||||||||||||" + string(CR) + string(ETXs)
	//fmt.Println(data)
	CheckSum = ASTMCheckSum(data)
	fullData = string(STX) + data + CheckSum + string(CR) + string(LF)
	time.Sleep(300 * time.Millisecond)
	conn.Write([]byte(fullData))

	_, err = bufio.NewReader(conn).ReadString(ACK)
	if err != nil {
		log.Print(err)
	}

	//******ORDER**********
	data = "3O|1|" + p.NumOT + "||^^^GLU\\^^^RBC|R||" + fechaOT + horaOT + "||||N||||||||||||||O|||||" + string(CR) + string(ETXs)
	//fmt.Println(data)
	CheckSum = ASTMCheckSum(data)
	fullData = string(STX) + data + CheckSum + string(CR) + string(LF)
	time.Sleep(300 * time.Millisecond)
	conn.Write([]byte(fullData))

	_, err = bufio.NewReader(conn).ReadString(ACK)
	if err != nil {
		log.Print(err)
	}

	//******LINE END**********
	data = "4L|1|N" + string(CR) + string(ETXs)
	//fmt.Println(data)
	CheckSum = ASTMCheckSum(data)
	fullData = string(STX) + data + CheckSum + string(CR) + string(LF)
	time.Sleep(300 * time.Millisecond)
	conn.Write([]byte(fullData))

	_, err = bufio.NewReader(conn).ReadString(ACK)
	if err != nil {
		log.Print(err)
	}

	//******EOT**********
	time.Sleep(300 * time.Millisecond)
	conn.Write([]byte{0x04})
}

func inactiveSample(conn net.Conn, p *hostQueryData, nacimiento string, fechaOT string, horaOT string) {
	t := time.Now()
	//******HEADER**********
	data := "1H|\\^&|||LIS||||||||LIS2-A2|" + t.Format("20060102150405") + string(CR) + string(ETXs)
	//fmt.Println(data)
	CheckSum := ASTMCheckSum(data)
	fullData := string(STX) + data + CheckSum + string(CR) + string(LF)
	time.Sleep(300 * time.Millisecond)
	conn.Write([]byte(fullData))

	_, err := bufio.NewReader(conn).ReadString(ACK)
	if err != nil {
		log.Print(err)
	}

	//******PERSON**********
	data = "2P|1|" + p.Nombres + "|" + p.CedulaActual + "||" + p.Nombres + "^" + p.Apellido1 + "||" + nacimiento + "|" + genero + "||||||OPOS|||||||||||||||||||||" + string(CR) + string(ETXs)
	//fmt.Println(data)
	CheckSum = ASTMCheckSum(data)
	fullData = string(STX) + data + CheckSum + string(CR) + string(LF)
	time.Sleep(300 * time.Millisecond)
	conn.Write([]byte(fullData))

	_, err = bufio.NewReader(conn).ReadString(ACK)
	if err != nil {
		log.Print(err)
	}

	////******ORDER**********
	//data = "3O|1|"+"0"+"||"+""+"|R||"+fechaOT+horaOT+"||||N||||||||||||||O|||||" + string(CR) + string(ETXs)
	////fmt.Println(data)
	//CheckSum = ASTMCheckSum(data)
	//fullData = string(STX) + data + CheckSum + string(CR) + string(LF)
	//conn.Write([]byte(fullData))
	//
	//time.Sleep(1 * time.Second)
	//
	//_, err = bufio.NewReader(conn).ReadString(ACK)
	//if err != nil {
	//	fmt.Print(err)
	//}

	//******LINE END**********
	data = "3L|1|N" + string(CR) + string(ETXs)
	//fmt.Println(data)
	CheckSum = ASTMCheckSum(data)
	fullData = string(STX) + data + CheckSum + string(CR) + string(LF)
	time.Sleep(300 * time.Millisecond)
	conn.Write([]byte(fullData))

	_, err = bufio.NewReader(conn).ReadString(ACK)
	if err != nil {
		log.Print(err)
	}

	//******EOT**********
	time.Sleep(300 * time.Millisecond)
	conn.Write([]byte{0x04})
}

func validMessage(db *sql.DB, conn net.Conn) {
	t := time.Now()
	soapMessage := "En la orden de trabajo " + check + " se esta tratando de Programar un Examen que NO esta ACTIVO :C210 Equipo: 071\n" +
		"No se Programaran los Examenes hasta que no se activen las Muestras.\n" +
		"Se deben revisar todos los Examenes que esten pendientes.\n" +
		"Fecha y Hora: " + t.Format("2006-01-02 15:04:05") + "\n" +
		"Equipo: 071 SYSMEX UN-2000"
	for _, p := range data {
		//qu ery
		if p.Sexo == "0" {
			genero = "F"
		} else {
			genero = "M"
		}
		nacimiento := strings.Replace(p.FechaNacimiento, "-", "", -1)
		fechaOT := strings.Replace(p.FechaOT, "-", "", -1)
		horaOT := strings.Replace(p.HoraOT, ":", "", -1)
		if horaOT == "000000" {
			inactiveSample(conn, p, nacimiento, fechaOT, horaOT)
			soapCrearReto(db, check, soapMessage)
			soapAlerta(check)
		} else {
			activeSample(conn, p, nacimiento, fechaOT, horaOT)
		}
		//crear mensaje
	}
}

func invalidMessage(conn net.Conn) {
	t := time.Now()
	//******HEADER**********
	data := "1H|\\^&|||LIS||||||||LIS2-A2|" + t.Format("20060102150405") + string(CR) + string(ETXs)
	//fmt.Println(data)
	CheckSum := ASTMCheckSum(data)
	fullData := string(STX) + data + CheckSum + string(CR) + string(LF)
	time.Sleep(300 * time.Millisecond)
	conn.Write([]byte(fullData))

	_, err := bufio.NewReader(conn).ReadString(ACK)
	if err != nil {
		log.Print(err)
	}

	////******PERSON**********
	//data = "2P|1||"+p.CedulaActual+"||"+p.Apellido1 +" "+ p.Apellido2+"^"+p.Nombres+"||"+nacimiento+"|"+genero+"||||||OPOS|||||||||||||||||||||" + string(CR) + string(ETXs)
	////fmt.Println(data)
	//CheckSum = ASTMCheckSum(data)
	//fullData = string(STX) + data + CheckSum + string(CR) + string(LF)
	//conn.Write([]byte(fullData))
	//
	//time.Sleep(1 * time.Second)
	//
	//_, err = bufio.NewReader(conn).ReadString(ACK)
	//if err != nil {
	//	fmt.Print(err)
	//}

	//******ORDER**********
	//data = "3O|1|"+"0"+"||"+""+"|R||"+""+"||||N||||||||||||||O|||||" + string(CR) + string(ETXs)
	////fmt.Println(data)
	//CheckSum = ASTMCheckSum(data)
	//fullData = string(STX) + data + CheckSum + string(CR) + string(LF)
	//conn.Write([]byte(fullData))
	//
	//time.Sleep(1 * time.Second)
	//
	//_, err = bufio.NewReader(conn).ReadString(ACK)
	//if err != nil {
	//	fmt.Print(err)
	//}

	//******LINE END**********
	data = "2L|1|N" + string(CR) + string(ETXs)
	//fmt.Println(data)
	CheckSum = ASTMCheckSum(data)
	fullData = string(STX) + data + CheckSum + string(CR) + string(LF)
	time.Sleep(300 * time.Millisecond)
	conn.Write([]byte(fullData))

	_, err = bufio.NewReader(conn).ReadString(ACK)
	if err != nil {
		log.Print(err)
	}

	//******EOT**********
	time.Sleep(300 * time.Millisecond)
	conn.Write([]byte{0x04})
}

func soapAlerta(numot string) {
	cli := soap.Client{
		URL:       soapURL,
		Namespace: hello.Namespace,
	}
	conn := hello.NewServiciosWebRPC(&cli)
	conn.WsAlertaMuestrasInactivas(numot, "071", "C210")
}
func soapCrearReto(db *sql.DB, numot string, soapMessage string) {
	t := time.Now()
	getDate := t.Format("2006-01-02")
	getTime := t.Format("15:04:05")
	dates := string(getDate)
	times := string(getTime)

	rows, err := db.Query("SELECT Hora FROM CalidadEnServicio where Orden = ? and fecha = ? order by Hora desc LIMIT 1", check, dates)
	if err != nil {
		log.Println(err) // Manejo de errores
	}
	for rows.Next() { // Almacena resultado del query en estructura revisado y liberado
		c := new(caseQueryData)
		err := rows.Scan(
			&c.hora)
		caseData = append(caseData, c)
		if err != nil {
			log.Fatal(err) // Manejo de errores
		}
	}
	if len(caseData) == 0 {
		cli := soap.Client{
			URL:       soapURL,
			Namespace: hello.Namespace,
		}
		conn := hello.NewServiciosWebRPC(&cli)
		conn.WsCrearReto(soapMessage, "SISTEMAS", hello.Date(dates), hello.Time(times), "", numot, false)
	} else {
		for _, d := range caseData {
			caseDate = d.hora
		}
		t = t.Add(-30 * time.Minute)
		getTime2 := t.Format("15:04:05")
		times2 := string(getTime2)
		timesCompareFormat, _ := time.Parse("15:04:05", times2)
		caseDateFormat, _ := time.Parse("15:04:05", caseDate)

		fmt.Printf("%v - %v", timesCompareFormat, caseDateFormat)
		if timesCompareFormat.After(caseDateFormat) {
			fmt.Println("ya paso media hora")
			cli := soap.Client{
				URL:       soapURL,
				Namespace: hello.Namespace,
			}
			conn := hello.NewServiciosWebRPC(&cli)
			conn.WsCrearReto(soapMessage, "SISTEMAS", hello.Date(dates), hello.Time(times), "", numot, false)
		} else {
			fmt.Println("no ha pasado media hora")
		}
	}
	caseData = caseData[:0]
}

func main() {
	//create your file with desired read/write permissions
	f, err := os.OpenFile("hostquery.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()
	//set output of logs to f
	log.SetOutput(f)
	//test case
	fmt.Println("Cargando Servidor...")

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error cargando archivo .env")
	}
	// archivo .env con la informacion de conexión
	dbDatabase := os.Getenv("DB_DATABASE")
	dbAddress := os.Getenv("DB_ADDRESS")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	tcpPort := os.Getenv("TCP_PORT")
	soapURL = os.Getenv("SOAP_URL")

	// String conexion MySQL
	dbConn := dbUser + ":" + dbPassword + "@tcp(" + dbAddress + ":" + dbPort + ")/" + dbDatabase

	ln, _ := net.Listen("tcp", ":"+tcpPort)

Retry:

	// Acepta condiciento en puerto indicado
	conn, err := ln.Accept()
	fmt.Println(conn.RemoteAddr().String())
	if err != nil {
		log.Println("error tcp", err)
	}

	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Println("error db")          // Manejo de errores
		c := time.Tick(10 * time.Second) // Reconexion TCP
		for now := range c {
			log.Println(now)
			goto Retry
		}
	}
	// Cerrar conexion a DB si la aplicacion termina para no bloquear puerto
	defer db.Close()

	// Open no abre una conexion. Validar datos DSN:
	if err := db.Ping(); err != nil {
		log.Println("error db", err)     // mensaje error
		c := time.Tick(10 * time.Second) // Reconexion TCP
		for now := range c {
			log.Println(now)
			goto Retry
		}
	}

	log.Println(err)
	if err != nil {
		conn.Close()
		goto Retry
	}
	// Loop infinito
	for {
	NewMessage:
		OT = ""
		check = ""
		fmt.Println("Inicio mensaje")
		// ENQ
		message, err := bufio.NewReader(conn).ReadString(ENQ)
		if err != nil {
			log.Println("timeout") // Manejo de errores
			if io.EOF == err {
				log.Println("connection dropped message", err)
				goto Retry
			}
			goto Retry // Sale del loop si se desconecta el cliente
		} else {
			fmt.Print("ENQ:\n")
			time.Sleep(300 * time.Millisecond)
			_, err = conn.Write([]byte{0x06})
			fmt.Print("ACK sent: ", err)
			for {
				// H Q L
				message, err = bufio.NewReader(conn).ReadString('\r')
				if err != nil {
					log.Println("desconectado") // Manejo de errores
					break                       // Sale del loop si se desconecta el cliente
				} else {
					OT, Q, L, response, err = verifyQuery(message)
				}
				if err != nil {
					time.Sleep(300 * time.Millisecond)
					conn.Write(response)
					log.Println(err)
				} else {
					time.Sleep(300 * time.Millisecond)
					conn.Write(response)
					fmt.Println(OT)
					check = OT
				}

				if L == true {
					// EOT
					message, err = bufio.NewReader(conn).ReadString(EOT)
					if err != nil {
						log.Println("desconectado") // Manejo de errores
						break                       // Sale del loop si se desconecta el cliente
					} else {
						fmt.Println("Fin mensaje")
						time.Sleep(300 * time.Millisecond)
						break
					}
				}
			}
			// enviar ENQ
			fmt.Println(check)
			fmt.Println("Envio orden")
			time.Sleep(300 * time.Millisecond)
			conn.Write([]byte{0x05})
			// respuesta
			_, err = bufio.NewReader(conn).ReadString(ACK)
			if err != nil {
				log.Print(err)
			}
			// OT vacia
			if check == "" {
				// crear examen sin OT
			} else { // OT existe
				p, err := hostQueryDB(db, check)
				if err != nil {
					log.Println(err)
				}
				if p != nil {
					data = append(data, p)
				}
				results, err := db.Query("SELECT NumOT,PendienteNumOT FROM OT WHERE PendienteNumOT=?", check)
				if err != nil {
					log.Println(err, "segundo error") // proper error handling instead of panic in your app
				}

				for results.Next() {
					var ot Orden
					err = results.Scan(&ot.NumOT, &ot.PendienteNumOT)
					if err != nil {
						log.Println(err, "Error scan")
					}

					a, err := hostQueryDB(db, ot.NumOT)
					if err != nil {
						log.Println(err, "cuarto error")
					}

					if a != nil {
						data = append(data, a)
					}
				}

				if len(data) == 0 {
					t := time.Now()
					soapMessage := "En la orden de trabajo " + check + " se esta tratando de Realizar un Examen que NO esta PROGRAMADO: C210 Equipo: 071\n" +
						"No se Programaran los Examenes hasta que no se verifique el Examen A Practicar.\n" +
						"Se deben revisar que la orden de trabajo tenga un examen: C210 programado.\n" +
						"Fecha y Hora: " + t.Format("2006-01-02 15:04:05") + "\n" +
						"Equipo: 071 SYSMEX UN-2000"
					invalidMessage(conn)
					soapCrearReto(db, check, soapMessage)
					soapAlerta(check)
				} else {
					validMessage(db, conn)
				}
			}
			data = data[:0] // limpiar slice datos
			goto NewMessage

		}

	}
	conn.Close() // Cierra conexion TCP
	goto Retry   // Reinicia la conexion TCP
}
