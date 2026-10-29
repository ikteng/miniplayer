import asyncio
import subprocess
import os
from loguru import logger
from typing import Callable, Optional

class AudioDetector:
    def __init__(self, assets_dir: str):
        self.assets_dir = assets_dir
        self.proc = None
        self.status = 'None'
        self.window_title = ''
        self.progress_seconds = -1
        self.total_seconds = -1
        self._callback: Optional[Callable] = None

    async def start(self, callback: Callable):
        self._callback = callback
        exe_path = os.path.join(self.assets_dir, 'AudioService', 'GetMusicStatus.exe')
        if not os.path.exists(exe_path):
            logger.error(f'GetMusicStatus.exe not found at {exe_path}')
            return
        # Run from assets dir like Java does
        try:
            self.proc = await asyncio.create_subprocess_exec(
                exe_path,
                '--device-id', 'default',
                '--platform', 'default',
                '--smtc', 'true',
                '--poll-interval', '1000',
                '--splayer-next-port', '37685',
                cwd=os.path.join(self.assets_dir, 'AudioService'),
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE
            )
            asyncio.create_task(self._read_loop())
            logger.info('Started GetMusicStatus.exe')
        except Exception as e:
            logger.error(f'Failed to start detector: {e}')

    async def _read_loop(self):
        while self.proc and self.proc.stdout:
            try:
                line = await self.proc.stdout.readline()
                if not line:
                    break
                text = line.decode('utf-8', errors='ignore').rstrip()
                if not text:
                    continue
                if text in ['Playing', 'Paused', 'None']:
                    self.status = text
                    if text == 'None':
                        self.window_title = ''
                        self.progress_seconds = -1
                        self.total_seconds = -1
                elif text.startswith('Progress:'):
                    try:
                        parts = text[9:].split('|')
                        if len(parts) == 2:
                            self.progress_seconds = int(parts[0])
                            self.total_seconds = int(parts[1])
                    except Exception:
                        pass
                else:
                    new_title = text
                    if new_title != self.window_title:
                        self.progress_seconds = -1
                        self.total_seconds = -1
                    self.window_title = new_title
                if self._callback:
                    await self._callback()
            except Exception as e:
                logger.error(f'Error reading detector: {e}')
                break

    async def stop(self):
        if self.proc:
            try:
                self.proc.terminate()
                await asyncio.wait_for(self.proc.wait(), timeout=2)
            except Exception:
                try:
                    self.proc.kill()
                except:
                    pass
            self.proc = None
