import streamlit as st
import time
import os
from collections.abc import Collection
from langchain.memory import ChatMessageHistory
from langchain_community.chat_message_histories import (
    StreamlitChatMessageHistory,
)
from langchain_groq import ChatGroq
from langchain_core.prompts import ChatPromptTemplate, MessagesPlaceholder
from langchain_community.vectorstores import MongoDBAtlasVectorSearch
from langchain_community.embeddings import HuggingFaceEmbeddings
from langchain.chains import create_history_aware_retriever, create_retrieval_chain
from langchain.chains.combine_documents import create_stuff_documents_chain
from langchain.output_parsers import ResponseSchema, StructuredOutputParser
from langchain_core.runnables.history import RunnableWithMessageHistory
from langchain_core.chat_history import BaseChatMessageHistory
from langchain.chains import RetrievalQA
import pymongo
import logging
import asyncio
import nest_asyncio
from langchain.docstore.document import Document
import redis
import threading

# Config
nest_asyncio.apply()
logging.basicConfig(level=logging.INFO)
database = "AlertSimAndRemediation"
collection = "alert_embed"
stream_name = "alerts"
index_name = "alert_index"

# Embedding model
embedding_args = {
    "model_name": "BAAI/bge-large-en-v1.5",
    "model_kwargs": {"device": "cpu"},
    "encode_kwargs": {"normalize_embeddings": True}
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
    logging.info("Listening to alerts...")
    try:
        last_id = '$'
        while True:
            entries = r.xread({stream_name: last_id}, block=0, count=None)
            if entries:
                stream, new_entries = entries[0]
                for entry_id, entry_data in new_entries:
                    description, entry_dict = await create_textual_description(entry_data)
                    await save({"content": description, "metadata": entry_dict})
                    # Update the last ID read
                    last_id = entry_id
                    st.toast(description, icon='🔔')
    except KeyboardInterrupt:
        print("Exiting...")

# Start Redis listener in a separate thread
def start_redis_listener():
    try:
        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        loop.run_until_complete(listen_to_alerts(r))
    except Exception as e:
        print(f"Error in Redis listener: {e}")
    finally:
        loop.close()

# Streamlit Application
st.set_page_config(
    page_title="ASMR Query Bot 🔔",
    page_icon="🔔",
    layout="wide",
    initial_sidebar_state="auto",
    menu_items={
        'About': "https://github.com/ankush-003/alerts-simulation-and-remediation"
    }
)

st.title('ASMR Query Bot 🔔')

# Start Redis listener in a separate thread
redis_listener_thread = threading.Thread(target=start_redis_listener)
redis_listener_thread.start()

# vector search
vector_search = MongoDBAtlasVectorSearch.from_connection_string(
    os.environ["MONGO_URI"],
    f"{database}.{collection}",
    embedding_model,
    index_name=index_name,
)

# contextualising prev chats
contextualize_q_system_prompt = """Given a chat history and the latest user question \
which might reference context in the chat history, formulate a standalone question \
which can be understood without the chat history. Do NOT answer the question, \
just reformulate it if needed and otherwise return it as is."""
contextualize_q_prompt = ChatPromptTemplate.from_messages(
    [
        ("system", contextualize_q_system_prompt),
        MessagesPlaceholder("chat_history"),
        ("human", "{input}"),
    ]
)

# prompt
system_prompt = """
You are a helpful query assistant for Alertmanager, an open-source system for monitoring and alerting on system metrics. Your goal is to accurately answer questions related to alerts triggered within the Alertmanager system based on the alert information provided to you. \
You will be given details about specific alerts, including the alert source, severity, category, and any other relevant metadata. Using this information, you should be able to respond to queries about the nature of the alert, what it signifies, potential causes, and recommended actions or troubleshooting steps. \
Your responses should be clear, concise, and tailored to the specific alert details provided, while also drawing from your broader knowledge about Alertmanager and monitoring best practices when relevant. If you cannot provide a satisfactory answer due to insufficient information, politely indicate that and ask for any additional context needed. \
<context>
{context}
</context>
"""

qa_prompt = ChatPromptTemplate.from_messages(
    [
        ("system", system_prompt),
        MessagesPlaceholder("chat_history"),
        ("human", "{input}"),
    ]
)

if "chat_messages" not in st.session_state:
    st.session_state.chat_messages = []

# streamlit history
history = StreamlitChatMessageHistory(key="chat_messages")

# Initialize chat history
if len(history.messages) == 0:
    history.add_ai_message("Hey I am ASMR Query Bot, how can i help you ?")

with st.sidebar:
    st.title('Settings ⚙️')
    st.subheader('Models and parameters')
    selected_model = st.sidebar.selectbox('Choose a model', ['Llama3-8B', 'Llama3-70B', 'Mixtral-8x7B'], key='selected_model')
    if selected_model == 'Mixtral-8x7B':
        model_name="mixtral-8x7b-32768"
    elif selected_model == 'Llama3-70B':
        model_name='Llama3-70b-8192'
    elif selected_model == 'Llama3-8B':
        model_name='Llama3-8b-8192'
    temp = st.sidebar.slider('temperature', min_value=0.01, max_value=1.0, value=0.0, step=0.01)
    k = st.sidebar.slider('number of docs retrieved', min_value=1, max_value=20, value=2, step=1)

def get_response(query, config):
  chat = ChatGroq(temperature=temp, model_name=model_name)
  qa_retriever = vector_search.as_retriever(
      search_type="similarity",
      search_kwargs={"k": k},
  )
  history_aware_retriever = create_history_aware_retriever(
      chat, qa_retriever, contextualize_q_prompt
  )
  question_answer_chain = create_stuff_documents_chain(chat, qa_prompt)
  rag_chain = create_retrieval_chain(history_aware_retriever, question_answer_chain)
  conversational_rag_chain = RunnableWithMessageHistory(
      rag_chain,
      lambda session_id: history,
      input_messages_key="input",
      history_messages_key="chat_history",
      output_messages_key="answer",
  )
  return conversational_rag_chain.invoke({"input": prompt}, config=config)

def clear_chat_history():
    st.session_state.chat_messages = []
    history.add_ai_message("Hey I am ASMR Query Bot, how can i help you ?")

st.sidebar.button('Clear Chat History', on_click=clear_chat_history)  

for msg in history.messages:
    st.chat_message(msg.type).write(msg.content)

# preprocessing context
def format_docs_with_metadata(docs):
    formatted_docs = []
    for i, doc in enumerate(docs, start=1):
        metadata_str = "\n".join([f"**{key}**: `{value}`\n" for key, value in doc.metadata.items() if key != "embedding"])
        formatted_doc = f"- {doc.page_content}\n\n**Metadata:**\n{metadata_str}"
        formatted_docs.append(formatted_doc)   
    return "\n\n".join(formatted_docs)

def stream_data(response):
  for word in response.split(" "):
        yield word + " "
        time.sleep(0.05)

if prompt := st.chat_input():
    with st.chat_message("Human"):
        st.markdown(prompt)

    # As usual, new messages are added to StreamlitChatMessageHistory when the Chain is called.
    config = {"configurable": {"session_id": "any"}}
    res = get_response(prompt, config)

    with st.chat_message("AI"):
      st.write_stream(stream_data(res['answer']))
      with st.popover("View Source"):
        st.markdown("### Source Alerts 📢")
        st.markdown(format_docs_with_metadata(res['context']))