import os
import asyncio
from app.core.config import settings
from app.services.settings_service import settings_service

class OutputService:
    def __init__(self):
        self.settings_service = settings_service

    async def output_async(self, track):
        try:
            out_dir = os.path.join(os.getcwd(), 'Outputs')
            os.makedirs(out_dir, exist_ok=True)
            # Write title and author
            try:
                with open(os.path.join(out_dir, 'title.txt'), 'w', encoding='utf-8') as f:
                    f.write(track.title or '')
            except Exception:
                pass
            try:
                with open(os.path.join(out_dir, 'author.txt'), 'w', encoding='utf-8') as f:
                    f.write(track.author or '')
            except Exception:
                pass
            # Write custom
            try:
                settings_out = self.settings_service.get_settings().output
                template = settings_out.template
                content = template.replace('{title}', track.title or '').replace('{author}', track.author or '').replace('{album}', track.album or '')
                with open(os.path.join(out_dir, 'custom.txt'), 'w', encoding='utf-8') as f:
                    f.write(content)
            except Exception:
                pass
        except Exception:
            pass

output_service = OutputService()
