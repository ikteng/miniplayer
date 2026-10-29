from fastapi import APIRouter
from app.models.entities import RespData

router = APIRouter()

@router.get('/api/system/info')
@router.get('/system/info')
async def system_info():
    return RespData(code=200, msg='success', data={'version': '1.0.0', 'platform': 'windows', 'name': 'now-playing-python'}).model_dump()

@router.get('/api/server/info')
@router.get('/server/info')
async def server_info():
    return RespData(code=200, msg='success', data={'version': '1.0.0'}).model_dump()

@router.get('/api')
async def api_root():
    return {'message': 'NowPlaying API'}
