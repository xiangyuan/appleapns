package apns

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type apnConn struct {
	tlsConn          *tls.Conn
	tlsCofig         tls.Config
	ReadTimeout      time.Duration
	ConnectLocation  string
	identifier       uint32
	MAX_PAYLOAD_SIZE int
	isConnected      bool
}

func (client *apnConn) ConnectApns() (err error) {
	if client.isConnected {
		return nil
	}
	if client.tlsConn != nil {
		client.Shutdown()
	}
	conn, err := net.Dial("tcp", client.ConnectLocation)
	if err != nil {
		return err
	}
	client.tlsConn = tls.Client(conn, &client.tlsCofig)
	err = client.tlsConn.Handshake()
	if err == nil {
		client.isConnected = true
	}
	return err
}

func (client *apnConn) Shutdown() (err error) {
	err = nil
	if client.tlsConn != nil {
		err = client.tlsConn.Close()
		client.isConnected = false
	}
	return err
}

func NewApnsClient(apnsURL string, certificate string, keyFile string) (apns *apnConn, err error) {
	cert, err := tls.LoadX509KeyPair(certificate, keyFile)
	if err != nil {
		return
	}
	apns = &apnConn{tlsConn: nil,
		tlsCofig: tls.Config{
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{cert}},
		ConnectLocation:  apnsURL,
		ReadTimeout:      200 * time.Millisecond,
		MAX_PAYLOAD_SIZE: 256,
		isConnected:      false,
	}

	return apns, nil
}

func (client *apnConn) SendPayload(token string, payload string, expaire time.Duration) (err error) {
	state := client.tlsConn.ConnectionState()
	fmt.Println("conn state %v \n", state)
	bpayload, err := json.Marshal(payload)
	btoken, err := hex.DecodeString(token)

	buffer := bytes.NewBuffer([]byte{})
	//1--3:identifier--3time-- token len -- token -- paylen -- payload
	binary.Write(buffer, binary.BigEndian, uint8(1))
	client.identifier++
	binary.Write(buffer, binary.BigEndian, uint32(client.identifier))
	binary.Write(buffer, binary.BigEndian, uint32(client.ReadTimeout))

	binary.Write(buffer, binary.BigEndian, uint16(len(btoken)))
	binary.Write(buffer, binary.BigEndian, btoken)

	binary.Write(buffer, binary.BigEndian, uint16(len(bpayload)))
	binary.Write(buffer, binary.BigEndian, bpayload)

	pdu := buffer.Bytes()

	//write to server
	_, err = client.tlsConn.Write(pdu)
	if err != nil {
		return
	}
	defer client.tlsConn.Close()
	client.tlsConn.SetReadDeadline(time.Now().Add(client.ReadTimeout))
	rbuff := [6]byte{}
	n, err := client.tlsConn.Read(rbuff[:])
	if n > 0 {
		fmt.Println("receive error %s \n", hex.EncodeToString(rbuff[:n]))
	}
	if err != nil {
		return
	}
	return
}

func CheckError(err error) {
	if err != nil {
		fmt.Println()
	}
}
