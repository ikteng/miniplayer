# Miniplay

A music now-playing information display tool with Python backend and React frontend.

## Quick Start

### Windows
Double-click start.bat or run start.ps1 in PowerShell

### Manual Start

#### Backend
`ash
cd backend
pip install -r requirements.txt
python -m uvicorn app.main:app --host 127.0.0.1 --port 9863
`

#### Frontend (Development)
`ash
cd frontend
npm install
npm run dev
`

#### Frontend (Production - served by backend)
`ash
cd frontend
npm install
npm run build
# Copy dist contents to backend/resources/static/
`

## Features
- Real-time music detection via Windows SMTC
- REST API and WebSocket support
- Lyrics display
- Multiple music platform support
- OBS streaming output
- Customizable themes and widgets
