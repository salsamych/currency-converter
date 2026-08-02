// Сервис конвертации валют с использованием курсов ЦБ РФ.
package main

import (
	"log"
	"os"
	"time"

	"currency-converter/internal/api"
	"currency-converter/internal/cache"
	"currency-converter/internal/converter"
	"currency-converter/internal/provider"
)

func main() {
	addr := ":8080"
	if env := os.Getenv("ADDR"); env != "" {
		addr = env
	}

	rateProvider := provider.NewCBRProvider()
	ratesCache := cache.NewCache(10 * time.Minute)
	conv := converter.NewConverter(rateProvider, ratesCache)
	handler := api.NewHandler(conv)
	server := api.NewServer(addr, handler)

	log.Printf("Сервис конвертации валют запущен на %s", addr)
	log.Fatal(server.ListenAndServe())
}
