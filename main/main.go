package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/fatih/color"
)

func main() {
	org := "ContinuumLLC"
	usernames := getUsersFromOrg(org)
	fmt.Println("1. Incorrect Username ; 2. Empty Fullname")
	var choice string
	fmt.Scan(&choice)
	switch choice {
	case "1":
		checkUsername(usernames)
	case "2":
		checkEmptyFullName(usernames)
	}
}

type user struct {
	Username string `json:"login"`
}

type users struct {
	User []user `json:"Usernames"`
}

type name struct {
	CName string `json:"name"`
}

func getUsersFromOrg(organisation string) map[int]users {
	client := &http.Client{}
	result := make(map[int]users)
	for i := 1; i <= 4; i++ {
		user := users{}
		url := fmt.Sprintf("https://api.github.com/orgs/"+organisation+"/members?page=%d&per_page=362", i)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Fatalf("getUsersFromOrg : Error in fetching Username. Reason : %v", err)
		}
		request.Header.Set("Authorization", "token 3477e10a3b8d35261184a97a8c848dfa2fa5a02d")
		resp, err := client.Do(request)
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			log.Fatalf("getUsersFromOrg : Error in Reading Response Body. Reason : %v", err)
		}
		err = json.Unmarshal(body, &user.User)
		if err != nil {
			log.Fatalf("getUsersFromOrg : Error in Unmarsalling. Reason : %v", err)
		}
		result[i] = user
	}
	return result
}

func checkUsername(usernames map[int]users) {
	fmt.Println("********** White : Contains Continuum ; Red : Does not contain Continuum in the Username ********")
	for i := 1; i <= len(usernames); i++ {
		for j := 0; j < len(usernames[i].User); j++ {
			if strings.Contains(strings.ToLower(usernames[i].User[j].Username), "continuum") {
				color.White(usernames[i].User[j].Username)
			} else {
				color.Red(usernames[i].User[j].Username)
			}
		}
	}
}

func checkEmptyFullName(usernames map[int]users) {
	fmt.Println("********** White : Updated Full Name ; Red : Requires Full Name ********")
	for i := 0; i < len(usernames); i++ {
		for j := 0; j < len(usernames[i].User); j++ {
			client := &http.Client{}
			names := name{}
			url := fmt.Sprintf("https://api.github.com/users/%s", usernames[i].User[j].Username)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				log.Fatalf("checkEmptyFullName : Error in fetching user's name. Reason : %v", err)
			}
			request.Header.Set("Authorization", "token 3477e10a3b8d35261184a97a8c848dfa2fa5a02d")
			resp, err := client.Do(request)
			if err != nil {
				log.Fatalf("checkEmptyFullName : Error in request. Reason : %v", err)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("checkEmptyFullName : Error in Reading Body. Reason : %v", err)
			}
			err = json.Unmarshal(body, &names)
			if err != nil {
				log.Fatalf("checkEmptyFullName : Error in Unmarshalling.Reason:%v", err)
			}
			if names.CName == "" {
				color.Red(usernames[i].User[j].Username)
			} else {
				color.White(usernames[i].User[j].Username)
			}
		}
	}
}