'use client';
import { useEffect, useRef, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

export default function Lobby() {
    const [currentWord, setCurrentWord] = useState('');
    const [lobbyState, setLobbyState] = useState('waiting');
    const [revealedWord, setRevealedWord] = useState('');
    const [attemptsLeft, setAttemptsLeft] = useState(6);
    const [instruction, setInstruction] = useState('Enter a word for your opponent to guess:');
    const [showRestart, setShowRestart] = useState(false);
    const ws = useRef<WebSocket | null>(null);
    const router = useRouter();
    const params = useSearchParams();
    const lobbyId = params.get('lobby');

    useEffect(() => {
        ws.current = new WebSocket(`ws://localhost:8080/ws?lobby=${lobbyId}`);
        if (ws.current) {
            ws.current.onopen = () => {
                fetchLobbyState();
            };
            ws.current.onmessage = (event: MessageEvent) => {
                const msg = JSON.parse(event.data);
                if (msg.type === 'start_game') {
                    setLobbyState('playing');
                    setInstruction('Game started! Type a letter to guess.');
                    setRevealedWord('');
                    setAttemptsLeft(6);
                } else if (msg.type === 'update') {
                    if (msg.revealed) {
                        setRevealedWord(msg.revealed.split('').join(' '));
                        setAttemptsLeft(msg.attempts);
                    }
                } else if (msg.type === 'win') {
                    setInstruction('You win!');
                } else if (msg.type === 'lost') {
                    setInstruction('You lost! The word was:');
                    setRevealedWord(msg.payload.split('').join(' '));
                } else if (msg.type === 'close') {
                    if (ws.current) ws.current.close();
                    router.push('/');
                } else if (msg.type === 'end') {
                    setLobbyState('ended');
                    setShowRestart(true);
                }
            };
        }
        return () => {};
        // eslint-disable-next-line
    }, [lobbyId]);

    const fetchLobbyState = async () => {
        const res = await fetch(`/api/lobby-state?lobby=${lobbyId}`);
        if (res.ok) {
            const data = await res.json();
            setLobbyState(data.state);
            if (data.state === 'playing') setInstruction('Game started! Type a letter to guess.');
        }
    };

    const handleLeave = async () => {
        await fetch('/api/leave-lobby', { method: 'POST' });
    };

    const handleRestart = () => {
        setShowRestart(false);
        setInstruction('Enter a word for your opponent to guess:');
        setLobbyState('waiting');
        if (ws.current) {ws.current.send(JSON.stringify({ type: 'restart', payload: 'r' }));}
    };

    const handleSubmitWord = () => {
        if (!currentWord) return alert('Enter a word first.');
        if (ws.current) {ws.current.send(JSON.stringify({ type: 'submit', payload: currentWord }));}
        setInstruction('Waiting for the other player to submit their word...');
        setCurrentWord('');
    };

    const handleKeyDown = (e: KeyboardEvent) => {
        const key: string = e.key.toLowerCase();
        if (lobbyState === 'waiting' || lobbyState === 'ready') {
            if (/^[a-z]$/.test(key)) setCurrentWord((w: string) => w + key);
            else if (e.key === 'Backspace') setCurrentWord((w: string) => w.slice(0, -1));
        } else if (lobbyState === 'playing') {
            if (ws.current) {ws.current.send(JSON.stringify({ type: 'guess', payload: key }));}
        }
    };

    useEffect(() => {
        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
        // eslint-disable-next-line
    }, [lobbyState, currentWord]);

    return (
        <main>
            <h2>{instruction}</h2>
            <h3 style={{ display: lobbyState === 'playing' ? '' : 'none' }}>
                Attempts Left: {attemptsLeft}
            </h3>
            {(lobbyState === 'waiting' || lobbyState === 'ready') && (
                <div>
                <h3>
                    <span>{currentWord}</span><span style={{ display: 'inline-block', width: '1ch', animation: 'blink 1s steps(2, start) infinite', color: 'black' }}>|</span>
                </h3>
                <button onClick={handleSubmitWord}>Submit Word</button>
                </div>
            )}
            <h2>{revealedWord}</h2>
            {showRestart && <button onClick={handleRestart}>Restart Game</button>}
            <button onClick={handleLeave}>Leave Lobby</button>
            <style>{`
                @keyframes blink {
                0%, 100% { opacity: 1; }
                50% { opacity: 0; }
                }
            `}</style>
        </main>
    )
}