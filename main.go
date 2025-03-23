package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

const Version = "0.5.0"
const Port = "8880"

type MacAddress struct {
	Mac string `json:"mac"`
}

func main() {
	fmt.Printf("styxut/WoL version %s\n", Version) 
	setupServer()
}

func setupServer() {
	http.HandleFunc("/wol", wolHandler)

	fmt.Printf("Go Server starting on :%s\n", Port)
	err := http.ListenAndServe(":"+Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func wolHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unsupported request method. Request method must be POST.", http.StatusMethodNotAllowed)
		return
	}

	var macAddress MacAddress

	err := json.NewDecoder(r.Body).Decode(&macAddress)
	if err != nil {
		response := "Invalid mac address: " + macAddress.Mac
		http.Error(w, response, http.StatusBadRequest)
		return
	}

	err = sendWOL(strings.ToUpper(macAddress.Mac))
	if err != nil {
		response := "Invalid mac address: " + macAddress.Mac
		http.Error(w, response, http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Println("Magic packet sent")
	}
}

func sendWOL(mac string) error {
	// Parse MAC address
	hwAddr, err := net.ParseMAC(mac)
	if err != nil {
		return fmt.Errorf("invalid MAC address: %w", err)
	}

	// Build magic packet: 6 x 0xFF followed by 16 repetitions of MAC
	packet := make([]byte, 102)
	copy(packet[:6], []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
	for i := 1; i <= 16; i++ {
		copy(packet[i*6:], hwAddr)
	}

	// Send packet via UDP broadcast
	fmt.Printf("Sending magic packet to %s\n", mac)
	addr := &net.UDPAddr{
		IP:   net.IPv4bcast, // 255.255.255.255
		Port: 9,
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to dial UDP: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write(packet)
	if err != nil {
		return fmt.Errorf("failed to send packet: %w", err)
	}

	return nil
}
