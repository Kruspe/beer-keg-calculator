package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	kegSize = 18
)

type helper struct {
	beerAmount       float64
	gradPlato        float64
	desiredCo2       float64
	beerGravity      float64
	requiredPressure float64
}

func newHelper() *helper {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("--- How much beer is in the keg? ---")
	input, err := reader.ReadString('\n')
	input = strings.Replace(input, "\n", "", -1)
	if err != nil {
		panic(err)
	}
	beerAmount, err := strconv.ParseFloat(input, 64)
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

	fmt.Println("--- What is your gravity after fermentation? ---")
	input, err = reader.ReadString('\n')
	input = strings.Replace(input, "\n", "", -1)
	if err != nil {
		panic(err)
	}
	gravity, err := strconv.ParseFloat(input, 64)
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

	return &helper{
		beerAmount:       beerAmount,
		gradPlato:        gradPlato,
		desiredCo2:       desiredCo2,
		beerGravity:      gravity,
		requiredPressure: desiredCo2/(10*math.Pow(math.E, -10.73797+(2617.25/(temperature+273.15)))) - 1.013,
	}
}

func (h *helper) neededCo2() float64 {
	return h.neededCo2Gas() + h.neededCo2Beer()
}

func (h *helper) neededCo2Gas() float64 {
	return (kegSize - h.beerAmount) * 1.83 * h.requiredPressure
}

func (h *helper) neededCo2Beer() float64 {
	return h.beerAmount * 1.66 * h.requiredPressure
}

func (h *helper) neededSugar() float64 {
	return h.neededCo2() * 2 * 0.957
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

func (h *helper) neededWort() int {
	return int(math.Round((h.neededSugar() / ((h.gradPlato - h.realGravity()) * 10)) * 1000))

}

func main() {
	h := newHelper()
	fmt.Println()
	fmt.Printf("For this keg you need %d ml of Wort\n", h.neededWort())
	fmt.Printf("The keg needs %.1f of pressure\n", h.requiredPressure)
	fmt.Printf("Your alcohol content will be roughly %.2f\n", h.alcoholContent())
}
