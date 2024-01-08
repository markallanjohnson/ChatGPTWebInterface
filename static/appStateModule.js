import apiModule from './apiModule.js';
import uiModule from './uiModule.js';

const AppState = (() => {
    let state = {
      currentSessionId: null,
      sessionNames: {},
      // Add other state properties as needed
    };
  
    function getCurrentSessionId() {
      return state.currentSessionId;
    }
  
    function setCurrentSessionId(sessionId) {
      state.currentSessionId = sessionId;
      // Any additional logic or side effects when setting the session ID
    }
  
    function getSessionName(sessionId) {
      return state.sessionNames[sessionId] || 'Unnamed Session';
    }
  
    function setSessionName(sessionId, name) {
      state.sessionNames[sessionId] = name;
      // Any additional logic or side effects when setting the session name
    }
  
    function getAllSessions() {
      return state.sessionNames;
      // Depending on the structure, you might need a more complex function to return sessions
    }
    
    function customLog(message, isError = false) {
        let logObject = { message, isError };
        fetch('/log', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(logObject),
        });
    }

    function selectSession(sessionId) {
        // Set the current session ID
        state.currentSessionId = sessionId;
        // Fetch and load the session history
        apiModule.loadSession(sessionId);
        customLog(`Selected session ${sessionId}`);
        // Close the modal
        uiModule.closeModal();
    }

    function storeSessionNames(sessions) {
      state.sessionNames = {};
        console.log('Storing session names', sessions);
        sessions.forEach(session => {
            state.sessionNames[session.session_id] = session.name;
        });
        console.log('Updated session names', state.sessionNames);
    }
  
    // Export public functions
    return {
      getCurrentSessionId,
      setCurrentSessionId,
      getSessionName,
      setSessionName,
      getAllSessions,
      customLog,
      selectSession,
      storeSessionNames,
    };
  })();
  
  // Export the AppState
  export default AppState;
  