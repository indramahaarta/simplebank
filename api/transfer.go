package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/indramhrt/simplebank/db/sqlc"
	"github.com/indramhrt/simplebank/token"
)

type transferMoneyRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) transferMoney(ctx *gin.Context) {
	var req transferMoneyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	account, valid := server.validAccount(ctx, arg.FromAccountID, req.Currency)
	if !valid {
		return
	}

	_, valid = server.validAccount(ctx, arg.ToAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != account.Owner {
		error := errors.New("you dont't have any authorization for this resource")
		ctx.JSON(http.StatusUnauthorized, errorResponse(error))
		return
	}

	if account.Balance < req.Amount {
		error := errors.New("you don't have enough balance")
		ctx.JSON(http.StatusBadRequest, errorResponse(error))
		return
	}

	result, err := server.store.TransferTx(ctx, db.TransferTxParams(arg))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (*db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, err)
			return nil, false
		}

		ctx.JSON(http.StatusInternalServerError, err)
		return nil, false
	}

	if account.Currency != currency {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("currency mismatch error")))
		return nil, false
	}

	return &account, true
}
