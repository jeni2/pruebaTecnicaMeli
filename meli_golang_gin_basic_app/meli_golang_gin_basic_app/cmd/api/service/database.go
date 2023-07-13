package service

import (
	"database/sql"
	"fmt"
	"meli_golang_gin_basic_app/cmd/api/model"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type IDatabase interface {
	Persist(requestBody *model.PersistRequest) int
	Scan(idDb string)
	GetClassification(id string) model.EsquemasData
	AddNewWord(word string) int
	GetWordList() model.WordListResponse
}

type Database struct {
}

func NewDatabase() *Database {
	return &Database{}
}

//funcion que almacena las credenciales de una base de datos

func (s *Database) Persist(requestBody *model.PersistRequest) int {
	db, err := sql.Open("mysql", "root:MySQLPassword2023@tcp(localhost:3306)/databasecredentials")
	if err != nil {
		fmt.Println("Error al conectar a la base de datos:", err)
		return -1
	}
	defer db.Close()

	// Insertar datos en la tabla
	insertQuery := "INSERT INTO databasecredentials.credentials (dbhost, dbport, dbusername, dbpassword) VALUES (?, ?,?,?) "
	result, err := db.Exec(insertQuery, requestBody.Host, requestBody.Port, requestBody.Username, requestBody.Password)
	if err != nil {
		fmt.Println("Error al insertar datos:", err)
		return -1
	}

	// Obtener el ID del último registro insertado
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error al obtener el último ID insertado:", err)
		return -1
	}
	fmt.Println("Último ID insertado:", lastInsertID)
	return int(lastInsertID)

}

//funcion que escanea la base de datos , la clasifica y la persiste la informacion clasificada

func (s *Database) Scan(idDb string) {
	fmt.Print("id a consultar: ", idDb)
	dbCredentials, err := sql.Open("mysql", "root:MySQLPassword2023@tcp(localhost:3306)/")
	if err != nil {
		fmt.Println("Error al conectar a la base de datos:", err)
		return
	}
	defer dbCredentials.Close()

	// Obtener credenciales de la base de datos especificada
	credenciales, err := dbCredentials.Query("SELECT * FROM databasecredentials.credentials WHERE id=" + idDb)
	if err != nil {
		fmt.Println("Error al obtener las credenciales:", err)
		return
	}
	defer credenciales.Close()
	// Leer los valores de las columnas en la fila
	var id int
	var dbHost string
	var dbPort int
	var dbUsername string
	var dbPassword string

	for credenciales.Next() {

		err := credenciales.Scan(&id, &dbHost, &dbPort, &dbUsername, &dbPassword)
		if err != nil {
			fmt.Println("Error al escanear los valores de la fila:", err)
			return
		}

		// imprimimos los valores leidos
		fmt.Println("Valores de la fila:", id, dbHost, dbPort, dbUsername, dbPassword)
	}

	// Verificar si hubo errores durante el recorrido de filas
	if err = credenciales.Err(); err != nil {
		fmt.Println("Error al recorrer las filas:", err)
		return
	}

	//nos conectamos a la db con las credenciales obtenidas

	db, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@tcp("+dbHost+":"+strconv.Itoa(dbPort)+")/")
	if err != nil {
		fmt.Println("Error al conectar a la base de datos con las credenciales obtenidas de la base de datos:", err)
		return
	}
	defer db.Close()
	fmt.Println("SE CONECTO A LA BASE DE DATOS")

	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE db_id = ?", "databaseClasification.dbschemas")

	// Ejecutar la sentencia SQL con el parámetro idDB (eliminacion de escaneo anterior)
	_, err = db.Exec(deleteQuery, idDb)
	if err != nil {
		fmt.Print("Error al borrar los datos de la tabla", "databaseClasification.dbschemas", err)
		return
	}

	fmt.Println("se borro correctamente el scaner anterior, scaneando de nuevo")

	informacionDB := model.EsquemasData{}
	// Obtener información sobre los esquemas
	schemas, err := db.Query("SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')")
	if err != nil {
		fmt.Println("Error al obtener los esquemas:", err)
		return
	}
	defer schemas.Close()

	for schemas.Next() {

		var schemaName string
		err := schemas.Scan(&schemaName)
		if err != nil {
			fmt.Println("Error al escanear los esquemas:", err)
			return
		}
		fmt.Println("Esquema:", schemaName)
		esquema1 := model.Esquema{EsquemaName: schemaName}

		// Obtener información sobre las tablas en el esquema actual
		tables, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = ?", schemaName)
		if err != nil {
			fmt.Println("Error al obtener las tablas:", err)
			return
		}
		defer tables.Close()

		for tables.Next() {

			var tableName string
			err := tables.Scan(&tableName)
			if err != nil {
				fmt.Println("Error al escanear las tablas:", err)
				return
			}
			fmt.Println("  Tabla:", tableName)
			tabla1 := model.Tabla{TableName: tableName}

			// Obtener información sobre las columnas de la tabla actual
			columns, err := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = ?", tableName)
			if err != nil {
				fmt.Println("Error al obtener las columnas:", err)
				return
			}
			defer columns.Close()

			for columns.Next() {
				var columnName, dataType string
				err := columns.Scan(&columnName, &dataType)
				if err != nil {
					fmt.Println("Error al escanear las columnas:", err)
					return
				}
				clasific := clasificarColumna(s.GetWordList().WordList, columnName)
				fmt.Println("    Columna:", columnName)
				fmt.Println("    Tipo de dato:", dataType)
				columna1 := model.Columna{ColumnName: columnName, Tipo: dataType, Clasificacion: clasific}
				tabla1.Columnas = append(tabla1.Columnas, columna1)
			}
			esquema1.Tablas = append(esquema1.Tablas, tabla1)
		}
		informacionDB.Esquemas = append(informacionDB.Esquemas, esquema1)
		informacionDB.DatabaseId = idDb
	}

	//fmt.Printf("informacionDB: %v\n", informacionDB)
	saveData(&informacionDB)

}

//Funcion que almacena la data en la base de datos a partir de un objeto esquemaData

func saveData(esquemasData *model.EsquemasData) {
	// Conexión a la base de datos MySQL
	db, err := sql.Open("mysql", "root:MySQLPassword2023@tcp(localhost:3306)/")
	if err != nil {
		fmt.Println("Error al conectar a la base de datos:", err)
		return
	}
	defer db.Close()
	// Insertar datos en la tabla dbSchemas
	for _, esquema := range esquemasData.Esquemas {
		insertSchemaQuery := "INSERT INTO databaseclasification.dbSchemas (schema_name, db_id, last_scan) VALUES (?,? ,?) ON DUPLICATE KEY UPDATE db_id = VALUES(db_id)"
		// Obtiene la fecha y hora actual
		currentTime := time.Now()
		schemaResult, err := db.Exec(insertSchemaQuery, esquema.EsquemaName, esquemasData.DatabaseId, currentTime)
		if err != nil {
			fmt.Println("Error al insertar datos en la tabla dbSchemas:", err)
			return
		}

		schemaID, err := schemaResult.LastInsertId()
		if err != nil {
			fmt.Println("Error al obtener el ID del esquema insertado:", err)
			return
		}
		// Insertar datos en la tabla dbTables
		for _, tabla := range esquema.Tablas {
			insertTableQuery := "INSERT INTO databaseclasification.dbTables (schema_id, table_name) VALUES (?, ?)"
			tableResult, err := db.Exec(insertTableQuery, schemaID, tabla.TableName)
			if err != nil {
				fmt.Println("Error al insertar datos en la tabla dbTables:", err)
				return
			}

			tableID, err := tableResult.LastInsertId()
			if err != nil {
				fmt.Println("Error al obtener el ID de la tabla insertada:", err)
				return
			}
			// Insertar datos en la tabla dbColumns
			for _, columna := range tabla.Columnas {
				insertColumnQuery := "INSERT INTO databaseclasification.dbColumns (table_id, column_name, data_type, classification) VALUES (?, ?, ?, ?)"
				_, err := db.Exec(insertColumnQuery, tableID, columna.ColumnName, columna.Tipo, columna.Clasificacion)
				if err != nil {
					fmt.Println("Error al insertar datos en la tabla dbColumns:", err)
					return
				}
			}

		}

	}
}

//Funcion que obtiene la informacion de la clasificacion de la base de datos

func (s *Database) GetClassification(id string) model.EsquemasData {

	response := model.EsquemasData{}
	// Conexión a la base de datos MySQL
	db, err := sql.Open("mysql", "root:MySQLPassword2023@tcp(localhost:3306)/")
	if err != nil {
		fmt.Println("Error al conectar a la base de datos:", err)
		return response
	}
	defer db.Close()

	// Consultar datos de la tabla dbSchemas
	schemasData := model.EsquemasData{}
	schemasQuery := "SELECT schema_name, last_scan FROM databaseclasification.dbSchemas WHERE db_id = ?"
	// Convertir cadena a entero
	idInt, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error al convertir el id en int, entero no valido :", err)
		return response
	}

	schemasResult, err := db.Query(schemasQuery, idInt)
	if err != nil {
		fmt.Println("Error al obtener datos de la tabla dbSchemas:", err)
		return response
	}
	defer schemasResult.Close()
	LastScan := []byte{}
	for schemasResult.Next() {
		esquema := model.Esquema{}

		err := schemasResult.Scan(&esquema.EsquemaName, &LastScan)
		if err != nil {
			fmt.Println("Error al escanear datos de la tabla databaseclasification.dbSchemas:", err)
			return response
		}

		// Consultar datos de la tabla dbTables para el esquema actual
		tablesQuery := "SELECT table_name FROM databaseclasification.dbTables WHERE schema_id = (SELECT schema_id FROM databaseclasification.dbSchemas WHERE schema_name = ? )"
		tablesResult, err := db.Query(tablesQuery, esquema.EsquemaName)
		if err != nil {
			fmt.Println("Error al obtener datos de la tabla dbTables:", err)
			return response
		}
		defer tablesResult.Close()

		for tablesResult.Next() {
			tabla := model.Tabla{}
			err := tablesResult.Scan(&tabla.TableName)
			if err != nil {
				fmt.Println("Error al escanear datos de la tabla databaseclasification.dbTables:", err)
				return response
			}

			// Consultar datos de la tabla dbColumns para la tabla actual
			columnsQuery := "SELECT column_name, data_type, classification FROM databaseclasification.dbColumns WHERE table_id = (SELECT table_id FROM databaseclasification.dbTables WHERE table_name = ?)"
			columnsResult, err := db.Query(columnsQuery, tabla.TableName)
			if err != nil {
				fmt.Println("Error al obtener datos de la tabla dbColumns:", err)
				return response
			}
			defer columnsResult.Close()

			for columnsResult.Next() {
				columna := model.Columna{}
				err := columnsResult.Scan(&columna.ColumnName, &columna.Tipo, &columna.Clasificacion)
				if err != nil {
					fmt.Println("Error al escanear datos de la tabla dbColumns:", err)
					return response
				}

				tabla.Columnas = append(tabla.Columnas, columna)
			}

			esquema.Tablas = append(esquema.Tablas, tabla)
		}

		schemasData.Esquemas = append(schemasData.Esquemas, esquema)
		schemasData.DatabaseId = id
		lastScanTime, err := time.Parse("2006-01-02 15:04:05", string(LastScan))
		if err != nil {
			fmt.Println("Error al convertir el valor en date:", err)

		}

		schemasData.LastScan = lastScanTime
	}

	// Convertir datos a formato JSON
	/*jsonData, err := json.Marshal(schemasData)
	if err != nil {
		fmt.Println("Error al convertir datos a JSON:", err)
		return response
	}*/
	if schemasData.DatabaseId == "" {
		fmt.Print("no se encontraron scaners de la base de datos")
		schemasData.DatabaseId = id
		schemasData.Esquemas = []model.Esquema{}
	}
	//fmt.Println("Datos consultados:\n", string(jsonData))
	return schemasData
}

func (s *Database) GetWordList() model.WordListResponse {
	response := model.WordListResponse{}
	wordList := []string{}
	// Conexión a la base de datos MySQL
	db, err := sql.Open("mysql", "root:MySQLPassword2023@tcp(localhost:3306)/")
	if err != nil {
		fmt.Println("Error al conectar a la base de datos:", err)
		return response
	}
	defer db.Close()

	query := "SELECT word FROM privateData.privateWord "
	wordsResult, err := db.Query(query)
	if err != nil {
		fmt.Println("Error al obtener datos de la tabla privateWord:", err)
		return response
	}
	defer wordsResult.Close()
	var word string

	for wordsResult.Next() {

		err := wordsResult.Scan(&word)
		if err != nil {
			fmt.Println("Error al escanear los valores de la fila:", err)
			return response
		}

		// imprimimos los valores leidos
		//fmt.Println("Valores de la fila:", word)
		wordList = append(wordList, word)
	}
	response.WordList = wordList
	return response

}

func (s *Database) AddNewWord(word string) int {
	// Conexión a la base de datos MySQL
	db, err := sql.Open("mysql", "root:MySQLPassword2023@tcp(localhost:3306)/")
	if err != nil {
		fmt.Println("Error al conectar a la base de datos:", err)

	}
	defer db.Close()

	// Insertar datos en la tabla
	insertQuery := "INSERT INTO privateData.privateWord (word) VALUES (?)"
	result, err := db.Exec(insertQuery, word)
	if err != nil {
		fmt.Println("Error al insertar datos:", err)
		return -1
	}

	// Obtener el ID del último registro insertado
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error al obtener el último ID insertado:", err)
		return -1
	}
	fmt.Println("Último ID insertado:", lastInsertID)
	return int(lastInsertID)
}

func clasificarColumna(wordList []string, columName string) string {

	columNamel := strings.ToLower(columName)
	resp := "N/A"
	for _, nombre := range wordList {
		nombre := strings.ToLower(nombre)
		lista := strings.Split(nombre, "_")

		/*if contienePalabras(columNamel, lista) {
			resp = nombre
		}*/

		if contienePalabras2(lista, columNamel) {
			fmt.Print("LA clasificacion de : " + columName + " seria : " + nombre)
			resp = strings.ToUpper(nombre)
		}

	}
	return resp
}

func contienePalabras2(list []string, p string) bool {
	// Verificar si la lista tiene solo una palabra
	if len(list) == 1 {
		return strings.Contains(p, list[0])
	}

	// Verificar si la palabra p contiene al menos list.size - 1 palabras de la lista
	count := 0
	for _, word := range list {
		if strings.Contains(p, word) {
			count++
		}
	}

	// Verificar la condición adicional
	if contieneName(list) {
		return count == len(list)
	}

	return count >= len(list)-1
}

func contieneName(list []string) bool {
	for _, word := range list {
		if strings.ToLower(word) == "name" {
			return true
		}
	}
	return false
}
