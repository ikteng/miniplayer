import base64
import aiohttp
import os
from app.core.config import settings as app_settings

class CoverService:
    def __init__(self):
        self.should_read_local_image = False
        self.prev_song_title = None
        self.prev_song_author = None
        self.cached_video_url = None
        self.video_url_cached = False

    async def convert_to_base64(self, cover_url: str):
        from app.models.entities import Base64Img
        if not cover_url:
            # Try to read default
            try:
                default_path = os.path.join(app_settings.assets_dir, 'no_cover_base64.txt')
                if os.path.exists(default_path):
                    with open(default_path, 'r', encoding='utf-8') as f:
                        data = f.read().strip()
                    return Base64Img(base64Img=data)
            except Exception:
                pass
            return Base64Img(base64Img='')
        try:
            async with aiohttp.ClientSession() as session:
                async with session.get(cover_url, timeout=10) as resp:
                    if resp.status == 200:
                        data = await resp.read()
                        content_type = resp.headers.get('Content-Type', 'image/jpeg')
                        b64 = base64.b64encode(data).decode('utf-8')
                        return Base64Img(base64Img=f'data:{content_type};base64,{b64}')
        except Exception:
            pass
        try:
            default_path = os.path.join(app_settings.assets_dir, 'no_cover_base64.txt')
            if os.path.exists(default_path):
                with open(default_path, 'r', encoding='utf-8') as f:
                    data = f.read().strip()
                return Base64Img(base64Img=data)
        except Exception:
            pass
        return Base64Img(base64Img='')

    async def get_video_url(self, song_title: str, song_author: str):
        if self.video_url_cached and self.prev_song_title == song_title and self.prev_song_author == song_author:
            return self.cached_video_url
        self.prev_song_title = song_title
        self.prev_song_author = song_author
        self.cached_video_url = None
        self.video_url_cached = True
        return None

cover_service = CoverService()
