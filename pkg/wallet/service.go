package wallet

import (
	"strings"
	"io"
	"strconv"
	"os"
	"log"
//	"fmt"
	"github.com/google/uuid"
	"errors"
	"github.com/Iftikhor99/wallet/v1/pkg/types"
)

	//ErrPhoneRegistered for
	var ErrPhoneRegistered = errors.New("phone already registered")
	//ErrAmountMustBePositive for
	var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
	//ErrAccountNotFound for
	var ErrAccountNotFound = errors.New("account not found")
	//ErrNotEnoughBalance for
	var ErrNotEnoughBalance = errors.New("not enough balance")
	//ErrPaymentNotFound for
	var ErrPaymentNotFound = errors.New("payment not found")
	//ErrFavoriteNotFound for
	var ErrFavoriteNotFound = errors.New("favorite not found")

	//Service for
	type Service struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
		favorites     []*types.Favorite
	}

	//RegisterAccount for
	func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
		for _, account := range s.accounts {
			if account.Phone == phone {
				return nil, ErrPhoneRegistered
			}
		}

		s.nextAccountID++
		account := &types.Account{
			ID:      s.nextAccountID,
			Phone:   phone,
			Balance: 0,
		}
		s.accounts = append(s.accounts, account)

		return account, nil
	}

	//FindAccountByID for
	func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
		for _, account := range s.accounts {
			if account.ID == accountID {
				return account, nil
			}
		}

		return nil, ErrAccountNotFound
	}

	//Deposit for
	func (s *Service) Deposit(accountID int64, amount types.Money) error {
		if amount <= 0 {
			return ErrAmountMustBePositive
		}

		account, err := s.FindAccountByID(accountID)
		if err != nil {
			return ErrAccountNotFound
		}

		// зачисление средств пока не рассматриваем как платёж
		account.Balance += amount
		return nil
	}

	//Pay for
	func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
		if amount <= 0 {
			return nil, ErrAmountMustBePositive
		}

		var account *types.Account
		for _, acc := range s.accounts {
			if acc.ID == accountID {
				account = acc
				break
			}
		}
		if account == nil {
			return nil, ErrAccountNotFound
		}

		if account.Balance < amount {
			return nil, ErrNotEnoughBalance
		}

		account.Balance -= amount
		paymentID := uuid.New().String()
		payment := &types.Payment{
			ID:        paymentID,
			AccountID: accountID,
			Amount:    amount,
			Category:  category,
			Status:    types.PaymentStatusInProgress,
		}
		s.payments = append(s.payments, payment)
		return payment, nil
	}

	//FindPaymentByID for
	func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
		for _, payment := range s.payments {
			if payment.ID == paymentID {
				return payment, nil
			}
		}

		return nil, ErrPaymentNotFound
	}

	//Reject for
	func (s *Service) Reject(paymentID string) error {
		payment, err := s.FindPaymentByID(paymentID)
		if err != nil {
			return err
		}
		account, err := s.FindAccountByID(payment.AccountID)
		if err != nil {
			return err
		}

		payment.Status = types.PaymentStatusFail
		account.Balance += payment.Amount
		return nil
	}

	//Repeat for
	func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
		payment, err := s.FindPaymentByID(paymentID)
		if err != nil {
			return nil, err
		}

		return s.Pay(payment.AccountID, payment.Amount, payment.Category)
	}

	//FavoritePayment for
	func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
		payment, err := s.FindPaymentByID(paymentID)
		if err != nil {
			return nil, err
		}

		favorite := &types.Favorite{
			ID:        uuid.New().String(),
			AccountID: payment.AccountID,
			Amount:    payment.Amount,
			Name:      name,
			Category:  payment.Category,
		}

		//s.favorites[len(s.favorites)] = favorite
		s.favorites = append(s.favorites, favorite)
		return favorite, nil
	}

	//FindFavoriteByID for
	func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
		for _, favorite := range s.favorites {
			if favorite.ID == favoriteID {
				return favorite, nil
			}
		}

		return nil, ErrFavoriteNotFound
	}

	//PayFromFavorite for
	func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
		favorite, err := s.FindFavoriteByID(favoriteID)
		if err != nil {
			return nil, err
		}

		return s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	}
//ExportToFile for
func (s *Service) ExportToFile(path string) error{

	
	
	fileNew, err := os.Create(path)
	if err != nil {
		log.Print(err)
		
	}

	defer func() {
		
		if cerr := fileNew.Close(); err != nil {
			log.Print(cerr)
		}
	}()
	for index, account := range s.accounts {
	//	account, err = s.FindAccountByID(int64(ind))
	// fmt.Println(newP2)
	// fmt.Println(ee3)
	 if index != 0 {
		_, err = fileNew.Write([]byte("|"))
		if err != nil {
			log.Print(err)
			
		}
		
	}

	_, err = fileNew.Write([]byte(strconv.FormatInt((account.ID), 10)))
	if err != nil {
		log.Print(err)
		
	}

	_, err = fileNew.Write([]byte(";"))
	if err != nil {
		log.Print(err)
		
	}
	_, err = fileNew.Write([]byte(string(account.Phone)))
	if err != nil {
		log.Print(err)
		
	}

	_, err = fileNew.Write([]byte(";"))
	if err != nil {
		log.Print(err)
		
	}

	_, err = fileNew.Write([]byte(strconv.FormatInt(int64(account.Balance), 10)))
	if err != nil {
		log.Print(err)
		
	}


	}

	return err

}

//ImportFromFile for
func (s *Service) ImportFromFile(path string) error{

	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	log.Printf("%#v", file)

	content := make([]byte, 0)
	buf := make([]byte, 4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		content = append(content, buf[:read]...)
	}
	
	data := string(content)
	newData := strings.Split(data, "|")
	//log.Print(data)
	//log.Print(newData)
	
	
	for ind1, stroka := range newData {
		//log.Print(stroka)
		account := &types.Account {
			}	
		newData2 := strings.Split(stroka, ";")
		for ind, stroka2 := range newData2 {
			// if stroka2 == "" {
			// 	return ErrPhoneRegistered
			// }
			//log.Print(stroka2)
			if ind == 0{
				id, _ := strconv.ParseInt(stroka2, 10, 64)
				account.ID = id
			}
			if ind == 1{
				account.Phone = types.Phone(stroka2)
			}
			if ind == 2{
				balance, _ := strconv.ParseInt(stroka2, 10, 64)
				account.Balance = types.Money(balance)
					
			}
					
			// if (ind1 == 0) && (ind ==2) {
				log.Print(ind1)
			// 	s.accounts = append(s.accounts, account)		
			// }  

			// if (ind1 == 1) && (ind ==2) {
			// 	log.Print(account)
			// 	s.accounts = append(s.accounts, account)		
			// } 

		}
		for _, accountCheck := range s.accounts {
			if accountCheck.Phone == account.Phone {
				return ErrPhoneRegistered
			}
			if accountCheck.ID == account.ID {
				return ErrPhoneRegistered
			}
			

		 }
		s.accounts = append(s.accounts, account)
	}
	for _, account := range s.accounts {
	//	if account.Phone == phone {
			log.Print(account)
	//	}
	}
	
	return nil

}