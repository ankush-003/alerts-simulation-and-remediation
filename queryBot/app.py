from chainlit.server import app
from fastapi import Request, File, UploadFile
import os
import chainlit as cl
import requests
import random
from huggingface_hub import login
from llama_index.core import VectorStoreIndex,SimpleDirectoryReader,ServiceContext,PromptTemplate
from llama_index.llms.huggingface import HuggingFaceLLM
# from llama_index.llms.ollama import Ollama
from llama_index.core.settings import Settings
import torch
from llama_index.core.output_parsers import PydanticOutputParser
from typing import List
from transformers import BitsAndBytesConfig
from langchain.embeddings.huggingface import HuggingFaceEmbeddings
from llama_index.core import ServiceContext
from llama_index.embeddings.langchain import LangchainEmbedding
from llama_index.core.memory import ChatMemoryBuffer
from huggingface_hub import login
from chainlit.input_widget import TextInput, Slider
from llama_index.readers.mongodb import SimpleMongoReader

login(token=os.environ["HF_TOKEN"])

system_prompt="""
You are a helpful query assistant for Alertmanager, an open-source system for monitoring and alerting on system metrics. Your goal is to accurately answer questions related to alerts triggered within the Alertmanager system based on the alert information provided to you.
You will be given details about specific alerts, including the alert source, severity, category, and any other relevant metadata. Using this information, you should be able to respond to queries about the nature of the alert, what it signifies, potential causes, and recommended actions or troubleshooting steps.
Your responses should be clear, concise, and tailored to the specific alert details provided, while also drawing from your broader knowledge about Alertmanager and monitoring best practices when relevant. If you cannot provide a satisfactory answer due to insufficient information, politely indicate that and ask for any additional context needed.
"""
query_wrapper_prompt=PromptTemplate("<|USER|>{query_str}<|ASSISTANT|>")
query_engine = None
documents = None
index = None

# Loading LLM
llm = HuggingFaceLLM(
    context_window=4096,
    max_new_tokens=256,
    generate_kwargs={"temperature": 0.0, "do_sample": False},
    system_prompt=system_prompt,
    query_wrapper_prompt=query_wrapper_prompt,
    tokenizer_name="meta-llama/Llama-2-7b-chat-hf",
    model_name="meta-llama/Llama-2-7b-chat-hf",
    device_map="auto",
    # uncomment this if using CUDA to reduce memory usage
    model_kwargs={"torch_dtype": torch.float16 , "load_in_4bit":True}
)

embed_model=LangchainEmbedding(
    HuggingFaceEmbeddings(model_name="sentence-transformers/all-mpnet-base-v2"))

Settings.llm = llm
Settings.embed_model = embed_model
Settings.chunk_size = 1024

print("Models Loaded Successfully!")

@cl.on_chat_start
async def on_chat_start():
    global index
    # setup
    settings = await cl.ChatSettings(
        [
            TextInput(id="MONGOURI", label="MongoDB Uri", initial=str(os.environ["MONGO_URI"])),
            TextInput(id="DB_NAME", label="DataBase Name", initial="users"),
            TextInput(id="COLLECTION", label="Collection Name", initial="alerts"),
        ]
    ).send()
    if settings["MONGOURI"] == "" or settings["DB_NAME"] == "" or settings["COLLECTION"] == "":
      await cl.Message(content="Please setup the configurations", disable_feedback=True).send()
    cl.user_session.set("MONGOURI", settings["MONGOURI"])
    cl.user_session.set("DB_NAME", settings["DB_NAME"])
    cl.user_session.set("COLLECTION", settings["COLLECTION"])

    uri = os.environ["MONGO_URI"]
    db_name = "users"
    collection_name = "alerts"
    # query_dict is passed into db.collection.find()
    query_dict = {}
    field_names = ["Source", "Category", "CreatedAt", "Remedy", "Severity", "Source"]

    # using mongoreader
    reader = SimpleMongoReader(uri=uri)
    documents = reader.load_data(
        db_name, collection_name, field_names, query_dict=query_dict
    )

    print(f"Fetched Documents...\n{documents}")

    # creating vector store
    index = VectorStoreIndex.from_documents(documents)
    # query_engine = index.as_query_engine()
    # cl.user_session.set("query_engine", query_engine)
    await cl.Message(content="Hi! I am ASMR Query Bot how can i help you ?").send()


@cl.on_settings_update
async def setup_agent(settings):
    cl.user_session.set("MONGOURI", settings["MONGOURI"])
    cl.user_session.set("DB_NAME", settings["DB_NAME"])
    cl.user_session.set("COLLECTION", settings["COLLECTION"])
    print("on_settings_update", settings)


@cl.on_message
async def main(message: cl.Message):
    global index
    user_query = message.content
    query_engine = index.as_query_engine()
    response = query_engine.query(user_query)
    elements = [
        cl.Text(name="response", content=response.response, display="inline")
    ]
    print(response)
    await cl.Message(content="Response to your query!", elements=elements).send()