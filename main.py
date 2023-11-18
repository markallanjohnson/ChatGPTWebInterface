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

def query_openai(client, model, user_input):
    """Send a query to OpenAI and return the response."""
    try:
        completion = client.chat.completions.create(
            model=model, 
            messages=[{"role": "user", "content": user_input}]
        )
        return completion.choices[0].message.content
    except Exception as e:
        logging.error(f"An error occurred: {e}")
        return None

def main():
    client = get_openai_client(api_key=os.environ["OPENAI_API_KEY"])
    model = "gpt-4-1106-preview" 

    while True:
        user_input = get_user_input()

        if user_input.lower() in ['exit', 'quit']:
            print("Exiting...")
            break

        response = query_openai(client, model, user_input)
        if response:
            print("AI Response:", response)

if __name__ == "__main__":
    main()
