package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/aler9/goroslib"
	"github.com/aler9/goroslib/pkg/msgs/sensor_msgs"
	"github.com/aler9/goroslib/pkg/msgs/std_msgs"
)

func initROS() {
	node1, err := goroslib.NewNode(goroslib.NodeConf{
		Name:          "wheelchair_remote",
		MasterAddress: Arduino.rosMasterAdress,
	})
	if err != nil {
		panic(err)
	}
	defer node1.Close()

	// Listen to the topic
	sub, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:     node1,
		Topic:    "wheelchair_move_command",
		Callback: onMessage,
	})
	if err != nil {
		panic(err)
	}
	defer sub.Close()

	// Publish the Joystick
	pub, err := goroslib.NewPublisher(goroslib.PublisherConf{
		Node:  node1,
		Topic: "joystick_user_input",
		Msg:   &std_msgs.UInt8MultiArray{} ,
	})

	r := node1.TimeRate(5 * time.Millisecond)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for {
			select {
			// publish a message every second
			case <-r.SleepChan():
				msg := &std_msgs.UInt8MultiArray{
					Data: []uint8{Arduino.forward, Arduino.right},
				}
				// fmt.Printf("Outgoing: %+v\n", msg)
				pub.Write(msg)

			// handle CTRL-C
			case <-c:
				return
			}
		}
	}()
}

func onMessage(msg *sensor_msgs.Imu) {
	fmt.Printf("Incoming: %+v\n", msg)
}
