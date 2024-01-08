// uiModule.js
import apiModule from './apiModule.js';
import AppState from './appStateModule.js';

// Get UI elements
const uiModule = (() => {
    const sidebar = document.getElementById('sidebar');
    const chatBox = document.getElementById('chat-box');
    const inputBox = document.getElementById('input-box');
    const projectSelectionModal = document.getElementById('project-selection-modal');
    const projectList = document.getElementById('project-list');
    const loadingIndicator = document.getElementById('loadingIndicator');
    const newChatButton = document.getElementById('new-chat-button');

        // Event listeners setup for UI elements
        function setupEventListeners() {
            newChatButton.addEventListener('click', openModal);
            inputBox.addEventListener('keypress', apiModule.sendMessage);

            document.querySelector('.close-button').addEventListener('click', closeModal);
            document.getElementById('project-selection-modal').addEventListener('click', handleModalClick);
            document.querySelector('#project-selection-modal button').addEventListener('click', apiModule.createNewChat);
            // Event delegation for session items and dropdowns
            document.getElementById('project-list').addEventListener('click', handleProjectListClick);
            document.addEventListener('click', closeAllDropdownsOutsideClick);
            // input box event listener
            document.getElementById('input-box').addEventListener('keypress', apiModule.sendMessage);
            // Modal event listeners
            projectList.addEventListener('click', function(event) {
                if (event.target && event.target.matches('.project-list-item span')) {
                    const sessionId = event.target.closest('.project-list-item').dataset.sessionId;
                    AppState.selectSession(sessionId); // Use selectSession to handle session selection
                }
            });
    
            // Dropdown event listeners
            document.querySelectorAll('.dropbtn').forEach(function(dropbtn) {
                dropbtn.onclick = function(event) {
                    // Prevent the event from closing dropdowns
                    event.stopPropagation();
                    // Close all other dropdowns
                    closeAllDropdowns();
                    // Toggle this dropdown
                    this.nextElementSibling.classList.toggle('show');
                };
            });

            // Handle all document click events here and delegate as necessary
            document.addEventListener('click', function(event) {
                
                if (projectList.contains(event.target)) {
                    if (event.target.matches('.session-entry span')) {
                        const sessionId = event.target.dataset.sessionId;
                        loadSession(sessionId);
                    }
                    closeAllDropdowns();
                }
            });
        
            // Event listener for text selection
            document.addEventListener('mouseup', function(event) {
                setTimeout(showSelectionDropdown, 10); // Add a short delay before showing the dropdown
            }, false);
        
        
            // Hide dropdown when clicking outside
            document.addEventListener('click', function(event) {
                const dropdown = document.getElementById('text-selection-dropdown');
                if (!dropdown.contains(event.target) && !event.target.matches('.dropbtn')) {
                    dropdown.style.display = 'none'; // Hide the dropdown if the click is outside
                }
            }, false);
    
        }

    function openModal() {
        projectSelectionModal.style.display = 'block';
        apiModule.loadSessionsIntoModal();
    }

    function closeModal() {
        projectSelectionModal.style.display = 'none';
    }

    function showLoadingIndicator(show) {
        loadingIndicator.style.display = show ? 'block' : 'none';
    }

    // MAKE SURE WE HAVE THIS FUNCTIONALITY
    //<input type="text" id="input-box" placeholder="Type your message here..." onkeypress="sendMessage(event)"></input>
    //<button onclick="createNewChat()">+ New Project</button>

    function appendMessage(message, sender) {
        AppState.customLog("Appending message...");
        const codeBlockRegex = /```[\s\S]*?```/g;
        let match;
        let lastIndex = 0;
        let formattedMessage = '';

        // Insert a line break if the last message wasn't from the same sender
        let lastChild = chatBox.lastElementChild;
        if (lastChild && !lastChild.classList.contains(sender + '-message')) {
            formattedMessage += '<br>'; // Add a single line break
        }

        while ((match = codeBlockRegex.exec(message)) !== null) {
            let textBeforeCode = message.substring(lastIndex, match.index);
            if (textBeforeCode) {
                formattedMessage += applyTextStyles(textBeforeCode);
            }
            let codeContent = match[0].slice(3, -3);
            formattedMessage += `<pre class="code-block"><code>${codeContent}</code></pre>`;
            lastIndex = match.index + match[0].length;
        }

        if (lastIndex < message.length) {
            let remainingText = message.substring(lastIndex);
            formattedMessage += applyTextStyles(remainingText);
        }

        // Append the formatted message to the chat box
        let messageDiv = document.createElement('div');
        messageDiv.className = 'message ' + sender + '-message'; // Ensure to use sender-message for styling
        messageDiv.innerHTML = formattedMessage;
        chatBox.appendChild(messageDiv);
        AppState.customLog(`Appending message: ${message} from ${sender}`);
        // Scroll to the bottom of the chat box
        chatBox.scrollTop = chatBox.scrollHeight;
    }

    function createSessionItem(session) {
        AppState.customLog(`Creating session item with session_id: ${session.session_id} and name: ${session.name}`);
        AppState.customLog("Creating session item...");
        const sessionItem = document.createElement('div');
        sessionItem.className = 'project-list-item';
        sessionItem.dataset.sessionId = session.session_id;
        AppState.customLog("Session name: " + session.name);

        const sessionName = document.createElement('span');
        sessionName.className = 'session-name';
        sessionName.textContent = session.name || 'Session ' + session.session_id;
        AppState.customLog("Session name text content: " + sessionName.textContent);

        // Dropdown and its contents
        const dropdown = document.createElement('div');
        dropdown.className = 'dropdown';
        const dropbtn = document.createElement('button');
        dropbtn.className = 'dropbtn';
        dropbtn.textContent = '...';
        dropbtn.onclick = function(event) {
            event.stopPropagation();
            this.nextElementSibling.classList.toggle('show');
        };
        dropdown.appendChild(dropbtn);

        const dropdownContent = document.createElement('div');
        dropdownContent.className = 'dropdown-content';
        dropdownContent.appendChild(createDropdownOption('Rename', () => promptForNewSessionName(session.session_id, sessionName)));

        dropdownContent.appendChild(createDropdownOption('Delete', () => apiModule.deleteChat(session.session_id)));
        dropdown.appendChild(dropdownContent);

        sessionItem.appendChild(sessionName);
        sessionItem.appendChild(dropdown);

        return sessionItem;
    }

    function promptForNewSessionName(sessionId, sessionNameElement) {
        const currentName = sessionNameElement.textContent;
        const newName = prompt('Enter new session name:', currentName);
        if (newName && newName !== currentName) {
            renameSession(sessionId, newName);
        }
    }

    function createDropdownOption(text, action) {
        const option = document.createElement('a');
        option.textContent = text;
        option.href = '#';
        option.onclick = function(event) {
            event.preventDefault();
            action();
        };
        return option;
    }

    // Toggles a dropdown menu
    function toggleDropdown(buttonElement) {
        closeAllDropdowns();
        buttonElement.nextElementSibling.classList.toggle('show');
    }

    function closeAllDropdowns() {
        var dropdowns = document.getElementsByClassName("dropdown-content");
        for (var i = 0; i < dropdowns.length; i++) {
            var openDropdown = dropdowns[i];
            if (openDropdown.classList.contains('show')) {
                openDropdown.classList.remove('show');
            }
        }
    }

    function closeAllDropdownsOutsideClick(event) {
        if (!event.target.matches('.dropbtn')) {
            closeAllDropdowns();
        }
    }

    function showSelectionDropdown() {
        const selection = window.getSelection();
        const selectionText = selection.toString().trim();

        const dropdown = document.getElementById('text-selection-dropdown');

        if (selectionText.length > 0) {
            const rangeRect = selection.getRangeAt(0).getBoundingClientRect();
            dropdown.style.top = `${rangeRect.bottom + window.scrollY}px`;
            dropdown.style.left = `${rangeRect.left + window.scrollX}px`;
            dropdown.style.display = 'block';
        } else {
            dropdown.style.display = 'none'; // Hide dropdown if no text is selected
        }
    }
    
    function updateActiveSessionDisplay() {
        // Remove the existing active session display if it exists
        const existingActiveSessionDisplay = document.getElementById('active-session-display');
        if (existingActiveSessionDisplay) {
            existingActiveSessionDisplay.remove();
        }
    
        // Only proceed if there is an active session ID
        const currentSessionId = AppState.getCurrentSessionId(); // Use getter for consistency
        if (currentSessionId) {
            const sessionName = AppState.getSessionName(currentSessionId); // Use getter for session name
            if(sessionName) { // Ensure sessionName is not undefined
                const sessionDisplay = createSessionItem({
                    session_id: currentSessionId,
                    name: sessionName
                });
                sessionDisplay.id = 'active-session-display';
    
                // Find the 'Directory' heading and insert the session button after it
                const directoryHeading = document.querySelector('#sidebar h2'); // Assuming 'Directory' is an <h2> inside #sidebar
                if (directoryHeading) {
                    directoryHeading.after(sessionDisplay); // Insert after the 'Directory' heading
                } else {
                    // If the 'Directory' heading cannot be found, append to the end of sidebar
                    sidebar.appendChild(sessionDisplay);
                }
            }
        }
    }
    

    // Handles modal click to close it on outside click
    function handleModalClick(event) {
        AppState.customLog("Handling modal click...");
        if (event.target.classList.contains('modal')) {
            closeModal();
        }
    }

    function handleProjectListClick(event) {
        event.stopPropagation();
        const sessionItem = event.target.closest('.project-list-item');
        AppState.customLog("Clicked on Session item: " + sessionItem);
        if (!sessionItem) return;

        const sessionId = sessionItem.dataset.sessionId;
        if (event.target.classList.contains('rename-option')) {
            const sessionNameElement = sessionItem.querySelector('span.session-name'); // Ensure this selector targets the correct span
            AppState.customLog("Session name element: " + sessionNameElement);
            const currentName = sessionNameElement.textContent;
            const newName = prompt('Enter new session name:', currentName);
            if (newName && newName !== currentName) {
                renameSession(sessionId, newName);
            }
        } else if (event.target.classList.contains('delete-option')) {
            if (confirm('Are you sure you want to delete this session?')) {
                apiModule.deleteChat(sessionId);
            }
        }
        closeAllDropdowns();
    }

    function copySelectedText() {
        const selection = window.getSelection();
        const selectionText = selection.toString();
        navigator.clipboard.writeText(selectionText).then(() => {
            alert('Copied to clipboard');
        }).catch(err => {
            console.error('Failed to copy text', err);
        });
    }

    function applyTextStyles(text) {
        // Headers
        text = text.replace(/(###\s(.*))/g, '<h3>$2</h3>');
        text = text.replace(/(##\s(.*))/g, '<h2>$2</h2>');
        text = text.replace(/(#\s(.*))/g, '<h1>$2</h1>');

        // Lists
        text = text.replace(/(\n- (.*))/g, '<li>$2</li>').replace(/<li>(.*?)<\/li>/g, '<ul>$&</ul>');
        text = text.replace(/(\n\d+\.\s(.*))/g, '<li>$2</li>').replace(/<li>(.*?)<\/li>/g, '<ol>$&</ol>');
        
        // Links
        text = text.replace(/(http[s]?:\/\/[^\s]+)/g, '<a href="$1">$1</a>');

        // Line Breaks
        text = text.replace(/\n\n/g, '<p></p>');
        text = text.replace(/\n/g, '<br>');

        // Bold and Italics
        text = text.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
        text = text.replace(/\*(.*?)\*/g, '<em>$1</em>');
        // text = text.replace(/```[\s\S]*?```/g, '<code-block>$1</code-block>');
        
        return text;
    }

    function populateSessionsInModal(sessions) {
        console.log('Sessions fetched:', JSON.stringify(sessions, null, 2));
        console.log('Current Session ID:', AppState.getCurrentSessionId());
        projectList.innerHTML = ''; // Clear the list first
        sessions.forEach(session => {
            const sessionItem = createSessionItem(session);
            projectList.appendChild(sessionItem);
            console.log("session " + session.session_id + " added to project list.") // string
            // Add active session class if this session is the current session
            if (session.session_id === parseInt(AppState.getCurrentSessionId(), 10)) {
                sessionItem.classList.add('active-session');
                console.log("session " + session.session_id + " is active.");
            }
        });
        console.log("current session: " + AppState.getCurrentSessionId()) // number
        AppState.storeSessionNames(sessions);
    }

    function displaySessions(sessions) {
        // Assuming you have some container in your HTML to display sessions
        const sessionsContainer = document.getElementById('sessions-container');
        sessionsContainer.innerHTML = ''; // Clear the container

        sessions.forEach(session => {
            // You might create a new DOM element for each session here
            const sessionElement = document.createElement('div');
            sessionElement.textContent = session.name; // Just an example, use actual session properties
            sessionsContainer.appendChild(sessionElement);
        });
    }
    
    function initializeApp() {
        setupEventListeners();
    }

    // Publicly exposed methods
    return {
        setupEventListeners: setupEventListeners,
        openModal: openModal,
        closeModal: closeModal,
        showLoadingIndicator: showLoadingIndicator,
        appendMessage: appendMessage,
        createSessionItem: createSessionItem,
        toggleDropdown: toggleDropdown,
        closeAllDropdownsOutsideClick: closeAllDropdownsOutsideClick,
        showSelectionDropdown: showSelectionDropdown,
        updateActiveSessionDisplay: updateActiveSessionDisplay,
        handleModalClick: handleModalClick,
        handleProjectListClick: handleProjectListClick,
        copySelectedText: copySelectedText,
        applyTextStyles: applyTextStyles,
        initializeApp: initializeApp,
        populateSessionsInModal: populateSessionsInModal,
        displaySessions: displaySessions,
        chatBox: chatBox,
        inputBox: inputBox,
    };
})();

// Export the UI Module
export default uiModule;
