package main

import (
	"crypto-selling-algorithm/feargreedindex"
	"crypto-selling-algorithm/pricehistory"
	"fmt"
	"time"
)

const timeLayout = "2006-01-02"

type SellingCriteria struct {
	SellingPercentage            int
	ProfitFactor                 float64
	FearAndGreedSellingThreshold int
	ChangePercentage             float64
	ChangePercentageCount        int
}

type CoinInfo struct {
	CoinSymbol      string
	AmountBought    float64
	HighestBuyPrice float64
}

type RelativePrice struct {
	Date               time.Time
	EthPrice           float64
	CoinPrice          float64
	RelativeToEthPrice float64
	ChangePercentage   float64
}

func main() {
	sc := &SellingCriteria{
		SellingPercentage:            25,
		ProfitFactor:                 2.0, // due to taxes
		FearAndGreedSellingThreshold: 45,  // chosen by science
		ChangePercentage:             5.0, // chosen by science
		ChangePercentageCount:        3,   // chosen by science
	}

	coinInfos := []*CoinInfo{
		{ // 100 USD invested @ 25 USD per coin -> 4 coins bought
			CoinSymbol:      "SOL",
			HighestBuyPrice: 4,
			AmountBought:    25.0,
		},
		{ // 100 USD invested @ 0.40 USD per coin -> 250 coins bought
			CoinSymbol:      "ADA",
			HighestBuyPrice: 0.4,
			AmountBought:    250.0,
		},
	}

	fgiData := readFearAndGreed()
	ethPrices := readHistory("ETH")

	for _, ci := range coinInfos {
		coinPrices := readHistory(ci.CoinSymbol)
		relativePrices := calcRelativePrices(ethPrices, coinPrices)

		calcChangePercentages(relativePrices)
		calcSellingOpportunities(fgiData, relativePrices, sc, ci)
	}
}

func readFearAndGreed() []*feargreedindex.FearAndGreedDatapoint {
	fgi, err := feargreedindex.New().FromFile("data/fgi.csv")
	if err != nil {
		panic(err)
	}
	return fgi
}

func readHistory(symbol string) []*pricehistory.PriceHistoryEntry {
	phr := pricehistory.NewReader()

	coinPrices, err := phr.FromFile(fmt.Sprintf("data/%s-usd-max.csv", symbol))
	if err != nil {
		panic(err)
	}

	return coinPrices
}

func calcRelativePrices(ethPrices, coinPrices []*pricehistory.PriceHistoryEntry) []*RelativePrice {
	if len(ethPrices) != len(coinPrices) {
		panic("pricehistory slices must have the same length")
	}

	prices := make([]*RelativePrice, 0, len(ethPrices))

	for idx, ethEntry := range ethPrices {
		ethPrice := ethEntry.Price
		coinPrice := coinPrices[idx].Price

		prices = append(prices, &RelativePrice{
			Date:               ethEntry.Date,
			EthPrice:           ethPrice,
			CoinPrice:          coinPrice,
			RelativeToEthPrice: coinPrice / ethPrice,
		})
	}

	return prices
}

func calcChangePercentages(relativePrices []*RelativePrice) {
	for idx, p := range relativePrices {
		if idx == 0 {
			p.ChangePercentage = 0.0
			continue
		}

		// ((newvalue / oldvalue) - 1) * 100

		// examples:

		// 200 / 100 = 2
		// 2 - 1 = 1
		// 1 * 100 = 100 %

		// 100 / 200 = 0,5
		// 0,5 - 1 = -0,5
		// -0,5 * 100 = -50 %

		prevVal := relativePrices[idx-1]

		p.ChangePercentage = ((p.RelativeToEthPrice / prevVal.RelativeToEthPrice) - 1) * 100.0
	}
}

func calcSellingOpportunities(fgiData []*feargreedindex.FearAndGreedDatapoint, relativePrices []*RelativePrice, sellingCriteria *SellingCriteria, coinInfo *CoinInfo) {

	if len(fgiData) != len(relativePrices) {
		panic("fearAndGreed slice must have the same length as relaivePrices slice")
	}

	fmt.Printf("%s - SELLING OPPORTUNITIES\n", coinInfo.CoinSymbol)

	coinAmount := coinInfo.AmountBought
	changesAboveThreshold := 0
	totalGains := 0.0

	for i := 0; i < len(fgiData); i++ {
		fgiDatapoint := fgiData[i]
		relativePrice := relativePrices[i]

		// price of coin must be at least HighestBuyingPrice * ProfitFactor defined in SellingCriteria
		if relativePrice.CoinPrice < coinInfo.HighestBuyPrice*sellingCriteria.ProfitFactor {
			continue
		}

		// fear and greed index must be at least FearAndGreedSellingThreshold defined in SellingCriteria
		if fgiDatapoint.Value < sellingCriteria.FearAndGreedSellingThreshold {
			continue
		}

		// price change percentage must be at least ChangePercentage defined in SellingCriteria
		if relativePrice.ChangePercentage >= sellingCriteria.ChangePercentage {
			changesAboveThreshold++
		} else {
			changesAboveThreshold = 0
		}

		// positive changes must be at least the same as ChangePercentageCount defined in SellingCriteria
		if changesAboveThreshold >= sellingCriteria.ChangePercentageCount {
			sellAmount := coinAmount * float64(sellingCriteria.SellingPercentage) / 100.0
			coinAmount = coinAmount - sellAmount
			totalGains = totalGains + float64(sellAmount)*relativePrice.CoinPrice

			fmt.Printf("\t%s - sell: %.1f coins - coins left: %.1f\n", relativePrice.Date.Format(timeLayout), sellAmount, coinAmount)
		}
	}

	fmt.Printf("TOTAL GAINS: %f\n", totalGains)
	fmt.Println("")
	fmt.Println("")

}
