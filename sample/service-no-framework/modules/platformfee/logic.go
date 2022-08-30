package platformfee

import (
    "github.com/jamestrandung/go-cte/sample/dependencies/configsfetcher"
    "github.com/jamestrandung/go-cte/sample/dto"
)

func AddPlatformFee(quote *dto.Quote, costConfigs configsfetcher.MergedCostConfigs) {
    quote.TotalCost = quote.TotalCost + costConfigs.PlatformFee
}
