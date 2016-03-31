package communication

import (
	"testing"
	"time"
	"github.com/satori/go.uuid"
)

func TestSend(t *testing.T) {
	sendCh := make(chan CommData)
	receiveCh, connStatus := Run(sendCh)
	go RunPrintConn(connStatus)
	go RunPrintMsg(receiveCh)
	time.Sleep(1 * time.Second)
	count := 0
	for{
		msgID := uuid.NewV4()
		msg := ResolveMsg("10.20.78.108", msgID.String(), "Test", count) 
		sendCh <- *msg
		time.Sleep(10 * time.Second)
		count += 1
	}
}


func TestListen(t *testing.T) {
	sendCh := make(chan CommData)
	receiveCh, connStatus := Run(sendCh)
	go RunPrintConn(connStatus)
	go RunPrintMsg(receiveCh)
	time.Sleep(1 * time.Second)
	count := 0
	for{
		time.Sleep(10 * time.Second)
		count += 1
	}
}

func RunPrintMsg(receiveCh <-chan CommData) {
	for{
		PrintMessage(<- receiveCh)
	}
}
func RunPrintConn(connStatus <-chan ConnData) {
	for{
		PrintConnData(<- connStatus)
	}
}

