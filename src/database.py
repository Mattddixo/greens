import sqlite3
from sqlite3 import Error
import os
import datetime
import json
from models import MicrogreenBatch, WateringEvent

DATABASE_NAME = "microgreens_tracker.db"

def create_connection():
    conn = None
    try:
        conn = sqlite3.connect(DATABASE_NAME)
    except Error as e:
        print(e)
    return conn

def setup_database():
    conn = create_connection()
    if conn is not None:
        with conn:
            conn.execute("""CREATE TABLE IF NOT EXISTS microgreen_batches (
                            id INTEGER PRIMARY KEY AUTOINCREMENT,
                            batchID TEXT NOT NULL,
                            type TEXT NOT NULL,
                            seedWeight REAL NOT NULL,
                            plantDate TEXT NOT NULL,
                            germinateDate TEXT,
                            harvestDate TEXT,
                            harvestWeight REAL,
                            status TEXT NOT NULL
                        );""")
            
            conn.execute("""CREATE TABLE IF NOT EXISTS watering_events (
                            id INTEGER PRIMARY KEY AUTOINCREMENT,
                            batchID TEXT NOT NULL,
                            waterDate TEXT NOT NULL,
                            FOREIGN KEY (batchID) REFERENCES microgreen_batches (batchID)
                        );""")
            print("Database setup completed.")
    else:
        print("Error! Cannot create the database connection.")

# Existing CRUD functions...

def export_to_json(file_path='database_export.json'):
    """Exports the database tables to a JSON file for easier readability."""
    conn = create_connection()
    data = {'microgreen_batches': [], 'watering_events': []}
    
    with conn:
        # Export microgreen_batches
        cur = conn.cursor()
        cur.execute("SELECT * FROM microgreen_batches")
        rows = cur.fetchall()
        columns = [column[0] for column in cur.description]
        for row in rows:
            data['microgreen_batches'].append(dict(zip(columns, row)))
        
        # Export watering_events
        cur.execute("SELECT * FROM watering_events")
        rows = cur.fetchall()
        columns = [column[0] for column in cur.description]
        for row in rows:
            data['watering_events'].append(dict(zip(columns, row)))
    
    # Handle datetime serialization
    def datetime_converter(o):
        if isinstance(o, datetime.datetime):
            return o.__str__()

    with open(file_path, 'w') as json_file:
        json.dump(data, json_file, default=datetime_converter, indent=4)

def insert_microgreen_batch(batch: MicrogreenBatch):
    """Insert a new microgreen batch record."""
    conn = create_connection()
    with conn:
        conn.execute('''INSERT INTO microgreen_batches(batchID, type, seedWeight, plantDate, germinateDate, harvestDate, harvestWeight, status)
                        VALUES(?,?,?,?,?,?,?,?)''', 
                     (batch.batchID, batch.type, batch.seedWeight, batch.plantDate.isoformat(), 
                      batch.germinateDate.isoformat() if batch.germinateDate else None, 
                      batch.harvestDate.isoformat() if batch.harvestDate else None, 
                      batch.harvestWeight, batch.status))

def update_microgreen_batch(batchID: str, updates: dict):
    """Update details of a microgreen batch record identified by batchID."""
    conn = create_connection()
    with conn:
        cur = conn.cursor()
        set_clause = ', '.join([f"{key} = ?" for key in updates])
        values = list(updates.values())
        values.append(batchID)
        cur.execute(f"UPDATE microgreen_batches SET {set_clause} WHERE batchID = ?", values)

def query_microgreen_batches(condition: str = "", params: tuple = ()):
    """Query microgreen batches with optional conditions and include their watering events."""
    conn = create_connection()
    batches = []
    with conn:
        cur = conn.cursor()
        cur.execute(f"SELECT * FROM microgreen_batches {condition}", params)
        batch_rows = cur.fetchall()
        for batch_row in batch_rows:
            batchID = batch_row[1]
            # Fetch watering events for this batch
            cur.execute("SELECT waterDate FROM watering_events WHERE batchID = ?", (batchID,))
            watering_event_rows = cur.fetchall()
            watering_events = [WateringEvent(waterDate=datetime.datetime.fromisoformat(row[0])) for row in watering_event_rows]
            
            batch = MicrogreenBatch(
                batchID=batch_row[1],
                type=batch_row[2],
                seedWeight=batch_row[3],
                plantDate=datetime.datetime.fromisoformat(batch_row[4]),
                status=batch_row[8],
                germinateDate=datetime.datetime.fromisoformat(batch_row[5]) if batch_row[5] else None,
                harvestDate=datetime.datetime.fromisoformat(batch_row[6]) if batch_row[6] else None,
                harvestWeight=batch_row[7],
                wateringEvents=watering_events
            )
            batches.append(batch)
    return batches

def insert_watering_event(batchID: str, watering_event: WateringEvent):
    conn = create_connection()
    with conn:
        cur = conn.cursor()
        cur.execute('''INSERT INTO watering_events(batchID, waterDate)
                       VALUES(?,?)''', 
                     (batchID, watering_event.waterDate.isoformat()))

def export_to_json(file_path='database_export.json'):
    """Exports the database tables to a JSON file for easier readability."""
    conn = create_connection()
    data = {'microgreen_batches': [], 'watering_events': []}
    
    with conn:
        # Export microgreen_batches
        cur = conn.cursor()
        cur.execute("SELECT * FROM microgreen_batches")
        rows = cur.fetchall()
        columns = [column[0] for column in cur.description]
        for row in rows:
            data['microgreen_batches'].append(dict(zip(columns, row)))
        
        # Export watering_events
        cur.execute("SELECT * FROM watering_events")
        rows = cur.fetchall()
        columns = [column[0] for column in cur.description]
        for row in rows:
            data['watering_events'].append(dict(zip(columns, row)))
    
    # Handle datetime serialization
    def datetime_converter(o):
        if isinstance(o, datetime.datetime):
            return o.__str__()

    with open(file_path, 'w') as json_file:
        json.dump(data, json_file, default=datetime_converter, indent=4)

setup_database()  # Automatically sets up database on module import
