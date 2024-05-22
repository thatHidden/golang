package yoopay

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Yoopay struct {
	shopID    string
	secretKey string
}

type Response struct {
	PaymentID    string `json:"id"`
	Status       string `json:"status"`
	Confirmation struct {
		Url string `json:"confirmation_url"`
	} `json:"confirmation"`
}

func NewYoopayService(shopId string, secretKey string) Yoopay {
	return Yoopay{
		shopID:    shopId,
		secretKey: secretKey,
	}
}

// CreatePayment returns map[string][string] with keys "id", "status" and "confirmation_url"
func (y *Yoopay) CreatePayment(price uint64, obj string, idempotenceKey string) (data map[string]string, err error) {
	url := os.Getenv("YOOPAY_PAYMENTS")

	priceStr := strconv.FormatUint(price, 10)

	requestBody, err := json.Marshal(map[string]interface{}{
		"amount": map[string]string{
			"value":    priceStr,
			"currency": "RUB",
		},
		"payment_method_data": map[string]string{
			"type": "bank_card",
		},
		"confirmation": map[string]string{
			"type":       "redirect",
			"return_url": "https://www.gosuslugi.ru/",
		},
		"description": obj,
		"capture":     false,
		"test":        true,
	})
	if err != nil {
		return data, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return data, err
	}
	req.SetBasicAuth(os.Getenv("YOOPAY_ID"), os.Getenv("test_AcRBlp608-UOtGMkSHtP1Oe_At2XptvTC6YgNXs1PzA")) // Авторизация
	req.Header.Set("Idempotence-Key", idempotenceKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	var response Response

	json.Unmarshal(body, &response)
	if err != nil {
		return data, err
	}

	data = map[string]string{
		"payment_id":       response.PaymentID,
		"status":           response.Status,
		"confirmation_url": response.Confirmation.Url,
	}

	return data, nil
}

func (y *Yoopay) CapturePayment(paymentId, idempotenceKey string) (err error) {
	url := "https://api.yookassa.ru/v3/payments/" + paymentId + "/capture"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth("378495", "test_AcRBlp608-UOtGMkSHtP1Oe_At2XptvTC6YgNXs1PzA")
	req.Header.Set("Idempotence-Key", idempotenceKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("YooKassa Status: " + resp.Status)
	}
	return
}

func (y *Yoopay) CancelPayment(paymentIds []string, idempotenceKey string) (err error) {
	for _, paymentId := range paymentIds {
		url := "https://api.yookassa.ru/v3/payments/" + paymentId + "/cancel"

		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return err
		}
		//это выносится??
		req.SetBasicAuth("378495", "test_AcRBlp608-UOtGMkSHtP1Oe_At2XptvTC6YgNXs1PzA") // Авторизация
		req.Header.Set("Idempotence-Key", idempotenceKey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			fmt.Println("err " + paymentId) //log
		}
	}
	return
}
