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
	LicensePlate  string  // PascalCase for struct fields instead of License_plate
	PricePerDay  float64
	IsAvailable      bool  // Changed Available to IsAvailable for better readability
	Bookings       []BookingPeriod
	CarType        string
}

type BookingPeriod struct {
	StartDate time.Time
	EndDate   time.Time
}

type Customer struct {
	ID      int
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

// changed functions to lowercase since they're not being exported
// Changed the receiver's name from rental to crs to keep it simple
func (crs *CarRentalSystem) addCar(make string, model string, year int, licensePlate string, pricePerDay float64, carType string) { // Changed license_plate to licensePlate to follow camelCase
	id := crs.nextCarId
	crs.nextCarId++
	crs.cars[id] = &Car{
		ID:            id,
		Make:          make,
		Model:         model,
		Year:          year,
		LicensePlate: licensePlate,
		PricePerDay: pricePerDay,
		IsAvailable:     true,
		// Bookings:      []BookingPeriod{}, removed unecessary declaration since go handles nil values properly
		CarType:       carType,
	}
}

func (crs *CarRentalSystem) addCustomer(name, contact, license string) {
	id := crs.nextCarId // apparently rental.nextCustomerId++ should happen after storing its value to avoid potential issues in concurrent scenarios.
	crs.nextCustomerId++
	crs.customers[id] = &Customer{
		ID:      id,
		Name:    name,
		Contact: contact,
		License: license,
	}
}

func (crs *CarRentalSystem) makeReservation(carId, customerId int, startDate, endDate time.Time) {
	car, exists := crs.cars[carId]
	if !exists {
		fmt.Println("Car not found")
		return
	}

	if !car.isAvailable(startDate, endDate) {
		fmt.Println("Car is not available for the selected dates")
		return
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
}

func (crs *CarRentalSystem) modifyReservation(reservationId int, startDate, endDate time.Time) {
	reservation, exists := crs.reservations[reservationId]
	if !exists {
		fmt.Println("Reservation not found")
		return
	}

	car, exists := crs.cars[reservation.CarId] // Added go's ok idiom for this map lookup
	if !exists {
		fmt.Println("Car not found")
		return 
	}

	if !car.isAvailable(startDate,endDate) { // Early return if the modification is not possible
		fmt.Println("Car not available on these days")
		return
	}

	for i, booking := range car.Bookings {
		if booking.StartDate.Equal(reservation.StartDate) && booking.EndDate.Equal(reservation.EndDate) { // Used Equal instead of == for better time zone comparison
			reservation.StartDate = startDate
	        reservation.EndDate = endDate
			car.Bookings[i] = BookingPeriod{StartDate: startDate, EndDate: endDate}
			fmt.Printf("Reservation dates successfully changed to %v - %v\n", startDate, endDate)
			return
		}
	}
}

func (crs *CarRentalSystem) cancelReservation(reservationId int) {
	reservation, exists := crs.reservations[reservationId]
	if !exists {
		fmt.Println("Reservation not found")
		return
	}

	car, exists := crs.cars[reservation.CarId]
	if !exists {
		fmt.Println("Car not found")
		return
	}

	for i, b := range car.Bookings {
		if b.StartDate == reservation.StartDate && b.EndDate == reservation.EndDate {
			car.Bookings = append(car.Bookings[:i], car.Bookings[i+1:]...)
			break
		}
	}

	delete(crs.reservations, reservationId)

	fmt.Println("Reservation canceled successfully")
}

func (crs *CarRentalSystem) search(carType string, price float64, startDate, endDate time.Time) (searchResult []Car) {
	for _, car := range crs.cars {
		if car.CarType == carType && car.PricePerDay <= price && car.isAvailable(startDate, endDate) {
			searchResult = append(searchResult, *car)
		}
	}
	return searchResult
}

func (c *Car) isAvailable(startDate, endDate time.Time) bool {
	for _, b := range c.Bookings {    // b instead of booking since it's already evident
		if startDate.Before(b.EndDate) && endDate.After(b.StartDate) {
			return false
		}
	}
	return true
}

