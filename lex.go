package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	content string
	data []string
	qURL = "https://sqs.us-east-1.amazonaws.com/263331787521/procrypt"
	svcSqs = sqs.New(session.New())

)

func sendsSQS(msg []string)  {
	svcSqs.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Location": &sqs.MessageAttributeValue{
				DataType: aws.String("String"),
				StringValue: aws.String(data[0]),
			},
			"Time": &sqs.MessageAttributeValue {
				DataType: aws.String("String"),
				StringValue: aws.String(data[1]),
			},
			"Phone": &sqs.MessageAttributeValue {
				DataType: aws.String("String"),
				StringValue: aws.String(data[2]),
			},
		},
		MessageBody: aws.String("Information about reservation."),
		QueueUrl: &qURL,

	})
	data = []string{}
}

type BotRequest struct {
	Messages []Message `json:"messages,omitempty"`
}

type BotResponce struct {
	Messages []Message `json:"messages,omitempty"`
}

type Message struct {
	Type         string        `json:"type,omitempty"`
	Unstructured *UnstructuredMessage `json:"unstructured,omitempty"`
}

type UnstructuredMessage struct {
	ID        string `json:"id,omitempty"`
	Text      string `json:"text,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type WeatherData struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int     `json:"sunrise"`
		Sunset  int     `json:"sunset"`
	} `json:"sys"`
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cod  int    `json:"cod"`
}

func weather() string {
	weather := WeatherData{}

	url := "http://api.openweathermap.org/data/2.5/weather?q=Manhattan,us&appid=b1cdcefb69ed835d369791fbdf91b536"

	client := http.Client{}

	req, err :=  http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return "Error Fetching Data"
	}
	body, geterr := ioutil.ReadAll(res.Body)
	if geterr != nil {
		log.Fatal(geterr)
		return "Error Fetching Data"
	}
	json.Unmarshal([]byte(body), &weather)
	kelvin := weather.Main.Temp
	temp := kelvin - 273.15
	temperature := fmt.Sprint("Temperature in New York is ", int(temp), "°C and it's, " + weather.Weather[0].Description + ".")
	return temperature
}

func handelRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	botRequest := BotRequest{}
	req := request.Body
	json.Unmarshal([]byte(req),&botRequest)
	responseSlice := []Message{}

	for _,v := range botRequest.Messages {
		if v.Unstructured.Text == "Hi" || v.Unstructured.Text == "hi" {
			body := Message{
				Type: request.HTTPMethod,
				Unstructured: &UnstructuredMessage{
					ID:        request.RequestContext.RequestID,
					Text:      "Hi! How are you doing today?",
					Timestamp: time.Now().String(),
				},
			}
			responseSlice = append(responseSlice, body)
			bs, _ := json.Marshal(BotResponce{Messages: responseSlice})
			return events.APIGatewayProxyResponse{Body: string(bs), StatusCode: 200}, nil
		} else if v.Unstructured.Text == "Hey" || v.Unstructured.Text == "hey" {
			body := Message{
				Type: request.HTTPMethod,
				Unstructured: &UnstructuredMessage{
					ID:        request.RequestContext.RequestID,
					Text:      "Hey! There.",
					Timestamp: time.Now().String(),
				},
			}
			responseSlice = append(responseSlice, body)
			bs, _ := json.Marshal(BotResponce{Messages: responseSlice})
			return events.APIGatewayProxyResponse{Body: string(bs), StatusCode: 200}, nil
		} else if v.Unstructured.Text == "Sup" || v.Unstructured.Text == "sup" {
			body := Message{
				Type: request.HTTPMethod,
				Unstructured: &UnstructuredMessage{
					ID:        request.RequestContext.RequestID,
					Text:      "What up, Homie!!",
					Timestamp: time.Now().String(),
				},
			}
			responseSlice = append(responseSlice, body)
			bs, _ := json.Marshal(BotResponce{Messages: responseSlice})
			return events.APIGatewayProxyResponse{Body: string(bs), StatusCode: 200}, nil
		} else if v.Unstructured.Text == "Time" || v.Unstructured.Text == "time" {
			location, _ := time.LoadLocation("America/New_York")
			body := Message{
				Type: request.HTTPMethod,
				Unstructured: &UnstructuredMessage{
					ID:        request.RequestContext.RequestID,
					Text:      "Time is " + time.Now().In(location).Format("3:04PM"),
					Timestamp: time.Now().String(),
				},
			}
			responseSlice = append(responseSlice, body)
			bs, _ := json.Marshal(BotResponce{Messages: responseSlice})
			return events.APIGatewayProxyResponse{Body: string(bs), StatusCode: 200}, nil
		} else if v.Unstructured.Text == "Day" || v.Unstructured.Text == "day" {
			location, _ := time.LoadLocation("America/New_York")
			body := Message{
				Type: request.HTTPMethod,
				Unstructured: &UnstructuredMessage{
					ID:        request.RequestContext.RequestID,
					Text:      "Today is " + time.Now().In(location).Weekday().String() + ".",
					Timestamp: time.Now().String(),
				},
			}
			responseSlice = append(responseSlice, body)
			bs, _ := json.Marshal(BotResponce{Messages: responseSlice})
			return events.APIGatewayProxyResponse{Body: string(bs), StatusCode: 200}, nil
		} else if v.Unstructured.Text == "Month" || v.Unstructured.Text == "month" {
			location, _ := time.LoadLocation("America/New_York")
			body := Message{
				Type: request.HTTPMethod,
				Unstructured: &UnstructuredMessage{
					ID:        request.RequestContext.RequestID,
					Text:      "This is " + time.Now().In(location).Month().String() + ".",
					Timestamp: time.Now().String(),
				},
			}
			responseSlice = append(responseSlice, body)
			bs, _ := json.Marshal(BotResponce{Messages: responseSlice})
			return events.APIGatewayProxyResponse{Body: string(bs), StatusCode: 200}, nil
		} else if v.Unstructured.Text == "Date" || v.Unstructured.Text == "date" {
			body := Message{
				Type: request.HTTPMethod,
				Unstructured: &UnstructuredMessage{
					ID:        request.RequestContext.RequestID,
					Text:      "Today is " + time.Now().Format("Mon Jan 2 2006"),
					Timestamp: time.Now().String(),
				},
			}
			responseSlice = append(responseSlice, body)
			bs, _ := json.Marshal(BotResponce{Messages: responseSlice})
			return events.APIGatewayProxyResponse{Body: string(bs), StatusCode: 200}, nil
		} else if v.Unstructured.Text == "Weather" || v.Unstructured.Text == "weather" {
			temperature := weather()
			body := Message{
				Type: request.HTTPMethod,
				Unstructured: &UnstructuredMessage{
					ID:        request.RequestContext.RequestID,
					Text:      temperature,
					Timestamp: time.Now().String(),
				},
			}
			responseSlice = append(responseSlice, body)
			bs, _ := json.Marshal(BotResponce{Messages: responseSlice})
			return events.APIGatewayProxyResponse{Body: string(bs), StatusCode: 200}, nil

		} else if v.Unstructured.Text == "Hello" {
			body := Message{
				Type: request.HTTPMethod,
				Unstructured: &UnstructuredMessage{
					ID:        request.RequestContext.RequestID,
					Text:      "Hi there, how can I help?",
					Timestamp: time.Now().String(),
				},
			}
			responseSlice = append(responseSlice, body)
			bs, _ := json.Marshal(BotResponce{Messages: responseSlice})
			return events.APIGatewayProxyResponse{Body: string(bs), StatusCode: 200}, nil

		} else if v.Unstructured.Text == "I need some restaurant suggestions." {
			body := Message{
				Type: request.HTTPMethod,
				Unstructured: &UnstructuredMessage{
					ID:        request.RequestContext.RequestID,
					Text:      "Great. I can help you with that. What city or city area are you looking to dine in?",
					Timestamp: time.Now().String(),
				},
			}
			responseSlice = append(responseSlice, body)
			bs, _ := json.Marshal(BotResponce{Messages: responseSlice})
			return events.APIGatewayProxyResponse{Body: string(bs), StatusCode: 200}, nil

		} else if strings.Contains(v.Unstructured.Text, "Please") {
			v.Unstructured.Text = strings.Replace(v.Unstructured.Text, "Please", "", -1)
			data = append(data, v.Unstructured.Text)
			body := Message{
				Type: request.HTTPMethod,
				Unstructured: &UnstructuredMessage{
					ID:        request.RequestContext.RequestID,
					Text:      "Got it, "+v.Unstructured.Text+" What cuisine would you like to try?",
					Timestamp: time.Now().String(),
				},
			}
			responseSlice = append(responseSlice, body)
			bs, _ := json.Marshal(BotResponce{Messages: responseSlice})
			return events.APIGatewayProxyResponse{Body: string(bs), StatusCode: 200}, nil

		} else if v.Unstructured.Text == "Japanese" {
			body := Message{
				Type: request.HTTPMethod,
				Unstructured: &UnstructuredMessage{
					ID:        request.RequestContext.RequestID,
					Text:      "Ok, how many people are in your party?",
					Timestamp: time.Now().String(),
				},
			}
			responseSlice = append(responseSlice, body)
			bs, _ := json.Marshal(BotResponce{Messages: responseSlice})
			return events.APIGatewayProxyResponse{Body: string(bs), StatusCode: 200}, nil

		} else if v.Unstructured.Text == "Two People" {
			body := Message{
				Type: request.HTTPMethod,
				Unstructured: &UnstructuredMessage{
					ID:        request.RequestContext.RequestID,
					Text:      "A few more to go. What date?",
					Timestamp: time.Now().String(),
				},
			}
			responseSlice = append(responseSlice, body)
			bs, _ := json.Marshal(BotResponce{Messages: responseSlice})
			return events.APIGatewayProxyResponse{Body: string(bs), StatusCode: 200}, nil

		} else if v.Unstructured.Text == "Today" {
			body := Message{
				Type: request.HTTPMethod,
				Unstructured: &UnstructuredMessage{
					ID:        request.RequestContext.RequestID,
					Text:      "What time?",
					Timestamp: time.Now().String(),
				},
			}
			responseSlice = append(responseSlice, body)
			bs, _ := json.Marshal(BotResponce{Messages: responseSlice})
			return events.APIGatewayProxyResponse{Body: string(bs), StatusCode: 200}, nil

		} else if strings.Contains(v.Unstructured.Text, "please") {
			v.Unstructured.Text = strings.Replace(v.Unstructured.Text, "please", "", -1)
			data = append(data, v.Unstructured.Text)
			body := Message{
				Type: request.HTTPMethod,
				Unstructured: &UnstructuredMessage{
					ID:        request.RequestContext.RequestID,
					Text:      "Great. Lastly, I need your phone number so I can send you my findings.",
					Timestamp: time.Now().String(),
				},
			}
			responseSlice = append(responseSlice, body)
			bs, _ := json.Marshal(BotResponce{Messages: responseSlice})
			return events.APIGatewayProxyResponse{Body: string(bs), StatusCode: 200}, nil

		} else if strings.Contains(v.Unstructured.Text, "+1") {
			data = append(data, v.Unstructured.Text)
			body := Message{
				Type: request.HTTPMethod,
				Unstructured: &UnstructuredMessage{
					ID:        request.RequestContext.RequestID,
					Text:      "You’re all set. Expect my recommendations shortly! Have a good day.",
					Timestamp: time.Now().String(),
				},
			}
			responseSlice = append(responseSlice, body)
			bs, _ := json.Marshal(BotResponce{Messages: responseSlice})
			if len(data) == 3 {
				sendsSQS(data)
			}
			return events.APIGatewayProxyResponse{Body: string(bs), StatusCode: 200}, nil

		} else if v.Unstructured.Text == "Thank You" {
			body := Message{
				Type: request.HTTPMethod,
				Unstructured: &UnstructuredMessage{
					ID:        request.RequestContext.RequestID,
					Text:      "You’re welcome.",
					Timestamp: time.Now().String(),
				},
			}
			responseSlice = append(responseSlice, body)
			bs, _ := json.Marshal(BotResponce{Messages: responseSlice})
			return events.APIGatewayProxyResponse{Body: string(bs), StatusCode: 200}, nil
		}

	}
	body := Message{
		Type: request.HTTPMethod,
		Unstructured: &UnstructuredMessage{
			ID: request.RequestContext.RequestID,
			Text: `Hi, I can not read your mind. Although, my creator is tyring hard to make me intelligent. In the mean time, type something instead, like 'Weather', 'time', 'month' and I'll try to answer :)`,
			Timestamp: time.Now().String(),
		},
	}
	responseSlice = append(responseSlice, body)
	bs, _ := json.Marshal(BotResponce{Messages:responseSlice})
	return events.APIGatewayProxyResponse{Body: string(bs), StatusCode:200} , nil
}

func main() {
	lambda.Start(handelRequest)
}
