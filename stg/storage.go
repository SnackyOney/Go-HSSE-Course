package stg

import "bank-service/exceptions"

type StorageInterface interface {
	Get(userID int) (int, error)

	Deposit(userID int, amount int) (int, error)

	Transfer(from, to int, amount int) (int, int, error)
}

type BalanceService struct {
	storage StorageInterface
}

func NewBalanceService(new_storage StorageInterface) *BalanceService {
	return &BalanceService{storage: new_storage}
}

func (service *BalanceService) GetBalance(userid int) (int, error) {
	return service.storage.Get(userid)
}

func (service *BalanceService) DepositBalance(userid int, amount int) (int, error) {
	if amount < 0 {
		return 0, exceptions.ErrInvalidAmount
	}
	return service.storage.Deposit(userid, amount)
}

func (service *BalanceService) TransferMoney(from, to int, amount int) (int, int, error) {
	if amount < 0 {
		return 0, 0, exceptions.ErrInvalidAmount
	}
	if from == to {
		return 0, 0, exceptions.ErrSelfTransfer
	}
	fromBalance, toBalance, err := service.storage.Transfer(from, to, amount)
	return fromBalance, toBalance, err
}
