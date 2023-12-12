package services

import (
	"consumer/repositories"
	"encoding/json"
	"events"
	"log"
	"reflect"
)

type EventHandler interface {
	Handle(topic string, eventByte []byte)
}

type accountEventHandler struct {
	accountRepo repositories.AccountRepository
}

func NewEventHandler(accountRepo repositories.AccountRepository) EventHandler {
	return accountEventHandler{accountRepo}
}

func (obj accountEventHandler) Handle(topic string, eventByte []byte) {
	switch topic {
	case reflect.TypeOf(events.OpenAccountEvent{}).Name():
		event := &events.OpenAccountEvent{}
		err := json.Unmarshal(eventByte, event)
		if err != nil {
			log.Println(err)
			return
		}

		//create a new bank account in the form of repositories.BankAccount
		backAccount := repositories.BankAccount{
			ID:            event.ID,
			AccountHolder: event.AccountHolder,
			AccountType:   event.AccountType,
			Balance:       event.OpeningBalance,
		}

		//save the bank account to the database
		err = obj.accountRepo.Save(backAccount)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("topic: %v, %#v", topic, event)
	case reflect.TypeOf(events.DepositFundEvent{}).Name():
		event := &events.DepositFundEvent{}
		err := json.Unmarshal(eventByte, event)
		if err != nil {
			log.Println(err)
			return

		}

		bankAccount, err := obj.accountRepo.FindById(event.ID)
		if err != nil {
			log.Println(err)
			return
		}

		bankAccount.Balance += event.Amount

		//update the bank account to the database
		err = obj.accountRepo.Save(bankAccount)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("topic: %v, %#v", topic, event)

	case reflect.TypeOf(events.WithdrawFundEvent{}).Name():
		// create a new instance of WithdrawFundEvent
		event := &events.WithdrawFundEvent{}
		// unmarshal the eventByte to event
		err := json.Unmarshal(eventByte, event)
		if err != nil {
			log.Println(err)
			return
		}

		//find the bank account by id
		bankAccount, err := obj.accountRepo.FindById(event.ID)
		if err != nil {
			log.Println(err)
			return
		}

		bankAccount.Balance -= event.Amount

		//update the bank account to the database
		err = obj.accountRepo.Save(bankAccount)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("topic: %v, %#v", topic, event)

	case reflect.TypeOf(events.CloseAccountEvent{}).Name():
		// create a new instance of CloseAccountEvent
		event := &events.CloseAccountEvent{}

		// unmarshal the eventByte to event
		err := json.Unmarshal(eventByte, event)
		if err != nil {
			log.Println(err)
			return
		}

		//delete the bank account from the database
		err = obj.accountRepo.Delete(event.ID)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("topic: %v, %#v", topic, event)
	default:
		log.Println("No handler for topic", topic)
	}
}
