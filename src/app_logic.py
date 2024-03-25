import datetime
from typing import List
from models import MicrogreenBatch, WateringEvent
import database

def plant_seeds(microgreen_type: str, seed_weight: float) -> MicrogreenBatch:
    now = datetime.datetime.now()
    batch_id = f"MG-{now.strftime('%Y%m%d%H%M%S')}"
    batch = MicrogreenBatch(
        batchID=batch_id,
        type=microgreen_type,
        seedWeight=seed_weight,
        plantDate=now,
        status="Planted",
        germinateDate=None,
        harvestDate=None,
        harvestWeight=None,
        wateringEvents=[]
    )
    database.insert_microgreen_batch(batch)
    database.export_to_json()  # Export database to JSON after update
    return batch

def update_germination(batch_id: str, germinate_date: datetime.datetime):
    updates = {'germinateDate': germinate_date.isoformat(), 'status': 'Germinated'}
    database.update_microgreen_batch(batch_id, updates)
    database.export_to_json()  # Export database to JSON after update

def record_watering(batch_id: str, water_date: datetime.datetime):
    watering_event = WateringEvent(waterDate=water_date)
    database.insert_watering_event(batch_id, watering_event)
    database.export_to_json()  # Export database to JSON after update

def update_to_harvested(batch_id: str, harvest_date: datetime.datetime, harvest_weight: float):
    updates = {'harvestDate': harvest_date.isoformat(), 'harvestWeight': harvest_weight, 'status': 'Harvested'}
    database.update_microgreen_batch(batch_id, updates)
    database.export_to_json()  # Export database to JSON after update

def query_batches_by_status(status: str) -> List[MicrogreenBatch]:
    condition = "WHERE status = ?"
    params = (status,)
    return database.query_microgreen_batches(condition, params)

def query_batches_by_type(microgreen_type: str) -> List[MicrogreenBatch]:
    condition = "WHERE type = ?"
    params = (microgreen_type,)
    return database.query_microgreen_batches(condition, params)

def query_batches_by_date_range(start_date: datetime.datetime, end_date: datetime.datetime) -> List[MicrogreenBatch]:
    condition = "WHERE plantDate >= ? AND plantDate <= ?"
    params = (start_date.isoformat(), end_date.isoformat())
    return database.query_microgreen_batches(condition, params)

def query_batch(batch_id: str) -> MicrogreenBatch:
    condition = "WHERE batchID = ?"
    params = (batch_id,)
    batch_records = database.query_microgreen_batches(condition, params)
    if batch_records:
        return batch_records[0]  # Assuming query_microgreen_batches returns MicrogreenBatch objects
    return None
