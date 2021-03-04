package tvscanner

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	API_URL          = "https://scanner.tradingview.com/"
	DEFAULT_SCREENER = "crypto"
	API_POSTFIX      = "scan"
	//
	INTERVAL_1_MINUTE   = "1m"
	INTERVAL_5_MINUTES  = "5m"
	INTERVAL_15_MINUTES = "15m"
	INTERVAL_1_HOUR     = "1h"
	INTERVAL_4_HOURS    = "4h"
	INTERVAL_1_DAY      = "1d"
	INTERVAL_1_WEEK     = "1W"
	INTERVAL_1_MONTH    = "1M"
)

var (
	ContextLogger = logrus.WithFields(logrus.Fields{
		"client_name": "tvscanner",
	})

	//indicators = []string{"Recommend.Other", "Recommend.All", "Recommend.MA", "RSI", "RSI[1]", "Stoch.K", "Stoch.D", "Stoch.K[1]", "Stoch.D[1]", "CCI20", "CCI20[1]", "ADX", "ADX+DI", "ADX-DI", "ADX+DI[1]", "ADX-DI[1]", "AO", "AO[1]", "Mom", "Mom[1]", "MACD.macd", "MACD.signal", "Rec.Stoch.RSI", "Stoch.RSI.K", "Rec.WR", "W.R", "Rec.BBPower", "BBPower", "Rec.UO", "UO", "close", "EMA5", "SMA5", "EMA10", "SMA10", "EMA20", "SMA20", "EMA30", "SMA30", "EMA50", "SMA50", "EMA100", "SMA100", "EMA200", "SMA200", "Rec.Ichimoku", "Ichimoku.BLine", "Rec.VWMA", "VWMA", "Rec.HullMA9", "HullMA9"}
	recommendsList  = []string{"Recommend.Other", "Recommend.All", "Recommend.MA"}
	oscillatorsList = []string{"RSI", "RSI[1]", "Stoch.K", "Stoch.D", "Stoch.K[1]", "Stoch.D[1]", "CCI20", "CCI20[1]", "ADX", "ADX+DI", "ADX-DI", "ADX+DI[1]", "ADX-DI[1]", "AO", "AO[1]", "Mom", "Mom[1]", "MACD.macd", "MACD.signal", "Rec.Stoch.RSI", "Stoch.RSI.K", "Rec.WR", "W.R", "Rec.BBPower", "BBPower", "Rec.UO", "UO", "close"}
	maList          = []string{"EMA10", "SMA10", "EMA20", "SMA20", "EMA30", "SMA30", "EMA50", "SMA50", "EMA100", "SMA100", "EMA200", "SMA200"}
	maSimpleList    = []string{"Rec.Ichimoku", "Ichimoku.BLine", "Rec.VWMA", "VWMA", "Rec.HullMA9", "HullMA9"}
	pivots          = []string{"Pivot.M.Classic.S3", "Pivot.M.Classic.S2", "Pivot.M.Classic.S1", "Pivot.M.Classic.Middle", "Pivot.M.Classic.R1", "Pivot.M.Classic.R2", "Pivot.M.Classic.R3", "Pivot.M.Fibonacci.S3", "Pivot.M.Fibonacci.S2", "Pivot.M.Fibonacci.S1", "Pivot.M.Fibonacci.Middle", "Pivot.M.Fibonacci.R1", "Pivot.M.Fibonacci.R2", "Pivot.M.Fibonacci.R3", "Pivot.M.Camarilla.S3", "Pivot.M.Camarilla.S2", "Pivot.M.Camarilla.S1", "Pivot.M.Camarilla.Middle", "Pivot.M.Camarilla.R1", "Pivot.M.Camarilla.R2", "Pivot.M.Camarilla.R3", "Pivot.M.Woodie.S3", "Pivot.M.Woodie.S2", "Pivot.M.Woodie.S1", "Pivot.M.Woodie.Middle", "Pivot.M.Woodie.R1", "Pivot.M.Woodie.R2", "Pivot.M.Woodie.R3", "Pivot.M.Demark.S1", "Pivot.M.Demark.Middle", "Pivot.M.Demark.R1"}
)

// Scanner represent a Scanner client
type Scanner struct {
	client           *client
	data             DataResponse
	recommendSummary RecommendSummary
}

type Data struct {
	Symbols struct {
		Tickers []string `json:"tickers"`
		Query   struct {
			Types []string `json:"types"`
		} `json:"query"`
	} `json:"symbols"`
	Columns []string `json:"columns"`
}

type DataResponse struct {
	Data []struct {
		Symbol string    `json:"s"`
		Data   []float64 `json:"d"`
	} `json:"data"`
	TotalCount int `json:"totalCount"`
}

type RecommendSummary struct {
	Recommend    Recommend
	BuyCount     int
	SellCount    int
	NeutralCount int
}

type Recommend struct {
	Summary        string
	Oscillators    string
	MovingAverages string
}

// New returns an instantiated Scanner struct
func New() *Scanner {
	return &Scanner{client: NewClient()}
}

// NewWithCustomHttpClient returns an instantiated Scanner struct with custom http client
func NewWithCustomHttpClient(httpClient *http.Client) *Scanner {
	return &Scanner{client: NewClientWithCustomHttpConfig(httpClient)}
}

// set enable/disable http request/response dump
func (c *Scanner) SetDebug(enable bool) {
	c.client.debug = enable
}

func concatAppend(slices [][]string) []string {
	var tmp []string
	for _, s := range slices {
		tmp = append(tmp, s...)
	}
	return tmp
}

// PrepareData prepare payload for request
func (c *Scanner) PrepareData(symbol, interval string, indicators []string) ([]byte, error) {
	// Default, 1 Day
	dataInterval := ""

	if interval == INTERVAL_1_MINUTE {
		// 1 Minute
		dataInterval = "|1"
	} else if interval == INTERVAL_5_MINUTES {
		// 5 Minutes
		dataInterval = "|5"
	} else if interval == INTERVAL_15_MINUTES {
		// 15 Minutes
		dataInterval = "|15"
	} else if interval == INTERVAL_1_HOUR {
		// 1 Hour
		dataInterval = "|60"
	} else if interval == INTERVAL_4_HOURS {
		// 4 Hour
		dataInterval = "|240"
	} else if interval == INTERVAL_1_WEEK {
		// 1 Week
		dataInterval = "|1W"
	} else if interval == INTERVAL_1_MONTH {
		// 1 Month
		dataInterval = "|1M"
	} else {
		if interval != INTERVAL_1_DAY {
			fmt.Println("Interval is empty or not valid, defaulting to 1 day.")
			// Default, 1 Day
			dataInterval = ""
		}
	}

	//indicators := concatAppend([][]string{recommendsList, oscillatorsList, maList, maSimpleList})

	data := Data{}
	data.Symbols.Tickers = []string{symbol}
	for _, ind := range indicators {
		data.Columns = append(data.Columns, fmt.Sprintf("%s%s", ind, dataInterval))
	}
	return json.Marshal(data)
}

func (c *Scanner) PrepareRecommendations() (err error) {
	// RECOMMENDATIONS
	c.recommendSummary.Recommend.Oscillators, err = c.ComputeRecommend(c.data.Data[0].Data[0])
	if err != nil {
		ContextLogger.Error(err)
		return err
	}
	c.recommendSummary.Recommend.Summary, err = c.ComputeRecommend(c.data.Data[0].Data[1])
	if err != nil {
		ContextLogger.Error(err)
		return err
	}
	c.recommendSummary.Recommend.MovingAverages, err = c.ComputeRecommend(c.data.Data[0].Data[2])
	if err != nil {
		ContextLogger.Error(err)
		return err
	}

	return nil
}

func (c *Scanner) GetRecommendations(screener, exchange, symbol, interval string) (RecommendSummary, error) {
	indicators := concatAppend([][]string{recommendsList})
	payload, err := c.PrepareData(fmt.Sprintf("%s:%s", exchange, symbol), interval, indicators)
	if err != nil {
		ContextLogger.Error(err, exchange, symbol)
		return RecommendSummary{}, err
	}
	r, err := c.client.do("POST", string(payload), false)
	if err != nil {
		ContextLogger.Errorf("Exchange (%s) or symbol (%s) not found %v", exchange, symbol, err)
		return RecommendSummary{}, err
	}
	err = json.Unmarshal(r, &c.data)
	if err != nil {
		ContextLogger.Error(err)
		return RecommendSummary{}, err
	}

	err = c.PrepareRecommendations()
	if err != nil {
		ContextLogger.Error(err)
	}

	if c.client.debug {
		fmt.Println(c.recommendSummary.Recommend.Summary,
			c.recommendSummary.Recommend.Oscillators,
			c.recommendSummary.Recommend.MovingAverages,
		)
	}

	return RecommendSummary{
		Recommend: c.recommendSummary.Recommend,
	}, nil
}

func (c *Scanner) GetIchimoku(screener, exchange, symbol, interval string) (Ichimoku string, value float64, err error) {
	indicators := []string{"Rec.Ichimoku", "Ichimoku.BLine"}
	payload, err := c.PrepareData(fmt.Sprintf("%s:%s", exchange, symbol), interval, indicators)
	if err != nil {
		ContextLogger.Error(err, exchange, symbol)
		return Ichimoku, value, err
	}
	r, err := c.client.do("POST", string(payload), false)
	if err != nil {
		ContextLogger.Errorf("Exchange (%s) or symbol (%s) not found %v", exchange, symbol, err)
		return Ichimoku, value, err
	}
	err = json.Unmarshal(r, &c.data)
	if err != nil {
		ContextLogger.Error(err)
		return Ichimoku, value, err
	}

	Ichimoku = c.ComputeSimple(c.data.Data[0].Data[0])
	value = c.data.Data[0].Data[1]

	if c.client.debug {
		fmt.Println(Ichimoku, value)
	}

	return Ichimoku, value, nil
}

func (c *Scanner) GetAnalysis(screener, exchange, symbol, interval string) (RecommendSummary, error) {
	indicators := concatAppend([][]string{recommendsList, oscillatorsList, maList, maSimpleList})
	payload, err := c.PrepareData(fmt.Sprintf("%s:%s", exchange, symbol), interval, indicators)
	if err != nil {
		ContextLogger.Error(err, exchange, symbol)
		return RecommendSummary{}, err
	}
	r, err := c.client.do("POST", string(payload), false)
	if err != nil {
		ContextLogger.Errorf("Exchange (%s) or symbol (%s) not found %v", exchange, symbol, err)
		return RecommendSummary{}, err
	}
	err = json.Unmarshal(r, &c.data)
	if err != nil {
		ContextLogger.Error(err)
		return RecommendSummary{}, err
	}

	oscillatorsCounter := map[string]int{"BUY": 0, "SELL": 0, "NEUTRAL": 0}
	maCounter := map[string]int{"BUY": 0, "SELL": 0, "NEUTRAL": 0}
	computedOscillators := map[string]string{}
	computedMa := map[string]string{}

	//
	// RECOMMENDATIONS
	err = c.PrepareRecommendations()
	if err != nil {
		ContextLogger.Error(err)
	}

	if c.client.debug {
		fmt.Println(c.recommendSummary.Recommend.Summary,
			c.recommendSummary.Recommend.Oscillators,
			c.recommendSummary.Recommend.MovingAverages,
		)
	}

	// TODO: Add checking for None

	// OSCILLATORS
	// RSI (14)
	computedOscillators["RSI"] = c.ComputeRSI(c.data.Data[0].Data[3], c.data.Data[0].Data[4])
	oscillatorsCounter[computedOscillators["RSI"]] += 1

	// Stoch %K
	computedOscillators["STOCH.K"] = c.ComputeStoch(c.data.Data[0].Data[5], c.data.Data[0].Data[6], c.data.Data[0].Data[7], c.data.Data[0].Data[8])
	oscillatorsCounter[computedOscillators["STOCH.K"]] += 1

	// CCI (20)
	computedOscillators["CCI"] = c.ComputeCCI20(c.data.Data[0].Data[9], c.data.Data[0].Data[10])
	oscillatorsCounter[computedOscillators["CCI"]] += 1

	// ADX (14)
	computedOscillators["ADX"] = c.ComputeADX(c.data.Data[0].Data[11], c.data.Data[0].Data[12], c.data.Data[0].Data[13], c.data.Data[0].Data[14], c.data.Data[0].Data[15])
	oscillatorsCounter[computedOscillators["ADX"]] += 1

	// AO
	computedOscillators["AO"] = c.ComputeAO(c.data.Data[0].Data[16], c.data.Data[0].Data[17])
	oscillatorsCounter[computedOscillators["AO"]] += 1

	// Mom (10)
	computedOscillators["Mom"] = c.ComputeMOM(c.data.Data[0].Data[18], c.data.Data[0].Data[19])
	oscillatorsCounter[computedOscillators["Mom"]] += 1

	// MACD
	computedOscillators["MACD"] = c.ComputeMACD(c.data.Data[0].Data[20], c.data.Data[0].Data[21])
	oscillatorsCounter[computedOscillators["MACD"]] += 1

	// Stoch RSI
	computedOscillators["Stoch.RSI"] = c.ComputeSimple(c.data.Data[0].Data[22])
	oscillatorsCounter[computedOscillators["Stoch.RSI"]] += 1

	// W%R
	computedOscillators["W%R"] = c.ComputeSimple(c.data.Data[0].Data[24])
	oscillatorsCounter[computedOscillators["W%R"]] += 1

	// BBP
	computedOscillators["BBP"] = c.ComputeSimple(c.data.Data[0].Data[26])
	oscillatorsCounter[computedOscillators["BBP"]] += 1

	// UO
	computedOscillators["UO"] = c.ComputeSimple(c.data.Data[0].Data[28])
	oscillatorsCounter[computedOscillators["UO"]] += 1

	// MOVING AVERAGES
	Close := c.data.Data[0].Data[30]
	maListCounter := 0
	for index := 31; index < 43; index++ {
		if &c.data.Data[0].Data[index] != nil {
			computedMa[maList[maListCounter]] = c.ComputeMA(c.data.Data[0].Data[index], Close)
			maCounter[computedMa[maList[maListCounter]]] += 1
			maListCounter += 1
		}
	}

	// MOVING AVERAGES, pt 2
	// ICHIMOKU
	computedMa["Ichimoku"] = c.ComputeSimple(c.data.Data[0].Data[43])
	maCounter[computedMa["Ichimoku"]] += 1

	// VWMA
	computedMa["VWMA"] = c.ComputeSimple(c.data.Data[0].Data[45])
	maCounter[computedMa["VWMA"]] += 1

	// HullMA (9)
	computedMa["HullMA"] = c.ComputeSimple(c.data.Data[0].Data[47])
	maCounter[computedMa["HullMA"]] += 1

	if c.client.debug {
		fmt.Printf("Summary - \"RECOMMENDATION\": %s, \"BUY\": %d, \"SELL\": %d, \"NEUTRAL\": %d\n",
			c.recommendSummary.Recommend.Summary,
			oscillatorsCounter["BUY"]+maCounter["BUY"],
			oscillatorsCounter["SELL"]+maCounter["SELL"],
			oscillatorsCounter["NEUTRAL"]+maCounter["NEUTRAL"],
		)
		fmt.Printf("Oscillators - \"RECOMMENDATION\": %s, \"BUY\": %d, \"SELL\": %d, \"NEUTRAL\": %d, \"COMPUTE\": %s\n",
			c.recommendSummary.Recommend.Oscillators,
			oscillatorsCounter["BUY"],
			oscillatorsCounter["SELL"],
			oscillatorsCounter["NEUTRAL"],
			computedOscillators,
		)
		fmt.Printf("MovingAverages - \"RECOMMENDATION\": %s, \"BUY\": %d, \"SELL\": %d, \"NEUTRAL\": %d, \"COMPUTE\": %s\n",
			c.recommendSummary.Recommend.MovingAverages,
			maCounter["BUY"],
			maCounter["SELL"],
			maCounter["NEUTRAL"],
			computedMa,
		)
	}

	return RecommendSummary{
		Recommend:    c.recommendSummary.Recommend,
		BuyCount:     oscillatorsCounter["BUY"] + maCounter["BUY"],
		SellCount:    oscillatorsCounter["SELL"] + maCounter["SELL"],
		NeutralCount: oscillatorsCounter["NEUTRAL"] + maCounter["NEUTRAL"],
	}, nil
}
