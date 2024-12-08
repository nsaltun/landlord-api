package model

type CalculationResponse struct {
	OnePlusOneCount   []Flat
	TwoPlusOneCount   int
	ThreePlusOneCount int
}

type CalculationRequest struct {
	Emsal               float64 `json:"emsal"`
	LandSquareMeter     float64 `json:"landSquareMeter"`
	MaxAllowedFlatCount int     `json:"maxAllowedFlatCount"`
	ExtendFactor        float64 `json:"extendFactor"`
}

type Flat struct {
	//Kaç tane daire çıkar
	Count int
	//Bir daire kaç m2'dir
	Size float64
}
