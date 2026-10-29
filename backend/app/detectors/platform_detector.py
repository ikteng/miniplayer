from loguru import logger

# Maps a session's source_app_user_model_id (AUMID) to our platform id.
# Browser AUMIDs (MSEdge/Chrome) are ambiguous - they could be any web-based
# platform (YouTube Music, browser fallback, etc.) so they're left unmapped.
AUMID_PLATFORM_MAP = {
    'spotify.exe': 'spotify',
    'tencent.qqmusic': 'qq',
    'kugou.exe': 'kugou',
    'kuwo.exe': 'kuwo',
    'cloudmusic.exe': 'netease',
    'cider': 'cider',
    'potplayermini64.exe': 'potplayer',
    'potplayer.exe': 'potplayer',
    'foobar2000.exe': 'foobar',
    'aimp.exe': 'aimp',
}


async def detect_platform() -> dict:
    try:
        from winsdk.windows.media.control import (
            GlobalSystemMediaTransportControlsSessionManager as Manager,
        )
    except ImportError:
        return {'platform': None, 'reason': 'winsdk not available'}

    try:
        mgr = await Manager.request_async()
    except Exception as e:
        logger.error(f'Failed to request SMTC session manager: {e}')
        return {'platform': None, 'reason': 'smtc unavailable'}

    for session in mgr.get_sessions():
        aumid = (session.source_app_user_model_id or '').lower()
        platform = AUMID_PLATFORM_MAP.get(aumid)
        if platform:
            try:
                info = await session.try_get_media_properties_async()
                title = info.title
            except Exception:
                title = ''
            return {'platform': platform, 'aumid': session.source_app_user_model_id, 'title': title}

    return {'platform': None, 'reason': 'no matching session found'}
