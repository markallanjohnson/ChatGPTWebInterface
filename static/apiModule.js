// apiModule.js

// This module will depend on the UI module to update the UI based on API responses
import AppState from './appStateModule.js';
import uiModule from './uiModule.js';

const apiModule = (() => {
    const baseURL = '/api'; // The base URL for your API

    function sendRequest(endpoint, method = 'GET', data = null) {
        const options = {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
        };
    
        if (data) {
            options.body = JSON.stringify(data);
        }
    
        return fetch(`${baseURL}${endpoint}`, options)
            .then(response => {
                if (!response.ok) {
                    throw new Error(`Server returned ${response.status}: ${response.statusText}`);
                }
                // Return the raw response object
                return response; // Do not parse here
            })
            .catch(error => console.error('API request failed:', error));
    }
    
    function sendMessage(event) {
        if (event.key === "Enter") {
            AppState.customLog('Enter key pressed');
            // Prevents the default action to ensure form is not submitted traditionally
            event.preventDefault();

            let message = uiModule.inputBox.value.trim();
            if (message === '') return; // Prevent sending empty messages
            uiModule.inputBox.value = '';

            // Append message and show loading indicator
            uiModule.appendMessage('You: ' + message, 'user');
            uiModule.showLoadingIndicator(true);
            AppState.customLog(`You: ${message}`);
            AppState.customLog('Current session ID: ' + AppState.getCurrentSessionId());
            // If there is no session, create one, otherwise send the message.
            if (AppState.getCurrentSessionId() === null) {
                createAndSendNewChat(message);
            } else {
                sendAndDisplayMessage(message);
            }
        }
    }
    
    function createNewChat() {
        // Prepare for new chat session
        AppState.customLog("Preparring for new chat session...");
        AppState.setCurrentSessionId(null);
        uiModule.updateActiveSessionDisplay();
    }
    
    function createAndSendNewChat(message) {
        // If there is no session, create one, otherwise send the message.
        AppState.customLog("Creating new chat session...");
        if (AppState.getCurrentSessionId() === null) {
            return sendRequest('/new-session', 'POST')
                .then(response => response.json())
                .then(data => {
                    AppState.setCurrentSessionId(data.session_id);
                    AppState.customLog(`Created new session ${data.session_id}`);
                    AppState.customLog(`Current session ID: ${AppState.getCurrentSessionId()}`);
                    return sendAndDisplayMessage(message); 
                });
        } else {
            return sendAndDisplayMessage(message);
        }
    }

    function sendAndDisplayMessage(message) {
        
        AppState.customLog(`Sending and displaying message, session id: ${AppState.getCurrentSessionId()}`);
        const endpoint = `/query?session_id=${encodeURIComponent(AppState.getCurrentSessionId())}&query=${encodeURIComponent(message)}`;
        return sendRequest(endpoint, 'POST', { query: message })
            .then(response => response.text())
            .then(data => {
                AppState.customLog(`AI: ${data}`);
                uiModule.showLoadingIndicator(false); // Hide loading indicator
                uiModule.appendMessage('AI: ' + data, 'ai');
            }).catch(error => {
                AppState.customLog('Error sending message:' + error, true);
                uiModule.showLoadingIndicator(false); // Hide loading indicator
            });
    }
    
    function loadSessions() {
        return sendRequest('/get-sessions')
            .then(sessions => uiModule.displaySessions(sessions))
            .catch(error => console.error('Error loading sessions:', error));
    }

    function deleteChat(sessionId) {
        return sendRequest(`/delete-session?session_id=${sessionId}`, 'DELETE')
            .then(() => uiModule.removeSessionItem(sessionId))
            .catch(error => console.error('Error deleting chat:', error));
    }

    function renameSession(sessionId, newName) {
        return sendRequest('/rename-session', 'PUT', { sessionId, newName })
            .then(() => uiModule.updateSessionName(sessionId, newName))
            .catch(error => console.error('Error renaming session:', error));
    }

    function loadSessionsIntoModal() {
        return sendRequest('/get-sessions', 'GET')
            .then(sessions => {
                if (!sessions || !Array.isArray(sessions)) {
                    throw new Error('Sessions is not an array');
                }
                uiModule.populateSessionsInModal(sessions); // Assuming you have this function in uiModule
            })
            .catch(error => console.error('Error fetching sessions:', error));
    }

    function loadSession(sessionId) {
        if (!sessionId) {
            console.error('No session ID provided to loadSession');
            return;
        }
    
        return sendRequest(`/get-session-history?session_id=${sessionId}`, 'GET')
            .then(response => {
                // Check if the response is JSON or plain text and handle accordingly
                const contentType = response.headers.get('Content-Type');
                return contentType && contentType.includes('application/json') ? response.json() : response.text();
            })
            .then(history => {
                if (Array.isArray(history)) {
                    uiModule.chatBox.innerHTML = '';
                    history.forEach(item => {
                        if (item && item.content && item.role) {
                            uiModule.appendMessage(item.content, item.role === 'user' ? 'user' : 'ai');
                        }
                    });
                } else {
                    console.warn('Expected history to be an array, but got:', history);
                }
                AppState.setCurrentSessionId(sessionId);
                uiModule.updateActiveSessionDisplay();
            })
            .catch(error => {
                console.error('Error loading session:', error);
                uiModule.appendMessage(`Error: ${error.message}`, 'system');
            })
            .finally(() => uiModule.showLoadingIndicator(false));
    }
    
    return {
        sendMessage,
        createNewChat,
        loadSessions,
        deleteChat,
        renameSession,
        loadSessionsIntoModal,
        loadSession,
        createAndSendNewChat,
    };
})();

// Export the API Module
export default apiModule;
