import os
import logging
from openai import OpenAI

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

def get_openai_client(api_key):
    """Create and return an OpenAI client with the given API key."""
    return OpenAI(api_key=api_key)

def get_user_input():
    """Prompt the user for input and return the input."""
    return input("Enter your query (or type 'exit' to quit): ")

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
    model = "gpt-4-1106-preview"
    conversation_history = []

    while True:
        user_input = get_user_input()

        if user_input.lower() in ['exit', 'quit']:
            print("Exiting...")
            break

        conversation_history.append({"role": "user", "content": user_input})

        response = query_openai(client, model, conversation_history)
        if response:
            print("AI Response:", response)
            conversation_history.append({"role": "assistant", "content": response})

if __name__ == "__main__":
    main()
