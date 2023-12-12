package repositories

import "gorm.io/gorm"

// Entity
type BankAccount struct {
	ID            string
	AccountHolder string
	AccountType   int
	Balance       float64
}

// interface
type AccountRepository interface {
	Save(bankAccount BankAccount) error
	Delete(id string) error
	FindAll() (bankAccount []BankAccount, err error)
	FindById(id string) (bankAccount BankAccount, err error)
}

type accountRepository struct {
	db *gorm.DB
}

// This constructor is used to create a new instance of accountRepository which implements AccountRepository interface
func NewAccountRepository(db *gorm.DB) AccountRepository {
	db.Table("Bank").AutoMigrate(&BankAccount{})
	return accountRepository{db}
}

// below methods are the implementation of AccountRepository interface
func (obj accountRepository) Save(bankAccount BankAccount) error {
	return obj.db.Table("Bank").Save(bankAccount).Error
}
func (obj accountRepository) Delete(id string) error {
	return obj.db.Table("Bank").Where("id = ?", id).Delete(&BankAccount{}).Error
}

func (obj accountRepository) FindAll() (bankAccount []BankAccount, err error) {
	err = obj.db.Table("Bank").Find(&bankAccount).Error
	return bankAccount, err
}

func (obj accountRepository) FindById(id string) (bankAccount BankAccount, err error) {
	err = obj.db.Table("Bank").Where("id = ?", id).First(&bankAccount).Error
	return bankAccount, err
}
