import os
from fastapi import APIRouter
from fastapi.responses import HTMLResponse, FileResponse
from app.core.config import settings as app_settings

router = APIRouter()

@router.get('/', response_class=HTMLResponse)
@router.get('/index', response_class=HTMLResponse)
@router.get('/page/{path:path}', response_class=HTMLResponse)
@router.get('/settings/{path:path}', response_class=HTMLResponse)
async def page(path: str = ''):
    index_path = os.path.join(app_settings.static_dir, 'index.html')
    if os.path.exists(index_path):
        return FileResponse(index_path)
    return '<html><head><title>MiniPlayer</title></head><body>MiniPlayer Python Backend</body></html>'
