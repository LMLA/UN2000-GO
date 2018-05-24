package main

import (
	"time"
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"log"
	"github.com/secsy/goftp"
)

// Carga datos a un servidor FTP
func UploadFTP(client *goftp.Client, filename string, location string) error {
	bigFile, err := os.Open(filename) //location + filename
	if err != nil {
		return err
	}

	err = client.Store(location, bigFile) //location
	if err != nil {
		return err
	}

	return nil
}

// Almacena los archivos contenidos en test/valid en MySQL y el servidor TCP
func store(ftpconnerr error, client *goftp.Client, filename string, location string) (err error){
	// almacenar en MySQL con valor valido, pero sin procesar hasta que se guarde en el FTP
	err = nil
	if ftpconnerr != nil {
		log.Println("error ftp")
	} else {
		// Enviar a FTP
		ftpUpErr := UploadFTP(client, "results/processed/" + filename, location+filename)
		if ftpUpErr != nil {
			log.Println(ftpUpErr)
			log.Println("no se proceso")
			err = ftpUpErr
		}

	}
	return err
}

// Funcion Principal
// Reenvia los archivos TXT en failed/valid para que sean leidos por 4D, tambien genera archivos TXT si en MySQL la
// orden aparece como no procesada
func main() {
	//Crea archivo con permisos
	f, err := os.OpenFile("store.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//cerrar archivo luego de escribir
	defer f.Close()
	//escribir a f
	log.SetOutput(f)
	// Carga archivo .env con variables de conexion
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error cargando archivo .env")
	}

	// archivo .env con la informacion de conexión
	ftpAddress := os.Getenv("FTP_ADDRESS")
	ftpPath := os.Getenv("FTP_PATH")
	ftpUser := os.Getenv("FTP_USER")
	ftpPassword := os.Getenv("FTP_PASSWORD")

	// Estructura con informacion de conexion a FTP
	config := goftp.Config{
		User:               ftpUser,
		Password:           ftpPassword,
		ConnectionsPerHost: 10,
		Timeout:            5 * time.Second,
		Logger:             os.Stderr,
	}

Retry:


// Crea conexion MySQL
	client, ftpconnerr := goftp.DialConfig(config, ftpAddress)
	if ftpconnerr != nil {
		log.Println(ftpconnerr)
		c := time.Tick(10 * time.Second) // Reconexion TCP
		for now := range c {
			log.Println(now)
			goto Retry
		}
	}


	// Verificacion de no errores TCP
	if ftpconnerr == nil {
		// LOOP INFINITO
		c := time.Tick(5 * time.Second) // busqueda resultados
		for now := range c {
			fmt.Println(now)
			// Carga de directorio valido
			directoryValid := "results/processed/"
			// Se abre el directorio
			outputDirReadValid, _ := os.Open(directoryValid)
			// Se leen los archivos encontrados
			outputDirFilesValid, _ := outputDirReadValid.Readdir(0)
			// Loop sobre los archivos
			for outputIndex := range outputDirFilesValid {
				outputFileHere := outputDirFilesValid[outputIndex]
				// Se obtiene el nombre del archivo
				filename := outputFileHere.Name()
				if err != nil {
					log.Println(err)
				} else {
					errstore := store(ftpconnerr, client, filename, ftpPath) // Funcion almacenamiento
					if errstore != nil{ // error almacenamiento ftp
						log.Println(err)
						client.Close()
						goto Retry
					} else{
						err = os.Remove("results/processed/"+filename) // borrar archivo procesado
						if err != nil {
							log.Println(err)
						}
					}
				}
			}
			if len(outputDirFilesValid) == 0 {
				fmt.Println("no hay archivos en el directorio results") // Directorio vacio
			}
		}
	}else {
		goto Retry
	}
}