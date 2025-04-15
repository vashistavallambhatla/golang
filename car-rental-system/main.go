package main

import (
	"fmt"
	"math/rand"
	"time"

	"crs/services"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	rentalSystem := services.NewCarRentalSystem()

	sedan, err := rentalSystem.EnrollCar("Toyota", "Camry", 2022, "ABC123", 50.0, "Sedan")
	if err != nil {
		fmt.Printf("Error enrolling car: %v\n", err)
		return
	}
	fmt.Printf("Enrolled car: %s %s (ID: %d)\n", sedan.Make, sedan.Model, sedan.ID)

	suv, err := rentalSystem.EnrollCar("Honda", "CR-V", 2023, "XYZ789", 75.0, "SUV")
	if err != nil {
		fmt.Printf("Error enrolling car: %v\n", err)
		return
	}
	fmt.Printf("Enrolled car: %s %s (ID: %d)\n", suv.Make, suv.Model, suv.ID)

	
	customer, err := rentalSystem.RegisterCustomer("John Doe", "john@example.com", "DL12345")
	if err != nil {
		fmt.Printf("Error registering customer: %v\n", err)
		return
	}
	fmt.Printf("Registered customer: %s (ID: %d)\n", customer.Name, customer.ID)

	startDate := time.Now().AddDate(0, 0, 1)
	endDate := time.Now().AddDate(0, 0, 4)

	reservation, err := rentalSystem.MakeReservation(sedan.ID, customer.ID, startDate, endDate)
	if err != nil {
		fmt.Printf("Error making reservation: %v\n", err)
		return
	}
	fmt.Printf("Created reservation (ID: %d) for %s from %s to %s\n",
		reservation.ID,
		sedan.Make+" "+sedan.Model,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))
	fmt.Printf("Total cost: $%.2f\n", reservation.TotalCost)

	availableCars, err := rentalSystem.FindAvailableCarsByFilters("SUV", 100.0, startDate, endDate)
	if err != nil {
		fmt.Printf("Error searching for cars: %v\n", err)
	} else {
		fmt.Printf("\nAvailable cars matching criteria:\n")
		for _, car := range availableCars {
			fmt.Printf("- %s %s (%s): $%.2f per day\n", car.Make, car.Model, car.CarType, car.PricePerDay)
		}
	}


	newStartDate := time.Now().AddDate(0, 0, 5)
	newEndDate := time.Now().AddDate(0, 0, 8)

	modifiedReservation, err := rentalSystem.ModifyReservation(reservation.ID, newStartDate, newEndDate)
	if err != nil {
		fmt.Printf("\nError modifying reservation: %v\n", err)
	} else {
		fmt.Printf("\nReservation modified successfully:\n")
		fmt.Printf("New dates: %s to %s\n",
			modifiedReservation.StartDate.Format("2006-01-02"),
			modifiedReservation.EndDate.Format("2006-01-02"))
	}

	cancelResult, err := rentalSystem.CancelReservation(reservation.ID)
	if err != nil {
		fmt.Printf("\nError cancelling reservation: %v\n", err)
	} else {
		fmt.Printf("\n%s\n", cancelResult)
	}

	invalidStartDate := time.Now().AddDate(0, 0, 10)
	invalidEndDate := time.Now().AddDate(0, 0, 5)

	_, err = rentalSystem.MakeReservation(suv.ID, customer.ID, invalidStartDate, invalidEndDate)
	fmt.Printf("\nExpected error with invalid dates: %v\n", err)

	
	cars, err := rentalSystem.FindAvailableCarsByFilters("", 0, time.Time{}, time.Time{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("\nAll available cars:")
		for _, car := range cars {
			fmt.Printf("- %s %s (%s): $%.2f per day\n", car.Make, car.Model, car.CarType, car.PricePerDay)
		}
	}
}