package tvscanner

import (
	"errors"
)

const (
	Buy        = "BUY"
	StrongBuy  = "STRONG_BUY"
	Sell       = "SELL"
	StrongSell = "STRONG_SELL"
	Neutral    = "NEUTRAL"
	Error      = "ERROR"
)

// ComputeRecommend return "STRONG_BUY", "BUY", "NEUTRAL", "SELL", "STRONG_SELL", or "ERROR"
func (c *Scanner) ComputeRecommend(value float64) (string, error) {
	switch {
	case value >= -1 && value < -0.5:
		return StrongSell, nil
	case value >= -0.5 && value < 0:
		return Sell, nil
	case value == 0:
		return Neutral, nil
	case value > 0 && value <= 0.5:
		return Buy, nil
	case value > 0.5 && value <= 1:
		return StrongBuy, nil
	default:
		return "", errors.New("Failed ComputeRecommend ")
	}
}

// ComputeSimple return "BUY", "SELL", or "NEUTRAL"
func (c *Scanner) ComputeSimple(value float64) string {
	switch {
	case value == -1:
		return Sell
	case value == 1:
		return Buy
	default:
		return Neutral
	}
}

// ComputeMA
func (c *Scanner) ComputeMA(ma, close float64) string {
	if ma < close {
		return Buy
	} else if ma > close {
		return Sell
	}
	return Neutral
}

// ComputeRSI
func (c *Scanner) ComputeRSI(rsi, rsi1 float64) string {
	if rsi < 30 && rsi1 > rsi {
		return Buy
	} else if rsi > 70 && rsi1 < rsi {
		return Sell
	}
	return Neutral

}

// ComputeStoch
func (c *Scanner) ComputeStoch(k, d, k1, d1 float64) string {
	if k < 20 && d < 20 && k > d && k1 < d1 {
		return Buy
	} else if k > 80 && d > 80 && k < d && k1 > d1 {
		return Sell
	}
	return Neutral
}

// ComputeCCI20
func (c *Scanner) ComputeCCI20(cci20, cci201 float64) string {
	if cci20 < -100 && cci20 > cci201 {
		return Buy
	} else if cci20 > 100 && cci20 < cci201 {
		return Sell
	}
	return Neutral
}

// ComputeADX
func (c *Scanner) ComputeADX(adx, adxpdi, adxndi, adxpdi1, adxndi1 float64) string {
	if adx > 20 && adxpdi1 < adxndi1 && adxpdi > adxndi {
		return Buy
	} else if adx > 20 && adxpdi1 > adxndi1 && adxpdi < adxndi {
		return Sell
	}
	return Neutral
}

// ComputeAO
func (c *Scanner) ComputeAO(ao, ao1 float64) string {
	if ao > 0 && ao1 < 0 || ao > 0 && ao1 > 0 && ao > ao1 {
		return Buy
	} else if ao < 0 && ao1 > 0 || ao < 0 && ao1 < 0 && ao < ao1 {
		return Sell
	}
	return Neutral
}

// ComputeMOM
func (c *Scanner) ComputeMOM(mom, mom1 float64) string {
	if mom > mom1 {
		return Buy
	} else if mom < mom1 {
		return Sell
	}
	return Neutral
}

// ComputeMACD
func (c *Scanner) ComputeMACD(macd, signal float64) string {
	if macd > signal {
		return Buy
	} else if macd < signal {
		return Sell
	}
	return Neutral
}

// ComputeBBBuy
func (c *Scanner) ComputeBBBuy(close, bblower float64) string {
	if close < bblower {
		return Buy
	}
	return Neutral
}

// ComputeBBSell
func (c *Scanner) ComputeBBSell(close, bbupper float64) string {
	if close > bbupper {
		return Sell
	}
	return Neutral
}

// ComputePSAR
func (c *Scanner) ComputePSAR(psar, open float64) string {
	if psar < open {
		return Buy
	} else if psar > open {
		return Sell
	}
	return Neutral
}
