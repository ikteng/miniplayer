import os
from pydantic_settings import BaseSettings, SettingsConfigDict

class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_prefix='MINIPLAY_',
        protected_namespaces=()
    )
    
    app_name: str = 'Miniplay'
    server_port: int = 9863
    assets_dir: str = os.path.join(os.path.dirname(os.path.dirname(os.path.dirname(__file__))), 'assets')
    settings_dir: str = os.path.join(os.path.dirname(os.path.dirname(os.path.dirname(__file__))), 'settings')
    static_dir: str = os.path.join(os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(__file__)))), 'frontend', 'dist')

settings = Settings()
