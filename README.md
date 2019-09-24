# UN-2000
El sistema SYSMEX UN-2000 comprende 3 equipos UC-3500, UF-5000 y UD-10. 
Se realizó una interfaz utilizando el lenguaje GO

## Contexto

El equipo UC-3500 es un analizador de química de la orina totalmente automatizado, la muestra de orina es pipeteada sobre cada almohadilla de una tira de prueba dedicada dentro del analizador. La secuencia completa, comenzando desde la aspiración de la muestra, hasta la comparación del color y el envío final de resultados es completamente automática.

El equipo UF-5000 es un citómetro completamente automatizado, realiza citometría de flujo fluorescente con láser semiconductor azul y enfoque hidrodinámico en conductimetría de dos canales diferentes.

El equipo UD-10 realiza captura de imágenes de las partículas en las muestras, las imágenes de partículas del UD-10 le ayudarán a diferenciar los resultados anómalos o confusos. 

## Configuración

La informacion de configuracion se agrega en un archivo .env de la siguiente manera:

```
DB_ADDRESS  : 131.1.18.106
DB_DATABASE : 4dlab
DB_PORT     : 3306
DB_USER     : 4duser
DB_PASSWORD : LavAmerikx09
TCP_ADDRESS : localhost
TCP_PORT    : 10002
TCP_PORT_SERVER : 10001
FTP_ADDRESS : 131.1.18.111
SOAP_URL    : http://131.1.18.106:8081/4DSOAP
FTP_PATH    : /home/conlab97/Ftp/iib/071/
FTP_USER    : conlab97
FTP_PASSWORD : lab3000                 
```

## Compilación

Para LINUX : env GOOS=linux GOARCH=amd64 go build main.go
Para Mac OS: env GOOS=darwin GOARCH=amd64 go build main.go

## Instalación

1. Configurar archivo .env.
2. Correr binario.
```
nohup store.go &
```

## Interfaz

La consulta y transformación del mensaje se realizó utilizando **GO** y consultando una base de datos **MySQL**.
El comportamiento de la interfaz comprende 3 partes:

### Recepción 

El sistema UN-2000 envía una serie de resultados de las pruebas realizadas a las muestras de orina (UF y Chemistry) cuando se hace un envió al HOST se comunica usando protocolo TCP/IP al puerto 10001 al equipo principal donde está ubicado la interfaz de comunicación creada en GO, el equipo se comunica usando el estándar ASTM el cual envía tramas al HOST indicando información del usuario, orden y resultados, están información es almacenada en cadenas de texto

### Almacenamiento

El binario store se encarga de enviar las cadenas de texto que genera la recepción de resultados, genera una conexión a un servidor FTP y al enviar exitosamente el archivo de texto, lo elimina de la carpeta temporal

### Host Query

Cuando el equipo es alimentado con una muestra, lee el código de barras en el tubo de muestra, si la lectura es exitosa realiza una conexión con la interfaz al puerto 10002, esta comunicación comprende dos partes, primero el envío a la interfaz de la información almacenada en el código de barras (información del usuario y orden de trabajo). La segunda parte de la comunicación es iniciada por la interfaz (GO);

Al obtener la orden de trabajo la interfaz inicia una consulta a la base de datos MySQL consultado los datos de usuario (demográficos) y si se tiene programada una prueba (Código: C210 - Citoquímico de orina), se comprenden 3 casos posibles en este punto:

1. La orden de trabajo efectivamente tiene un examen a practicar asociado al equipo en cuestión, se genera un mensaje ASTM con la información de la orden y el examen a practicar a la muestra.

2. La orden de trabajo no tiene un examen a practicar asociado, se genera un mensaje ASTM indicando que no se debe realizar prueba a la muestra, se crea un reto en el sistema **4D** asociado a la OT describiendo el problema.

3. La orden de trabajo tiene un examen a practicar válido pero la muestra esta inactiva, se genera un mensaje ASTM indicando que no se debe realizar prueba a la muestra, se crea un reto en el sistema **4D** asociado a la OT describiendo el problema.

Una vez termina la comunicación la interfaz queda en modo escucha para la comunicación de la próxima muestra.

## Mantenimiento

La interfaz puede correr en entornos windows, mac o linux por medio de un binario que se compila.

|GOOS (OS) |	GOARCH (arquitectura)|
| ------------- | ------------- |
|android|arm|
|darwin|386|
|darwin|amd64|
|darwin|arm|
|darwin|arm64|
|dragonfly|amd64|
|freebsd|386|
|freebsd|amd64|
|freebsd|arm|
|linux|386|
|linux|amd64|
|linux|arm|
|linux|arm64|
|linux|ppc64|
|linux|ppc64le|
|linux|mips|
|linux|mipsle|
|linux|mips64|
|linux|mips64le|
|netbsd|386|
|netbsd|amd64|
|netbsd|arm|
|openbsd|386|
|openbsd|amd64|
|openbsd|arm|
|plan9|386|
|plan9|amd64|
|solaris|amd64|
|windows|386|
|windows|amd64|

*Para compilar se usa el siguiente código:*
```
env GOOS=linux GOARCH=amd64 go store.go
```
*Mas información en este [enlace](https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04).*

**1. El equipo UN-2000 muestra mensaje de error 00CC en la pantalla de resultados.**

* *Este mensaje significa error de comunicación con el HOST con la orden de trabajo específico, significa que esta muestra no puede ser procesada por no tener un examen asignado a este equipo o la muestra se encuentra inactiva se debe verificar los exámenes a realizar a la orden de trabajo con el error.*

**2. El equipo UN-2000 muestra alerta roja en la casilla SERVER.**

* *Este error indica problemas de comunicación con el UWAM, reportar el problema al encargado de mantenimiento del equipo.*

**3. El equipo UN-2000 muestra alerta roja en la casilla HOST.**

* *Este error indica problemas de comunicación con la interfaz del equipo, verificar si hay conexión a la ip del servidor de la interfaz, esta prueba se puede realizar haciendo ping a la direccion IP del servidor de la interfaz, se puede verificar la IP desde Vsphere en el equipo **interfaces-go** en el cluster **PURE_MDE**, se puede verificar la credenciales de acceso en el documento networking en google drive.*

* *El debug se realiza sobre el codigo sin compilar, no sobre el binario.*

* *Se recomienda usar [DELVE](https://github.com/derekparker/delve), ya que puede integrarse con SublimeText ([documentacion](https://github.com/dishmaev/GoDebug)) o VSCode ([documentacion](https://github.com/Microsoft/vscode-go/wiki/Debugging-Go-code-using-VS-Code)).*


**4. El equipo servidor no es accesible.**

* *Realizar un reinicio del equipo desde la plataforma Vsphere. Se realiza de la siguiente manera, ingresar a la dirección https://vcenter.mde.lmla.co/vsphere-client/?csp seleccionar el equipo **interfaces-go**  ubicado en **MDE** en el **cluster PURE_MDE** (se puede verificar la credenciales de acceso en el documento networking en google drive), para reiniciar el equipo seleccionar la opción **all actions** > **power** > **reset** en la parte superior de la pantalla de administración*

**5. Como bajar y levantar la interfaz.**

* *Se usa el comando*
```
ps aux | grep -i ./store
```

* *Se verifica el PID del primer proceso que se esta ejecutando, obviar el segundo resultado*
```
kill -9 <PID>
```
* *Verificar logs de errores.*

**6. Verificar que la interfaz este activa.**
* *Se usa el comando*
```
ps aux | grep -i ./store
```

# TRAMA ASTM UN-2000
| Trama | Identificador | Descripcion | 
| ------------- | ------------- | ------------- |
|Header|H|Contains sender and receiver information.|
|Patient Information|P|Contains patient information.|
|Analysis order or query|O-Q|Contains analysis order information.|
|Analysis result|R|Contains analysis result information.|
|Manufacturer information|M|Not used.|
|Message termination|L|Indicates the end of the message.|

## HEADER

```
<STX>1H|\^&|||U-WAM^00-06_Build007^A1159^^^^AU501736||||||||LIS2-A2|20180120044334
```
| Trama | Ejemplo |
| ------------- | ------------- |
|ASTM_HEADER_DELIMITER|1H|
|ASTM_HEADER_MESSAGE_CONTROL_ID|\\^&|
|ASTM_HEADER_PASSWORD|vacio|
|ASTM_HEADER_SENDERID|U-WAM^00-06_Build007^A1159^^^^AU501736|
|ASTM_HEADER_SENDERADDR|vacio|
|ASTM_HEADER_RESERVED|vacio|
|ASTM_HEADER_PHONENUMBER|vacio|
|ASTM_HEADER_SENDERDETAILS|vacio|
|ASTM_HEADER_RECEIVERID|vacio|
|ASTM_HEADER_COMMENTS|vacio|
|ASTM_HEADER_PROCESSINGID|vacio|
|ASTM_HEADER_VERSION|LIS2-A2|
|ASTM_HEADER_TIMESTAMP|20180120044334|

## PATIENT

```
<STX>2P|1||3666354||HORACIO DE JESUS^TORO||19421223|M||||||OPOS
```
  
| Trama | Ejemplo |
| ------------- | ------------- |
|ASTM_PATIENT_DELIMITER|2P|
|ASTM_PATIENT_SEQUENCE|1|
|ASTM_PATIENT_PRACTICED_PATIENT_ID|vacio|
|ASTM_PATIENT_LAB_PATIENT_ID|3666354|
|ASTM_PATIENT_ID3|vacio|
|ASTM_PATIENT_NAME|HORACIO DE JESUS^TORO|
|ASTM_PATIENT_MOTHERMAIDENNAME|vacio|
|ASTM_PATIENT_BIRTHDATE|19421223|
|ASTM_PATIENT_SEX|M|
|ASTM_PATIENT_RACE|vacio|
|ASTM_PATIENT_ADDRESS|vacio|
|ASTM_PATIENT_RESERVED_FIELD|vacio|
|ASTM_PATIENT_TELEPHONE|vacio|
|ASTM_PATIENT_PHYSICIAN_ID|vacio|
|ASTM_PATIENT_SPECIAL_FIELD1|OPOS|
|ASTM_PATIENT_SPECIAL_FIELD2|vacio|
|ASTM_PATIENT_HEIGHT|vacio|
|ASTM_PATIENT_WEIGHT|vacio|
|ASTM_PATIENT_DIAGNOSTIC|vacio|
|ASTM_PATIENT_MEDICATIONS|vacio|
|ASTM_PATIENT_DIET|vacio|
|ASTM_PATIENT_PRACTICE_F1|vacio|
|ASTM_PATIENT_PRACTICE_F2|vacio|
|ASTM_PATIENT_ADMISSION_DISCHARGE_DATES|vacio|
|ASTM_PATIENT_ADMISSION_STATUS|vacio|
|ASTM_PATIENT_LOCATION|vacio|
|ASTM_PATIENT_NATURE_ALT_DIAG_CODE|vacio|
|ASTM_PATIENT_ALTERNATIVE_DIAG_CODE|vacio|
|ASTM_PATIENT_PATIENT_RELIGION|vacio|
|ASTM_PATIENT_MARITAL_STATUS|vacio|
|ASTM_PATIENT_ISOLATION_STATUS|vacio|
|ASTM_PATIENT_LANGUAGE|vacio|
|ASTM_PATIENT_HOSPITAL_SERVICE|vacio|
|ASTM_PATIENT_HOSPITAL_INSTITUTION|vacio|
|ASTM_PATIENT_DOSAGE_CATEGORY|vacio|

## ORDER

```
<STX>3O|1|2523670||^^^Path_CAST\^^^BACT\|R||20170623091034||||N|||20170623110851|*||||||||||F
```
| Trama | Ejemplo |
| ------------- | ------------- |
|ASTM_ORDER_DELIMITER|3O|
|ASTM_ORDER_SECUENCE|1|
|ASTM_ORDER_SPECIMEN_ID|2523670|
|ASTM_ORDER_INSTRUMENT_SPECIMEN_ID|vacio|
|ASTM_ORDER_UNIVERSAL_TEST_ID|^^^Path_CAST\^^^BACT...|
|ASTM_ORDER_PRIORITY|R|
|ASTM_ORDER_REQUESTED_DATE_TIME|vacio|
|ASTM_ORDER_COLLECTION_DATE_TIME|20170623091034|
|ASTM_ORDER_COLLECTION_END_TIME|vacio|
|ASTM_ORDER_COLLECTION_VOLUME|vacio|
|ASTM_ORDER_COLLECTOR_ID|vacio|
|ASTM_ORDER_ACTION_CODE|N|
|ASTM_ORDER_DANGER_CODE|vacio|
|ASTM_ORDER_RELEVANT_CLINICAL_INFO|vacio|
|ASTM_ORDER_DATE_TIME_SPECIMEN_RECEIVED|20170623110851
|ASTM_ORDER_SPECIMEN_DESCRIPTOR|\*|
|ASTM_ORDER_ORDERING_PHYSICIAN|vacio|
|ASTM_ORDER_PHYSICIAN_PHONENUMBER|vacio|
|ASTM_ORDER_USER_FIELD1|vacio|
|ASTM_ORDER_USER_FIELD2|vacio|
|ASTM_ORDER_LABORATORY_FIELD1|vacio|
|ASTM_ORDER_LABORATORY_FIELD2|vacio|
|ASTM_ORDER_DATE_TIME|vacio|
|ASTM_ORDER_INSTRUMENT_CHARGE|vacio|
|ASTM_ORDER_INSTRUMENT_ID|vacio|
|ASTM_ORDER_REPORT_TYPE|F|
|ASTM_ORDER_RESERVED_FIELD|vacio|
|ASTM_ORDER_LOCATION_OF_SPECIMEN|vacio|
|ASTM_ORDER_NOSOCIMIAL_INFECTION_FLAG|vacio|
|ASTM_ORDER_SPECIMEN_SERVICE|vacio|
|ASTM_ORDER_SPECIMEN_INSTITUTION|vacio|

## RESULT

```
<STX>7R|2|^^^URO^A^1^S^  0010^02|2.0^MAINFORMAT|mg/dL||H||||^^admin^administrator||20180119233905|UC-3500
```

| Trama | Ejemplo |
| ------------- | ------------- |
|ASTM_RESULT_DELIMITER|7R|
|ASTM_RESULT_SECUENCE|2|
|ASTM_RESULT_TEST_ID|^^^URO^A^1^S^ 0010^02|
|ASTM_RESULT_DATA_MEASURE|2_0^MAINFORMAT|
|ASTM_RESULT_UNITS|mg/dL|
|ASTM_RESULT_REFERENCE_RANGES|vacio|
|ASTM_RESULT_RESULT_ABNORMAL_FLAGS|H|
|ASTM_RESULT_NATURE_OF_ABNORMALITY|vacio|
|ASTM_RESULT_RESULT_STATUS|vacio|
|ASTM_RESULT_DATE_OF_CHANGE|vacio|
|ASTM_RESULT_OPERATOR_IDENTIFICATION|^^admin^administrator|
|ASTM_RESULT_DATE_TIME_TEST_STARTED|vacio|
|ASTM_RESULT_DATE_TIME_TEST_COMPLETED|20180119233905|
|ASTM_RESULT_INSTRUMENT_IDENTIFICATION|UC-3500|

## TERMINATION

```
<STX>7L|1|N
```
  
| Trama | Ejemplo |
| ------------- | ------------- |
|ASTM_TERMINATOR_DELIMITER|7L|
|ASTM_TERMINATOR_SECUENCE|1|
|ASTM_TERMINATOR_TERMINATOR_CODE|N|
