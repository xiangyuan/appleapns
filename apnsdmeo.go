package main

import (
	"apns"
	"fmt"
	"time"
)

func main() {

	conn, err := apns.NewApnsClient(apns.APPLE_DEVELOPMENT_PUSH, "cert.pem", "key.pem")
	err = conn.ConnectApns()
	if err != nil {
		fmt.Println(err.Error())
	}

	payload := make(map[string]interface{})
	payload["aps"] = map[string]string{"alert": "hello push"}
	err = conn.SendPayload(apns.DEVICE_TOKEN, payload, time.Now().Add(10*60))
	if err != nil {
		fmt.Println("error : ", err.Error())
	}

}
