def get_formatted_duration(seconds: int) -> str:
    if seconds < 0:
        seconds = 0
    m, s = divmod(seconds, 60)
    return f'{m}:{s:02d}'
