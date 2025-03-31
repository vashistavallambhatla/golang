package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}

type UserResponse struct {
	Data      []User `json:"data"`
	Page      int    `json:"page"`
	PerPage   int    `json:"per_page"`
	Total     int    `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func main() {
	baseURL := "https://reqres.in/api/users"  // Corrected URL from regres.in to reqres.in
	maxUsers := 100
	perPage := 12
	var allUsers []User
	currentPage := 1

	for len(allUsers) < maxUsers {
		url := baseURL + "?page=" + strconv.Itoa(currentPage) + "&per_page=" + strconv.Itoa(perPage)
		fmt.Println("Fetching URL:", url)
		
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Error fetching users: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Response failed with code %d", resp.StatusCode)
		}
		
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error while reading response body: %v", err)
		}
		
		if len(body) > 100 {
			fmt.Printf("Response Body (preview): %s...\n", body[:100])
		} else {
			fmt.Printf("Response Body: %s\n", body)
		}
		
		var userResp UserResponse
		err = json.Unmarshal(body, &userResp)
		if err != nil {
			log.Fatalf("Error occurred while unmarshalling the data: %v", err)
		}
		
		if len(userResp.Data) == 0 {
			break
		}
		
		allUsers = append(allUsers, userResp.Data...)
	
		if currentPage >= userResp.TotalPages {
			break
		}
		
		currentPage++
	}
	
	if len(allUsers) > maxUsers {
		allUsers = allUsers[:maxUsers]
	}
	
	fmt.Printf("Retrieved %d users\n", len(allUsers))
	
	file, err := os.OpenFile("module7/users.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()
	
	usersJSON, err := json.MarshalIndent(allUsers, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling users to JSON: %v", err)
	}
	
	_, err = file.Write(usersJSON)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
	
	fmt.Println("Successfully saved users to users.json")
}