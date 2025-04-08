package main
import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)
type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}
type UserResponse struct {
	Data       []User `json:"data"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	Total      int    `json:"total"`
	TotalPages int    `json:"total_pages"`
}
func fetchUserDetails(id int, wg *sync.WaitGroup, sem chan struct{}, userCh chan<- User) {
	defer wg.Done()
	sem <- struct{}{} 
	defer func() { <-sem }()
	url := "https://reqres.in/api/users/" + strconv.Itoa(id)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching user %d: %v", id, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Request for user %d failed with code %d", id, resp.StatusCode)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response for user %d: %v", id, err)
		return
	}
	var userWrapper struct {
		Data User `json:"data"`
	}
	err = json.Unmarshal(body, &userWrapper)
	if err != nil {
		log.Printf("Error unmarshalling response for user %d: %v", id, err)
		return
	}
	userCh <- userWrapper.Data
}
func main() {
	url := "https://reqres.in/api/users?page=1"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching users: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Initial request failed with code %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading initial response body: %v", err)
	}
	var userResp UserResponse
	err = json.Unmarshal(body, &userResp)
	if err != nil {
		log.Fatalf("Error unmarshalling initial response: %v", err)
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)
	userCh := make(chan User)

	var allUsers []User
	go func() { // Fan-out, Fan-in approach? The go func() block runs in the background as a goroutine, and it waits for data to be sent through the userCh channel.
		for user := range userCh {
			allUsers = append(allUsers, user)
		}
	}() // This function will run indefinitely until the channel userCh is closed or the main exits.
	for _, user := range userResp.Data {
		wg.Add(1)
		go fetchUserDetails(user.ID, &wg, sem, userCh)
	}
	wg.Wait()
	close(userCh) 

	fmt.Printf("Fetched details for %d users:\n", len(allUsers))
	for _, user := range allUsers {
		fmt.Printf("ID: %d, Name: %s %s\n", user.ID, user.FirstName, user.LastName)
	}
}

















