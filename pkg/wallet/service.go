package wallet

import (
	"math"
	"sync"

	//	"path/filepath"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	//	"fmt"
	"errors"

	"github.com/Iftikhor99/wallet/v1/pkg/types"
	"github.com/google/uuid"
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

//ErrFileNotFound for
var ErrFileNotFound = errors.New("file not found")

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
func (s *Service) ExportToFile(path string) error {

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
func (s *Service) ImportFromFile(path string) error {

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
		account := &types.Account{}
		newData2 := strings.Split(stroka, ";")
		for ind, stroka2 := range newData2 {
			// if stroka2 == "" {
			// 	return ErrPhoneRegistered
			// }
			//log.Print(stroka2)
			if ind == 0 {
				id, _ := strconv.ParseInt(stroka2, 10, 64)
				account.ID = id
			}
			if ind == 1 {
				account.Phone = types.Phone(stroka2)
			}
			if ind == 2 {
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

//Export for
func (s *Service) Export(dir string) error {

	// dir, err := filepath.Abs(dir)

	// if err != nil {
	// 	log.Print(err)

	// }
	var err = errors.New("Error")
	err = nil
	lenAcou := len(s.accounts)

	if lenAcou != 0 {

		dirAccount := dir + "/accounts.dump"
		log.Print(dirAccount)

		fileAccounts, err := os.Create(dirAccount)
		if err != nil {
			log.Print(err)

		}

		defer func() {
			if cerr := fileAccounts.Close(); err != nil {
				log.Print(cerr)
			}
		}()

		for index, account := range s.accounts {
			//	account, err = s.FindAccountByID(int64(ind))
			// fmt.Println(newP2)
			// fmt.Println(ee3)
			if index != 0 {
				_, err = fileAccounts.Write([]byte("\n"))
				if err != nil {
					log.Print(err)

				}

			}
			_, err = fileAccounts.Write([]byte(strconv.FormatInt((account.ID), 10)))
			if err != nil {
				log.Print(err)

			}

			_, err = fileAccounts.Write([]byte(";"))
			if err != nil {
				log.Print(err)

			}
			_, err = fileAccounts.Write([]byte(string(account.Phone)))
			if err != nil {
				log.Print(err)

			}

			_, err = fileAccounts.Write([]byte(";"))
			if err != nil {
				log.Print(err)

			}

			_, err = fileAccounts.Write([]byte(strconv.FormatInt(int64(account.Balance), 10)))
			if err != nil {
				log.Print(err)

			}

		}
	}

	lenPay := len(s.payments)
	if lenPay != 0 {

		dirPayment := dir + "/payments.dump"
		filePayments, err := os.Create(dirPayment)
		if err != nil {
			log.Print(err)

		}

		defer func() {
			if cerr := filePayments.Close(); err != nil {
				log.Print(cerr)
			}
		}()

		for index, payment := range s.payments {
			//	account, err = s.FindAccountByID(int64(ind))
			// fmt.Println(newP2)

			if index != 0 {
				_, err = filePayments.Write([]byte("\n"))
				if err != nil {
					log.Print(err)

				}

			}

			_, err = filePayments.Write([]byte(string(payment.ID)))
			if err != nil {
				log.Print(err)

			}

			_, err = filePayments.Write([]byte(";"))
			if err != nil {
				log.Print(err)

			}

			_, err = filePayments.Write([]byte(strconv.FormatInt(int64(payment.AccountID), 10)))
			if err != nil {
				log.Print(err)

			}

			_, err = filePayments.Write([]byte(";"))
			if err != nil {
				log.Print(err)

			}

			_, err = filePayments.Write([]byte(strconv.FormatInt(int64(payment.Amount), 10)))
			if err != nil {
				log.Print(err)

			}

			_, err = filePayments.Write([]byte(";"))
			if err != nil {
				log.Print(err)

			}

			_, err = filePayments.Write([]byte(string(payment.Category)))
			if err != nil {
				log.Print(err)

			}

			_, err = filePayments.Write([]byte(";"))
			if err != nil {
				log.Print(err)

			}

			_, err = filePayments.Write([]byte(string(payment.Status)))
			if err != nil {
				log.Print(err)

			}

		}
	}

	lenFav := len(s.favorites)

	if lenFav != 0 {
		dirFavorite := dir + "/favorites.dump"
		fileFavorites, err := os.Create(dirFavorite)
		if err != nil {
			log.Print(err)

		}

		defer func() {
			if cerr := fileFavorites.Close(); err != nil {
				log.Print(cerr)
			}
		}()

		for index, favorite := range s.favorites {
			//	account, err = s.FindAccountByID(int64(ind))
			// fmt.Println(newP2)

			if index != 0 {
				_, err = fileFavorites.Write([]byte("\n"))
				if err != nil {
					log.Print(err)

				}

			}

			_, err = fileFavorites.Write([]byte(string(favorite.ID)))
			if err != nil {
				log.Print(err)

			}

			_, err = fileFavorites.Write([]byte(";"))
			if err != nil {
				log.Print(err)

			}

			_, err = fileFavorites.Write([]byte(strconv.FormatInt(int64(favorite.AccountID), 10)))
			if err != nil {
				log.Print(err)

			}

			_, err = fileFavorites.Write([]byte(";"))
			if err != nil {
				log.Print(err)

			}

			_, err = fileFavorites.Write([]byte(string(favorite.Name)))
			if err != nil {
				log.Print(err)

			}

			_, err = fileFavorites.Write([]byte(";"))
			if err != nil {
				log.Print(err)

			}

			_, err = fileFavorites.Write([]byte(strconv.FormatInt(int64(favorite.Amount), 10)))
			if err != nil {
				log.Print(err)

			}

			_, err = fileFavorites.Write([]byte(";"))
			if err != nil {
				log.Print(err)

			}

			_, err = fileFavorites.Write([]byte(string(favorite.Category)))
			if err != nil {
				log.Print(err)

			}

		}
	}

	return err

}

//Import for
func (s *Service) Import(dir string) error {

	dirAccount := dir + "/accounts.dump"
	fileAccount, err := os.Open(dirAccount)
	log.Print(dirAccount)
	if err != nil {
		log.Print(err)
		err = ErrFileNotFound
	}
	if err != ErrFileNotFound {
		defer func() {
			err := fileAccount.Close()
			if err != nil {
				log.Print(err)
			}
		}()

		log.Printf("%#v", fileAccount)

		content := make([]byte, 0)
		buf := make([]byte, 4)
		for {
			read, err := fileAccount.Read(buf)
			if err == io.EOF {
				break
			}
			content = append(content, buf[:read]...)
		}

		data := string(content)
		newData := strings.Split(data, "\n")
		//log.Print(data)
		//log.Print(newData)

		for ind1, stroka := range newData {
			//log.Print(stroka)
			account := &types.Account{}
			newData2 := strings.Split(stroka, ";")
			for ind, stroka2 := range newData2 {
				// if stroka2 == "" {
				// 	return ErrPhoneRegistered
				// }
				//log.Print(stroka2)
				if ind == 0 {
					id, _ := strconv.ParseInt(stroka2, 10, 64)
					account.ID = id
				}
				if ind == 1 {
					account.Phone = types.Phone(stroka2)
				}
				if ind == 2 {
					balance, _ := strconv.ParseInt(stroka2, 10, 64)
					account.Balance = types.Money(balance)

				}

				log.Print(ind1)

			}
			errExist := 1
			for _, accountCheck := range s.accounts {

				if accountCheck.ID == account.ID {
					accountCheck.Phone = account.Phone
					accountCheck.Balance = account.Balance
					errExist = 0
				}

			}
			if errExist == 1 {
				s.accounts = append(s.accounts, account)
			}
		}
		for _, account := range s.accounts {
			//	if account.Phone == phone {
			log.Print(account)
			//	}
		}
	}

	dirPayment := dir + "/payments.dump"
	filePayments, err := os.Open(dirPayment)
	if err != nil {
		log.Print(err)
		return ErrFileNotFound
	}
	if err != ErrFileNotFound {
		defer func() {
			err := filePayments.Close()
			if err != nil {
				log.Print(err)
			}
		}()

		log.Printf("%#v", filePayments)

		contentPayment := make([]byte, 0)
		bufPayment := make([]byte, 4)
		for {
			read, err := filePayments.Read(bufPayment)
			if err == io.EOF {
				break
			}
			contentPayment = append(contentPayment, bufPayment[:read]...)
		}

		dataPayment := string(contentPayment)
		newDataPayment := strings.Split(dataPayment, "\n")
		//log.Print(data)
		//log.Print(newData)

		for ind1, stroka := range newDataPayment {
			//log.Print(stroka)
			payment := &types.Payment{}
			newData2 := strings.Split(stroka, ";")
			for ind, stroka2 := range newData2 {
				// if stroka2 == "" {
				// 	return ErrPhoneRegistered
				// }
				//log.Print(stroka2)
				if ind == 0 {
					//id, _ := stroka2
					payment.ID = stroka2
				}
				if ind == 1 {
					accountID, _ := strconv.ParseInt(stroka2, 10, 64)
					payment.AccountID = int64(accountID)
				}

				if ind == 2 {
					balance, _ := strconv.ParseInt(stroka2, 10, 64)
					payment.Amount = types.Money(balance)
				}

				if ind == 3 {
					payment.Category = types.PaymentCategory(stroka2)
				}

				if ind == 4 {
					payment.Status = types.PaymentStatus(stroka2)
				}

				log.Print(ind1)

			}
			errExist := 1
			for _, paymentCheck := range s.payments {

				if paymentCheck.ID == payment.ID {
					paymentCheck.AccountID = payment.AccountID
					paymentCheck.Amount = payment.Amount
					paymentCheck.Category = payment.Category
					paymentCheck.Status = payment.Status
					errExist = 0
				}

			}
			if errExist == 1 {
				s.payments = append(s.payments, payment)
			}
		}
		for _, payment := range s.payments {
			//	if account.Phone == phone {
			log.Print(payment)
			//	}
		}
	}

	dirFavorite := dir + "/favorites.dump"
	fileFavorites, err := os.Open(dirFavorite)
	if err != nil {
		log.Print(err)
		err = ErrFileNotFound
	}
	if err != ErrFileNotFound {
		defer func() {
			err := fileFavorites.Close()
			if err != nil {
				log.Print(err)
			}
		}()

		log.Printf("%#v", fileFavorites)

		contentFavorite := make([]byte, 0)
		bufFavorite := make([]byte, 4)
		for {
			read, err := fileFavorites.Read(bufFavorite)
			if err == io.EOF {
				break
			}
			contentFavorite = append(contentFavorite, bufFavorite[:read]...)
		}

		dataFavorite := string(contentFavorite)
		newDataFavorite := strings.Split(dataFavorite, "\n")
		//log.Print(data)
		//log.Print(newData)

		for ind1, stroka := range newDataFavorite {
			//log.Print(stroka)
			favorite := &types.Favorite{}
			newData2 := strings.Split(stroka, ";")
			for ind, stroka2 := range newData2 {
				// if stroka2 == "" {
				// 	return ErrPhoneRegistered
				// }
				//log.Print(stroka2)
				if ind == 0 {
					//id, _ := stroka2
					favorite.ID = stroka2
				}
				if ind == 1 {
					accountID, _ := strconv.ParseInt(stroka2, 10, 64)
					favorite.AccountID = int64(accountID)
				}

				if ind == 2 {
					favorite.Name = stroka2
				}
				if ind == 3 {
					balance, _ := strconv.ParseInt(stroka2, 10, 64)
					favorite.Amount = types.Money(balance)
				}

				if ind == 4 {
					favorite.Category = types.PaymentCategory(stroka2)
				}

				log.Print(ind1)

			}
			errExist := 1
			for _, favoriteCheck := range s.favorites {

				if favoriteCheck.ID == favorite.ID {
					favoriteCheck.AccountID = favorite.AccountID
					favoriteCheck.Name = favorite.Name
					favoriteCheck.Amount = favorite.Amount
					favoriteCheck.Category = favorite.Category
					errExist = 0
				}

			}
			if errExist == 1 {
				s.favorites = append(s.favorites, favorite)
			}
		}
		for _, favorite := range s.favorites {
			//	if account.Phone == phone {
			log.Print(favorite)
			//	}
		}
	}
	return nil

}

//HistoryToFiles for
func (s *Service) HistoryToFiles(payments []types.Payment, dir string, record1 int) error {

	// dir, err := filepath.Abs(dir)

	// if err != nil {
	// 	log.Print(err)

	// }
	var err = errors.New("Error")
	err = nil
	record := float64(record1)

	lenPay := len(payments)
	if lenPay != 0 {
		//	for i := 1; i < 3; i++ {
		//		str := strconv.FormatInt(int64(i), 10)
		dirPayment := ""
		if lenPay <= record1 {
			dirPayment = dir + "/payments.dump"
		} else {
			dirPayment = dir + "/payments1.dump"
		}
		filePayments, err := os.Create(dirPayment)
		if err != nil {
			log.Print(err)

		}

		defer func() {
			if cerr := filePayments.Close(); err != nil {
				log.Print(cerr)
			}
		}()
		fileNumber1 := 1
		for index, payment := range payments {
			//	account, err = s.FindAccountByID(int64(ind))
			// fmt.Println(newP2)

			fileNumber := int(math.Ceil(float64(index+1) / record))
			//log.Print(fileNumber)
			log.Printf("fileNumber %v", fileNumber)
			if fileNumber > fileNumber1 {
				log.Printf("fileNumber1 %v", fileNumber1)
				str := strconv.FormatInt(int64(fileNumber), 10)
				dirPayment = dir + "/payments" + str + ".dump"
				filePayments, err = os.Create(dirPayment)
				if err != nil {
					log.Print(err)

				}

				defer func() {
					if cerr := filePayments.Close(); err != nil {
						log.Print(cerr)
					}
				}()
			}

			_, err = filePayments.Write([]byte(string(payment.ID)))
			if err != nil {
				log.Print(err)

			}

			_, err = filePayments.Write([]byte(";"))
			if err != nil {
				log.Print(err)

			}

			_, err = filePayments.Write([]byte(strconv.FormatInt(int64(payment.AccountID), 10)))
			if err != nil {
				log.Print(err)

			}

			_, err = filePayments.Write([]byte(";"))
			if err != nil {
				log.Print(err)

			}

			_, err = filePayments.Write([]byte(strconv.FormatInt(int64(payment.Amount), 10)))
			if err != nil {
				log.Print(err)

			}

			_, err = filePayments.Write([]byte(";"))
			if err != nil {
				log.Print(err)

			}

			_, err = filePayments.Write([]byte(string(payment.Category)))
			if err != nil {
				log.Print(err)

			}

			_, err = filePayments.Write([]byte(";"))
			if err != nil {
				log.Print(err)

			}

			_, err = filePayments.Write([]byte(string(payment.Status)))
			if err != nil {
				log.Print(err)

			}

			if (fileNumber >= fileNumber1) || (fileNumber1 == 1) {
				_, err = filePayments.Write([]byte("\n"))
				if err != nil {
					log.Print(err)

				}

			}

			fileNumber1 = fileNumber
			//		}
		}
	}

	for _, payment := range payments {
		//	if account.Phone == phone {
		log.Print(payment)
		//	}
	}

	return err

}

//ExportAccountHistory for
func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {
	var foundPayments []types.Payment

	for _, payment := range s.payments {
		log.Print(accountID)
		log.Print(payment.AccountID)
		if payment.AccountID == accountID {
			foundPayments = append(foundPayments, *payment)
			//	return foundPayments, nil
		}
	}
	if foundPayments == nil {
		return nil, ErrAccountNotFound
	}
	return foundPayments, nil
}

// //Simple for
// func (s *Service) Simple() types.Money {
// 	lenPay := len(s.payments)
// 	allPayments := s.payments
// 	log.Print(s.payments)
// 	log.Print(lenPay)
// 	log.Print(allPayments[0].Amount)
// 	return allPayments[0].Amount
// }

//SumPayments for
func (s *Service) SumPayments(goroutines int) types.Money {
	// err := types.Money(5)
	// return err
	wg := sync.WaitGroup{}

	wg.Add(goroutines) // cKonbKOo ropyTMH péM

	mu := sync.Mutex{}
	sum := types.Money(0)
	lenPay := len(s.payments)
	numberOfPaymentPerRoutine := lenPay / goroutines
	//timesOfPayments := 1
	allPayments := s.payments
	index := 0
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done() // cooOwaem, 4TO 3aKkoHUunN
			val := types.Money(0)

			for ; index < numberOfPaymentPerRoutine; index++ {
				if index < lenPay {
					val += allPayments[index].Amount
				}
			}
			mu.Lock()
			defer mu.Unlock()
			sum += val // TOMbKO B KOHUE 3anvcbiBaeM CYMMY
			numberOfPaymentPerRoutine += numberOfPaymentPerRoutine
		}()
	}
	// go func() {
	// 	defer wg.Done() // coo6waem, 4TO 3aKoH4UNN
	// 	val := types.Money(0)
	// 	for index, payment := range s.payments {
	// 		if index > numberOfPaymentPerRoutine{
	// 			break
	// 		}
	// 		val += payment.Amount
	// 	}
	// 	mu.Lock()
	// 	defer mu.Unlock()
	// 	sum += val // TOMbKO B KOHWE 3anucbiBaeM CyMMy

	// }()

	wg.Wait()
	return sum
}

//FilterPayments for
func (s *Service) FilterPayments(accountID int64, goroutines int) ([]types.Payment, error) {

	var foundPayments []types.Payment
	//	var newPayments []types.Payment
	var allfoundPayments []types.Payment
	for _, payment := range s.payments {
		if payment.AccountID == accountID {
			foundPayments = append(foundPayments, *payment)
		}
	}
	if foundPayments == nil {
		return nil, ErrAccountNotFound
	}
	wg := sync.WaitGroup{}

	wg.Add(goroutines) // cKonbKOo ropyTMH péM

	mu := sync.Mutex{}

	lenPay := len(foundPayments)
	numberOfPaymentPerRoutine := int(math.Ceil(float64((lenPay + 1) / goroutines)))
	//numberOfPaymentPerRoutine := lenPay / goroutines
	//timesOfPayments := 1
	allPayments := foundPayments
	index := 0
	for i := 0; i < goroutines; i++ {
		
		go func() {
			defer wg.Done() // cooOwaem, 4TO 3aKkoHUunN
			var newPayments []types.Payment
			for ; index < numberOfPaymentPerRoutine; index++ {
				if index < lenPay {
					newPayments = append(newPayments, allPayments[index])

				}
			}
			numberOfPaymentPerRoutine += numberOfPaymentPerRoutine
			mu.Lock()
			defer mu.Unlock()
			allfoundPayments = append(allfoundPayments, newPayments...)

		}()
	}
	// go func() {
	// 	defer wg.Done() // coo6waem, 4TO 3aKoH4UNN
	// 	val := types.Money(0)
	// 	for index, payment := range s.payments {
	// 		if index > numberOfPaymentPerRoutine{
	// 			break
	// 		}
	// 		val += payment.Amount
	// 	}
	// 	mu.Lock()
	// 	defer mu.Unlock()
	// 	sum += val // TOMbKO B KOHWE 3anucbiBaeM CyMMy

	// }()

	wg.Wait()
	return allfoundPayments, nil
}

// //FilterPaymentsNormal for
// func (s *Service) FilterPaymentsNormal(accountID int64, goroutines int) ([]types.Payment, error) {

// 	var foundPayments []types.Payment

// 	var allfoundPayments []types.Payment
// 	for _, payment := range s.payments {
// 		if payment.AccountID == accountID {
// 			foundPayments = append(foundPayments, *payment)
// 		}
// 	}
// 	if foundPayments == nil {
// 		return nil, ErrAccountNotFound
// 	}

// 	lenPay := len(foundPayments)
// 	numberOfPaymentPerRoutine := int(math.Ceil(float64((lenPay + 1) / goroutines)))
// 	//numberOfPaymentPerRoutine := lenPay / goroutines
// 	//timesOfPayments := 1
// 	allPayments := foundPayments
// 	index := 0
// 	for i := 0; i < goroutines; i++ {
// 		newPayments := []types.Payment{}

// 		for ; index < numberOfPaymentPerRoutine; index++ {
// 			if index < lenPay {
// 				newPayments = append(newPayments, allPayments[index])

// 			}
// 		}
// 		numberOfPaymentPerRoutine += numberOfPaymentPerRoutine

// 		allfoundPayments = append(allfoundPayments, newPayments...)

// 	}

// 	return allfoundPayments, nil
// }
