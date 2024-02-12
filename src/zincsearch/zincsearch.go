package zincsearch

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func QueryByEmailField(fieldName, value string) ([]byte, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando el archivo .env")
	}

	// Obtener las variables de entorno
	url := os.Getenv("API_URL") + "/emails/_search"
	user := os.Getenv("API_USER")
	password := os.Getenv("API_PASSWORD")

	query := fmt.Sprintf(`{
		"search_type": "match",
		"query": {
			"term": "%s",
			"field": "%s"
		},
		"sort_fields": ["-@timestamp"],
		"from": 0,
		"max_results": 3,
		"_source": [
			"to","from","subject"
		]
	}`, value, fieldName)
	fmt.Println(query)
	req, err := http.NewRequest("POST", url, strings.NewReader(query))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(user, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body, err
}

func QueryByEmailId(id string) ([]byte, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando el archivo .env")
	}

	// Obtener las variables de entorno
	url := os.Getenv("API_URL") + "/emails/_doc/" + id
	user := os.Getenv("API_USER")
	password := os.Getenv("API_PASSWORD")
	log.Println("url:", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(user, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body, err
}
