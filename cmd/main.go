package main

import (
	"sync"
	//	"io/ioutil"
	"path/filepath"
	//	"github.com/Iftikhor99/wallet/v1/pkg/types"
	//	"strings"
	//	"strconv"
	//	"io"
	"fmt"
	"log"
	"os"

	"github.com/Iftikhor99/wallet/v1/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	accountTest1, err1 := svc.RegisterAccount("+992000000001")
	if err1 != nil {
		fmt.Println(err1)
		return
	}

	err1 = svc.Deposit(accountTest1.ID, 100_000_00)
	if err1 != nil {
		switch err1 {
		case wallet.ErrAmountMustBePositive:
			fmt.Println("Сумма должна быть положительной")
		case wallet.ErrAccountNotFound:
			fmt.Println("Аккаунт пользователя не найден")
		}
		return
	}
	fmt.Println(accountTest1.Balance)

	newP1, ee21 := svc.Pay(accountTest1.ID, 6_000_00, "car")
	fmt.Println(newP1)
	fmt.Println(ee21)

	accountTest, err := svc.RegisterAccount("+992000000002")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = svc.Deposit(accountTest.ID, 200_000_00)
	if err != nil {
		switch err {
		case wallet.ErrAmountMustBePositive:
			fmt.Println("Сумма должна быть положительной")
		case wallet.ErrAccountNotFound:
			fmt.Println("Аккаунт пользователя не найден")
		}
		return
	}
	fmt.Println(accountTest.Balance)

	newP, ee2 := svc.Pay(accountTest.ID, 1_000_00, "food")
	newP, ee2 = svc.Pay(accountTest.ID, 2_000_00, "food")
	newP, ee2 = svc.Pay(accountTest.ID, 3_000_00, "food")
	newP, ee2 = svc.Pay(accountTest.ID, 4_000_00, "food")
	newP, ee2 = svc.Pay(accountTest.ID, 5_000_00, "food")
	//newP, ee2 = svc.Pay(accountTest.ID, 1_000_00, "auto")

	fmt.Println(accountTest.Balance)
	fmt.Println(newP)
	fmt.Println(ee2)

	// newP2, ee3 := svc.FindPaymentByID(newP.ID)
	// fmt.Println(newP2)
	// fmt.Println(ee3)

	// //ee4 := svc.Reject(newP.ID)
	// //fmt.Println(account.Balance)
	// //fmt.Println(ee4)

	// newP3, ee5 := svc.Repeat(newP.ID)
	// fmt.Println(ee5)
	// fmt.Println(newP3.Amount)
	// fmt.Println(account.Balance)

	// fav, errFv := svc.FavoritePayment(newP.ID, "Tcell")
	// fmt.Println(errFv)
	// fmt.Println(fav)

	// newP4, eeFv2 := svc.PayFromFavorite("fav.ID")
	// fmt.Println(eeFv2)
	// fmt.Println(newP4)

	// fmt.Println(account.Balance)

	abs, err := filepath.Abs("data/readme.txt")

	if err != nil {
		log.Print(err)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Print(err)
		return
	}

	log.Print(wd)
	log.Print(abs)

	// err = ioutil.WriteFile("c:/projects/wallet/data/readme1.txt", []byte(accountTest.ID), 0600)
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// err = svc.ImportFromFile("c:/projects/wallet/data/accounts.dump")
	// if err != nil {
	//  	log.Print(err)
	//  	return
	//  }

	// err = svc.Export(wd)
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// err = svc.Import(wd)
	// if err != nil {
	//  	log.Print(err)
	//  	return
	// }
	
	// paymentsFound, err := svc.ExportAccountHistory(newP1.AccountID)
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// err = svc.HistoryToFiles(paymentsFound, wd, 2)
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }
	//svc.Simple()
	pay, err8 := svc.FilterPayments(1,2) 
	log.Print(pay)
	log.Print(err8)

}

//Concurrently for
func Concurrently() int64 {

	wg := sync.WaitGroup{}

	wg.Add(2) // cKonbKOo ropyTMH péM

	mu := sync.Mutex{}
	sum := int64(0)

	go func() {
		defer wg.Done() // cooOwaem, 4TO 3aKkoHUunN
		val := int64(0)
		for i := 0; i < 1000; i++ {
			val++
		}
		mu.Lock()
		defer mu.Unlock()
		sum += val // TOMbKO B KOHUE 3anvcbiBaeM CYMMY
	}()

	go func() {
		defer wg.Done() // coo6waem, 4TO 3aKoH4UNN
		val := int64(0)
		for i := 0; i < 1000; i++ {
			val++
		}
		mu.Lock()
		defer mu.Unlock()
		sum += val // TOMbKO B KOHWE 3anucbiBaeM CyMMy

	}()

	wg.Wait()
	return sum
}
