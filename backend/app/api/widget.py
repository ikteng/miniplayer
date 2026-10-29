import os
from fastapi import APIRouter, Request
from fastapi.responses import HTMLResponse, FileResponse
from app.core.config import settings as app_settings
from app.services.settings_service import settings_service

router = APIRouter()

templates_dir = os.path.join(os.path.dirname(os.path.dirname(os.path.dirname(__file__))), 'resources', 'templates')

def _read_template(name: str) -> str:
    with open(os.path.join(templates_dir, name), 'r', encoding='utf-8') as f:
        return f.read()

@router.get('/widget', response_class=HTMLResponse)
@router.get('/widget/{profile_id}', response_class=HTMLResponse)
async def widget_page(profile_id: str = ''):
    widget_path = os.path.join(app_settings.static_dir, 'widget.html')
    if os.path.exists(widget_path):
        return FileResponse(widget_path)
    return _read_template('widget.html')

@router.get('/settings-widget', response_class=HTMLResponse)
async def settings_widget_page():
    return _read_template('settings-widget.html')

@router.get('/widgets/amuse/settings')
async def widget_settings():
    return {'settings': settings_service.get_settings_widget().model_dump()}

@router.post('/widgets/amuse/settings')
async def update_widget_settings(request: Request):
    data = await request.json()
    settings_service.update_settings_widget(data)
    return {'settings': settings_service.get_settings_widget().model_dump()}
