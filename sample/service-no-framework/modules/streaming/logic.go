package streaming

import (
	"github.com/jamestrandung/go-cte/sample/config"
	"github.com/jamestrandung/go-cte/sample/dto"
)

func StreamQuote(quote *dto.Quote) {
	config.Print("Streaming calculated cost:", quote.TotalCost)
}
