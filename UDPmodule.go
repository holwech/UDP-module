package UDPmodule

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

type Packet struct {
	SourceIP, DestinationIP, ID, Response string
	Content                               []byte
}

func Init(readPort string, writePort string) (<-chan Packet, chan<- Packet) {
	receive := make(chan Packet, 10)
	send := make(chan Packet, 10)
	go listen(receive, readPort)
	go broadcast(send, writePort)
	return receive, send
}

func broadcast(send chan CommData, localIP string, port string) {
	fmt.Printf("COMM: Broadcasting message to: %s%s\n", broadcast_addr, port)
	broadcastAddress, err := net.ResolveUDPAddr("udp", broadcast_addr+port)
	printError("ResolvingUDPAddr in Broadcast failed.", err)
	localAddress, err := net.ResolveUDPAddr("udp", GetLocalIP())
	connection, err := net.DialUDP("udp", localAddress, broadcastAddress)
	printError("DialUDP in Broadcast failed.", err)

	localhostAddress, err := net.ResolveUDPAddr("udp", "localhost"+port)
	printError("ResolvingUDPAddr in Broadcast localhost failed.", err)
	lConnection, err := net.DialUDP("udp", localAddress, localhostAddress)
	printError("DialUDP in Broadcast localhost failed.", err)
	defer connection.Close()

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	for {
		message := <-send
		err := encoder.Encode(message)
		printError("Encode error in broadcast: ", err)
		_, err = connection.Write(buffer.Bytes())
		if err != nil {
			_, err = lConnection.Write(buffer.Bytes())
			printError("Write in broadcast localhost failed", err)
		}
		buffer.Reset()
	}
}

func listen(receive chan CommData, port string) {
	localAddress, err := net.ResolveUDPAddr("udp", port)
	connection, err := net.ListenUDP("udp", localAddress)
	defer connection.Close()
	var message CommData

	for {
		inputBytes := make([]byte, 4096)
		length, _, err := connection.ReadFromUDP(inputBytes)
		buffer := bytes.NewBuffer(inputBytes[:length])
		decoder := gob.NewDecoder(buffer)
		err = decoder.Decode(&message)
		if message.Key == com_id {
			receive <- message
		}
	}
}

func PrintMessage(data CommData) {
	fmt.Printf("=== Data received ===\n")
	fmt.Printf("Key: %s\n", data.Key)
	fmt.Printf("SenderIP: %s\n", data.SenderIP)
	fmt.Printf("ReceiverIP: %s\n", data.ReceiverIP)
	fmt.Printf("Message ID: %s\n", data.MsgID)
	fmt.Printf("= Data = \n")
	fmt.Printf("Data type: %s\n", data.Response)
	fmt.Printf("Content: %v\n", data.Content)
}

func printError(errMsg string, err error) {
	if err != nil {
		fmt.Println(errMsg)
		fmt.Println(err.Error())
	}
}

func GetLocalIP() string {
	var localIP string
	addr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("GetLocalIP in communication failed")
		return "localhost"
	}
	for _, val := range addr {
		if ip, ok := val.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				localIP = ip.IP.String()
			}
		}
	}
	return localIP
}

func ResolveMsg(senderIP string, receiverIP string, msgID string, response string, content map[string]interface{}) (commData *CommData) {
	message := CommData{
		Key:        com_id,
		SenderIP:   senderIP,
		ReceiverIP: receiverIP,
		MsgID:      msgID,
		Response:   response,
		Content:    content,
	}
	return &message
}
