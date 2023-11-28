import os
import logging
import requests
from bs4 import BeautifulSoup
from openai import OpenAI

logging.basicConfig(filename='application.log', filemode='w', level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

class OpenAIChatbot:
    def __init__(self, api_key, model="gpt-4-1106-preview"):
        self.client = OpenAI(api_key=api_key)
        self.model = model
        self.conversation_history = []
        self.context = ""  # Context storage

    def get_user_input(self):
        return input("\nEnter query or type 'exit' to quit: ")

    def fetch_html_content(self, url):
        try:
            response = requests.get(url)
            response.raise_for_status()
            logging.info(f"Fetched URL: {url} with status: {response.status_code}")
            logging.info(f"Snippet of HTML content: {response.text[:500]}")  # Log first 500 characters
            return response.text
        except requests.RequestException as e:
            logging.error(f"Error fetching URL: {e}")
            return None

    def extract_text_from_html(self, html_content):
        soup = BeautifulSoup(html_content, 'html.parser')
        extracted_text = soup.get_text()
        logging.info(f"Snippet of extracted text: {extracted_text[:500]}")  # Log first 500 characters
        return extracted_text

    def update_context(self, new_context):
        self.context = new_context

    def query_openai(self, user_input):
        try:
            print("\nProcessing...", end='\r')
            
            # Format the query with context only if context is not empty
            if self.context.strip():
                query = "Context: " + self.context + "\n\n" + "Query: " + user_input
            else:
                query = user_input  # If no context, use user input directly

            self.conversation_history.append({"role": "user", "content": query})
            completion = self.client.chat.completions.create(
                model=self.model, 
                messages=self.conversation_history
            )
            print(" " * 50, end='\r')  # Clear the processing message
            response = completion.choices[0].message.content
            self.conversation_history.append({"role": "assistant", "content": response})
            return response
        except Exception as e:
            logging.error(f"An error occurred: {e}")
            return None

    def run(self):
        while True:
            user_input = self.get_user_input()

            if user_input.lower() in ['exit', 'quit']:
                print("\nExiting...")
                break

            if user_input.lower().startswith('http'):
                html_content = self.fetch_html_content(user_input)
                if html_content:
                    text = self.extract_text_from_html(html_content)
                    self.update_context(text)
                    print("\nContext updated with extracted text from URL.")
                    continue  # Skip querying after updating context

            response = self.query_openai(user_input)
            if response:
                print("\nAI Response:", response)

def main():
    api_key = os.environ["OPENAI_API_KEY"]
    chatbot = OpenAIChatbot(api_key)

    # Example URL for initial context
    context_url = "https://www.sec.gov/Archives/edgar/data/1113027/000119312504116689/dn23c3a.htm"
    html_content = chatbot.fetch_html_content(context_url)
    if html_content:
        text = chatbot.extract_text_from_html(html_content)
        chatbot.update_context(text)
        print("Initial context set from URL:", context_url)
    print(chatbot.context)
    #chatbot.run()

if __name__ == "__main__":
    main()
