package main

import (
	"log"
	"net/http"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"gitlab.com/gomidi/midi/writer"
	driver "gitlab.com/gomidi/portmididrv"
)

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// This example expects the first input and output port to be connected
// somehow (are either virtual MIDI through ports or physically connected).
// We write to the out port and listen to the in port.
func main() {
	drv, err := driver.New()
	must(err)

	// make sure to close all open ports at the end
	defer drv.Close()

	ins, err := drv.Ins()
	must(err)

	outs, err := drv.Outs()
	must(err)

	log.Println(ins)
	log.Println(outs)

	in, out := ins[0], outs[0]

	must(in.Open())
	must(out.Open())

	wr := writer.New(out)

	// listen for MIDI
	// go mid.NewReader().ReadFrom(in)

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())

		if err := writer.NoteOn(wr, 61, 100); err != nil {
			return err
		}
		time.Sleep(time.Second)
		if err := writer.NoteOff(wr, 61); err != nil {
			return err
		}

		return nil
	})
	type drawingMsg struct {
		Color string
		X0    float64
		X1    float64
		Y0    float64
		Y1    float64
	}
	server.OnEvent("/", "drawing", func(s socketio.Conn, msg drawingMsg) {
		if err := func() error {
			//s.Emit("reply", "have "+msg)
			note := uint8(msg.X1 * 127)
			if note < 20 {
				note = 20
			}
			velocity := uint8(msg.Y1 * 127)
			if velocity < 20 {
				velocity = 20
			}
			log.Println("note:", note, "velocity:", velocity, "input:", msg)
			if err := writer.NoteOn(wr, note, velocity); err != nil {
				return err
			}
			time.Sleep(time.Nanosecond * 1000000)
			if err := writer.NoteOff(wr, note); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			log.Printf("error on drawing event: %v", err)
		}
	})
	server.OnError("/", func(e error) {
		log.Println("error:", e)

		if err := func() error {
			if err := writer.NoteOn(wr, 62, 100); err != nil {
				return err
			}
			time.Sleep(time.Second)
			if err := writer.NoteOff(wr, 62); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			log.Printf("socket.io error event: %v", err)
		}
	})
	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		if err := func() error {
			log.Println("closed", msg)

			if err := writer.NoteOn(wr, 63, 100); err != nil {
				return err
			}
			time.Sleep(time.Second)
			if err := writer.NoteOff(wr, 63); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			log.Printf("error on disconnect event: %v", err)
		}
	})
	go func() {
		if err := server.Serve(); err != nil {
			panic(err)
		}
	}()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
