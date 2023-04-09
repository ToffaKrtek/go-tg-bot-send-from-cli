package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	Token  string
	ChatId string
)

func getEnv(variable string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	env, ok := os.LookupEnv(variable)
	if !ok {
		return "", fmt.Errorf("%s не существует.", variable)
	}
	if len(strings.TrimSpace(env)) == 0 {
		return "", fmt.Errorf("%s пустой", env)
	}
	return env, nil
}

func getUrl() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s", Token)
}

func sendMessage(message string) (bool, error) {
	var response *http.Response
	var err error
	url := fmt.Sprintf("%s/sendMessage", getUrl())
	body, _ := json.Marshal(map[string]string{
		"chat_id": ChatId,
		"text":    message,
	})
	response, err = http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}
	fmt.Printf("Отправлено сообщение '%s'", message)
	fmt.Printf("Ответ : %s", string(body))
	return true, nil
}

func main() {
	var err error
	Token, err = getEnv("TOKEN")
	if err != nil {
		log.Fatal("%s", err)
	}
	ChatId, err = getEnv(("CHAT_ID"))
	if err != nil {
		log.Fatal("%s", err)
	}
	var message string
	flag.StringVar(&message, "msg", "Msg", "Сообщение в чат")
	flag.Parse()

	_, err = sendMessage(message)
	if err != nil {
		log.Fatal("%s", err)
	}
}
