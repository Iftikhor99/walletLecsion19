package main

import (
//	"io/ioutil"
	"path/filepath"
//	"github.com/Iftikhor99/wallet/v1/pkg/types"
//	"strings"
//	"strconv"
//	"io"
	"log"
	"os"
	"fmt"
	"github.com/Iftikhor99/wallet/v1/pkg/wallet"
)


func main() {
	svc := &wallet.Service{}
	accountTest , err := svc.RegisterAccount("+992000000001")
	if err != nil {
		fmt.Println(err)
		return
	} 

	err = svc.Deposit(accountTest.ID, 100_000_00)
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

	
	// accountTest , err = svc.RegisterAccount("+992000000002")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// } 

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


	newP, ee2 := svc.Pay(accountTest.ID,10_000_00,"food")

	
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

	fav, errFv := svc.FavoritePayment(newP.ID, "Tcell")
	fmt.Println(errFv)
	fmt.Println(fav)

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

	err = svc.Import("c:/projects/wallet/data/")
	if err != nil {
	 	log.Print(err)
	 	return
	}
	


	

}