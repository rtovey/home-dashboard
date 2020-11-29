package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	c "github.com/rtovey/astro/common"
	"github.com/rtovey/astro/lunar"
	astro "github.com/rtovey/astro/lunar"
	"github.com/rtovey/astro/solar"
)

func main() {
	http.HandleFunc("/", homePage)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	lunarPhase := astro.Phase(now)

	loc, _ := time.LoadLocation("Europe/London")
	observer := c.Observer{
		Latitude:  51,
		Longitude: 0,
		Location:  loc,
	}

	moonRise, moonSet := getNextMoonRiseSetTime(observer, now)

	sunRise, sunSet := getNextSunRiseSetTime(observer, now)

	fmt.Fprintf(w, "Current lunar phase: %.0f%%\n"+
		"Next moonrise: %v\n"+
		"Next moonset: %v\n\n"+
		"Next sunrise: %v\n"+
		"Next sunset: %v\n",
		lunarPhase*100, moonRise, moonSet, sunRise, sunSet)
}

func getNextMoonRiseSetTime(observer c.Observer, now time.Time) (time.Time, time.Time) {
	lunarRiseSetTimeForToday := lunar.RiseSetTime(observer, now)
	lunarRiseSetTimeForTomorrow := lunar.RiseSetTime(observer, now.Add(24*time.Hour))
	moonRise := lunarRiseSetTimeForToday.Rise
	moonSet := lunarRiseSetTimeForToday.Set

	if moonRise.Before(now) {
		moonRise = lunarRiseSetTimeForTomorrow.Rise
	}

	if moonSet.Before(moonRise) {
		moonSet = lunarRiseSetTimeForTomorrow.Set
	}

	return moonRise, moonSet
}

func getNextSunRiseSetTime(observer c.Observer, now time.Time) (time.Time, time.Time) {
	solarRiseSetTimeForToday := solar.RiseSetTime(observer, now)
	solarRiseSetTimeForTomorrow := solar.RiseSetTime(observer, now.Add(24*time.Hour))
	sunRise := solarRiseSetTimeForToday.Rise
	sunSet := solarRiseSetTimeForToday.Set

	if sunRise.Before(now) {
		sunRise = solarRiseSetTimeForTomorrow.Rise
	}

	if sunSet.Before(now) {
		sunSet = solarRiseSetTimeForTomorrow.Set
	}

	return sunRise, sunSet
}
