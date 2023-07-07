package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type helper struct {
	gradPlato        float64
	desiredCo2       float64
	beerGravity      float64
	requiredPressure float64
	kegSize          float64
}

func newHelper() *helper {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("--- What is the the keg size? ---")
	input, err := reader.ReadString('\n')
	input = strings.Replace(input, "\n", "", -1)
	if err != nil {
		panic(err)
	}
	kegSize, err := strconv.ParseFloat(input, 64)
	if err != nil {
		panic(err)
	}

	fmt.Println("--- What is the current temperature? ---")
	input, err = reader.ReadString('\n')
	input = strings.Replace(input, "\n", "", -1)
	if err != nil {
		panic(err)
	}
	temperature, err := strconv.ParseFloat(input, 64)
	if err != nil {
		panic(err)
	}

	fmt.Println("--- What is your desired CO2 amount? ---")
	input, err = reader.ReadString('\n')
	input = strings.Replace(input, "\n", "", -1)
	if err != nil {
		panic(err)
	}
	desiredCo2, err := strconv.ParseFloat(input, 64)
	if err != nil {
		panic(err)
	}

	fmt.Println("--- How much Grad Plato did it have? ---")
	input, err = reader.ReadString('\n')
	input = strings.Replace(input, "\n", "", -1)
	if err != nil {
		panic(err)
	}
	gradPlato, err := strconv.ParseFloat(input, 64)
	if err != nil {
		panic(err)
	}

	fmt.Println("--- What is your gravity after fermentation? (Use 2 for estimation before fermentation) ---")
	input, err = reader.ReadString('\n')
	input = strings.Replace(input, "\n", "", -1)
	if err != nil {
		panic(err)
	}
	gravity, err := strconv.ParseFloat(input, 64)
	if err != nil {
		panic(err)
	}

	return &helper{
		gradPlato:        gradPlato,
		desiredCo2:       desiredCo2,
		beerGravity:      gravity,
		requiredPressure: desiredCo2/(10*math.Pow(math.E, -10.73797+(2617.25/(temperature+273.15)))) - 1.013,
		kegSize:          kegSize,
	}
}

func (h *helper) neededCo2(beerAmount, estimatedWortAmount float64) float64 {
	return h.neededCo2Gas(beerAmount, estimatedWortAmount) + h.neededCo2Beer(beerAmount)
}

func (h *helper) neededCo2Gas(beerAmount, estimatedWortAmount float64) float64 {
	leftOverAir := h.kegSize - beerAmount - estimatedWortAmount
	if leftOverAir < 0 {
		leftOverAir = 0
	}
	return leftOverAir * 1.83 * h.requiredPressure
}

func (h *helper) neededCo2Beer(beerAmount float64) float64 {
	return beerAmount * 1.66 * h.requiredPressure
}

func (h *helper) neededSugar(beerAmount, estimatedWortAmount float64) float64 {
	return h.neededCo2(beerAmount, estimatedWortAmount) * 2 * 0.957
}

func (h *helper) sugarToGradPlato(gramsPerLiter float64) float64 {
	return gramsPerLiter
}

func (h *helper) realGravity() float64 {
	return 0.1808*h.gradPlato + 0.8192*h.beerGravity
}

func (h *helper) alcoholContent() float64 {
	return math.Round((1/0.795*(h.gradPlato-h.realGravity())/(2.0665-0.010665*h.gradPlato))*100) / 100
}

func (h *helper) neededWort(beerAmount float64) int {
	estimatedWortAmount := h.neededSugar(beerAmount, 0) / ((h.gradPlato - h.realGravity()) * 10)
	return int(math.Round((h.neededSugar(beerAmount, estimatedWortAmount) / ((h.gradPlato - h.realGravity()) * 10)) * 1000))

}

func main() {
	h := newHelper()
	fmt.Printf("--- If you fill this keg completely you would roughly need %dml of wort ---\n", h.neededWort(h.kegSize))
	fmt.Println("--- So make sure to leave enough space for your wort in your keg!! ---")
	fmt.Println()

	fmt.Println("--- How much beer do you want to put in your keg? ---")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	input = strings.Replace(input, "\n", "", -1)
	if err != nil {
		panic(err)
	}
	beerAmount, err := strconv.ParseFloat(input, 64)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Printf("For this keg you need %d ml of Wort\n", h.neededWort(beerAmount))
	fmt.Printf("The keg needs %.1f of pressure\n", h.requiredPressure)
	fmt.Printf("Your alcohol content will be roughly %.2f\n", h.alcoholContent())
}
