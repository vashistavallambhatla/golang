package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Car struct {
	ID           int
	Make         string
	Model        string
	Year         int
	LicensePlate string
	PricePerDay  float64
	IsAvailable  bool
	Bookings     []BookingPeriod
	CarType      string
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
	Payment   *Payment
}

type PaymentStage string

const (
	Pending    PaymentStage = "Pending"
	Processing PaymentStage = "Processing"
	Completed  PaymentStage = "Completed"
	Failed     PaymentStage = "Failed"
)

type Payment struct {
	ID            int
	ReservationID int
	PaymentMethod string
	Amount        float64
	PaymentStage  PaymentStage
	Refund        bool
	Timestamp     time.Time
	RefundAddress string // Placeholder for now
}

type CarRentalSystem struct {
	cars              map[int]*Car
	customers         map[int]*Customer
	reservations      map[int]*Reservation
	payments          map[int]*Payment
	licensePlates     map[string]bool
	customerLicense   map[string]bool
	nextCarId         int
	nextCustomerId    int
	nextReservationId int
	nextPaymentId     int
}

func NewCarRentalSystem() *CarRentalSystem {
	return &CarRentalSystem{
		cars:            make(map[int]*Car),
		customers:       make(map[int]*Customer),
		reservations:    make(map[int]*Reservation),
		licensePlates:   make(map[string]bool),
		customerLicense: make(map[string]bool),
		payments:        make(map[int]*Payment),
	}
}

func (crs *CarRentalSystem) enrollCar(make, model string, year int, licensePlate string, pricePerDay float64, carType string) (Car, error) {
	_, licensePlateExists := crs.licensePlates[licensePlate]
	if licensePlateExists {
		return Car{}, fmt.Errorf("a car with this license plate already exists")
	}

	if pricePerDay <= 0 {
		return Car{}, fmt.Errorf("price per day of a car can't be negative")
	}

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

	crs.licensePlates[licensePlate] = true
	return *crs.cars[id], nil
}

func (crs *CarRentalSystem) registerCustomer(name, contact, license string) (Customer, error) {
	_, licenseExists := crs.customerLicense[license]
	if licenseExists {
		return Customer{}, fmt.Errorf("Customer with license number %v already exists", license)
	}

	id := crs.nextCustomerId
	crs.nextCustomerId++
	crs.customers[id] = &Customer{
		ID:      id,
		Name:    name,
		Contact: contact,
		License: license,
	}

	crs.customerLicense[license] = true
	return *crs.customers[id], nil
}

func (crs *CarRentalSystem) makeReservation(carId, customerId int, startDate, endDate time.Time) (Reservation, error) {
	car, carExists := crs.cars[carId]
	if !carExists {
		return Reservation{}, fmt.Errorf("car with ID %v doesn't exist", carId)
	}

	_, customerExists := crs.customers[customerId]
	if !customerExists {
		return Reservation{}, fmt.Errorf("customer with ID %v doesn't exist", customerId)
	}

	checkErr := car.isAvailable(startDate, endDate)
	if checkErr != nil {
		return Reservation{}, checkErr
	}

	id := crs.nextReservationId
	crs.nextReservationId++
	daysCount := int(endDate.Sub(startDate).Hours() / 24)
	totalCost := car.PricePerDay * float64(daysCount)

	reservation := &Reservation{
		ID:        id,
		Customer:  customerId,
		CarId:     carId,
		StartDate: startDate,
		EndDate:   endDate,
		TotalDays: daysCount,
		TotalCost: totalCost,
	}

	payment, err := crs.processPayment(id, reservation.TotalCost)
	if err != nil {
		crs.nextReservationId-- // rolling back the counter
		return Reservation{}, fmt.Errorf("payment failed due to %v", err)
	}

	reservation.Payment = payment
	crs.reservations[id] = reservation
	car.Bookings = append(car.Bookings, BookingPeriod{StartDate: startDate, EndDate: endDate})
	return *reservation, nil
}

func (crs *CarRentalSystem) modifyReservation(reservationId int, startDate, endDate time.Time) (Reservation, error) {
	reservation, exists := crs.reservations[reservationId]
	if !exists {
		return Reservation{}, fmt.Errorf("reservation with the ID %v not found",reservationId)
	}

	originalDuration := reservation.EndDate.Sub(reservation.StartDate)
	newDuration := endDate.Sub(startDate)

	if originalDuration != newDuration { // Same window won't require a new payment - to keep things simple for now
		return Reservation{}, fmt.Errorf("modification not allowed: new date window must be the same as the original duration (%v)", originalDuration)
	}

	car, exists := crs.cars[reservation.CarId]
	if !exists {
		return Reservation{}, fmt.Errorf("car with ID %v not found",car.ID)
	}

	checkErr := car.isAvailable(startDate, endDate)
	if checkErr != nil {
		return Reservation{}, checkErr
	}

	for i, booking := range car.Bookings {
		if booking.StartDate.Equal(reservation.StartDate) && booking.EndDate.Equal(reservation.EndDate) {
			reservation.StartDate = startDate
			reservation.EndDate = endDate
			car.Bookings[i] = BookingPeriod{StartDate: startDate, EndDate: endDate}
			return *reservation, nil
		}
	}

	return Reservation{}, fmt.Errorf("original booking period not found")
}

func (crs *CarRentalSystem) cancelReservation(reservationId int) (string, error) {
	reservation, exists := crs.reservations[reservationId]
	if !exists {
		return "", fmt.Errorf("reservation with the ID %v not found",reservationId)
	}

	car, exists := crs.cars[reservation.CarId]
	if !exists {
		return "", fmt.Errorf("car with ID %v not found",car.ID)
	}

	for i, booking := range car.Bookings {
		if booking.StartDate.Equal(reservation.StartDate) && booking.EndDate.Equal(reservation.EndDate) {
			car.Bookings = append(car.Bookings[:i], car.Bookings[i+1:]...)
			break
		}
	}

	paymentId := reservation.Payment.ID
	err := crs.initiateRefund(paymentId)
	if err != nil {
		return "", fmt.Errorf("failed to initiate refund: %v", err)
	} // should ideally send a refundID incase customers want to follow up on this later,this is just a place holder logic for now.

	delete(crs.reservations, reservationId)
	return "Reservation cancelled successfully and initiated refund", nil
}

func (crs *CarRentalSystem) findAvailableCarsByFilters(carType string, price float64, startDate, endDate time.Time) ([]Car, error) {
	if carType == "" && price <= 0  && startDate.IsZero() && endDate.IsZero(){ // if no filters are provided, return all cars
		allCars := make([]Car, 0, len(crs.cars))
		for _, car := range crs.cars {
			allCars = append(allCars, *car)
		}
		return allCars, nil
	}

	var searchResult []Car

	for _, car := range crs.cars {
		checkErr := car.isAvailable(startDate, endDate)
		if checkErr != nil {
			continue
		}

		typeMatches := carType == "" || car.CarType == carType
		priceMatches := price <= 0 || car.PricePerDay <= price

		if typeMatches && priceMatches {
			searchResult = append(searchResult, *car)
		}
	}

	if len(searchResult) == 0 {
		return nil, fmt.Errorf("no cars found that match your requirements")
	}

	return searchResult, nil
}

func (c *Car) isAvailable(startDate, endDate time.Time) error {
	if startDate.IsZero() || endDate.IsZero() || !startDate.Before(endDate) {
		return fmt.Errorf("invalid booking window: start and end dates must be non-zero and start < end")
	}

	for _, booking := range c.Bookings {
		if startDate.Before(booking.EndDate) && endDate.After(booking.StartDate) {
			return fmt.Errorf("car is not available for the selected dates")
		}
	}

	return nil
}

func (crs *CarRentalSystem) processPayment(reservationID int, amount float64) (*Payment, error) {
	id := crs.nextPaymentId
	crs.nextPaymentId++

	payment := &Payment{
		ID:            id,
		ReservationID: reservationID,
		Amount:        amount,
		PaymentStage:  Processing,
		Timestamp:     time.Now(),
	}

	stage, err := callMockGateway(payment)
	payment.PaymentStage = stage

	if err != nil {
		crs.nextPaymentId--
		return nil, err
	}

	crs.payments[payment.ID] = payment

	return payment, nil
}

func (crs *CarRentalSystem) initiateRefund(paymentId int) error {
	payment, exists := crs.payments[paymentId]
	if !exists {
		return fmt.Errorf("payment with ID %d not found", paymentId)
	}
	payment.Refund = true
	return nil
} // This refund == true should be picked by a refund service ideally and processed

func callMockGateway(payment *Payment) (PaymentStage, error) {
	if rand.Intn(100) < 90 {
		return Completed, nil
	} else {
		return Failed, fmt.Errorf("server Issue - please try again")
	}
}


func main() {

	rand.New(rand.NewSource(time.Now().UnixNano()))

	rentalSystem := NewCarRentalSystem()

	sedan, err := rentalSystem.enrollCar("Toyota", "Camry", 2022, "ABC123", 50.0, "Sedan")
	if err != nil {
		fmt.Printf("Error enrolling car: %v\n", err)
		return
	}
	fmt.Printf("Enrolled car: %s %s (ID: %d)\n", sedan.Make, sedan.Model, sedan.ID)

	suv, err := rentalSystem.enrollCar("Honda", "CR-V", 2023, "XYZ789", 75.0, "SUV")
	if err != nil {
		fmt.Printf("Error enrolling car: %v\n", err)
		return
	}
	fmt.Printf("Enrolled car: %s %s (ID: %d)\n", suv.Make, suv.Model, suv.ID)


	customer, err := rentalSystem.registerCustomer("John Doe", "john@example.com", "DL12345")
	if err != nil {
		fmt.Printf("Error registering customer: %v\n", err)
		return
	}
	fmt.Printf("Registered customer: %s (ID: %d)\n", customer.Name, customer.ID)

	startDate := time.Now().AddDate(0, 0, 1)  
	endDate := time.Now().AddDate(0, 0, 4)   

	reservation, err := rentalSystem.makeReservation(sedan.ID, customer.ID, startDate, endDate)
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

	availableCars, err := rentalSystem.findAvailableCarsByFilters("SUV", 100.0, startDate, endDate)
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
	
	modifiedReservation, err := rentalSystem.modifyReservation(reservation.ID, newStartDate, newEndDate)
	if err != nil {
		fmt.Printf("\nError modifying reservation: %v\n", err)
	} else {
		fmt.Printf("\nReservation modified successfully:\n")
		fmt.Printf("New dates: %s to %s\n", 
			modifiedReservation.StartDate.Format("2006-01-02"), 
			modifiedReservation.EndDate.Format("2006-01-02"))
	}

	cancelResult, err := rentalSystem.cancelReservation(reservation.ID)
	if err != nil {
		fmt.Printf("\nError cancelling reservation: %v\n", err)
	} else {
		fmt.Printf("\n%s\n", cancelResult)
	}

	
	invalidStartDate := time.Now().AddDate(0, 0, 10)
	invalidEndDate := time.Now().AddDate(0, 0, 5) 
	
	_, err = rentalSystem.makeReservation(suv.ID, customer.ID, invalidStartDate, invalidEndDate)
	fmt.Printf("\nExpected error with invalid dates: %v\n", err)


	cars, err := rentalSystem.findAvailableCarsByFilters("", 0, startDate, endDate)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(cars)
	}
}