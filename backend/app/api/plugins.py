from fastapi import APIRouter, Request
from app.models.entities import RespData

router = APIRouter()

@router.get('/api/plugins')
async def get_plugins():
    return RespData(code=200, msg='success', data=[]).model_dump()

@router.post('/api/plugins')
async def post_plugins(request: Request):
    return RespData(code=200, msg='success', data=None).model_dump()
