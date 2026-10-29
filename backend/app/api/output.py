from fastapi import APIRouter, Request
from app.models.entities import RespData

router = APIRouter()

@router.get('/api/output')
async def get_output():
    return RespData(code=200, msg='success', data={}).model_dump()

@router.post('/api/output')
async def post_output(request: Request):
    try:
        data = await request.json()
        from app.services.settings_service import settings_service
        settings_service.get_settings().output = settings_service.get_settings().output.__class__(**data) if hasattr(data, 'items') else settings_service.get_settings().output
        return RespData(code=200, msg='success', data=None).model_dump()
    except Exception as e:
        return RespData(code=500, msg=str(e), data=None).model_dump()
