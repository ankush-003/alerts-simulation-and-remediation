# QueryBot 🤖

a real-time alert processing system with an interactive chatbot interface. It utilizes Redis Streams for data ingestion, MongoDB for vector storage, and LangChain for creating a conversational AI model.

- [Embedding Service on 🤗](https://ankush-003-asmr-embedding-service.hf.space/)
- [ChatBot on 🤗](https://ankush-003-asmr-query-bot.hf.space)

## Features

- Real-time alert ingestion using Redis Streams
- Alert vectorization using Sentence Transformers
- Vector storage in MongoDB
- Interactive chatbot interface with Streamlit
- Context-aware responses using LangChain and Groq
- Support for multiple language models (Llama3-8B, Llama3-70B, Mixtral-8x7B)
- Customizable model parameters (temperature, number of retrieved documents)

## System Architecture

### Data Pipeline:

- Listens to alerts from Redis Streams
- Processes and vectorizes alert data
- Stores vectorized data in MongoDB


### Chatbot:

- Utilizes Streamlit for the user interface
- Implements LangChain for conversational AI
- Uses Groq for language model inference
- Provides responses with referenced context


<!-- 
## Dependencies

Redis
MongoDB
Sentence Transformers
LangChain
PyMongo
AsyncIO
Streamlit
Groq -->

## Setup

Install required packages:
Copypip install redis pymongo sentence-transformers langchain streamlit groq

Set up environment variables:

- `MONGO_URI`: MongoDB connection string
- `REDIS_HOST`: Redis host address
- `REDIS_PWD`: Redis password

## Data Pipeline (Embedding Service)
The data pipeline performs the following steps:

- Connects to Redis and MongoDB
- Listens to the "alerts" stream in Redis
- Processes incoming alerts:
    - Creates a textual description of the alert
    - Vectorizes the description using the BAAI/bge-large-en-v1.5 model
    - Stores the vectorized alert in MongoDB



## Chatbot Details
The chatbot is implemented using Streamlit and provides the following features:

- Interactive chat interface
- Model selection (Llama3-8B, Llama3-70B, Mixtral-8x7B)
- Adjustable temperature and number of retrieved documents
- Context-aware responses using LangChain and Groq
- Ability to view source alerts for each response

## Screenshots

### Embedding Service

![image](https://github.com/ankush-003/alerts-simulation-and-remediation/assets/94037471/eea1cc6a-392b-437c-b4c5-53334d586bb4)

### ChatBot

![image](https://github.com/ankush-003/alerts-simulation-and-remediation/assets/94037471/100277cd-52cc-41f6-bcea-da7bc836d97e)


## Implementation Details

- Uses `MongoDBAtlasVectorSearch` for retrieving relevant alert contexts
- Implements a context-aware retriever to understand chat history
- Utilizes Groq's language models for generating responses
- Streams the response for a more interactive experience

## Configuration

- database: Name of the MongoDB database (default: "AlertSimAndRemediation")
- collection: Name of the MongoDB collection (default: "alert_embed")
- stream_name: Name of the Redis stream (default: "alerts")
- index_name: Name of the vector search index (default: "alert_index")
