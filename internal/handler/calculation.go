package handler

import (
	"math"
	"net/http"

	"github.com/nsaltun/landlord-api/internal/model"
	"github.com/nsaltun/landlord-api/pkg/middlewares"
)

const (
	MaxOnePlusOneSquareMeter = 45
	MinOnePlusOneSquareMeter = 35
	MaxTwoPlusOneSquareMeter = 100
	MinTwoPlusOneSquareMeter = 60
)

type LandCalculation interface {
	HandleFieldCalculation(ctx *middlewares.HttpContext)
}

type landCalculation struct {
}

func NewLandCalculationHandler() LandCalculation {
	return &landCalculation{}
}

func (l *landCalculation) HandleFieldCalculation(ctx *middlewares.HttpContext) {
	var req model.CalculationRequest
	err := ctx.CastRequestBody(&req)
	if err != nil {
		ctx.JSONErrResp(err)
		return
	}

	resp := &model.CalculationResponse{
		OnePlusOneCount: CalculateOnePlusOneCount(req.LandSquareMeter, req.Emsal, req.ExtendFactor),
		TwoPlusOneCount: CalculateTwoPlusOneCount(req.LandSquareMeter, req.Emsal, req.ExtendFactor),
		// ThreePlusOneCount: 8,
	}
	ctx.JSON(http.StatusOK, resp)
}

func CalculateOnePlusOneCount(landSquareMeter float64, emsal float64, extendFactor float64) []model.Flat {
	netConstructionArea := landSquareMeter * emsal * extendFactor

	flats := []model.Flat{}
	//Calculate min size flats
	flatSizeMin := float64(MinOnePlusOneSquareMeter)
	countMin := netConstructionArea / flatSizeMin
	if !isExactInteger(countMin) {
		//Daire sayısını net olarak dönmek istediğimiz için daire başına m2 yi yükseltiyoruz burada
		countMin = math.Floor(countMin)
		flatSizeMin = netConstructionArea / countMin
	}
	flats = append(flats, model.Flat{Count: int(countMin), Size: math.Round(flatSizeMin)})

	//Calculate max size flats
	countMax := netConstructionArea / float64(MaxOnePlusOneSquareMeter)
	//Daire sayısını artırıyoruz sayıyı virgüllü dönmek istemediğimiz için
	if !isExactInteger(countMax) {
		//Calculate by ceil
		countCeil := math.Ceil(countMax)
		flatSizeCeil := netConstructionArea / countCeil
		flats = append(flats, model.Flat{
			Count: int(countCeil),
			Size:  math.Round(flatSizeCeil),
		})

		//Calculate by floor
		countFloor := math.Floor(countMax)
		flatSizeFloor := netConstructionArea / countFloor
		flats = append(flats, model.Flat{
			Count: int(countFloor),
			Size:  math.Round(flatSizeFloor),
		})

	} else {
		flats = append(flats, model.Flat{
			Count: int(math.Ceil(countMax)),
			Size:  math.Round(float64(MaxOnePlusOneSquareMeter)),
		})
	}

	return flats
}

func CalculateTwoPlusOneCount(landSquareMeter float64, emsal float64, extendFactor float64) []model.Flat {
	netConstructionArea := landSquareMeter * emsal * extendFactor

	flats := []model.Flat{}
	//Calculate min size flats
	flatSizeMin := float64(MinTwoPlusOneSquareMeter)
	countMin := netConstructionArea / flatSizeMin
	if !isExactInteger(countMin) {
		//Daire sayısını net olarak dönmek istediğimiz için daire başına m2 yi yükseltiyoruz burada
		countMin = math.Floor(countMin)
		flatSizeMin = netConstructionArea / countMin
	}
	flats = append(flats, model.Flat{Count: int(countMin), Size: math.Round(flatSizeMin)})

	//Calculate max size flats
	countMax := netConstructionArea / float64(MaxTwoPlusOneSquareMeter)
	//Daire sayısını artırıyoruz sayıyı virgüllü dönmek istemediğimiz için
	if !isExactInteger(countMax) {
		//Calculate by ceil
		countCeil := math.Ceil(countMax)
		flatSizeCeil := netConstructionArea / countCeil
		flats = append(flats, model.Flat{
			Count: int(countCeil),
			Size:  math.Round(flatSizeCeil),
		})

		//Calculate by floor
		countFloor := math.Floor(countMax)
		flatSizeFloor := netConstructionArea / countFloor
		flats = append(flats, model.Flat{
			Count: int(countFloor),
			Size:  math.Round(flatSizeFloor),
		})

	} else {
		flats = append(flats, model.Flat{
			Count: int(math.Ceil(countMax)),
			Size:  float64(MaxTwoPlusOneSquareMeter),
		})
	}

	return flats
}

func isExactInteger(f float64) bool {
	return f == float64(int32(f))
}
