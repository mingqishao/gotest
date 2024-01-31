package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/go-shuttle/v2"
	fleetv1 "go.goms.io/aks/rp/protos/fleet/v1"
	"google.golang.org/protobuf/proto"

	"encoding/base64"
	// Replace with the actual path to your generated Go protobuf package
)

type MyError struct{}

func (m *MyError) Error() string {
	return "boom"
}

func main() {
	// azlog.SetListener(func(event azlog.Event, s string) {
	// 	fmt.Printf("%s, [%s] %s\n", time.Now(), event, s)
	// })

	// receiveDLQFromTopic()
	// GetQueueRuntime()
	// sendTopic()
	// runNetwork()
	// fmt.Println("start sleep")
	// time.Sleep(10 * time.Second)
	// sendQueue()
	receiveDLQFromQueue()
}

func runNetwork() {
	fmt.Println("start runNetwork")
	stopChan := make(chan bool, 1)
	go network(stopChan)
	fmt.Println("finish runNetwork")
}

func network(stopChan chan bool) {
	fmt.Println("start network")
	switch {
	case <-stopChan:
		fmt.Println("stop network")
	}
}

func fortest(row []int) {
	row = append(row, 4)
}

func receiveDLQFromQueue() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// handle error
	}
	client, err := azservicebus.NewClient("e2eakssvcbus-r2xh3d97h.servicebus.windows.net", cred, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to Service Bus")
	receiver, err := client.NewReceiverForQueue("fleetoperations", &azservicebus.ReceiverOptions{
		SubQueue: azservicebus.SubQueueDeadLetter,
	})
	// receiver, err := client.NewReceiverForSubscription("managedcluster", "fleet", &azservicebus.ReceiverOptions{
	// 	SubQueue: azservicebus.SubQueueDeadLetter,
	// })
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to Queue")
	// ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	// defer cancel()
	ctx := context.TODO()
	msgs, err := receiver.ReceiveMessages(ctx, 2, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("received messages")
	fmt.Println(len(msgs))

	marshaller := shuttle.DefaultProtoMarshaller{}
	for _, m := range msgs {
		fmt.Printf("SequenceNumber: %d, DeliveryCount: %d, Body: %s\n", m.SequenceNumber, m.DeliveryCount, string(m.Body))
		fmt.Println("begin AbandonMessage")
		receiver.AbandonMessage(ctx, m, nil)
		fmt.Println("endAbandonMessage")

		marshaller.Unmarshal()
	}
}

func receiveDLQFromTopic() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// handle error
	}
	client, err := azservicebus.NewClient("e2eakssvcbus-r6cy242w9.servicebus.windows.net", cred, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to Service Bus")
	receiver, err := client.NewReceiverForSubscription("managedcluster", "fleet", &azservicebus.ReceiverOptions{
		SubQueue: azservicebus.SubQueueDeadLetter,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to Queue")
	// ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	// defer cancel()
	ctx := context.TODO()
	msgs, err := receiver.ReceiveMessages(ctx, 2, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("received messages")
	fmt.Println(len(msgs))
	for _, m := range msgs {
		fmt.Printf("SequenceNumber: %d, DeliveryCount: %d, ID: %s, Body: %s\n", m.SequenceNumber, m.DeliveryCount, m.MessageID, string(m.Body))
		fmt.Println("begin AbandonMessage")
		receiver.CompleteMessage(ctx, m, nil)
		// receiver.AbandonMessage(ctx, m, nil)
		fmt.Println("endAbandonMessage")
	}
}

func sendQueue() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	client, err := azservicebus.NewClient("e2eakssvcbus-r2xh3d97h.servicebus.windows.net", cred, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	senderOptions := &shuttle.SenderOptions{
		Marshaller:               &shuttle.DefaultProtoMarshaller{},
		EnableTracingPropagation: true,
		SendTimeout:              -1,
	}

	fmt.Println("Connected to Service Bus")
	// sender, err := client.NewSender("managedcluster/subscriptions/fleet", nil)
	sender, err := client.NewSender("fleetoperations", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	shuttleSender := shuttle.NewSender(sender, senderOptions)

	message := &fleetv1.Fleet{
		Name:           "test",
		SubscriptionId: "sub1",
		ResourceGroup:  "rg1",
	}
	// serializedMessage, err := proto.Marshal(message)
	// if err != nil {
	// 	fmt.Println("Error serializing message:", err)
	// 	return
	// }
	shuttleSender.SendMessage(context.TODO(), message)
	// err = sender.SendMessage(context.TODO(), &azservicebus.Message{
	// 	Body: serializedMessage,
	// }, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func sendTopic() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// handle error
	}
	client, err := azservicebus.NewClient("e2eakssvcbus-r6cy242w9.servicebus.windows.net", cred, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to Service Bus")
	// sender, err := client.NewSender("managedcluster/subscriptions/fleet", nil)
	sender, err := client.NewSender("managedcluster", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = sender.SendMessage(context.TODO(), &azservicebus.Message{
		ApplicationProperties: map[string]interface{}{
			"type": "ManagedClusterDeleted",
		},
		Body: []byte("{}"),
	}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func GetQueueRuntime() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// handle error
	}
	client, err := admin.NewClient("e2eakssvcbus-r2xh3d97h.servicebus.windows.net", cred, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.TODO()
	p, err := client.GetQueueRuntimeProperties(ctx, "fleet", nil)
	fmt.Printf("p: %+v\n", p)
	fmt.Printf("err: %+v\n", err)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(p.DeadLetterMessageCount)
}

func GetTopicRuntime() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// handle error
	}
	client, err := admin.NewClient("e2eakssvcbus-r6cy242w9.servicebus.windows.net", cred, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.TODO()
	p, err := client.GetSubscriptionRuntimeProperties(ctx, "managedcluster", "fleet", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(p.ActiveMessageCount)
	fmt.Println(p.DeadLetterMessageCount)
}

func protocolbase64() {
	// Create an instance of ExampleMessage
	message := &fleetv1.Fleet{
		Name: "test",
	}

	// Serialize the message to bytes
	serializedMessage, err := proto.Marshal(message)
	if err != nil {
		fmt.Println("Error serializing message:", err)
		return
	}

	// Encode the bytes to base64
	encodedMessage := base64.StdEncoding.EncodeToString(serializedMessage)

	fmt.Println("Base64 Encoded Message:", encodedMessage)
}
