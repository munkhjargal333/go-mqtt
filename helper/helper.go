package helper

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

func ParseStateData(stateData string) (map[int]int, error) {
	stateMap := make(map[int]int)

	// Split the stateData string by commas
	pairs := strings.Split(stateData, ",")

	for _, pair := range pairs {
		// Split each pair into key and value
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid state data format")
		}

		// Convert key and value to int
		key, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		value, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		// Add key-value pair to the state map
		stateMap[key] = value
	}

	return stateMap, nil
}

func RemoveState(message string) (string, string) {
	// Find the start and end index of the "state" field
	start := strings.Index(message, `"state":{`)
	if start == -1 {
		return message, ""
	}
	end := strings.Index(message[start:], "}") + start + 2
	if end == -1 {
		return message, ""
	}

	// Extract the "state" field content
	stateData := message[start+len(`"state":{`) : end]

	// Remove the last two elements from the "state" data
	stateDataParts := strings.Split(stateData, ",")
	if len(stateDataParts) >= 2 {
		stateData = strings.Join(stateDataParts[:len(stateDataParts)-2], ",")
	}

	// Remove the "state" field from the message
	modifiedMessage := message[:start] + message[end:]

	return modifiedMessage, stateData
}

func RemoveCanItems(message string) (string, string) {
	// Find the start and end index of the "state" field
	start := strings.Index(message, `"can":{`)
	if start == -1 {
		return message, ""
	}
	end := strings.Index(message[start:], "}") + start + 2
	if end == -1 {
		return message, ""
	}

	// Extract the "state" field content
	stateData := message[start+len(`"can":{`) : end]

	// Remove the last two elements from the "state" data
	stateDataParts := strings.Split(stateData, ",")
	if len(stateDataParts) >= 2 {
		stateData = strings.Join(stateDataParts[:len(stateDataParts)-2], ",")
	}

	// Remove the "state" field from the message
	modifiedMessage := message[:start] + message[end:]

	return modifiedMessage, stateData
}

func ArrayToString(intSlice []int) string {
	stringSlice := make([]string, len(intSlice))
	for i, v := range intSlice {
		stringSlice[i] = strconv.Itoa(v)
	}

	return strings.Join(stringSlice, ",")
}

func MapToString(m map[int]int) string {
	// Get the keys and sort them
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// Create string representation of sorted key-value pairs
	pairs := make([]string, 0, len(m))
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%d:%d", k, m[k]))
	}
	return strings.Join(pairs, ",")
}

func WaitForShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Println("Shutting down...")
	os.Exit(0)
}
