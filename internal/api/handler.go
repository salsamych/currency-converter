// Package api реализует HTTP обработчики для сервиса конвертации валют.
package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"currency-converter/internal/converter"
)

// Handler обрабатывает HTTP запросы.
type Handler struct {
	converter *converter.Converter
}

// NewHandler создаёт новый обработчик.
func NewHandler(c *converter.Converter) *Handler {
	return &Handler{converter: c}
}

// Convert обрабатывает запрос на конвертацию валют.
func (h *Handler) Convert(w http.ResponseWriter, r *http.Request) {
	from := strings.ToUpper(r.URL.Query().Get("from"))
	to := strings.ToUpper(r.URL.Query().Get("to"))
	amountStr := r.URL.Query().Get("amount")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "некорректная сумма",
		})
		return
	}

	result, err := h.converter.Convert(from, to, amount)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"result": result,
		"from":   from,
		"to":     to,
		"amount": amount,
	})
}

// Currencies возвращает список доступных валют.
func (h *Handler) Currencies(w http.ResponseWriter, _ *http.Request) {
	currencies, err := h.converter.GetAvailableCurrencies()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "не удалось получить список валют",
		})
		return
	}

	writeJSON(w, http.StatusOK, currencies)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
