import app_logic
import datetime

def display_menu():
    print("\n=== Microgreen Tracker CLI ===")
    print("1. Plant Seeds")
    print("2. Update Germination")
    print("3. Record Watering")
    print("4. Update to Harvested")
    print("5. Query Batches")
    print("6. Exit")

def get_status_choice():
    status_options = {
        '1': 'Planted',
        '2': 'Germinated',
        '3': 'Harvested'
    }
    print("\nSelect a status option:")
    for key, value in status_options.items():
        print(f"{key}. {value}")
    choice = input("> ")
    return status_options.get(choice, None)

def format_and_print_batch(batch):
    print(f"\nBatch ID: {batch.batchID}")
    print(f"Type: {batch.type}")
    print(f"Seed Weight: {batch.seedWeight}g")
    print(f"Plant Date: {batch.plantDate.strftime('%Y-%m-%d %H:%M:%S')}")
    print(f"Status: {batch.status}")

    if batch.germinateDate:
        print(f"Germinate Date: {batch.germinateDate.strftime('%Y-%m-%d %H:%M:%S')}")
    if batch.harvestDate:
        print(f"Harvest Date: {batch.harvestDate.strftime('%Y-%m-%d %H:%M:%S')}")
        print(f"Harvest Weight: {batch.harvestWeight}g" if batch.harvestWeight else "Harvest Weight: Not Available")
    
    if batch.wateringEvents:
        print("Watering Events:")
        for event in batch.wateringEvents:
            print(f" - {event.waterDate.strftime('%Y-%m-%d %H:%M:%S')}")
    else:
        print("Watering Events: None")
    print("-" * 40)  # Divider for readability

def run_cli():
    while True:
        display_menu()
        option = input("Enter option: ")
        if option == "1":
            microgreen_type = input("Enter Microgreen Type: ")
            try:
                seed_weight = float(input("Enter Seed Weight: "))
                app_logic.plant_seeds(microgreen_type, seed_weight)
                print("Seeds planted successfully.")
            except ValueError:
                print("Invalid seed weight. Please enter a valid number.")
        elif option == "2":
            batch_id = input("Enter Batch ID: ")
            app_logic.update_germination(batch_id, datetime.datetime.now())
            print("Germination updated successfully.")
        elif option == "3":
            batch_id = input("Enter Batch ID: ")
            app_logic.record_watering(batch_id, datetime.datetime.now())
            print("Watering recorded successfully.")
        elif option == "4":
            batch_id = input("Enter Batch ID: ")
            try:
                harvest_weight = float(input("Enter Harvest Weight: "))
                app_logic.update_to_harvested(batch_id, datetime.datetime.now(), harvest_weight)
                print("Updated to harvested successfully.")
            except ValueError:
                print("Invalid harvest weight. Please enter a valid number.")
        elif option == "5":
            query_batches()
        elif option == "6":
            print("\nExiting... Goodbye!")
            break
        else:
            print("\nInvalid option. Please try again.")

def query_batches():
    print("\nQuery Options:")
    print("1. By Status")
    print("2. By Type")
    print("3. By Date Range")
    print("4. Single Batch")
    option = input("Select query type: ")
    try:
        batches = []
        if option == "1":
            status = get_status_choice()
            if status:
                batches = app_logic.query_batches_by_status(status)
            else:
                print("Invalid status option.")
                return
        elif option == "2":
            microgreen_type = input("Enter Type: ")
            batches = app_logic.query_batches_by_type(microgreen_type)
        elif option == "3":
            start_date = datetime.datetime.strptime(input("Enter Start Date (YYYY-MM-DD): "), "%Y-%m-%d")
            end_date = datetime.datetime.strptime(input("Enter End Date (YYYY-MM-DD): "), "%Y-%m-%d")
            batches = app_logic.query_batches_by_date_range(start_date, end_date)
        elif option == "4":
            batch_id = input("Enter Batch ID: ")
            batch = app_logic.query_batch(batch_id)
            if batch:
                batches.append(batch)
        else:
            print("Invalid query type.")
            return

        for batch in batches:
            format_and_print_batch(batch)
    except Exception as e:
        print(f"\nAn error occurred: {e}")

if __name__ == "__main__":
    run_cli()
