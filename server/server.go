package server

import (
	"bank-service/stg"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type BalanceServer struct {
	Src *stg.BalanceService
	Mux *http.ServeMux
}

type err struct {
	Error string `json:"error"`
}

func writeInvalidJSON(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(err{Error: "invalid json"})
}

func writeInvalidJSONError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	_ = json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}
func NewBalanceServer(service stg.BalanceService) *BalanceServer {
	return &BalanceServer{Src: &service, Mux: http.NewServeMux()}
}

func (server *BalanceServer) SetServer() {
	server.Mux.HandleFunc("GET /api/get_balance/{id}", server.GetBalanceHandler)
	server.Mux.HandleFunc("POST /api/deposit/{id}", server.DepositHandler)
	server.Mux.HandleFunc("POST /api/transfer/{id}", server.TransferHandler)
	server.Mux.HandleFunc("GET /live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func (server *BalanceServer) DepositHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var request struct {
		Amount int `json:"add_sum"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeInvalidJSON(w, http.StatusBadRequest)
		return
	}

	new_money_amount := request.Amount
	userid, _ := strconv.Atoi(r.PathValue("id"))
	balance, _ := server.Src.DepositBalance(userid, new_money_amount)

	response := map[string]int{
		"balance": balance,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (server *BalanceServer) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Валидация и преобразование ID
	idStr := r.PathValue("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Проверка что ID положительный
	if userID <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Получение баланса с обработкой ошибок
	balance, err := server.Src.GetBalance(userID)
	if err != nil {
		// Логируем ошибку для отладки
		log.Printf("GetBalance error for user %d: %v", userID, err)

		// В зависимости от типа ошибки возвращаем соответствующий статус
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Формирование ответа
	response := map[string]int{
		"balance": balance,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (server *BalanceServer) TransferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var request struct {
		OtherUserId int `json:"other_id"`
		Amount      int `json:"add_sum"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeInvalidJSON(w, http.StatusBadRequest)
		return
	}

	userid, _ := strconv.Atoi(r.PathValue("id"))
	amount := request.Amount
	other_userid := request.OtherUserId

	cur_user_balance, other_user_balance, err := server.Src.TransferMoney(userid, other_userid, amount)

	if err != nil {
		writeInvalidJSONError(w, err)
		return
	}

	response := map[string]int{
		"current_user_balance": cur_user_balance,
		"other_user_balance":   other_user_balance,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
