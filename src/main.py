import database  # Ensure database setup is performed
import cli_interface

def main():
    # Explicit database setup call could be here if setup was not automatic
    # database.setup()  # Uncomment if setup() is not called automatically on import
    
    print("Welcome to the Microgreen Tracker Application!")
    cli_interface.run_cli()

if __name__ == "__main__":
    main()
