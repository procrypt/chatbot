package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type LatLong struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Bounds struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"bounds"`
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

type Zomato struct {
	ResultsFound int `json:"results_found"`
	ResultsStart int `json:"results_start"`
	ResultsShown int `json:"results_shown"`
	Restaurants  []struct {
		Restaurant struct {
			R struct {
				ResID int `json:"res_id"`
			} `json:"R"`
			Apikey   string `json:"apikey"`
			ID       string `json:"id"`
			Name     string `json:"name"`
			URL      string `json:"url"`
			Location struct {
				Address         string `json:"address"`
				Locality        string `json:"locality"`
				City            string `json:"city"`
				CityID          int    `json:"city_id"`
				Latitude        string `json:"latitude"`
				Longitude       string `json:"longitude"`
				Zipcode         string `json:"zipcode"`
				CountryID       int    `json:"country_id"`
				LocalityVerbose string `json:"locality_verbose"`
			} `json:"location"`
			SwitchToOrderMenu  int           `json:"switch_to_order_menu"`
			Cuisines           string        `json:"cuisines"`
			AverageCostForTwo  int           `json:"average_cost_for_two"`
			PriceRange         int           `json:"price_range"`
			Currency           string        `json:"currency"`
			Offers             []interface{} `json:"offers"`
			OpentableSupport   int           `json:"opentable_support"`
			IsZomatoBookRes    int           `json:"is_zomato_book_res"`
			MezzoProvider      string        `json:"mezzo_provider"`
			IsBookFormWebView  int           `json:"is_book_form_web_view"`
			BookFormWebViewURL string        `json:"book_form_web_view_url"`
			BookAgainURL       string        `json:"book_again_url"`
			Thumb              string        `json:"thumb"`
			UserRating         struct {
				AggregateRating string `json:"aggregate_rating"`
				RatingText      string `json:"rating_text"`
				RatingColor     string `json:"rating_color"`
				Votes           string `json:"votes"`
			} `json:"user_rating"`
			PhotosURL                   string        `json:"photos_url"`
			MenuURL                     string        `json:"menu_url"`
			FeaturedImage               string        `json:"featured_image"`
			HasOnlineDelivery           int           `json:"has_online_delivery"`
			IsDeliveringNow             int           `json:"is_delivering_now"`
			IncludeBogoOffers           bool          `json:"include_bogo_offers"`
			Deeplink                    string        `json:"deeplink"`
			IsTableReservationSupported int           `json:"is_table_reservation_supported"`
			HasTableBooking             int           `json:"has_table_booking"`
			EventsURL                   string        `json:"events_url"`
			EstablishmentTypes          []interface{} `json:"establishment_types"`
		} `json:"restaurant"`
	} `json:"restaurants"`
}

var (
	msg         []string
	qURL        = "https://sqs.us-east-1.amazonaws.com/263331787521/procrypt"
	svcSqs      = sqs.New(session.New())
	svcSns      = sns.New(session.New())
	topicARN    = "arn:aws:sns:us-east-1:263331787521:procrypt"
	reservation = map[string]string{}
)

type Info struct {
	Phone    string            `json:"phone"`
	Location string            `json:"location"`
	Time     string            `json:"time"`
	Data     map[string]string `json:"data"`
}

func dynamo(time, phone, location string, data map[string]string) {
	svc := dynamodb.New(session.New())
	info := Info{
		Phone:    phone,
		Location: location,
		Time:     time,
		Data:     data,
	}

	av, _ := dynamodbattribute.MarshalMap(info)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("reservation"),
	}
	svc.PutItem(input)
	reservation = data
}

func getLatLong(s string) (lat, lng string, err error) {
	location := LatLong{}

	s = strings.Replace(s, " ", "+", -1)
	url := "https://maps.googleapis.com/maps/api/geocode/json?address=" + s + "&key=<key>"

	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", "", err
	}
	res, er := client.Do(req)
	if er != nil {
		return "", "", err
	}
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &location)

	return strconv.FormatFloat(location.Results[0].Geometry.Location.Lat, 'f', 6, 64), strconv.FormatFloat(location.Results[0].Geometry.Location.Lng, 'f', 6, 64), nil
}

func sendSMS() {
	count := 1
	findings := ""
	for k, v := range reservation {
		s := strconv.Itoa(count)
		findings += s + ")" + k + " " + v + " \n"
		count++
	}

	msgToSend := "Hello! Here are my Japanese restaurant suggestions for 2 people for today at, " + msg[0] + " in " + msg[2] + "\n" + findings + "Enjoy your meal!"

	params := &sns.PublishInput{
		Message:     aws.String(msgToSend),
		PhoneNumber: aws.String(msg[1]),
	}
	_, err := svcSns.Publish(params)
	if err != nil {
		log.Fatal(err)
	}
	// Empty the Reservation Map
	reservation = map[string]string{}
	// Empty the data Map
	msg = []string{}
}

func data(cuisine string, lat, lng string) map[string]string {
	z := Zomato{}
	m := make(map[string]string)

	url := "https://developers.zomato.com/api/v2.1/search?q=" + cuisine + "&lat=" + lat + "&lon=" + lng
	client := http.Client{}
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Add("user-key", "<key>")
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &z)
	for i := 0; i < len(z.Restaurants); i++ {
		if strings.Contains(z.Restaurants[i].Restaurant.Cuisines, cuisine) {
			m[z.Restaurants[i].Restaurant.Name] = z.Restaurants[i].Restaurant.Location.Address
		}
	}
	return m
}

func sqsWorker(event events.SQSEvent) {
	if len(event.Records) != 0 {
		for i := 0; i < len(event.Records); i++ {
			msg = append(msg, *event.Records[i].MessageAttributes["Time"].StringValue)
			msg = append(msg, *event.Records[i].MessageAttributes["Phone"].StringValue)
			msg = append(msg, *event.Records[i].MessageAttributes["Location"].StringValue)
		}
		lat, lng, _ := getLatLong(msg[2])
		info := data("Japanese", lat, lng)
		if len(info) < 3 {
			dynamo(msg[0], msg[1], msg[2], info)
		} else {
			count := 0
			m := make(map[string]string)
			for k, v := range info {
				if count < 4 {
					m[k] = v
					count++
				}
			}
			dynamo(msg[0], msg[1], msg[2], m)
			sendSMS()
		}
	}
}

func main() {
	lambda.Start(sqsWorker)
}
