package DB

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "postgres"
// 	dbname   = "footai"
// )
func ConnectPsql(user string, pass string, host string, port string, dbname string) (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		user, pass, host, port, dbname)

	db, err := sql.Open("pgx", psqlInfo)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("Connected to PostgreSQL!")
	return db, nil
}

func PrintAllRows(rows *sql.Rows) ([]map[string]interface{}, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[Print rows]Recovered from a potential panic:", err)
			// Log the error for further debugging
			// Potentially return a default value or an error
		}
	}()

    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    } 

    values := make([]interface{}, len(columns))
	var allRows []map[string]interface{}

    for i := range values {
        values[i] = new(interface{})
    }

    for rows.Next() {
		rowValues := make(map[string]interface{})

        err = rows.Scan(values...)
        if err != nil {
            return nil, err
        }

        i := 0
        for _, value := range values {
			rowValues[columns[i]] = *(value.(*interface{}))
			i++
        }
        fmt.Println()
		allRows = append(allRows, rowValues)

    }
    return allRows, nil 
}

func StoreLog(db *sql.DB, logType string, server string, description string) error {
	query := `
		INSERT INTO logs (timestamp, type, server, description) 
		VALUES ($1, $2, $3, $4) 
	`
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	_, err := db.Exec(query, timestamp, logType, server, description)
	if err != nil {
		return fmt.Errorf("error inserting log: %w", err)
	}

	return nil
}

func CreateSchema(db *sql.DB, schemaName string) error {
	query := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName)

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating schema: %v", err)
	}

	fmt.Println("Schema created successfully (if it didn't already exist).")
	return nil
}
