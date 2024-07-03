from langchain_mongodb import MongoDBAtlasVectorSearch
from langchain_huggingface import HuggingFaceEmbeddings
# from dotenv import load_dotenv
import os
import pymongo
import logging
import nest_asyncio
from langchain.docstore.document import Document
import redis
import threading
import asyncio
import gradio as gr
import pandas as pd
import plotly.express as px
from datetime import datetime, timedelta

# config
# nest_asyncio.apply()
logging.basicConfig(level = logging.INFO)
database = "AlertSimAndRemediation"
collection = "alert_embed"
stream_name = "alerts"

# Global variables to store alert information
latest_alert = "No alerts yet."
alert_count = 0

# embedding model
embedding_args = {
    "model_name" : "BAAI/bge-large-en-v1.5",
    "model_kwargs" : {"device": "cpu"},
    "encode_kwargs" : {"normalize_embeddings": True}
}
embedding_model = HuggingFaceEmbeddings(**embedding_args)

# Mongo Connection
connection = pymongo.MongoClient(os.environ["MONGO_URI"])
alert_collection = connection[database][collection]

# Redis connection
r = redis.Redis(host=os.environ['REDIS_HOST'], password=os.environ['REDIS_PWD'], port=16652)

# Preprocessing
async def create_textual_description(entry_data):
    entry_dict = {k.decode(): v.decode() for k, v in entry_data.items()}

    category = entry_dict["Category"]
    created_at = entry_dict["CreatedAt"]
    acknowledged = "Acknowledged" if entry_dict["Acknowledged"] == "1" else "Not Acknowledged"
    remedy = entry_dict["Remedy"]
    severity = entry_dict["Severity"]
    source = entry_dict["Source"]
    node = entry_dict["node"]

    description = f"A {severity} alert of category {category} was raised from the {source} source for node {node} at {created_at}. The alert is {acknowledged}. The recommended remedy is: {remedy}."

    return description, entry_dict

# Saving alert doc
async def save(entry):
    vector_search = MongoDBAtlasVectorSearch.from_documents(
        documents=[Document(
            page_content=entry["content"],
            metadata=entry["metadata"]
        )],
        embedding=embedding_model,
        collection=alert_collection,
        index_name="alert_index",
    )
    logging.info("Alerts stored successfully!")

# Listening to alert stream
async def listen_to_alerts(r):
    global latest_alert, alert_count
    try:
        last_id = '$'

        while True:
            entries = r.xread({stream_name: last_id}, block=0, count=None)

            if entries:
                stream, new_entries = entries[0]

                for entry_id, entry_data in new_entries:
                    description, entry_dict = await create_textual_description(entry_data)
                    await save({
                        "content" : description,
                        "metadata" : entry_dict
                    })
                    print(description)
                    latest_alert = description
                    alert_count += 1
                    # Update the last ID read
                    last_id = entry_id
                    await asyncio.sleep(1) 

    except KeyboardInterrupt:
        print("Exiting...")

def run_alert_listener():
    asyncio.run(listen_to_alerts(r))

# Start the alert listener thread
alert_thread = threading.Thread(target=run_alert_listener)
alert_thread.start()

# dashboard

# Function to get the latest alert and count
def get_latest_alert():
    global latest_alert, alert_count
    return latest_alert, alert_count

with gr.Blocks(theme=gr.themes.Soft()) as app:
    gr.Markdown("# 🚨 Alert Dashboard 🚨")
    
    with gr.Row():
        with gr.Column():
            latest_alert_box = gr.Textbox(label="Latest Alert", lines=3, interactive=False)
            alert_count_box = gr.Number(label="Alert Count", interactive=False)
            refresh_button = gr.Button("Refresh", variant="primary")
    
    refresh_button.click(get_latest_alert, inputs=None, outputs=[latest_alert_box, alert_count_box])
    
    app.load(get_latest_alert, inputs=None, outputs=[latest_alert_box, alert_count_box])
    
    # Auto-refresh every 30 seconds
    app.load(get_latest_alert, inputs=None, outputs=[latest_alert_box, alert_count_box], every=30)

# Launch the app

# Launch the app
# if __name__ == "__main__":
app.launch()