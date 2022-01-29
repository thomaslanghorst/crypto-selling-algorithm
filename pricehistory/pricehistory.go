package pricehistory

import (
	"crypto-selling-algorithm/common"
	"strconv"
	"time"
)

const timeLayout = "2006-01-02 15:04:05 UTC"

type Reader interface {
	FromFile(csvFile string) ([]*PriceHistoryEntry, error)
}

type reader struct {
}

type PriceHistoryEntry struct {
	Date  time.Time
	Price float64
}

func NewReader() Reader {
	return &reader{}
}

func (r *reader) FromFile(csvFile string) ([]*PriceHistoryEntry, error) {
	allLines, err := common.ReadFile(csvFile)
	if err != nil {
		return nil, err
	}

	priceHistoryEntries := make([]*PriceHistoryEntry, 0, len(allLines))

	for i, line := range allLines {

		// skip header
		if i == 0 {
			continue
		}

		date, err := time.Parse(timeLayout, line[0])
		if err != nil {
			return nil, err
		}
		price, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			return nil, err
		}

		priceHistoryEntries = append(priceHistoryEntries, &PriceHistoryEntry{
			Date:  date,
			Price: price,
		})
	}

	return priceHistoryEntries, nil
}
