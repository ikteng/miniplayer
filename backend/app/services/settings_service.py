import json
import os
from typing import Dict, Any
from app.core.config import settings
from app.models.entities import SettingsWrapper, SettingsGeneral, SettingsLyric, SettingsOutput, SettingsLyricCommon, SettingsWidget
from app.core.events import event_bus

class SettingsService:
    def __init__(self):
        self.settings_dir = settings.settings_dir
        os.makedirs(self.settings_dir, exist_ok=True)
        self.settings_file = os.path.join(self.settings_dir, 'settings.json')
        self.wrapper = self.load()

    def load(self) -> SettingsWrapper:
        if os.path.exists(self.settings_file):
            try:
                with open(self.settings_file, 'r', encoding='utf-8') as f:
                    data = json.load(f)
                wrapper = SettingsWrapper(**data)
                return wrapper
            except Exception:
                pass
        wrapper = SettingsWrapper()
        self.save_wrapper(wrapper)
        return wrapper

    def save_wrapper(self, wrapper: SettingsWrapper = None):
        if wrapper is None:
            wrapper = self.wrapper
        try:
            with open(self.settings_file, 'w', encoding='utf-8') as f:
                json.dump(wrapper.model_dump(by_alias=True), f, ensure_ascii=False, indent=2)
        except Exception:
            pass

    def save(self):
        self.save_wrapper()

    def get_settings_general(self) -> SettingsGeneral:
        return self.wrapper.general

    def get_settings(self) -> SettingsWrapper:
        return self.wrapper

    def get_settings_widget(self) -> SettingsWidget:
        return self.wrapper.widget

    def update_settings_widget(self, data: Dict[str, Any]):
        for k, v in data.items():
            if hasattr(self.wrapper.widget, k):
                setattr(self.wrapper.widget, k, v)
        self.save()

    def update_settings_general(self, general: Dict[str, Any] or SettingsGeneral):
        if isinstance(general, dict):
            for k, v in general.items():
                # Convert camelCase to snake_case if needed
                key = k
                if hasattr(self.wrapper.general, k):
                    setattr(self.wrapper.general, k, v)
                else:
                    # Try to map
                    pass
        elif hasattr(general, 'model_dump'):
            data = general.model_dump()
            for k, v in data.items():
                if hasattr(self.wrapper.general, k):
                    setattr(self.wrapper.general, k, v)
        else:
            for k, v in general.__dict__.items():
                if not k.startswith('_') and hasattr(self.wrapper.general, k):
                    setattr(self.wrapper.general, k, v)
        self.save()
        try:
            event_bus.publish('settings_general_changed', self.wrapper.general)
        except Exception:
            pass

    def set_settings_general(self, general: Dict[str, Any]):
        self.update_settings_general(general)

    def set_settings(self, wrapper: Dict[str, Any]):
        if 'general' in wrapper:
            self.update_settings_general(wrapper['general'])
        if 'lyric' in wrapper:
            for k, v in wrapper['lyric'].items():
                if hasattr(self.wrapper.lyric, k):
                    setattr(self.wrapper.lyric, k, v)
        if 'output' in wrapper:
            for k, v in wrapper['output'].items():
                if hasattr(self.wrapper.output, k):
                    setattr(self.wrapper.output, k, v)
        if 'widget' in wrapper:
            for k, v in wrapper['widget'].items():
                if hasattr(self.wrapper.widget, k):
                    setattr(self.wrapper.widget, k, v)
        if 'lyric_common' in wrapper:
            for k, v in wrapper['lyric_common'].items():
                if hasattr(self.wrapper.lyric_common, k):
                    setattr(self.wrapper.lyric_common, k, v)
        if 'plugins' in wrapper:
            if isinstance(wrapper['plugins'], dict):
                self.wrapper.plugins.update(wrapper['plugins'])
        self.save()

settings_service = SettingsService()
