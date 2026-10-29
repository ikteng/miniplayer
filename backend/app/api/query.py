from fastapi import APIRouter
from app.models.entities import Query, QueryProgress, Player, Track, RespData, CoverConvertDTO, CoverVideoDTO
from app.services.nowplaying_service import nowplaying_service
from app.services.cover_service import cover_service

router = APIRouter()

@router.get('/query')
@router.get('/api/query')
async def query():
    q = nowplaying_service.query()
    return q.model_dump(by_alias=True)

@router.get('/query/progress')
@router.get('/api/query/progress')
async def query_progress():
    qp = nowplaying_service.query_progress()
    return qp.model_dump(by_alias=True)

@router.get('/query/player')
@router.get('/api/query/player')
async def query_player():
    p = nowplaying_service.query_player()
    return p.model_dump(by_alias=True)

@router.get('/query/track')
@router.get('/api/query/track')
async def query_track():
    t = nowplaying_service.query_track()
    return t.model_dump(by_alias=True)

@router.post('/cover/convert')
@router.post('/api/cover/convert')
async def cover_convert(dto: CoverConvertDTO):
    from app.models.entities import Base64Img
    result = await cover_service.convert_to_base64(dto.cover_url if hasattr(dto, 'cover_url') else getattr(dto, 'coverUrl', ''))
    return result.model_dump(by_alias=True)

@router.post('/cover/videoUrl')
@router.post('/api/cover/videoUrl')
async def cover_video(dto: CoverVideoDTO):
    result = await cover_service.get_video_url(dto.songTitle, dto.songAuthor)
    return result

@router.get('/api/query/hasSong')
async def has_song():
    return nowplaying_service.has_song().model_dump()

@router.get('/api/query/isConnected')
async def is_connected():
    return nowplaying_service.is_connected().model_dump()
