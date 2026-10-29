import json
import os
import asyncio
from typing import Dict, Any
from app.core.config import settings
from app.models.entities import Lyric, SettingsLyricCommon
from app.core.events import event_bus

class LyricService:
    def __init__(self):
        self.cache_dir = os.path.join(settings.settings_dir, 'lyric_cache')
        os.makedirs(self.cache_dir, exist_ok=True)
        self.current_lyric = Lyric()
        self.current_lyric_window_title = None
        self.fetch_lyric_enabled = True
        self.settings_common = SettingsLyricCommon()

    def get_lyric(self):
        return self.current_lyric

    async def get_lyric_by_title_author(self, title: str = '', author: str = ''):
        key = f"{title}_{author}".replace('/', '_').replace('\\', '_')
        cache_file = os.path.join(self.cache_dir, f"{key}.json")
        if os.path.exists(cache_file):
            try:
                with open(cache_file, 'r', encoding='utf-8') as f:
                    data = json.load(f)
                lyric = Lyric(**data)
                self.current_lyric = lyric
                return lyric
            except Exception:
                pass
        lyric = Lyric(title=title, author=author, lyric='', tlyric='')
        try:
            with open(cache_file, 'w', encoding='utf-8') as f:
                json.dump(lyric.model_dump(), f, ensure_ascii=False)
        except Exception:
            pass
        self.current_lyric = lyric
        return lyric

    def force_refresh_lyric(self):
        try:
            event_bus.publish('lyric_changed', self.current_lyric)
        except Exception:
            pass

    def get_settings(self, id: str = None):
        from app.models.entities import SettingsLyric
        return SettingsLyric()

    def get_common_settings(self):
        return self.settings_common

    def update_settings(self, id: str, settings):
        pass

    def update_common_settings(self, settings):
        if isinstance(settings, dict):
            for k, v in settings.items():
                if hasattr(self.settings_common, k):
                    setattr(self.settings_common, k, v)
        elif hasattr(settings, 'model_dump'):
            data = settings.model_dump()
            for k, v in data.items():
                if hasattr(self.settings_common, k):
                    setattr(self.settings_common, k, v)
        try:
            event_bus.publish('lyric_changed', self.current_lyric)
        except Exception:
            pass

lyric_service = LyricService()
