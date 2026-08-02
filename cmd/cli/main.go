package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	var (
		serverAddr string
		from       string
		to         string
		amount     float64
		list       bool
	)

	flag.StringVar(&serverAddr, "server", "http://localhost:8080", "адрес сервера")
	flag.StringVar(&from, "from", "", "исходная валюта (например USD)")
	flag.StringVar(&to, "to", "", "целевая валюта (например RUB)")
	flag.Float64Var(&amount, "amount", 0, "сумма для конвертации")
	flag.BoolVar(&list, "list", false, "показать доступные валюты")
	flag.Parse()

	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	if list {
		showCurrencies(serverAddr)
		return
	}

	if from == "" || to == "" || amount <= 0 {
		fmt.Println("Использование: cli -from USD -to RUB -amount 100")
		os.Exit(1)
	}

	convertCurrency(serverAddr, from, to, amount)
}

func convertCurrency(serverAddr, from, to string, amount float64) {
	reqURL := fmt.Sprintf("%s/convert?from=%s&to=%s&amount=%.2f",
		serverAddr, url.QueryEscape(from), url.QueryEscape(to), amount)

	resp, err := http.Get(reqURL) //nolint:gosec,noctx
	if err != nil {
		fmt.Printf("Ошибка соединения: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var errResp map[string]string
		if json.Unmarshal(body, &errResp) == nil {
			fmt.Printf("Ошибка: %s\n", errResp["error"])
		} else {
			fmt.Printf("Ошибка сервера (статус %d)\n", resp.StatusCode)
		}
		os.Exit(1)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Ошибка парсинга ответа: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result["result"], to)
}

func showCurrencies(serverAddr string) {
	resp, err := http.Get(serverAddr + "/currencies") //nolint:gosec,noctx
	if err != nil {
		fmt.Printf("Ошибка соединения: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Ошибка сервера: %s\n", body)
		os.Exit(1)
	}

	var currencies map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&currencies); err != nil {
		fmt.Printf("Ошибка парсинга ответа: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Доступные валюты:")
	for code, name := range currencies {
		fmt.Printf("%s - %s\n", code, name)
	}
}
