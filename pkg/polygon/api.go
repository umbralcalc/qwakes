package polygon

import (
	"context"
	"fmt"
	"os"
	"time"

	rest "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
)

// FetchMinutelyAggs requests Polygon daily aggs from the API with
// the given ticker and date range.
func FetchMinutelyAggs(ticker string, from, to time.Time) ([]models.Agg, error) {
	apiKey := os.Getenv("POLYGON_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("missing POLYGON_API_KEY")
	}

	c := rest.New(apiKey)
	ctx := context.Background()

	params := models.ListAggsParams{
		Ticker:     ticker,
		Multiplier: 1,
		Timespan:   "minute",
		From:       models.Millis(from),
		To:         models.Millis(to),
	}.WithOrder(models.Asc).WithAdjusted(true)

	iter := c.ListAggs(ctx, params)
	var out []models.Agg
	for iter.Next() {
		out = append(out, iter.Item())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
