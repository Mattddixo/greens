package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Mattddixo/greens/chaincodeactions" // Adjust the import path based on your project's structure
)

func main() {
	// Assume actions is initialized and passed from main.go, this is just for placeholder
	var actions *chaincodeactions.ChaincodeActions

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("=== Microgreen Tracker CLI ===")
		fmt.Println("1. Plant Seeds")
		fmt.Println("2. Germinated")
		fmt.Println("3. Harvested")
		fmt.Println("4. Water")
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
			seedWeight, err := strconv.ParseFloat(scanner.Text(), 64)
			if err != nil {
				fmt.Println("Invalid input for seed weight. Please enter a valid number.")
				continue
			}

			err = actions.PlantSeeds(microgreenType, seedWeight)
			if err != nil {
				fmt.Printf("Failed to plant seeds: %v\n", err)
			} else {
				fmt.Println("Seeds planted successfully.")
			}
		case "2":
			fmt.Print("Enter Batch ID: ")
			scanner.Scan()
			batchID := scanner.Text()

			err := actions.UpdateGermination(batchID)
			if err != nil {
				fmt.Printf("Failed to update germination: %v\n", err)
			} else {
				fmt.Println("Germination updated successfully.")
			}
		case "3":
			fmt.Print("Enter Batch ID: ")
			scanner.Scan()
			batchID := scanner.Text()

			fmt.Print("Enter Harvest Weight: ")
			scanner.Scan()
			harvestWeight, err := strconv.ParseFloat(scanner.Text(), 64)
			if err != nil {
				fmt.Println("Invalid input for harvest weight. Please enter a valid number.")
				continue
			}

			err = actions.UpdateToHarvested(batchID, harvestWeight)
			if err != nil {
				fmt.Printf("Failed to update to harvested: %v\n", err)
			} else {
				fmt.Println("Updated to harvested successfully.")
			}
		case "4":
			fmt.Print("Enter Batch ID: ")
			scanner.Scan()
			batchID := scanner.Text()

			err := actions.RecordWatering(batchID)
			if err != nil {
				fmt.Printf("Failed to record watering: %v\n", err)
			} else {
				fmt.Println("Watering recorded successfully.")
			}
		case "5":
			queryCLI(scanner, actions)
		case "6":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option. Please enter a number between 1 and 6.")
		}
	}
}

func queryCLI(scanner *bufio.Scanner, actions *chaincodeactions.ChaincodeActions) {
	fmt.Println("Query Types:")
	fmt.Println("1. By Status")
	fmt.Println("2. By Type")
	fmt.Println("3. By Date Range")
	fmt.Println("4. Single Batch")
	fmt.Print("Select query type: ")

	scanner.Scan()
	queryType := scanner.Text()

	switch queryType {
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
	
