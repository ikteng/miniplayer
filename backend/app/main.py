import os
from fastapi import FastAPI
from fastapi.staticfiles import StaticFiles
from fastapi.middleware.cors import CORSMiddleware
from loguru import logger

from app.api import query, settings, lyric, audio, system, plugins, output, profile, page, public, server, widget
from app.ws import lyric_ws, websocket_lyric
from app.core.config import settings as app_settings
from app.detectors.audio_detector import AudioDetector
from app.services.nowplaying_service import nowplaying_service
from app.services.output_service import output_service

app = FastAPI(title='Miniplay', version='1.0.0')

app.add_middleware(
    CORSMiddleware,
    allow_origins=['*'],
    allow_credentials=True,
    allow_methods=['*'],
    allow_headers=['*'],
)

app.include_router(query.router)
app.include_router(settings.router)
app.include_router(lyric.router)
app.include_router(audio.router)
app.include_router(system.router)
app.include_router(plugins.router)
app.include_router(output.router)
app.include_router(profile.router)
app.include_router(lyric_ws.router)
app.include_router(websocket_lyric.router)
app.include_router(page.router)
app.include_router(public.router)
app.include_router(server.router)
app.include_router(widget.router)

static_dir = app_settings.static_dir
if os.path.exists(static_dir):
    app.mount('/static', StaticFiles(directory=static_dir), name='static_files')
    app.mount('/', StaticFiles(directory=static_dir, html=True), name='static')

detector = AudioDetector(app_settings.assets_dir)
query.nowplaying = nowplaying_service

async def on_update():
    status = detector.status
    window_title = detector.window_title
    
    import time
    now = time.time()
    
    if status == 'None':
        if nowplaying_service.prev_status == 'None':
            if now - nowplaying_service.none_occurs_time > 3:
                had_song = nowplaying_service.player.has_song
                from app.models.entities import Player, Track
                nowplaying_service.player = Player()
                nowplaying_service.track = Track()
                nowplaying_service.reset_timer()
                nowplaying_service.prev_window_title = ''
                nowplaying_service.prev_is_paused = True
                nowplaying_service.prev_csharp_progress_seconds = nowplaying_service.CSHARP_PROGRESS_RESET
        else:
            nowplaying_service.none_occurs_time = now
        nowplaying_service.prev_status = status
        return
    
    nowplaying_service.player.has_song = True
    is_paused = status != 'Playing'
    nowplaying_service.player.is_paused = is_paused
    
    if nowplaying_service.prev_window_title != window_title and window_title:
        nowplaying_service.reset_timer()
        nowplaying_service.start_timer()
        nowplaying_service.prev_csharp_progress_seconds = nowplaying_service.CSHARP_PROGRESS_RESET
        from app.utils.songutil import parse_window_title
        parts = parse_window_title(window_title)
        if parts[0]:
            nowplaying_service.track.title = parts[0]
        elif window_title:
            nowplaying_service.track.title = window_title
        if len(parts) > 1 and parts[1]:
            nowplaying_service.track.author = parts[1]
        if nowplaying_service.track.duration <= 0:
            nowplaying_service.track.duration = 300
            nowplaying_service.track.duration_human = '5:00'
        try:
            import asyncio
            asyncio.create_task(output_service.output_async(nowplaying_service.track))
        except Exception:
            pass
    else:
        if is_paused:
            nowplaying_service.pause_timer()
        else:
            nowplaying_service.start_timer()
    
    nowplaying_service.prev_window_title = window_title
    nowplaying_service.prev_status = status
    nowplaying_service.prev_is_paused = is_paused

@app.on_event('startup')
async def startup_event():
    logger.info('Starting Miniplay Backend on port 9863')
    await detector.start(on_update)

@app.on_event('shutdown')
async def shutdown_event():
    await detector.stop()
