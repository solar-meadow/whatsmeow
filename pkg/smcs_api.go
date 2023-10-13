package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

type UserMess struct {
	Status        int64  `json:"status"`
	LastDate      string `json:"last_date"`
	LastTimestamp int64  `json:"last_timestamp"`
	Err           string `json:"err"`
	SendDate      string `json:"send_date"`
	SendTimestamp int64  `json:"send_timestamp"`
	Phone         string `json:"phone"`
	Cost          string `json:"cost"`
	SenderID      string `json:"sender_id"`
	StatusName    string `json:"status_name"`
	Message       string `json:"message"`
	MCCMNC        string `json:"mccmnc"`
	Country       string `json:"country"`
	Operator      string `json:"operator"`
	Region        string `json:"region"`
	Type          int    `json:"type"`
	ID            int    `json:"id"`
	IntID         string `json:"int_id"`
	SMSCnt        int    `json:"sms_cnt"`
	Format        int    `json:"format"`
	CRC           int64  `json:"crc"`
}

func GetRequestSmcs(ctx context.Context, phone string) (*string, error) {

	login, err := CheckEnvExist("LOGIN")
	if err != nil {
		return nil, err
	}
	password, err := CheckEnvExist("PASSWORD")
	if err != nil {
		return nil, err
	}
	apiLink, err := CheckEnvExist("API_LINK")
	if err != nil {
		return nil, err
	}

	// url params
	data := url.Values{}
	data.Set("get_messages", "1")
	data.Set("login", login)
	data.Set("psw", password)
	data.Set("phone", phone)
	data.Set("fmt", "3") // format json

	resp, err := http.Post(apiLink, "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jsonData []UserMess
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, err
	}
	str := jsonData[0].Message
	re := regexp.MustCompile(`\d`)
	digits := re.FindAllString(str, -1)

	if len(digits) == 4 {
		str = digits[0] + digits[1] + digits[2] + digits[3]
	} else {
		return nil, fmt.Errorf("invalid count digits: %d ", len(digits))
	}

	return &str, nil
}
