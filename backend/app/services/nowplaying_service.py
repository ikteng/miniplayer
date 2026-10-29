import time
from loguru import logger
from app.models.entities import Player, Track, Query, QueryProgress, RespData
from app.utils.timeutil import get_formatted_duration
from app.utils.songutil import parse_window_title
from app.core.events import event_bus
from app.services.output_service import output_service

class NowPlayingService:
    def __init__(self):
        self.player = Player()
        self.track = Track()
        self.timer_start = time.time()
        self.timer_paused_at = None
        self.prev_window_title = ''
        self.prev_status = 'None'
        self.prev_is_paused = True
        self.none_occurs_time = 0
        self.CSHARP_PROGRESS_NONE = -1
        self.CSHARP_PROGRESS_RESET = -2
        self.prev_csharp_progress_seconds = self.CSHARP_PROGRESS_RESET

    def reset_timer(self):
        self.timer_start = time.time()
        self.timer_paused_at = None

    def start_timer(self):
        if self.timer_paused_at is not None:
            paused_for = time.time() - self.timer_paused_at
            self.timer_start += paused_for
            self.timer_paused_at = None

    def pause_timer(self):
        if self.timer_paused_at is None:
            self.timer_paused_at = time.time()

    def get_time_ms(self):
        if self.timer_paused_at is not None:
            return int((self.timer_paused_at - self.timer_start) * 1000)
        return int((time.time() - self.timer_start) * 1000)

    def query(self) -> Query:
        return Query(player=self.query_player(), track=self.track)

    def query_player(self) -> Player:
        progress_ms = self.get_time_ms()
        progress_sec = int(progress_ms / 1000)
        duration = self.track.duration
        if duration <= 0:
            duration = 300
        self.player.seekbar_current_position = progress_sec
        self.player.seekbar_current_position_human = get_formatted_duration(progress_sec)
        self.player.state_percent = float(progress_sec) / float(duration) if duration > 0 else 0.0
        self.player.has_song = self.player.has_song
        self.player.is_paused = self.player.is_paused
        return self.player

    def query_track(self) -> Track:
        return self.track

    def query_progress(self) -> QueryProgress:
        return QueryProgress(progress_ms=self.get_time_ms())

    def has_song(self) -> RespData:
        return RespData(code=200, msg='success', data=self.player.has_song)

    def is_connected(self) -> RespData:
        return RespData(code=200, msg='success', data=self.player.has_song)

nowplaying_service = NowPlayingService()
