from fastapi import APIRouter
from app.models.entities import RespData

router = APIRouter()

@router.get('/api/server/info')
async def server_info():
    return RespData(code=200, msg='success', data={'version': '1.0.0'}).model_dump()
