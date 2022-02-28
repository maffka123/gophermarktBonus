package app

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/maffka123/gophermarktBonus/internal/config"
	"github.com/maffka123/gophermarktBonus/internal/models"
	"github.com/maffka123/gophermarktBonus/internal/storage"
	"github.com/theplant/luhn"
	"go.uber.org/zap"
)

func ComparePass(ctx context.Context, expected string, actual string) bool {
	return (subtle.ConstantTimeCompare([]byte(expected), []byte(actual)) == 1)

}

func PrepOrderNumber(ctx context.Context, n []byte) (bool, int64, error) {
	id, err := strconv.Atoi(string(n))
	if err != nil {
		return false, 0, fmt.Errorf("convert to int fialed: %v", err)
	}

	return luhn.Valid(id), int64(id), nil
}

func UserIDFromContext(ctx context.Context) (int64, error) {
	_, uID, err := jwtauth.FromContext(ctx)
	if err != nil {
		return 0, err
	}
	return int64(uID["user_id"].(float64)), nil
}

func UpdateStatus(ctx context.Context, t <-chan time.Time, logger *zap.Logger, db storage.DBinterface, cfg *config.Config) {
	client := &http.Client{}
	for {
		select {
		case <-t:
			logger.Info("starting bonus update")
			oin := make(chan []models.Order)
			oout := make(chan models.Order)
			go db.SelectOrdersForUpdate(ctx, cfg, oin, oout)
			go getAccrual(ctx, cfg, oin, oout, logger, client)
		case <-ctx.Done():
			logger.Info("context canceled")
		}
	}
}

func getAccrual(ctx context.Context, cfg *config.Config, oin chan []models.Order, oout chan models.Order, logger *zap.Logger, client *http.Client) {
	url := fmt.Sprintf("http://%s/api/orders/", cfg.AccrualSystem)
	var intermOrder models.AccrualOrder
	orders := <-oin

	for _, order := range orders {
		url += fmt.Sprint(order.ID)
		request, err := http.NewRequest(http.MethodGet, url, nil)

		if err != nil {
			logger.Fatal("request creation failed", zap.Error(err))
		}

		response, requestErr := requestWithRetry(client, request, logger)
		if requestErr != nil {
			logger.Error(requestErr.Error())
		}
		defer response.Body.Close()

		decoder := json.NewDecoder(response.Body)
		decoder.Decode(&intermOrder)

		oout <- models.Order{ID: intermOrder.ID, Amount: intermOrder.Amount, Status: intermOrder.Status}

	}
	close(oout)
	logger.Info("bonus update finished")
}

func requestWithRetry(client *http.Client, request *http.Request, logger *zap.Logger) (*http.Response, error) {
	var response *http.Response
	var requestErr error
	for i := 0; i < 5; i++ {
		response, requestErr = client.Do(request)
		if requestErr != nil {
			logger.Info("Retrying: " + requestErr.Error())
		} else if response.StatusCode == 429 {
			logger.Info("Too many requests")
			time.Sleep(30 * time.Second)
		} else if response.StatusCode == 200 {
			return response, nil
		}
		logger.Info("Retrying...")
		time.Sleep(time.Duration(i*10) * time.Second)
	}

	return response, requestErr
}
