package features

import (
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/umbralcalc/qwakes/pkg/polygon"
)

func GetMinutelyAggsDataFrame(ticker string, from, to time.Time) dataframe.DataFrame {
	aggs, err := polygon.FetchMinutelyAggs(ticker, from, to)
	if err != nil {
		panic(err)
	}
	tickers := make([]string, 0)
	closes := make([]float64, 0)
	highs := make([]float64, 0)
	lows := make([]float64, 0)
	transactions := make([]int64, 0)
	opens := make([]float64, 0)
	unixMillisTimestamps := make([]float64, 0)
	volumes := make([]float64, 0)
	vWAPs := make([]float64, 0)
	oTCs := make([]bool, 0)
	for _, agg := range aggs {
		tickers = append(tickers, agg.Ticker)
		closes = append(closes, agg.Close)
		highs = append(highs, agg.High)
		lows = append(lows, agg.Low)
		transactions = append(transactions, agg.Transactions)
		opens = append(opens, agg.Open)
		unixMillisTimestamps = append(
			unixMillisTimestamps,
			float64(time.Time(agg.Timestamp).UnixMilli()),
		)
		volumes = append(volumes, agg.Volume)
		vWAPs = append(vWAPs, agg.VWAP)
		oTCs = append(oTCs, agg.OTC)
	}
	df := dataframe.New(
		series.New(tickers, series.String, "ticker"),
		series.New(closes, series.Float, "close"),
		series.New(highs, series.Float, "high"),
		series.New(lows, series.Float, "low"),
		series.New(transactions, series.Int, "transations"),
		series.New(opens, series.Float, "open"),
		series.New(unixMillisTimestamps, series.Float, "unix_millis_timestamp"),
		series.New(volumes, series.Float, "volume"),
		series.New(vWAPs, series.Float, "VWAP"),
		series.New(oTCs, series.Bool, "OTC"),
	)
	return df
}
