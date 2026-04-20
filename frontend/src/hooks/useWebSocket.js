import { useEffect, useRef, useCallback } from 'react';

export function useWebSocket(roomId, onMessage) {
  const wsRef = useRef(null);
  const onMessageRef = useRef(onMessage);
  onMessageRef.current = onMessage;

  const connect = useCallback(() => {
    if (!roomId) return;
    const token = localStorage.getItem('token');
    const wsUrl = `ws://localhost:8080/ws?room_id=${roomId}&token=${token}`;
    const ws = new WebSocket(wsUrl);

    ws.onopen = () => console.log(`[ws] connected to room ${roomId}`);

    ws.onmessage = (e) => {
      try {
        const data = JSON.parse(e.data);
        onMessageRef.current(data);
      } catch {}
    };

    ws.onclose = () => {
      console.log('[ws] disconnected');
      // Reconnect after 2s
      setTimeout(connect, 2000);
    };

    ws.onerror = (err) => console.error('[ws] error', err);

    wsRef.current = ws;
  }, [roomId]);

  useEffect(() => {
    connect();
    return () => {
      if (wsRef.current) {
        wsRef.current.onclose = null; // prevent reconnect on unmount
        wsRef.current.close();
      }
    };
  }, [connect]);

  const send = useCallback((content) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ content }));
    }
  }, []);

  return { send };
}
