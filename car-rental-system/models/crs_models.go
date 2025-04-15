package models

import (
	"fmt"
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
	RefundAddress string
}

func (c *Car) IsCarAvailable(startDate, endDate time.Time) error {
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

