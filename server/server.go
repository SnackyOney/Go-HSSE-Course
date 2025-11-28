package server

import (
	"bank-service/stg"
	"encoding/json"
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

func NewBalanceServer(service stg.BalanceService) *BalanceServer {
	return &BalanceServer{Src: &service, Mux: http.NewServeMux()}
}

func (server *BalanceServer) SetServer() {
	server.Mux.HandleFunc("api/get_balance/{id}", server.GetBalanceHandler)
	server.Mux.HandleFunc("api/deposit/{id}", server.DepositHandler)
	server.Mux.HandleFunc("api/transfer/{id}", server.TransferHandler)
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

	server.Src.TransferMoney(userid, other_userid, amount)

	cur_user_balance, _ := server.Src.GetBalance(userid)
	other_user_balance, _ := server.Src.GetBalance(userid)

	response := map[string]int{
		"current_user_balance": cur_user_balance,
		"other_user_balance":   other_user_balance,
	}

	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
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

	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}

func (server *BalanceServer) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	userid, _ := strconv.Atoi(r.PathValue("id"))
	balance, _ := server.Src.GetBalance(userid)

	response := map[string]int{
		"balance": balance,
	}

	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
