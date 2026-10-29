import re

def parse_window_title(window_title: str):
    # Simple parser similar to Java version
    if not window_title:
        return [None, None]
    # Try common patterns
    # e.g. "Artist - Title"
    if ' - ' in window_title:
        parts = window_title.split(' - ', 1)
        return [parts[1].strip(), parts[0].strip()]
    return [window_title.strip(), None]
