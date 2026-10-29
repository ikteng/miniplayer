from fastapi import APIRouter

router = APIRouter()

@router.get('/public/{path:path}')
async def public_resource(path: str):
    return {}
