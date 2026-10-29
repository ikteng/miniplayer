from fastapi import APIRouter, Request
from app.models.entities import RespData
from app.detectors.platform_detector import detect_platform

router = APIRouter()

@router.get('/api/audio/devices')
@router.get('/audio/devices')
async def get_audio_devices():
    return RespData(code=200, msg='success', data=[]).model_dump()

@router.get('/api/audio/volume')
@router.get('/audio/volume')
async def get_volume():
    return RespData(code=200, msg='success', data={'volume': 0}).model_dump()

@router.post('/api/audio/device')
async def set_device(request: Request):
    return RespData(code=200, msg='success', data=None).model_dump()

@router.get('/api/audio/detectPlatform')
async def detect_platform_endpoint():
    result = await detect_platform()
    return RespData(code=200, msg='success', data=result).model_dump()
