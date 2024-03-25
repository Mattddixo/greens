package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Mattddixo/greens/chaincodeactions" // Update this import path as necessary.
)

func runCLI(actions *chaincodeactions.ChaincodeActions) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n=== Microgreen Tracker CLI ===")
		fmt.Println("1. Plant Seeds")
		fmt.Println("2. Update Germination")
		fmt.Println("3. Update to Harvested")
		fmt.Println("4. Record Watering")
		fmt.Println("5. Query")
		fmt.Println("6. Exit")
		fmt.Print("Enter option: ")

		scanner.Scan()
		option := scanner.Text()

		switch option {
		case "1":
			fmt.Print("Enter Microgreen Type: ")
			scanner.Scan()
			microgreenType := scanner.Text()
			fmt.Print("Enter Seed Weight: ")
			scanner.Scan()
			seedWeight, _ := strconv.ParseFloat(scanner.Text(), 64)
			actions.PlantSeeds(microgreenType, seedWeight)
		case "2":
			fmt.Print("Enter Batch ID: ")
			scanner.Scan()
			batchID := scanner.Text()
			actions.UpdateGermination(batchID)
		case "3":
			fmt.Print("Enter Batch ID: ")
			scanner.Scan()
			batchID := scanner.Text()
			fmt.Print("Enter Harvest Weight: ")
			scanner.Scan()
			harvestWeight, _ := strconv.ParseFloat(scanner.Text(), 64)
			actions.UpdateToHarvested(batchID, harvestWeight)
		case "4":
			fmt.Print("Enter Batch ID: ")
			scanner.Scan()
			batchID := scanner.Text()
			actions.RecordWatering(batchID)
		case "5":
			handleQueryOptions(scanner, actions)
		case "6":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

func handleQueryOptions(scanner *bufio.Scanner, actions *chaincodeactions.ChaincodeActions) {
	fmt.Println("Query Options:")
	fmt.Println("1. By Status")
	fmt.Println("2. By Type")
	fmt.Println("3. By Date Range")
	fmt.Println("4. Single Batch")
	fmt.Print("Select query type: ")
	scanner.Scan()
	option := scanner.Text()

	switch option {
	case "1":
		fmt.Print("Enter Status: ")
		scanner.Scan()
		status := scanner.Text()
		fmt.Println(actions.QueryBatchesByStatus(status))
	case "2":
		fmt.Print("Enter Type: ")
		scanner.Scan()
		microgreenType := scanner.Text()
		fmt.Println(actions.QueryBatchesByType(microgreenType))
	case "3":
		fmt.Print("Enter Start Date: ")
		scanner.Scan()
		startDate := scanner.Text()
		fmt.Print("Enter End Date: ")
		scanner.Scan()
		endDate := scanner.Text()
		fmt.Println(actions.QueryBatchesByDateRange(startDate, endDate))
	case "4":
		fmt.Print("Enter Batch ID: ")
		scanner.Scan()
		batchID := scanner.Text()
		fmt.Println(actions.QueryBatch(batchID))
	default:
		fmt.Println("Invalid query type.")
	}
}

func main() {
	// Assuming that the actions object is correctly initialized and passed from main.go
	// This is a placeholder. Ensure that you adjust your main.go accordingly.
	var actions *chaincodeactions.ChaincodeActions
	runCLI(actions)
}
