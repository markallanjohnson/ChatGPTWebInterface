// app.js
import uiModule from './uiModule.js';
import AppState from './appStateModule.js';

document.addEventListener('DOMContentLoaded', () => {
    AppState.customLog('DOM fully loaded and parsed');
    uiModule.setupEventListeners();
    AppState.customLog('Event listeners setup')
    AppState.customLog('Session ID: ' + AppState.getCurrentSessionId());
    //apiModule.loadSessions();
});

