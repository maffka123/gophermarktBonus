package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/maffka123/gophermarktBonus/internal/app"
	"github.com/maffka123/gophermarktBonus/internal/models"
	"github.com/maffka123/gophermarktBonus/internal/storage"

	"strings"

	"go.uber.org/zap"
)

type Handler struct {
	db     storage.DBinterface
	logger *zap.Logger
	ctx    context.Context
}

func NewHandler(ctx context.Context, db storage.DBinterface, logger *zap.Logger) Handler {
	return Handler{
		db:     db,
		logger: logger,
		ctx:    ctx,
	}
}

func (h *Handler) HandlerPostRegister(tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var u models.User
		err := decoder.Decode(&u)
		h.logger.Debug("recieved new user: ", zap.String("login", u.Login))

		if err != nil {
			http.Error(w, fmt.Sprintf("400 - Register json cannot be decoded: %s", err), http.StatusBadRequest)
			return
		}

		exists, err := h.db.CreateNewUser(h.ctx, u)
		if exists == -1 {
			http.Error(w, fmt.Sprintf("409 - Login is already taken: %s", err), http.StatusConflict)
			return
		} else if err != nil {
			http.Error(w, fmt.Sprintf("500 - Internal error: %s", err), http.StatusInternalServerError)
			return
		} else if exists == 1 {

			_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": u.ID})

			h.logger.Debug("logged in: ", zap.String("login", u.Login))
			http.SetCookie(w, &http.Cookie{
				Name:  "jwt",
				Value: tokenString,
			})
			w.Header().Set("application-type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok}`))
		}
	}
}

func (h *Handler) HandlerPostLogin(tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var u models.User
		err := decoder.Decode(&u)
		h.logger.Debug("user is trying to login: ", zap.String("login", u.Login))

		if err != nil {
			http.Error(w, fmt.Sprintf("400 - Register json cannot be decoded: %s", err), http.StatusBadRequest)
			return
		}

		pass, err := h.db.SelectPass(h.ctx, &u)

		if pass == nil || !app.ComparePass(h.ctx, *pass, u.Password) {
			http.Error(w, fmt.Sprintf("401 - user or password are wrong: %s", err), http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, fmt.Sprintf("500 - Internal error: %s", err), http.StatusInternalServerError)
			return
		}

		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": u.ID})

		h.logger.Debug("logged in: ", zap.String("login", u.Login))
		http.SetCookie(w, &http.Cookie{
			Name:  "jwt",
			Value: tokenString,
		})
		w.Header().Set("application-type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok}`))

	}
}

func (h *Handler) HandlerPostOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var valid bool
		order := models.Order{}
		orderID, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("400 - could not parse order id: %s", err), http.StatusBadRequest)
			return
		}
		valid, order.ID, err = app.PrepOrderNumber(h.ctx, orderID)
		if err != nil {
			http.Error(w, fmt.Sprintf("400 - could not parse order id: %s", err), http.StatusBadRequest)
			return
		} else if !valid {
			http.Error(w, fmt.Sprintf("422 - wrong format of the order number: %s", err), http.StatusUnprocessableEntity)
			return
		}

		h.logger.Debug("adding new order: ", zap.String("order", string(orderID)))

		currUser, err := app.UserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("400 - could not parse user id from token: %s", err), http.StatusBadRequest)
			return
		}

		expectedUser, err := h.db.SelectUserForOrder(h.ctx, order)

		if expectedUser != 0 && currUser == expectedUser {
			w.Header().Set("application-type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok}`))
			return
		} else if expectedUser != 0 && currUser != expectedUser {
			http.Error(w, fmt.Sprintf("409 - Order was used by diferent user: %s", err), http.StatusConflict)
			return
		} else if err != nil {
			h.logger.Debug(err.Error())
			http.Error(w, fmt.Sprintf("500 - Internal error: %s", err), http.StatusInternalServerError)
			return
		}

		order.UserID = int(currUser)
		order.Status = "NEW"
		order.Type = "top_up"

		err = h.db.InsertOrder(h.ctx, order)
		if err != nil {
			http.Error(w, fmt.Sprintf("500 - Internal error: %s", err), http.StatusInternalServerError)
			return
		}

		h.logger.Debug("order accepted: ", zap.String("login", string(orderID)))
		w.Header().Set("application-type", "text/plain")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"status":"ok}`))
	}
}

func (h *Handler) HandlerGetOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		h.logger.Debug("getting list of orders")

		currUser, err := app.UserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("400 - could not parse user id from token: %s", err), http.StatusBadRequest)
			return
		}

		orders, err := h.db.SelectAllOrders(h.ctx, currUser)
		if err != nil {
			http.Error(w, fmt.Sprintf("500 - internal server error: %s", err), http.StatusBadRequest)
			return
		} else if len(orders) == 0 {
			w.Header().Set("application-type", "text/plain")
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte(`{"status":"ok}`))
			return
		}

		mJSON, err := json.Marshal(orders)
		if err != nil {
			http.Error(w, fmt.Sprintf("500 - could not prepare data for return: %s", err), http.StatusBadRequest)
			return
		}

		h.logger.Debug("list of orders for user: ", zap.String("login", fmt.Sprint(currUser)))
		w.Header().Set("application-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mJSON))
	}
}

func (h *Handler) HandlerGetBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		h.logger.Debug("getting balance")

		currUser, err := app.UserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("400 - could not parse user id from token: %s", err), http.StatusBadRequest)
			return
		}

		balance, err := h.db.SelectBalance(h.ctx, currUser)
		if err != nil {
			http.Error(w, fmt.Sprintf("500 - internal server error: %s", err), http.StatusBadRequest)
			return
		}

		mJSON, err := json.Marshal(balance)
		if err != nil {
			http.Error(w, fmt.Sprintf("500 - could not prepare data for return: %s", err), http.StatusBadRequest)
			return
		}

		h.logger.Debug("balance for user: ", zap.String("login", fmt.Sprint(currUser)))
		w.Header().Set("application-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mJSON))
	}
}

func (h *Handler) HandlerPostWithdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		h.logger.Debug("withdraw")

		currUser, err := app.UserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("400 - could not parse user id from token: %s", err), http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var o models.Withdrawal
		err = decoder.Decode(&o)

		if err != nil && strings.Contains(err.Error(), "valid") {
			http.Error(w, fmt.Sprintf("422 - internal server error: %s", err), http.StatusUnprocessableEntity)
			return

		} else if err != nil {
			http.Error(w, fmt.Sprintf("500 - internal server error: %s", err), http.StatusBadRequest)
			return
		}

		balance, err := h.db.SelectBalance(h.ctx, currUser)
		if err != nil {
			http.Error(w, fmt.Sprintf("500 - internal server error: %s", err), http.StatusBadRequest)
			return
		}

		if balance.Current < o.Amount {
			http.Error(w, fmt.Sprintf("402 - currenct balance is not enough: %s", err), http.StatusPaymentRequired)
			return
		}

		order := models.Order{ID: o.ID, Amount: -o.Amount, UserID: int(currUser), Status: "PROCESSED", Type: "withdraw"}
		err = h.db.InsertOrder(h.ctx, order)
		if err != nil {
			http.Error(w, fmt.Sprintf("500 - internal server error: %s", err), http.StatusBadRequest)
			return
		}

		h.logger.Debug("withdraw succesfull for user: ", zap.String("login", fmt.Sprint(currUser)))
		w.Header().Set("application-type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok}`))
	}
}

func (h *Handler) HandlerGetWithdrawals() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		h.logger.Debug("withdrawals")

		currUser, err := app.UserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("400 - could not parse user id from token: %s", err), http.StatusBadRequest)
			return
		}

		orders, err := h.db.SelectAllWithdrawals(h.ctx, currUser)
		if err != nil {
			http.Error(w, fmt.Sprintf("500 - internal server error: %s", err), http.StatusBadRequest)
			return
		} else if len(*orders) == 0 {
			w.Header().Set("application-type", "text/plain")
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte(`{"status":"ok}`))
			return
		}

		mJSON, err := json.Marshal(orders)
		if err != nil {
			http.Error(w, fmt.Sprintf("500 - could not prepare data for return: %s", err), http.StatusBadRequest)
			return
		}

		h.logger.Debug("list of orders for user: ", zap.String("login", fmt.Sprint(currUser)))
		w.Header().Set("application-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mJSON))
	}
}
