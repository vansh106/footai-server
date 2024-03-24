package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	DB "footai.com/m/Db"
	Gpt "footai.com/m/Gpt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	openai "github.com/sashabaranov/go-openai"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	err := godotenv.Load("cred.env")
	if err != nil {
		fmt.Println("[server] Error loading .env file")
	}

	ctx := context.Background()
	client := Gpt.Initialize(os.Getenv("OPENAI_API_KEY"))

	db, err := DB.ConnectPsql(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASS"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	router := gin.Default()

	router.GET("/", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"response": "working",
		})
	})

	router.POST("/post", func(c *gin.Context) {
		if c.Request.Method != "POST" {
			c.String(http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		requestBody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, "Error reading request body")
			return
		}

		var data map[string]interface{}
		err = json.Unmarshal(requestBody, &data)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid JSON format")
			return
		}

		prompt, ok := data["prompt"].(string)
		if !ok {
			c.String(http.StatusBadRequest, "Missing 'name' field")
			return
		}
		// prompt = "How many points does manchester united need to get on the top of the table?"

		res, err := Run(db, client, ctx, prompt)
		if err != nil {
			c.String(http.StatusBadRequest, "Couldn't perform operation!"+prompt)
		}
		c.JSON(http.StatusOK, gin.H{
			"response": res,
		})

	})

	fmt.Println("Server listening on port 3000")
	log.Fatal(router.Run("0.0.0.0:3000"))

}

func Run(db *sql.DB, client *openai.Client, ctx context.Context, prompt string) (string, error) {
	query, err := Gpt.GenerateChat(client, ctx, DB.GenSqlQuery(prompt))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(query)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[MAIN] Query didn't learn", err)
		}
	}()

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		DB.StoreLog(db, "error", prompt, "", err.Error())
		return "", err
	}
	defer rows.Close()

	data, err := DB.PrintAllRows(rows)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// fmt.Println(data)

	result, err := Gpt.GenerateChat(client, ctx, DB.GenPrompt(prompt, data))
	if err != nil {
		fmt.Println(err)
		DB.StoreLog(db, "error", prompt, query, err.Error())
		return "", err
	}
	fmt.Println(prompt + "\n" + result)
	err  = DB.StoreLog(db, "success", prompt, query, result)
	if err != nil {
		fmt.Println("[INSERTING LOGS]",err)
	}
	return result, nil
}
