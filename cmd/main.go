package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


type SibionicUser struct {
	Email string `json:"email"`
	Password string `json:"password"`
}
var API_URL string;
var USER_ID string;


type LoginResponse struct {
	Timestamp int64  `json:"timestamp"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Data      struct {
		AccessToken string `json:"access_token"`
		ExpiresIn  int    `json:"expires_in"`
	} `json:"data"`
	ErrorData interface{} `json:"errorData"`
	Success   bool        `json:"success"`
}


type GlucoseInfo struct {
	I    int     `json:"i"`
	T    string  `json:"t"` // record time 
	V    float64 `json:"v"` // blood glucose value mmol/l
	S    int     `json:"s"`
	Ast  int     `json:"ast"`
	Bl   float64 `json:"bl"`
	Name *string `json:"name"`
}

type Target struct {
	Upper int `json:"upper"`
	Lower int `json:"lower"`
	IsRec int `json:"isRec"`
	DrType int `json:"drType"`
}

type GlucoseData struct {
	UserId             string        `json:"userId"`
	UserNickName       string        `json:"userNickName"`
	DrType             int           `json:"drType"`
	DeviceId           string        `json:"deviceId"`
	DeviceName         string        `json:"deviceName"`
	BlueToothNum       string        `json:"blueToothNum"`
	DeviceEnableTime   string        `json:"deviceEnableTime"`
	DeviceStatus       int           `json:"deviceStatus"`
	DeviceAlarmStatus  int           `json:"deviceAlarmStatus"`
	LatestIndex       int           `json:"latestIndex"`
	LatestGlucoseValue float64      `json:"latestGlucoseValue"`
	BloodGlucoseTrend  int          `json:"bloodGlucoseTrend"`
	LatestGlucoseTime  string       `json:"latestGlucoseTime"`
	DeviceAbnormalTime *string      `json:"deviceAbnormalTime"`
	DeviceLastTime     string       `json:"deviceLastTime"`
	GlucoseInfos      []GlucoseInfo `json:"glucoseInfos"`
	Target            Target        `json:"target"`
}

type GlucoseDataResponse struct {
	Timestamp int64       `json:"timestamp"`
	Code      int        `json:"code"`
	Msg       string     `json:"msg"`
	Data      GlucoseData `json:"data"`
	ErrorData interface{} `json:"errorData"`
	Success   bool       `json:"success"`
}


func Login(user SibionicUser) (string, error) {

	client := &http.Client{}
	loginURL := fmt.Sprintf("%s/auth/app/user/login", API_URL)

	jsonData, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("error marshaling user data: %v", err)
	}

	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if !loginResp.Success {
		return "", fmt.Errorf("login failed: %s", loginResp.Msg)
	}
	return loginResp.Data.AccessToken, nil
}

func GetData(token string) (*GlucoseDataResponse, error) {

	client := &http.Client{}
	dataURL := fmt.Sprintf("%s/user/app/follow/deviceGlucose", API_URL)

	body, err := json.Marshal(map[string]string{
		"id": USER_ID,
		"range": "24",
	})

	if err != nil {
		return nil, fmt.Errorf("error marshaling body: %v", err)
	}

	req, err := http.NewRequest("POST", dataURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var dataResp GlucoseDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&dataResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &dataResp, nil
}

func main() {
	
	var user SibionicUser
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	user.Email = os.Getenv("USER_EMAIL")
	user.Password = os.Getenv("USER_PASSWORD")
	API_URL = os.Getenv("API_URL")
	USER_ID = os.Getenv("USER_ID")

	fmt.Println(user)

	token, err := Login(user)
	if err != nil {
		log.Fatal("Error logging in: %v", err)
	}

	data, err := GetData(token)
	if err != nil {
		log.Fatal("Error getting data: %v", err)
	}
	
	fmt.Printf("Glucose Data: %+v\n", data.Data.GlucoseInfos)
	
	r := gin.Default()

	
	r.GET("/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	
	r.Run(":8080")
}