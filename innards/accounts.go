package innards

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Account struct {
	Id       string            `json:"id"`
	Email    string            `json:"email"`
	Password string            `json:"password"`
	Links    map[string]string `json:"links"`
}

const accountsPath = "data/accounts/"

func findAccountWithId(id string) (Account, error) {
	var fileBytes, err = os.ReadFile(accountsPath + id + "/details.json")
	if err != nil {
		return Account{}, err
	}
	var acc Account
	err = json.Unmarshal(fileBytes, &acc)
	if err != nil {
		return Account{}, err
	}
	return acc, nil
}

func existsAccountWithId(id string) bool {
	// if we cannot stat, then we have an error, if we have no error, then such an account exists
	var _, err = os.Stat(accountsPath + id)
	return err == nil
}

// Should get implemented when we move to database
func existsAccountWithEmail(email string) bool {
	return false
}

func persistAccount(acc Account) error {
	if !existsAccountWithId(acc.Id) {
		var err = os.Mkdir(accountsPath+acc.Id, 0755)
		if err != nil {
			fmt.Println("ERROR: cannot make directory for account:", acc.Id)
			return err
		}
		if acc.Links == nil {
			acc.Links = make(map[string]string)
		}
	}
	var fileBytes, err = json.Marshal(acc)
	if err != nil {
		return err
	}
	err = os.WriteFile(accountsPath+acc.Id+"/details.json", fileBytes, 0755)
	/*
		if err != nil {
			return err
		}
		return nil
	*/
	return err
}

func validateAccount(id, password string) bool {
	var acc, err = findAccountWithId(id)
	if err != nil {
		//fmt.Println("ERROR no such account:", id)
		return false
	}
	//fmt.Println("acc.P: ", acc.Password, " // P: ", password)
	return acc.Password == password
}

func correspondingLink(id, link string) (string, error) {
	var acc, err = findAccountWithId(id)
	if err != nil {
		return "", err
	}
	if corresponding, ok := acc.Links[link]; ok {
		return corresponding, nil
	}
	return "", errors.New("correspondingLink(" + id + ", " + link + "): no such link for user")
}

func putLinkOnAccountAndPersist(alias, url string, acc Account) (string, error) {
	if !strings.Contains(url, "://") {
		url = "https://" + url
	}
	acc.Links[alias] = url
	if err := persistAccount(acc); err != nil {
		return "", err
	}
	return url, nil
}
