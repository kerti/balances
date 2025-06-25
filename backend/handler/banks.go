package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/kerti/balances/backend/handler/response"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/service"
	"github.com/kerti/balances/backend/util/cachetime"
	"github.com/kerti/balances/backend/util/ctxprops"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/logger"
)

// BankAccount is the handler interface for Bank Accounts
type BankAccount interface {
	Startup()
	Shutdown()
	HandleCreateBankAccount(w http.ResponseWriter, r *http.Request)
	HandleGetBankAccountByID(w http.ResponseWriter, r *http.Request)
	HandleGetBankAccountByFilter(w http.ResponseWriter, r *http.Request)
	HandleUpdateBankAccount(w http.ResponseWriter, r *http.Request)
	HandleDeleteBankAccount(w http.ResponseWriter, r *http.Request)
	HandleCreateBankAccountBalance(w http.ResponseWriter, r *http.Request)
	HandleGetBankAccountBalanceByID(w http.ResponseWriter, r *http.Request)
	HandleGetBankAccountBalanceByFilter(w http.ResponseWriter, r *http.Request)
	HandleUpdateBankAccountBalance(w http.ResponseWriter, r *http.Request)
	HandleDeleteBankAccountBalance(w http.ResponseWriter, r *http.Request)
}

// BankAccountImpl is the handler implementation for Bank Accounts
type BankAccountImpl struct {
	Service service.BankAccount `inject:"bankAccountService"`
}

// Startup performs startup functions
func (h *BankAccountImpl) Startup() {
	logger.Trace("Bank Account Handler starting up...")
}

// Shutdown cleans up everything and shuts down
func (h *BankAccountImpl) Shutdown() {
	logger.Trace("Bank Account Handler shutting down...")
}

// HandleCreateBankAccount handles the request
func (h *BankAccountImpl) HandleCreateBankAccount(w http.ResponseWriter, r *http.Request) {
	input, err := h.getInputFromRequest(w, r)
	if err != nil {
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
func (h *BankAccountImpl) HandleGetBankAccountByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(w, r)
	if err != nil {
		return
	}

	r.ParseForm()
	_, withBalances := r.Form["withBalances"]
	balanceStartDateStr, withBalanceStartDate := r.Form["balanceStartDate"]
	balanceEndDateStr, withBalanceEndDate := r.Form["balanceEndDate"]
	pageSizeStr, withPageSize := r.Form["pageSize"]

	var balanceStartDate cachetime.NCacheTime
	if withBalanceStartDate {
		balanceStartDate.Scan(balanceStartDateStr[0])
	}

	var balanceEndDate cachetime.NCacheTime
	if withBalanceEndDate {
		balanceEndDate.Scan(balanceEndDateStr[0])
	}

	var pageSize *int = nil
	if withPageSize {
		parsedPageSize, err := strconv.Atoi(pageSizeStr[0])
		if err == nil {
			pageSize = &parsedPageSize
		}
	}

	bankAccount, err := h.Service.GetByID(id, withBalances, balanceStartDate, balanceEndDate, pageSize)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, bankAccount.ToOutput())
}

// HandleGetBankAccountByFilter handles the request
func (h *BankAccountImpl) HandleGetBankAccountByFilter(w http.ResponseWriter, r *http.Request) {
	var input model.BankAccountFilterInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	bankAccounts, pageInfo, err := h.Service.GetByFilter(input)
	if err != nil {
		response.RespondWithError(w, err)
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
func (h *BankAccountImpl) HandleUpdateBankAccount(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(w, r)
	if err != nil {
		return
	}

	input, err := h.getInputFromRequest(w, r)
	if err != nil {
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

	response.RespondWithJSON(w, http.StatusOK, bankAccount.ToOutput())
}

// HandleDeleteBankAccount handles the request
func (h *BankAccountImpl) HandleDeleteBankAccount(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(w, r)
	if err != nil {
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
func (h *BankAccountImpl) HandleCreateBankAccountBalance(w http.ResponseWriter, r *http.Request) {
	input, err := h.getBalanceInputFromRequest(w, r)
	if err != nil {
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
func (h *BankAccountImpl) HandleGetBankAccountBalanceByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(w, r)
	if err != nil {
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
func (h *BankAccountImpl) HandleGetBankAccountBalanceByFilter(w http.ResponseWriter, r *http.Request) {
	var input model.BankAccountBalanceFilterInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	bankAccountBalances, pageInfo, err := h.Service.GetBalancesByFilter(input)
	if err != nil {
		response.RespondWithError(w, err)
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
func (h *BankAccountImpl) HandleUpdateBankAccountBalance(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(w, r)
	if err != nil {
		return
	}

	input, err := h.getBalanceInputFromRequest(w, r)
	if err != nil {
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
func (h *BankAccountImpl) HandleDeleteBankAccountBalance(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(w, r)
	if err != nil {
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

func (h *BankAccountImpl) getInputFromRequest(w http.ResponseWriter, r *http.Request) (input model.BankAccountInput, err error) {
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
	}

	return
}

func (h *BankAccountImpl) getBalanceInputFromRequest(w http.ResponseWriter, r *http.Request) (input model.BankAccountBalanceInput, err error) {
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
	}

	return
}
