package model

import "time"

type Database struct {
	Response string
	//dto
}

type PersistResponse struct {
	Id int
}

type PersistRequest struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Message string
}

// Estructuras que modelan la informacion de la base de datos escaneada
type Columna struct {
	ColumnName    string `json:"columnName"`
	Tipo          string `json:"tipo"`
	Clasificacion string `json:"Clasificacion"`
}

type Tabla struct {
	TableName string    `json:"TableName"`
	Columnas  []Columna `json:"columnas"`
}

type Esquema struct {
	EsquemaName string  `json:"EsquemaName"`
	Tablas      []Tabla `json:"Tablas"`
}

type EsquemasData struct {
	DatabaseId string    `json:"databaseId"`
	LastScan   time.Time `json:"last_scan"`
	Esquemas   []Esquema `json:"Esquemas"`
}

type WordListResponse struct {
	WordList []string `json:"wordList"`
}
