from fastapi import APIRouter, Request
from app.models.entities import SettingsWrapper, RespData, SettingsGeneral
from app.services.settings_service import settings_service

router = APIRouter()

@router.get('/api/settings')
async def get_settings():
    return settings_service.get_settings().model_dump(by_alias=True)

@router.post('/api/settings')
async def post_settings(request: Request):
    try:
        data = await request.json()
        settings_service.set_settings(data)
        return RespData(code=200, msg='success', data=None).model_dump()
    except Exception as e:
        return RespData(code=500, msg=str(e), data=None).model_dump()

@router.get('/api/settings/general')
async def get_settings_general():
    return settings_service.get_settings_general().model_dump(by_alias=True)

@router.post('/api/settings/general')
@router.put('/api/settings/general')
async def post_settings_general(request: Request):
    try:
        data = await request.json()
        settings_service.update_settings_general(data)
        return RespData(code=200, msg='success', data=None).model_dump()
    except Exception as e:
        return RespData(code=500, msg=str(e), data=None).model_dump()
