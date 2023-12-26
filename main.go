package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	expectedToken = "4321"
	updateURL     = "http://localhost:8000/answ/update_async/"
)

type HRResult struct {
	AnswID string `json:"answ_id"`
	Suite  int    `json:"suite"`
	Token  string `json:"token"`
}

func main() {
	http.HandleFunc("/async_task", handleProcess)
	fmt.Println("Server running at port :8088")
	http.ListenAndServe(":8088", nil)
}

func handleProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	answid := r.FormValue("answ_id")
	token := r.FormValue("token")
	fmt.Println(answid, token)

	if token == "" || token != expectedToken {
		http.Error(w, "Токены не совпадают", http.StatusForbidden)
		fmt.Println("Токены не совпадают")
		fmt.Println(token, expectedToken)
		return
	}

	w.WriteHeader(http.StatusOK)

	go func() {
		delay := 10
		time.Sleep(time.Duration(delay) * time.Second)

		result := rand.Intn(101)

		// Отправка результата на другой сервер
		expResult := HRResult{
			AnswID: answid,
			Suite:  result,
			// Token:  token,
		}
		fmt.Println("json", expResult)
		jsonValue, err := json.Marshal(expResult)
		if err != nil {
			fmt.Println("Ошибка при маршализации JSON:", err)
			return
		}

		req, err := http.NewRequest(http.MethodPost, updateURL, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Println("Ошибка при создании запроса на обновление:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		answ, err := client.Do(req)
		if err != nil {
			fmt.Println("Ошибка при отправке запроса на обновление:", err)
			return
		}
		defer answ.Body.Close()

		fmt.Println("Ответ от сервера обновления:", answ.Status)
	}()
}
