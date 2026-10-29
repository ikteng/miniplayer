from pydantic import BaseModel, Field
from typing import Optional, Any, Dict, List

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
    code: int = 200
    msg: str = ''
    data: Any = None

class Base64Img(BaseModel):
    base64_img: str = Field(default='', alias='base64Img')
    base64: str = ''

class Lyric(BaseModel):
    title: str = ''
    author: str = ''
    duration: int = 0
    lyric: str = ''
    tlyric: str = ''
    yrc_lyric: str = ''
    yrc_tlyric: str = ''
    qrc_lyric: str = ''
    qrc_tlyric: str = ''
    source: str = 'netease'
    has_lyric: bool = False
    has_translated_lyric: bool = False
    has_karaoke_lyric: bool = False

class SettingsGeneral(BaseModel):
    device_id: str = 'default'
    device_name: str = 'Default Audio Driver'
    platform: str = 'netease'
    auto_launch_home_page: bool = True
    run_at_startup: bool = False
    update_check_freq: int = 0
    smtc: bool = True
    fallback_platform_enabled: bool = False
    fallback_platform: str = 'netease'
    poll_interval: int = 100
    we_sing_cache_path: str = ''
    splayer_next_port: int = 14558

class SettingsLyric(BaseModel):
    use_custom_config: bool = False
    config_file_name: str = ''
    lyric_color: str = '#f5f5f5'
    lyric_font_size: int = 28
    lyric_font_weight: int = 700
    stroke_color: str = '#000000'
    stroke_size: int = 0
    translation_color: str = '#f5f5f5'
    translation_font_size: int = 28
    translation_font_weight: int = 700
    translation_stroke_color: str = '#000000'
    translation_stroke_size: int = 0
    gap: int = 0
    translation_gap: int = 0
    album_cover_lyric: bool = False
    animation: bool = False
    background_color: str = '#4d000000'
    position_type: str = 'bottom'
    type_lyric: bool = False
    karaoke_animation: bool = True
    vertical_align: bool = True
    show_interlude: bool = True
    show_translation: bool = True
    full_screen_lyric: bool = False
    enable_hover: bool = False
    offset_y: int = 0
    last_update_time: int = 0

class SettingsLyricCommon(BaseModel):
    lyric_source: str = 'netease'
    auto_select_best_lyric: bool = True
    romanization: bool = False
    translate: bool = True
    chinese_conversion: str = 'off'

class SettingsOutput(BaseModel):
    output_title: bool = True
    output_author: bool = True
    output_cover: bool = True
    output_custom: bool = True
    template: str = '{title} - {author}'
    enable_when_none: bool = False
    custom_none: str = ''

class SettingsWidget(BaseModel):
    cover: int = 0
    cover_blur: bool = False
    cover_glow: bool = False
    first_login: bool = False
    hide_delay: int = 10
    hide_equalizer: bool = False
    hide_on_pause: bool = False
    is_demo: bool = False
    player: int = 0
    player_colors: bool = True
    song_change_only: bool = False
    theme: int = 0
    tint_color: str = '#40a4d8'
    visible_duration: int = 5

class SettingsWrapper(BaseModel):
    general: SettingsGeneral = SettingsGeneral()
    lyric: SettingsLyric = SettingsLyric()
    lyric_common: SettingsLyricCommon = SettingsLyricCommon()
    output: SettingsOutput = SettingsOutput()
    widget: SettingsWidget = SettingsWidget()
    plugins: Dict[str, Any] = {}

class CoverConvertDTO(BaseModel):
    cover_url: str = Field(alias='cover_url')
    
    class Config:
        populate_by_name = True

class CoverVideoDTO(BaseModel):
    songTitle: str = ''
    songAuthor: str = ''
