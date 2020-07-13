package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kerti/balances/backend/handler/response"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/service"
	"github.com/kerti/balances/backend/util/ctxprops"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/logger"
	"github.com/satori/uuid"
)

// BankAccount handles all requests related to Bank Accounts
type BankAccount struct {
	Service *service.BankAccount `inject:"bankAccountService"`
}

// Startup perform startup functions
func (h *BankAccount) Startup() {
	logger.Trace("Bank Account Handler starting up...")
}

// Shutdown cleans up everything and shuts down
func (h *BankAccount) Shutdown() {
	logger.Trace("Bank Account Handler shutting down...")
}

// HandleCreateBankAccount handles the request
func (h *BankAccount) HandleCreateBankAccount(w http.ResponseWriter, r *http.Request) {
	var input model.BankAccountInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)
	bankAccount, err := h.Service.Create(input, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, bankAccount.ToOutput())
}

// HandleGetBankAccountByID handles the request
func (h *BankAccount) HandleGetBankAccountByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := uuid.FromString(vars["id"])
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	r.ParseForm()
	_, withBalances := r.Form["withBalances"]
	balanceCountStr, withBalanceCount := r.Form["balanceCount"]

	balanceCount := 0
	if withBalanceCount {
		balanceCount, err = strconv.Atoi(balanceCountStr[0])
		if err != nil {
			balanceCount = 10
		}
	}

	bankAccount, err := h.Service.GetByID(id, withBalances, balanceCount)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, bankAccount.ToOutput())
}

// HandleGetBankAccountByFilter handles the request
func (h *BankAccount) HandleGetBankAccountByFilter(w http.ResponseWriter, r *http.Request) {
	var input model.BankAccountFilterInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	bankAccounts, pageInfo, err := h.Service.GetByFilter(input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	outputs := make([]model.BankAccountOutput, 0)
	for _, bankAccount := range bankAccounts {
		output := bankAccount.ToOutput()
		outputs = append(outputs, output)
	}

	pageOutput := model.PageOutput{
		Items:    outputs,
		PageInfo: pageInfo,
	}

	response.RespondWithJSON(w, http.StatusOK, pageOutput)
}

// HandleUpdateBankAccount handles the request
func (h *BankAccount) HandleUpdateBankAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := uuid.FromString(vars["id"])
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	var input model.BankAccountInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	if input.ID.String() != id.String() {
		response.RespondWithError(w, failure.BadRequestFromString("id mismatch"))
		return
	}

	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)
	bankAccount, err := h.Service.Update(input, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, bankAccount.ToOutput())
}

// HandleDeleteBankAccount handles the request
func (h *BankAccount) HandleDeleteBankAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := uuid.FromString(vars["id"])
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)
	bankAccount, err := h.Service.Delete(id, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, bankAccount.ToOutput())
}

// HandleCreateBankAccountBalance handles the request
func (h *BankAccount) HandleCreateBankAccountBalance(w http.ResponseWriter, r *http.Request) {
	var input model.BankAccountBalanceInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)
	bankAccountBalance, err := h.Service.CreateBalance(input, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, bankAccountBalance.ToOutput())
}

// HandleGetBankAccountBalanceByID handles the request
func (h *BankAccount) HandleGetBankAccountBalanceByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := uuid.FromString(vars["id"])
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	bankAccountBalance, err := h.Service.GetBalanceByID(id)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, bankAccountBalance.ToOutput())
}

// HandleGetBankAccountBalanceByFilter handles the request
func (h *BankAccount) HandleGetBankAccountBalanceByFilter(w http.ResponseWriter, r *http.Request) {
	var input model.BankAccountBalanceFilterInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	bankAccountBalances, pageInfo, err := h.Service.GetBalancesByFilter(input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	outputs := make([]model.BankAccountBalanceOutput, 0)
	for _, bankAccountBalance := range bankAccountBalances {
		output := bankAccountBalance.ToOutput()
		outputs = append(outputs, output)
	}

	pageOutput := model.PageOutput{
		Items:    outputs,
		PageInfo: pageInfo,
	}

	response.RespondWithJSON(w, http.StatusOK, pageOutput)
}

// HandleUpdateBankAccountBalance handles the request
func (h *BankAccount) HandleUpdateBankAccountBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := uuid.FromString(vars["id"])
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	var input model.BankAccountBalanceInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	if input.ID.String() != id.String() {
		response.RespondWithError(w, failure.BadRequestFromString("id mismatch"))
		return
	}

	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)
	bankAccountBalance, err := h.Service.UpdateBalance(input, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, bankAccountBalance.ToOutput())
}

// HandleDeleteBankAccountBalance handles the request
func (h *BankAccount) HandleDeleteBankAccountBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := uuid.FromString(vars["id"])
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)
	bankAccountBalance, err := h.Service.DeleteBalance(id, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, bankAccountBalance.ToOutput())
}
