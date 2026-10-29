from fastapi import APIRouter, Query, Request
from app.models.entities import RespData, Lyric
from app.services.lyric_service import lyric_service

router = APIRouter()

@router.get('/api/lyric')
async def get_lyric(title: str = Query(''), author: str = Query('')):
    if title or author:
        result = await lyric_service.get_lyric_by_title_author(title, author)
        return result.model_dump(by_alias=True)
    result = lyric_service.get_lyric()
    return result.model_dump(by_alias=True)

@router.get('/lyric')
async def get_lyric_alt(title: str = Query(''), author: str = Query('')):
    if title or author:
        result = await lyric_service.get_lyric_by_title_author(title, author)
        return result.model_dump(by_alias=True)
    result = lyric_service.get_lyric()
    return result.model_dump(by_alias=True)

@router.get('/api/lyric/settings')
async def get_lyric_settings(id: str = Query('main')):
    result = lyric_service.get_settings(id)
    return result.model_dump(by_alias=True)

@router.get('/api/lyric/settings/common')
async def get_lyric_settings_common():
    result = lyric_service.get_common_settings()
    return result.model_dump(by_alias=True)

@router.put('/api/lyric/settings')
async def put_lyric_settings(request: Request, id: str = Query('main')):
    try:
        data = await request.json()
        lyric_service.update_settings(id, data)
        return RespData(code=200, msg='success', data=None).model_dump()
    except Exception as e:
        return RespData(code=500, msg=str(e), data=None).model_dump()

@router.put('/api/lyric/settings/common')
async def put_lyric_settings_common(request: Request):
    try:
        data = await request.json()
        lyric_service.update_common_settings(data)
        return RespData(code=200, msg='success', data=None).model_dump()
    except Exception as e:
        return RespData(code=500, msg=str(e), data=None).model_dump()
