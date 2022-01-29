# Crypto Selling Algorithm

The crypto selling algorithm is a small program that indicates selling opportunities for crypto coins.

It looks at the prices of Solana and Cardano from _2021-01-18_ to _2022-01-18_  and watches out for two criteria to find selling points:

1. the [fear and greed index](https://alternative.me/crypto/fear-and-greed-index/)
2. the [price of a coin relative](https://www.coingecko.com/en/coins/solana) to ethereum

## Selling criteria

Selling criteria can be changed within the `SellingCriteria` struct, where the values are the following:

`SellingPercentage` -  what percentage of the total coin amount should be sold at a selling opportunity

`ProfitFactor` - how high the minimum selling price should be in regards to the highest buying price

`FearAndGreedSellingThreshold` - the minimum value the fear and greed index should have

`ChangePercentage` - the minimum change percentage of the price of a coin relative to ethereum from one day to another 

`ChangePercentageCount` - how many days in a row the change percentage value should be above the defined value

## Buying information 

Buying information of a coin can be defined in the `coinInfos` slice. At the moment, only one buying time can be defined though. Each coin has its own struct containing the following values:

`CoinSymbol` - the symbol of the coin, e.g. SOL or ADA

`HighestBuyPrice` - the price at which the coin was bought

`AmountBought` - the number of coins that have been bought

## Running the algorithm

The code can be executed running `go run main.go` in the terminal

## Example output

```
SOL - SELLING OPPORTUNITIES
        2021-02-24 - sell: 6.2 coins - coins left: 18.8
        2021-02-25 - sell: 4.7 coins - coins left: 14.1
        2021-03-29 - sell: 3.5 coins - coins left: 10.5
        2021-08-18 - sell: 2.6 coins - coins left: 7.9
        2021-08-19 - sell: 2.0 coins - coins left: 5.9
        2021-08-29 - sell: 1.5 coins - coins left: 4.4
        2021-09-09 - sell: 1.1 coins - coins left: 3.3
TOTAL GAINS: 910.430528


ADA - SELLING OPPORTUNITIES
        2021-02-27 - sell: 62.5 coins - coins left: 187.5
        2021-02-28 - sell: 46.9 coins - coins left: 140.6
TOTAL GAINS: 140.760512
```
