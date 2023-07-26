//go:build !solution

package hotelbusiness

import "fmt"

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {

	load := make([]Load, 0)

	if len(guests) == 0 {
		return load

	}

	checkInCounter := make(map[int]int)
	checkOutCounter := make(map[int]int)

	min := guests[0].CheckInDate
	max := guests[0].CheckOutDate

	for _, guest := range guests {
		_, exists := checkInCounter[guest.CheckInDate]
		if exists {
			checkInCounter[guest.CheckInDate] += 1
		} else {
			checkInCounter[guest.CheckInDate] = 1
		}

		_, exists = checkOutCounter[guest.CheckOutDate]
		if exists {
			checkOutCounter[guest.CheckOutDate] += 1
		} else {
			checkOutCounter[guest.CheckOutDate] = 1
		}

		if guest.CheckInDate < min {
			min = guest.CheckInDate
		}
		if guest.CheckOutDate > max {
			max = guest.CheckOutDate
		}
	}

	guestsAmount := 0

	fmt.Println(checkInCounter)
	fmt.Println(checkOutCounter)

	for ; min <= max; min++ {
		fmt.Println(min)

		checkIn, checkInExists := checkInCounter[min]
		checkOut, checkOutExists := checkOutCounter[min]

		if !checkInExists && !checkOutExists {
			continue
		}
		if !checkInExists {
			checkIn = 0
		}
		if !checkOutExists {
			checkOut = 0
		}

		if checkIn == checkOut {
			continue
		}

		fmt.Println(checkIn, checkOut)

		guestsAmount += checkIn - checkOut

		loadLog := Load{
			StartDate:  min,
			GuestCount: guestsAmount,
		}
		load = append(load, loadLog)

	}

	return load
}
