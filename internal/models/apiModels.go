package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/theplant/luhn"
)

type User struct {
	Id       int64   `json:"id,omitempty"`
	Login    string  `json:"login"`
	Balance  float64 `json:"balance"`
	Password string  `json:"password"`
}

type Order struct {
	Id     int64     `json:"number,omitempty"`
	Status string    `json:"status,omitempty"`
	Amount int64     `json:"accrual,omitempty"`
	Date   time.Time `json:"uploaded_at,omitempty"`
	Type   string    `json:"type,omitempty"`
	UserID int       `json:"user_id,omitempty"`
}

type Withdrawal struct {
	Id     int64     `json:"order,omitempty"`
	Amount int64     `json:"sum,omitempty"`
	Date   time.Time `json:"processed_at,omitempty"`
}

type Balance struct {
	Current   int64 `json:"current"`
	Withdrawn int64 `json:"withdrawn"`
}

func (b *Balance) MarshalJSON() ([]byte, error) {
	type newBalance struct {
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}

	nb := newBalance{Current: float64(b.Current) / 100,
		Withdrawn: float64(b.Withdrawn) / 100}

	return json.Marshal(nb)
}

func (b *Balance) UnmarshalJSON(data []byte) error {
	type newBalance struct {
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}

	var nu newBalance

	if err := json.Unmarshal(data, &nu); err != nil {
		return err
	}

	b.Current = int64(nu.Current * 100)
	b.Withdrawn = int64(nu.Withdrawn * 100)

	return nil
}

func (u *User) UnmarshalJSON(data []byte) error {
	type newU User
	nu := (*newU)(u)

	if err := json.Unmarshal(data, &nu); err != nil {
		return err
	}

	np := sha256.Sum256([]byte(u.Password))
	u.Password = hex.EncodeToString((np[:]))

	return nil
}

func (u *Withdrawal) UnmarshalJSON(data []byte) error {
	type newU struct {
		Id     string  `json:"order,omitempty"`
		Amount float64 `json:"sum,omitempty"`
	}
	nu := newU{}

	if err := json.Unmarshal(data, &nu); err != nil {
		return err
	}

	s, err := strconv.Atoi(nu.Id)
	if err != nil {
		return fmt.Errorf("order amount is not valid")
	}

	if !luhn.Valid(int(s)) {
		return fmt.Errorf("order id is not valid")
	}

	u.Id = int64(s)
	u.Amount = int64(nu.Amount * 100)

	return nil
}

func (b *Withdrawal) MarshalJSON() ([]byte, error) {
	type newWithdrawal struct {
		Id     string  `json:"order,omitempty"`
		Amount float64 `json:"sum,omitempty"`
		Date   string  `json:"processed_at,omitempty"`
	}

	nb := newWithdrawal{
		Amount: math.Abs(float64(b.Amount)) / 100,
		Date:   b.Date.Format(time.RFC3339),
		Id:     fmt.Sprint(b.Id),
	}

	return json.Marshal(nb)
}

func (b *Order) MarshalJSON() ([]byte, error) {
	type newOrder struct {
		Id     string  `json:"number,omitempty"`
		Status string  `json:"status"`
		Amount float64 `json:"accrual"`
		Date   string  `json:"uploaded_at,omitempty"`
	}

	nb := newOrder{
		Amount: math.Abs(float64(b.Amount)) / 100,
		Date:   b.Date.Format(time.RFC3339),
		Id:     fmt.Sprint(b.Id),
		Status: b.Status,
	}

	return json.Marshal(nb)
}

func (u *Order) UnmarshalJSON(data []byte) error {
	type newU struct {
		Id     string    `json:"number,omitempty"`
		Status string    `json:"status"`
		Amount float64   `json:"accrual"`
		Date   time.Time `json:"uploaded_at,omitempty"`
	}
	nu := newU{}

	if err := json.Unmarshal(data, &nu); err != nil {
		return err
	}

	s, err := strconv.Atoi(nu.Id)
	if err != nil {
		return fmt.Errorf("order id is not valid")
	}

	if !luhn.Valid(int(s)) {
		return fmt.Errorf("order id is not valid")
	}

	u.Id = int64(s)
	u.Amount = int64(nu.Amount * 100)
	u.Date = nu.Date
	u.Status = nu.Status

	return nil
}
