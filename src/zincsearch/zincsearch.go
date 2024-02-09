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

func QueryByEmail(fieldName, value string) ([]byte, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando el archivo .env")
	}

	// Obtener las variables de entorno
	url := os.Getenv("API_URL")
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
		"max_results": 1,
		"_source": [
			"to","from","subject","body"
		]
	}`, value, fieldName)
	fmt.Println(query)
	req, err := http.NewRequest("POST", url, strings.NewReader(query))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(user, password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

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
