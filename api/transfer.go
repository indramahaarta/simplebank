package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/indramhrt/simplebank/db/sqlc"
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

	if !server.validAccount(ctx, arg.FromAccountID, req.Currency) {
		return
	}

	if !server.validAccount(ctx, arg.ToAccountID, req.Currency) {
		return
	}

	result, err := server.store.TransferTx(ctx, db.TransferTxParams(arg))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, err)
			return false
		}

		ctx.JSON(http.StatusInternalServerError, err)
		return false
	}

	if account.Currency != currency {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("currency mismatch error")))
		return false
	}

	return true
}
