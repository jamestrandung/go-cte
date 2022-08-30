package vat

import (
    "github.com/jamestrandung/go-cte/sample/dependencies/configsfetcher"
    "github.com/jamestrandung/go-cte/sample/dto"
)

func AddVAT(quote *dto.Quote, costConfigs configsfetcher.MergedCostConfigs) {
    vatAmount := quote.TotalCost * costConfigs.VATPercent / 100

    quote.VATAmount = vatAmount
    quote.TotalCost = quote.TotalCost + vatAmount
}
