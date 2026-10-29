from typing import Callable, List

class EventBus:
    def __init__(self):
        self._listeners = {}

    def subscribe(self, event_type: str, listener: Callable):
        if event_type not in self._listeners:
            self._listeners[event_type] = []
        self._listeners[event_type].append(listener)

    def publish(self, event_type: str, *args, **kwargs):
        for listener in self._listeners.get(event_type, []):
            try:
                listener(*args, **kwargs)
            except Exception:
                pass

event_bus = EventBus()
