import os
import sys
import json
import logging
from openai import OpenAI

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

def get_openai_client(api_key):
    """Create and return an OpenAI client with the given API key."""
    return OpenAI(api_key=api_key)

def query_openai(client, model, conversation_history):
    """Send a query to OpenAI and return the response."""
    try:
        completion = client.chat.completions.create(
            model=model, 
            messages=conversation_history
        )
        response = completion.choices[0].message.content
        return response
    except Exception as e:
        logging.error(f"An error occurred: {e}")
        return None

def main():
    client = get_openai_client(api_key=os.environ["OPENAI_API_KEY"])
    #model = "gpt-4-vision-preview"
    model = "gpt-4-1106-preview"

    # Read the conversation history from stdin
    conversation_history_json = sys.stdin.read()
    conversation_history = json.loads(conversation_history_json)

    response = query_openai(client, model, conversation_history)
    if response:
        # Write the response to stdout
        print(response)

if __name__ == "__main__":
    main()
