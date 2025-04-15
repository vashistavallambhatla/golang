package services

import (
	"fmt"
	"math/rand"
	"time"

	"crs/models"
)

type CarRentalSystem struct {
	cars              map[int]*models.Car
	customers         map[int]*models.Customer
	reservations      map[int]*models.Reservation
	payments          map[int]*models.Payment
	licensePlates     map[string]bool
	customerLicense   map[string]bool
	nextCarId         int
	nextCustomerId    int
	nextReservationId int
	nextPaymentId     int
}

// NewCarRentalSystem creates a new instance of the car rental system
func NewCarRentalSystem() *CarRentalSystem {
	return &CarRentalSystem{
		cars:            make(map[int]*models.Car),
		customers:       make(map[int]*models.Customer),
		reservations:    make(map[int]*models.Reservation),
		licensePlates:   make(map[string]bool),
		customerLicense: make(map[string]bool),
		payments:        make(map[int]*models.Payment),
	}
}

// EnrollCar adds a new car to the system
func (crs *CarRentalSystem) EnrollCar(make, model string, year int, licensePlate string, pricePerDay float64, carType string) (models.Car, error) {
	_, licensePlateExists := crs.licensePlates[licensePlate]
	if licensePlateExists {
		return models.Car{}, fmt.Errorf("a car with this license plate already exists")
	}

	if pricePerDay <= 0 {
		return models.Car{}, fmt.Errorf("price per day of a car can't be negative")
	}

	id := crs.nextCarId
	crs.nextCarId++
	crs.cars[id] = &models.Car{
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

// RegisterCustomer registers a new customer in the system
func (crs *CarRentalSystem) RegisterCustomer(name, contact, license string) (models.Customer, error) {
	_, licenseExists := crs.customerLicense[license]
	if licenseExists {
		return models.Customer{}, fmt.Errorf("customer with license number %v already exists", license)
	}

	id := crs.nextCustomerId
	crs.nextCustomerId++
	crs.customers[id] = &models.Customer{
		ID:      id,
		Name:    name,
		Contact: contact,
		License: license,
	}

	crs.customerLicense[license] = true
	return *crs.customers[id], nil
}

// MakeReservation creates a new reservation in the system
func (crs *CarRentalSystem) MakeReservation(carId, customerId int, startDate, endDate time.Time) (models.Reservation, error) {
	car, carExists := crs.cars[carId]
	if !carExists {
		return models.Reservation{}, fmt.Errorf("car with ID %v doesn't exist", carId)
	}

	_, customerExists := crs.customers[customerId]
	if !customerExists {
		return models.Reservation{}, fmt.Errorf("customer with ID %v doesn't exist", customerId)
	}

	checkErr := car.IsCarAvailable(startDate, endDate)
	if checkErr != nil {
		return models.Reservation{}, checkErr
	}

	id := crs.nextReservationId
	crs.nextReservationId++
	daysCount := int(endDate.Sub(startDate).Hours() / 24)
	totalCost := car.PricePerDay * float64(daysCount)

	reservation := &models.Reservation{
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
		return models.Reservation{}, fmt.Errorf("payment failed due to %v", err)
	}

	reservation.Payment = payment
	crs.reservations[id] = reservation
	car.Bookings = append(car.Bookings, models.BookingPeriod{StartDate: startDate, EndDate: endDate})
	return *reservation, nil
}

// ModifyReservation modifies an existing reservation
func (crs *CarRentalSystem) ModifyReservation(reservationId int, startDate, endDate time.Time) (models.Reservation, error) {
	reservation, exists := crs.reservations[reservationId]
	if !exists {
		return models.Reservation{}, fmt.Errorf("reservation with the ID %v not found", reservationId)
	}

	originalDuration := reservation.EndDate.Sub(reservation.StartDate)
	newDuration := endDate.Sub(startDate)

	if originalDuration != newDuration {
		return models.Reservation{}, fmt.Errorf("modification not allowed: new date window must be the same as the original duration (%v)", originalDuration)
	}

	car, exists := crs.cars[reservation.CarId]
	if !exists {
		return models.Reservation{}, fmt.Errorf("car with ID %v not found", car.ID)
	}

	checkErr := car.IsCarAvailable(startDate, endDate)
	if checkErr != nil {
		return models.Reservation{}, checkErr
	}

	for i, booking := range car.Bookings {
		if booking.StartDate.Equal(reservation.StartDate) && booking.EndDate.Equal(reservation.EndDate) {
			reservation.StartDate = startDate
			reservation.EndDate = endDate
			car.Bookings[i] = models.BookingPeriod{StartDate: startDate, EndDate: endDate}
			return *reservation, nil
		}
	}

	return models.Reservation{}, fmt.Errorf("original booking period not found")
}

// CancelReservation cancels an existing reservation
func (crs *CarRentalSystem) CancelReservation(reservationId int) (string, error) {
	reservation, exists := crs.reservations[reservationId]
	if !exists {
		return "", fmt.Errorf("reservation with the ID %v not found", reservationId)
	}

	car, exists := crs.cars[reservation.CarId]
	if !exists {
		return "", fmt.Errorf("car with ID %v not found", car.ID)
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
	}

	delete(crs.reservations, reservationId)
	return "Reservation cancelled successfully and initiated refund", nil
}

// FindAvailableCarsByFilters finds available cars based on filters
func (crs *CarRentalSystem) FindAvailableCarsByFilters(carType string, price float64, startDate, endDate time.Time) ([]models.Car, error) {
	if carType == "" && price <= 0 && startDate.IsZero() && endDate.IsZero() {
		allCars := make([]models.Car, 0, len(crs.cars))
		for _, car := range crs.cars {
			allCars = append(allCars, *car)
		}
		return allCars, nil
	}

	var searchResult []models.Car

	for _, car := range crs.cars {
		checkErr := car.IsCarAvailable(startDate, endDate)
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

// processPayment processes a payment for a reservation
func (crs *CarRentalSystem) processPayment(reservationID int, amount float64) (*models.Payment, error) {
	id := crs.nextPaymentId
	crs.nextPaymentId++

	payment := &models.Payment{
		ID:            id,
		ReservationID: reservationID,
		Amount:        amount,
		PaymentStage:  models.Processing,
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

// initiateRefund initiates a refund for a payment
func (crs *CarRentalSystem) initiateRefund(paymentId int) error {
	payment, exists := crs.payments[paymentId]
	if !exists {
		return fmt.Errorf("payment with ID %d not found", paymentId)
	}
	payment.Refund = true
	return nil
}

// callMockGateway simulates calling a payment gateway
func callMockGateway(payment *models.Payment) (models.PaymentStage, error) {
	if rand.Intn(100) < 90 {
		return models.Completed, nil
	} else {
		return models.Failed, fmt.Errorf("server Issue - please try again")
	}
}
