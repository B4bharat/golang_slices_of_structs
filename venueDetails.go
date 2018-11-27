package main

import (
	"fmt"
)

func main() {
	// CategoryDetails : ticket category details
	type CategoryDetails struct {
		PriceCode        string `json:"PriceCode"`
		PriceDescription string `json:"PriceDescription"`
		AlternateName    string `json:"alternateName"`
	}

	// ShowTimeDetails : show time details
	type ShowTimeDetails struct {
		SessionID       string            `json:"SessionId"`
		CutOffDateTime  string            `json:"CutOffDateTime"`
		ShowTimeDisplay string            `json:"ShowTimeDisplay"`
		SessionRealShow string            `json:"sessionRealShow"`
		ArrCategory     []CategoryDetails `json:"arrCategory"`
	}

	// ShowDateDetails : show date details
	type ShowDateDetails struct {
		ShowDateCode    string            `json:"ShowDateCode"`
		ShowDateDisplay string            `json:"ShowDateDisplay"`
		ArrShowTimes    []ShowTimeDetails `json:"arrShowTimes"`
	}

	// VenueDetails : festival event venue info
	type VenueDetails struct {
		VenueCode    string            `json:"VenueCode"`
		VenueName    string            `json:"VenueName"`
		ArrShowDates []ShowDateDetails `json:"arrShowDates"`
	}

	// FestivalEventDetails : festival event details schema
	type FestivalEventDetails struct {
		ArrVenues []VenueDetails `json:"arrVenues"`
	}

	uapiVenueDetails := VenueDetails{
		"RUTB",
		"Rann Utsav Tent City: Bhuj",
		[]ShowDateDetails{
			{
				"20181124",
				"Saturday, 24 Nov",
				[]ShowTimeDetails{
					{
						"10493",
						"201811211230",
						"12:30 PM",
						"24-11-2018 12:30:00",
						[]CategoryDetails{
							{
								"T011",
								"1N 2D NonAC SwisCottage(Rs. 12980.00)",
								"1N 2D NonAC SwisCottage",
							},
							{
								"T010",
								"1N 2D DeluxAC SwisCottage(Rs. 16756.00)",
								"1N 2D DeluxAC SwisCottage",
							},
						},
					},
				},
			},
			{
				"20181125",
				"Sunday, 25 Nov",
				[]ShowTimeDetails{
					{
						"10494",
						"201811221230",
						"12:30 PM",
						"25-11-2018 12:30:00",
						[]CategoryDetails{
							{
								"T011",
								"1N 2D NonAC SwisCottage(Rs. 12980.00)",
								"Half Ticket",
							},
							{
								"T010",
								"1N 2D DeluxAC SwisCottage(Rs. 16756.00)",
								"Half Ticket",
							},
						},
					},
				},
			},
		},
	}

	type dateRangeDetails struct {
		startDate string
		endDate   string
	}

	type ticketGroupsDetails struct {
		name           string
		description    string
		shouldValidate bool
		tickets        []string
		dates          []string
		venueCodes     []string
		times          []string
		dateRange      dateRangeDetails
	}

	ticketGroups := make([]ticketGroupsDetails, 1)
	etCode := "ET00085207"

	showDates := uapiVenueDetails.ArrShowDates
	// 1st loop showdates
	for i := 0; i < len(showDates); i++ {
		currShowTime := showDates[i].ArrShowTimes
		// - save showDateCode in a variable
		// showDateCode := showDates[i].ShowDateCode

		// for every showdate, loop through showTimes
		for j := 0; j < len(currShowTime); j++ {
			// - save sessionId in a variable
			sessionId := currShowTime[j].SessionID

			currCategoryDetail := currShowTime[j].ArrCategory
			// for every showtime, loop through categoryDetails
			for k := 0; k < len(currCategoryDetail); k++ {
				/*
					- showDateCode into dates array && check 'dateRange' struct if showDateCode <= startDate then startDate = showDateCode else if >= endDate then endDate = showDateCode
				*/
				// - save priceCode in a variable
				priceCode := currCategoryDetail[k].PriceCode

				// - check if alternateName is present, if not grab PriceDescription and push it into ticketGroupNames
				ticketGroupName := ""
				if currCategoryDetail[k].AlternateName != "" {
					ticketGroupName = currCategoryDetail[k].AlternateName
				} else {
					ticketGroupName = currCategoryDetail[k].PriceDescription
				}
				tgNameExists := false

				// if ticketGroupNames doesn't exists, then create a 'name' property in the 'tickGroups' struct and add the curr value
				fmt.Println(ticketGroupName)
				for idx, _ := range ticketGroups {
					if ticketGroups[idx].name == ticketGroupName {
						/*
							- switch (name) - set 'description', based on equivalent 'name'
								name: description
								'General Festival Pass': '18 Years and above',
								'Student': '(Aged 18- 25) | A student has to carry a valid Student ID',
								'Half Ticket': 'Eligible for children between the ages 5-17 years of age, accompanied by an adult. Screenings taking place on the 27th and 28th of October'
									&& set shouldValidate to true
								'Press': ''
						*/
						switch ticketGroupName {
						case "General Festival Pass":
							ticketGroups[idx].description = "18 Years and above"
						case "Student":
							ticketGroups[idx].description = "(Aged 18- 25) | A student has to carry a valid Student ID"
						case "Half Ticket":
							ticketGroups[idx].description = "Eligible for children between the ages 5-17 years of age, accompanied by an adult. Screenings taking place on the 27th and 28th of October"
							ticketGroups[idx].shouldValidate = true
						default:

						}

						ticketCode := etCode + "-" + sessionId + "-" + priceCode
						ticketGroups[idx].tickets = append(ticketGroups[idx].tickets, ticketCode)

						// No need to create a new ticket group
						tgNameExists = true
						break
					}
				}

				if tgNameExists == false {

					ticketCode := etCode + "-" + sessionId + "-" + priceCode
					newTickets := make([]string, 1)
					newTickets = append(newTickets, ticketCode)

					newTicketGroup := ticketGroupsDetails{name: ticketGroupName, description: "", tickets: newTickets}
					ticketGroups = append(ticketGroups, newTicketGroup)
				}
			}
		}
	}

	// removes the first element from the slice
	ticketGroups = append(ticketGroups[:0], ticketGroups[1:]...)
	fmt.Println(ticketGroups)
}