from fastapi import APIRouter, Request
from app.models.entities import RespData

router = APIRouter()

@router.get('/api/profile')
async def get_profile():
    return RespData(code=200, msg='success', data=[]).model_dump()

@router.post('/api/profile')
async def post_profile(request: Request):
    return RespData(code=200, msg='success', data=None).model_dump()
