from pydantic import BaseModel
from typing import Optional, Any

class Player(BaseModel):
    has_song: bool = False
    is_paused: bool = True
    seekbar_current_position: int = 0
    seekbar_current_position_human: str = '0:00'
    state_percent: float = 0.0

class Track(BaseModel):
    title: str = ''
    author: str = ''
    album: str = ''
    cover: str = ''
    duration: int = 0
    duration_human: str = '0:00'
    url: str = ''

class Query(BaseModel):
    player: Player
    track: Track

class QueryProgress(BaseModel):
    progress_ms: int = 0

class RespData(BaseModel):
    data: Any

class CoverConvertDTO(BaseModel):
    cover_url: str

class CoverVideoDTO(BaseModel):
    songTitle: str = ''
    songAuthor: str = ''
