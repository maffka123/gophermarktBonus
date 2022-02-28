package handlers

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/maffka123/gophermarktBonus/internal/storage"
	"go.uber.org/zap"
)

func BonusRouter(ctx context.Context, db storage.DBinterface, secret string, logger *zap.Logger) chi.Router {

	r := chi.NewRouter()
	mh := NewHandler(ctx, db, logger)
	tokenAuth := jwtauth.New("HS256", []byte(secret), nil)

	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(jwtauth.Verifier(tokenAuth))
	//r.Use(jwtauth.Authenticator)

	r.Route("/api/user/", func(r chi.Router) {
		r.Post("/register", Conveyor(mh.HandlerPostRegister(), unpackGZIP, checkForJSON))
		r.Post("/login", Conveyor(mh.HandlerPostLogin(tokenAuth), unpackGZIP, checkForJSON))
		r.With(jwtauth.Authenticator).Post("/orders", Conveyor(mh.HandlerPostOrders(), unpackGZIP, checkForText))
		r.Get("/orders", Conveyor(mh.HandlerGetOrders(), unpackGZIP, packGZIP))
		r.Route("/balance", func(r chi.Router) {
			r.Get("/", Conveyor(mh.HandlerGetBalance(), unpackGZIP))
			r.Post("/withdraw", Conveyor(mh.HandlerPostWithdraw(), unpackGZIP))
			r.Get("/withdrawals", Conveyor(mh.HandlerGetWithdrawals(), unpackGZIP))
		})

	})

	return r
}
