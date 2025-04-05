package main

import (
	"fmt"
	"time"
)

type Car struct {
	ID            int
	Make          string
	Model         string
	Year          int
	LicensePlate  string
	PricePerDay   float64
	IsAvailable   bool
	Bookings      []BookingPeriod
	CarType       string
}

type BookingPeriod struct {
	StartDate time.Time
	EndDate   time.Time
}

type Customer struct {
	ID      int
	Name    string
	Contact string
	License string
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

func (crs *CarRentalSystem) addCar(make, model string, year int, licensePlate string, pricePerDay float64, carType string) error {
	id := crs.nextCarId
	crs.nextCarId++
	crs.cars[id] = &Car{
		ID:           id,
		Make:         make,
		Model:        model,
		Year:         year,
		LicensePlate: licensePlate,
		PricePerDay:  pricePerDay,
		IsAvailable:  true,
		CarType:      carType,
	}
	return nil
}

func (crs *CarRentalSystem) addCustomer(name, contact, license string) error {
	id := crs.nextCustomerId
	crs.nextCustomerId++
	crs.customers[id] = &Customer{
		ID:      id,
		Name:    name,
		Contact: contact,
		License: license,
	}
	return nil
}

func (crs *CarRentalSystem) makeReservation(carId, customerId int, startDate, endDate time.Time) error {
	car, exists := crs.cars[carId]
	if !exists {
		return fmt.Errorf("car not found")
	}

	if !car.isAvailable(startDate, endDate) {
		return fmt.Errorf("car is not available for the selected dates")
	}

	id := crs.nextReservationId
	crs.nextReservationId++
	daysCount := int(endDate.Sub(startDate).Hours() / 24)
	totalCost := car.PricePerDay * float64(daysCount)

	crs.reservations[id] = &Reservation{
		ID:        id,
		Customer:  customerId,
		CarId:     carId,
		StartDate: startDate,
		EndDate:   endDate,
		TotalDays: daysCount,
		TotalCost: totalCost,
	}

	car.Bookings = append(car.Bookings, BookingPeriod{StartDate: startDate, EndDate: endDate})
	return nil
}

func (crs *CarRentalSystem) modifyReservation(reservationId int, startDate, endDate time.Time) error {
	reservation, exists := crs.reservations[reservationId]
	if !exists {
		return fmt.Errorf("reservation not found")
	}

	car, exists := crs.cars[reservation.CarId]
	if !exists {
		return fmt.Errorf("car not found")
	}

	if !car.isAvailable(startDate, endDate) {
		return fmt.Errorf("car is not available for the selected dates")
	}

	for i, b := range car.Bookings {
		if b.StartDate.Equal(reservation.StartDate) && b.EndDate.Equal(reservation.EndDate) {
			reservation.StartDate = startDate
			reservation.EndDate = endDate
			car.Bookings[i] = BookingPeriod{StartDate: startDate, EndDate: endDate}
			return nil
		}
	}
	return fmt.Errorf("original booking period not found")
}

func (crs *CarRentalSystem) cancelReservation(reservationId int) error {
	reservation, exists := crs.reservations[reservationId]
	if !exists {
		return fmt.Errorf("reservation not found")
	}

	car, exists := crs.cars[reservation.CarId]
	if !exists {
		return fmt.Errorf("car not found")
	}

	for i, b := range car.Bookings {
		if b.StartDate.Equal(reservation.StartDate) && b.EndDate.Equal(reservation.EndDate) {
			car.Bookings = append(car.Bookings[:i], car.Bookings[i+1:]...)
			break
		}
	}

	delete(crs.reservations, reservationId)
	return nil
}

func (crs *CarRentalSystem) search(carType string, price float64, startDate, endDate time.Time) []Car {
	var searchResult []Car
	for _, car := range crs.cars {
		if car.CarType == carType && car.PricePerDay <= price && car.isAvailable(startDate, endDate) {
			searchResult = append(searchResult, *car)
		}
	}
	return searchResult
}

func (c *Car) isAvailable(startDate, endDate time.Time) bool {
	for _, b := range c.Bookings {
		if startDate.Before(b.EndDate) && endDate.After(b.StartDate) {
			return false
		}
	}
	return true
}
