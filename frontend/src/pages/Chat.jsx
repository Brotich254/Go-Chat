import { useEffect, useState, useRef, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../api';
import { useAuth } from '../context/AuthContext';
import { useWebSocket } from '../hooks/useWebSocket';
import { format } from '../utils/date';

export default function Chat() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [rooms, setRooms] = useState([]);
  const [activeRoom, setActiveRoom] = useState(null);
  const [messages, setMessages] = useState([]);
  const [onlineUsers, setOnlineUsers] = useState([]);
  const [input, setInput] = useState('');
  const [newRoom, setNewRoom] = useState('');
  const [showNewRoom, setShowNewRoom] = useState(false);
  const bottomRef = useRef(null);

  // Fetch rooms
  useEffect(() => {
    api.get('/rooms').then(({ data }) => {
      setRooms(data);
      if (data.length > 0 && !activeRoom) setActiveRoom(data[0]);
    });
  }, []);

  // Fetch message history when room changes
  useEffect(() => {
    if (!activeRoom) return;
    setMessages([]);
    setOnlineUsers([]);
    api.get(`/messages?room_id=${activeRoom.id}`).then(({ data }) => setMessages(data));
  }, [activeRoom?.id]);

  // Scroll to bottom on new messages
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleWsMessage = useCallback((data) => {
    if (data.type === 'message' && data.message) {
      setMessages((prev) => [...prev, data.message]);
    } else if (data.type === 'online_users') {
      setOnlineUsers(data.users || []);
      // Refresh room online counts
      setRooms((prev) => prev.map((r) =>
        r.id === data.room_id ? { ...r, online_count: data.users.length } : r
      ));
    }
  }, []);

  const { send } = useWebSocket(activeRoom?.id, handleWsMessage);

  const handleSend = (e) => {
    e.preventDefault();
    if (!input.trim()) return;
    send(input.trim());
    setInput('');
  };

  const handleCreateRoom = async (e) => {
    e.preventDefault();
    if (!newRoom.trim()) return;
    try {
      const { data } = await api.post('/rooms', { name: newRoom.trim() });
      setRooms((prev) => [...prev, data]);
      setActiveRoom(data);
      setNewRoom('');
      setShowNewRoom(false);
    } catch {
      alert('Room name taken');
    }
  };

  const handleLogout = () => { logout(); navigate('/login'); };

  return (
    <div className="flex h-screen bg-gray-950">
      {/* Sidebar */}
      <div className="w-64 bg-gray-900 border-r border-gray-800 flex flex-col">
        <div className="p-4 border-b border-gray-800">
          <h1 className="text-lg font-bold text-indigo-400">💬 GoChat</h1>
          <p className="text-xs text-gray-500 mt-0.5">{user?.name}</p>
        </div>

        <div className="flex-1 overflow-y-auto p-2">
          <div className="flex items-center justify-between px-2 py-1 mb-1">
            <span className="text-xs text-gray-500 uppercase tracking-wide">Rooms</span>
            <button onClick={() => setShowNewRoom(!showNewRoom)}
              className="text-gray-500 hover:text-indigo-400 text-lg leading-none">+</button>
          </div>

          {showNewRoom && (
            <form onSubmit={handleCreateRoom} className="px-2 mb-2">
              <input value={newRoom} onChange={(e) => setNewRoom(e.target.value)}
                placeholder="Room name" autoFocus
                className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-1.5 text-sm focus:outline-none focus:ring-1 focus:ring-indigo-500" />
            </form>
          )}

          {rooms.map((room) => (
            <button key={room.id} onClick={() => setActiveRoom(room)}
              className={`w-full text-left px-3 py-2 rounded-lg text-sm transition flex items-center justify-between ${
                activeRoom?.id === room.id ? 'bg-indigo-600 text-white' : 'text-gray-400 hover:bg-gray-800 hover:text-white'
              }`}>
              <span># {room.name}</span>
              {room.online_count > 0 && (
                <span className={`text-xs px-1.5 py-0.5 rounded-full ${activeRoom?.id === room.id ? 'bg-indigo-500' : 'bg-gray-700'}`}>
                  {room.online_count}
                </span>
              )}
            </button>
          ))}
        </div>

        <div className="p-3 border-t border-gray-800">
          <button onClick={handleLogout} className="text-xs text-gray-500 hover:text-red-400 transition">Logout</button>
        </div>
      </div>

      {/* Main chat area */}
      <div className="flex-1 flex flex-col">
        {activeRoom ? (
          <>
            {/* Header */}
            <div className="px-5 py-3 border-b border-gray-800 bg-gray-900 flex items-center justify-between">
              <div>
                <h2 className="font-semibold"># {activeRoom.name}</h2>
                {activeRoom.description && <p className="text-xs text-gray-500">{activeRoom.description}</p>}
              </div>
              <div className="flex items-center gap-2 text-xs text-gray-500">
                <span className="w-2 h-2 bg-green-500 rounded-full inline-block"></span>
                {onlineUsers.length} online
              </div>
            </div>

            {/* Messages */}
            <div className="flex-1 overflow-y-auto p-5 space-y-3">
              {messages.map((msg) => (
                <div key={msg.id} className={`flex gap-3 ${msg.user_name === user?.name ? 'flex-row-reverse' : ''}`}>
                  <div className="w-8 h-8 rounded-full bg-indigo-600 flex items-center justify-center text-xs font-bold shrink-0">
                    {msg.user_name?.[0]?.toUpperCase()}
                  </div>
                  <div className={`max-w-xs lg:max-w-md ${msg.user_name === user?.name ? 'items-end' : 'items-start'} flex flex-col`}>
                    <div className="flex items-baseline gap-2 mb-1">
                      <span className="text-xs font-semibold text-gray-300">{msg.user_name}</span>
                      <span className="text-xs text-gray-600">{format(msg.created_at)}</span>
                    </div>
                    <div className={`px-3 py-2 rounded-2xl text-sm ${
                      msg.user_name === user?.name
                        ? 'bg-indigo-600 text-white rounded-tr-sm'
                        : 'bg-gray-800 text-gray-100 rounded-tl-sm'
                    }`}>
                      {msg.content}
                    </div>
                  </div>
                </div>
              ))}
              <div ref={bottomRef} />
            </div>

            {/* Input */}
            <form onSubmit={handleSend} className="p-4 border-t border-gray-800 bg-gray-900">
              <div className="flex gap-3">
                <input
                  value={input}
                  onChange={(e) => setInput(e.target.value)}
                  placeholder={`Message #${activeRoom.name}`}
                  className="flex-1 bg-gray-800 border border-gray-700 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                />
                <button type="submit" disabled={!input.trim()}
                  className="bg-indigo-600 hover:bg-indigo-700 text-white px-5 py-2.5 rounded-xl text-sm font-semibold disabled:opacity-40 transition">
                  Send
                </button>
              </div>
            </form>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center text-gray-600">
            Select a room to start chatting
          </div>
        )}
      </div>

      {/* Online users panel */}
      {activeRoom && (
        <div className="w-48 bg-gray-900 border-l border-gray-800 p-3">
          <p className="text-xs text-gray-500 uppercase tracking-wide mb-3">Online — {onlineUsers.length}</p>
          <div className="space-y-2">
            {onlineUsers.map((name) => (
              <div key={name} className="flex items-center gap-2 text-sm text-gray-300">
                <span className="w-2 h-2 bg-green-500 rounded-full"></span>
                {name}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
