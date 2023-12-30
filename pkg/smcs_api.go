package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

type UserMess struct {
	Status        int64  `json:"status"`
	LastDate      string `json:"last_date"`
	LastTimestamp int64  `json:"last_timestamp"`
	Err           int `json:"err"`
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

type ErrMess struct {
	Text string `json:"error"`
	Code int    `json:"error_code"`
}

var (
	ErrNoUserHistory string = "пользователь еще не отправлял запрос для получение кода"
	ErrForbidden     string = "нельзя отправлять запрос на номер сотрудника"
)

const (
	ErrInvalidNumber1             = "Абонент не существует: Указанный номер телефона не существует."
	ErrSubscriberNotInNet         = "Абонент не в сети: Телефон абонента отключен или находится вне зоны действия сети."
	ErrServiceNotConnected        = "Не подключена услуга: Абонент не может принять SMS-сообщение."
	ErrInvalidPhoneNumber         = "Ошибка в телефоне абонента: Не удается доставить сообщение из-за ошибки в телефонном аппарате или SIM-карте."
	ErrSubscriberBlocked          = "Абонент заблокирован: Нулевой или отрицательный баланс, заблокирован оператором или добровольная блокировка."
	ErrNoServiceSupport           = "Нет поддержки сервиса: Аппарат абонента не поддерживает работу с данной услугой."
	ErrVirtualSending             = "Виртуальная отправка: Уведомление появляется при отправке сообщения в режиме тестирования."
	ErrSimCardReplacement         = "Замена SIM-карты: Ошибка отправки сообщения из-за замены абонентом SIM-карты."
	ErrOperatorQueueOverflow      = "Переполнена очередь у оператора: Абонент недоступен, но сообщения продолжают поступать."
	ErrSubscriberNotAnswering     = "Абонент не отвечает: Во время дозвона абонент не взял трубку."
	ErrNoTemplate                 = "Нет шаблона: Отправка возможна только по определенному шаблону."
	ErrForbiddenIPAddress         = "Запрещенный IP-адрес: Попытка отправки сообщения с неразрешенного IP-адреса."
	ErrSubscriberBusy             = "Абонент занят: Линия занята или абонент отменил вызов."
	ErrConversionError            = "Ошибка конвертации: Произошла ошибка конвертации текста или звукового файла."
	ErrAnsweringMachineDetected   = "Зафиксирован автоответчик: Был зафиксирован автоответчик абонента."
	ErrUnregisteredSenderID       = "Незарегистрированный Sender ID: Попытка отправки сообщения от незарегистрированного имени отправителя."
	ErrRejectedByOperator         = "Отклонено оператором: Оператор отклонил сообщение без указания точного кода ошибки."
	ErrInvalidFormatNumber        = "Неверный формат номера: Мобильный код и длина номера неверны."
	ErrNumberNotAllowedBySettings = "Номер запрещен настройками: Попадание номера под ограничения, установленные клиентом."
	ErrDailyMessageLimitExceeded  = "Превышен лимит сообщений: Превышен суточный лимит сообщений, указанный клиентом."
	ErrNoRoute                    = "Нет маршрута: На данный номер отправка сообщений недоступна в нашем сервисе."
	ErrInvalidFormatNumber249     = "Неверный формат номера: Мобильный код указанного номера и длина номера неверны."
	ErrNumberProhibitedBySettings = "Номер запрещен настройками: Номер попал под ограничения, установленные клиентом."
	ErrExceedDailyLimitPerNumber  = "Превышен лимит на один номер: Превышен суточный лимит сообщений на один номер."
	ErrNumberProhibited           = "Номер запрещен: Попытка указания клиентом запрещенного номера."
	ErrSpamFilterForbidden        = "Запрещено спам-фильтром: Текст сообщения содержит запрещенные выражения или ссылки."
	ErrUnregisteredSenderID255    = "Незарегистрированный Sender ID: Попытка отправки сообщения от незарегистрированного имени отправителя."
	ErrOperatorRejected           = "Отклонено оператором: Оператор отклонил сообщение без указания точного кода ошибки."
)

func GetRequestSmcs(phone string) (*string, error) {
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
	req, err := http.NewRequest(http.MethodPost, apiLink, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("body: ", string(body))
	//fmt.Println(string(body))
	var jsonData []UserMess
	var errData ErrMess
	if err := json.Unmarshal(body, &jsonData); err != nil {
		if jsonErr := json.Unmarshal(body, &errData); jsonErr != nil {
			fmt.Println("json err: ", jsonErr)
			return nil, err
		} else if errData.Code == 3 {
			return nil, fmt.Errorf(ErrNoUserHistory)
		}
		return nil, fmt.Errorf(errData.Text)
	}
	str := jsonData[0].Message
	re := regexp.MustCompile(`\d`)
	digits := re.FindAllString(str, -1)

	if len(digits) == 4 {
		str = digits[0] + digits[1] + digits[2] + digits[3]
	} else {
		return nil, fmt.Errorf("invalid count digits: %d ", len(digits))
	}
	if str == "2990"{ // для рэббит хол
		str = jsonData[0].Message
	}
	mes := getErrorMessage(jsonData[0].Err)
	if mes != ""{
		str = mes
	}

	result := fmt.Sprintf("Запрос для номера: [%s] - Cообщения: %s", phone, str)

	return &result, nil
}

