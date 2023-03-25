package ds

import (
	log "github.com/sirupsen/logrus"
	"strconv"
)

type ExecutorOrder struct {
	ExecutorID int
	OrderID    int
	Price      int
}

func Unmarshal(payload string) (execOrder ExecutorOrder, err error) {
	var flagReading bool
	var exec, order, price string
	var indicator int
	payload = payload[1 : len(payload)-1]
	for _, symbol := range payload {
		if symbol == ',' {
			flagReading = false
			indicator++
			continue
		}
		if symbol == ':' {
			flagReading = true
			continue
		}
		if flagReading {
			switch indicator {
			case 0:
				exec += string(symbol)
			case 1:
				order += string(symbol)
			case 2:
				price += string(symbol)
			}
		}
	}
	execOrder.OrderID, err = strconv.Atoi(order)
	if err != nil {
		log.Println("couldn't convert string to OrderID field")
		return execOrder, err
	}
	execOrder.ExecutorID, err = strconv.Atoi(exec)
	if err != nil {
		log.Println("couldn't convert string to ExecutorID field")
		return execOrder, err
	}
	execOrder.Price, err = strconv.Atoi(price)
	if err != nil {
		log.Println("couldn't convert string to ExecutorID field")
		return execOrder, err
	}
	return execOrder, err
}
