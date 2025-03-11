package main

import (
	"fmt"
	"time"
)

type Car struct {
	ID             int
	Make           string
	Model          string
	Year           int
	License_plate  string
	Price_per_day  float64
	Available      bool
	Bookings       []BookingPeriod
	CarType        string
}

type BookingPeriod struct {
	StartDate time.Time
	EndDate   time.Time
}

type Customer struct {
	ID       int
	Name     string
	Contact  string
	License  string
}

type Reservation struct {
	ID        int
	Customer  int
	CarId     int
	StartDate time.Time
	EndDate   time.Time
	TotalDays int
	TotalCost float64
}

type CarRentalSystem struct {
	cars              map[int]*Car
	customers         map[int]*Customer
	reservations      map[int]*Reservation
	nextCarId         int
	nextCustomerId    int
	nextReservationId int
}

func NewCarRentalSystem() *CarRentalSystem {
	return &CarRentalSystem{
		cars:         make(map[int]*Car),
		customers:    make(map[int]*Customer),
		reservations: make(map[int]*Reservation),
	}
}

func (rental *CarRentalSystem) AddCar(make string, model string, year int, license_plate string, price_per_day float64, carType string) {
	rental.nextCarId++
	rental.cars[rental.nextCarId] = &Car{
		ID:            rental.nextCarId,
		Make:          make,
		Model:         model,
		Year:          year,
		License_plate: license_plate,
		Price_per_day: price_per_day,
		Available:     true,
		Bookings:      []BookingPeriod{},
		CarType:       carType,
	}
}

func (rental *CarRentalSystem) AddCustomer(name, contact, license string) {
	rental.nextCustomerId++
	rental.customers[rental.nextCustomerId] = &Customer{
		ID:      rental.nextCustomerId,
		Name:    name,
		Contact: contact,
		License: license,
	}
}

func (rental *CarRentalSystem) MakeReservation(carId, customerId int, startDate, endDate time.Time) {
	car, exists := rental.cars[carId]
	if !exists {
		fmt.Println("Car not found")
		return
	}

	if !car.isAvailable(startDate, endDate) {
		fmt.Println("Car is not available for the selected dates")
		return
	}

	rental.nextReservationId++
	daysCount := int(endDate.Sub(startDate).Hours() / 24)
	totalCost := car.Price_per_day * float64(daysCount)

	rental.reservations[rental.nextReservationId] = &Reservation{
		ID:        rental.nextReservationId,
		Customer:  customerId,
		CarId:     carId,
		StartDate: startDate,
		EndDate:   endDate,
		TotalDays: daysCount,
		TotalCost: totalCost,
	}

	car.Bookings = append(car.Bookings, BookingPeriod{StartDate: startDate, EndDate: endDate})
}

func (rental *CarRentalSystem) ModifyReservation(reservationId int, startDate, endDate time.Time) {
	reservation, exists := rental.reservations[reservationId]
	if !exists {
		fmt.Println("Reservation not found")
		return
	}

	car := rental.cars[reservation.CarId]

	for i, booking := range car.Bookings {
		if booking.StartDate == reservation.StartDate && booking.EndDate == reservation.EndDate {
			if !car.isAvailable(startDate, endDate) {
				fmt.Println("Car not available on these dates")
				return
			}

			reservation.StartDate = startDate
			reservation.EndDate = endDate

			car.Bookings[i] = BookingPeriod{StartDate: startDate, EndDate: endDate}

			fmt.Printf("Reservation dates successfully changed to %v - %v\n", startDate, endDate)
			return
		}
	}
}

func (rental *CarRentalSystem) CancelReservation(reservationId int) {
	reservation, exists := rental.reservations[reservationId]
	if !exists {
		fmt.Println("Reservation not found")
		return
	}

	car := rental.cars[reservation.CarId]

	for i, booking := range car.Bookings {
		if booking.StartDate == reservation.StartDate && booking.EndDate == reservation.EndDate {
			car.Bookings = append(car.Bookings[:i], car.Bookings[i+1:]...)
			break
		}
	}

	delete(rental.reservations, reservationId)

	fmt.Println("Reservation canceled successfully")
}

func (rental *CarRentalSystem) Search(carType string, price float64, startDate, endDate time.Time) (searchResult []Car) {
	for _, car := range rental.cars {
		if car.CarType == carType && car.Price_per_day <= price && car.isAvailable(startDate, endDate) {
			searchResult = append(searchResult, *car)
		}
	}
	return searchResult
}

func (c *Car) isAvailable(startDate, endDate time.Time) bool {
	for _, booking := range c.Bookings {
		if startDate.Before(booking.EndDate) && endDate.After(booking.StartDate) {
			return false
		}
	}
	return true
}
