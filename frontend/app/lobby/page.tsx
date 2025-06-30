'use client';
import { useEffect, useRef, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import Game from '../../components/Game';

export default function Lobby() {
    const [currentWord, setCurrentWord] = useState('');
    const [lobbyState, setLobbyState] = useState('waiting');
    const isHostRef = useRef(false); // mose isHost into a useRef to always use the updated value
    const [isHost, setIsHost] = useState(false);
    // Variables for player1 game
    const [revealedWord, setRevealedWord] = useState('');
    const [attemptsLeft, setAttemptsLeft] = useState(6);
    const [instruction, setInstruction] = useState('Enter a word for your opponent to guess:');
    // Variables for player2 game
    const [opponentRevealed, setOpponentRevealed] = useState('');
    const [opponentAttempts, setOpponentAttempts] = useState(6);
    const [opponentInstruction, setOpponentInstruction] = useState('Waiting...');
    const [showRestart, setShowRestart] = useState(false);

    const ws = useRef<WebSocket | null>(null);
    const router = useRouter();
    const params = useSearchParams();
    const lobbyId = params.get('lobby');
    const playerId = params.get('playerID')


    useEffect(() => {
        ws.current = new WebSocket(`ws://localhost:8080/ws?lobby=${lobbyId}&id=${playerId}`);
        if (ws.current) {
            ws.current.onopen = () => {
                fetchLobbyState();
            };
            ws.current.onmessage = (event: MessageEvent) => {
                const msg = JSON.parse(event.data);
                if (msg.type === 'start_game') {
                    setLobbyState('playing');
                    setInstruction('Type a letter to guess.');
                    setRevealedWord(msg.revealed.split('').join(' '));
                    setAttemptsLeft(6);
                    setOpponentInstruction('');
                    setOpponentRevealed(msg.opponent_revealed.split('').join(' '));
                    setOpponentAttempts(6);
                } else if (msg.type === 'update') {
                    if (msg.revealed) {
                        setRevealedWord(msg.revealed.split('').join(' '));
                        setOpponentRevealed(msg.opponent_revealed.split('').join(' '));
                        setAttemptsLeft(msg.attempts);
                        setOpponentAttempts(msg.opponent_attempts);
                    }
                } else if (msg.type === 'win') {
                    console.log(`Player ${msg.player} won. You are the host: ${isHostRef.current}.`)
                    if ((msg.player == "1" && isHostRef.current) || (msg.player == "2" && !isHostRef.current)) {
                        setInstruction('You win!');
                        setRevealedWord(msg.word.split('').join(' '));
                    } else if ((msg.player == "1" && !isHostRef.current) || (msg.player == "2" && isHostRef.current)) {
                        setOpponentInstruction("Opponent won!");
                        setOpponentRevealed(msg.word.split('').join(' '));
                    }

                } else if (msg.type === 'lost') {
                    console.log(`Player ${msg.player} lost. You are the host: ${isHostRef.current}.`)
                    if ((msg.player == "1" && isHostRef.current) || (msg.player == "2" && !isHostRef.current)) {
                        setInstruction('Game Over! The word was:');
                        setRevealedWord(msg.word.split('').join(' '));
                    } else if ((msg.player == "1" && !isHostRef.current) || (msg.player == "2" && isHostRef.current)) {
                        setOpponentInstruction("Game Over! The word was:");
                        setOpponentRevealed(msg.word.split('').join(' '));
                    }
                } else if (msg.type === 'close') {
                    if (ws.current) ws.current.close();
                    router.push('/');
                } else if (msg.type === 'end') {
                    setLobbyState('ended');
                    setShowRestart(true);
                }
            };
        }
        return () => { };
        // eslint-disable-next-line
    }, [lobbyId]);

    const fetchLobbyState = async () => {
        const res = await fetch(`/api/lobby-state?lobby=${lobbyId}`, {
            method: 'GET',
            credentials: 'include',
        });
        if (res.ok) {
            const data = await res.json();
            setLobbyState(data.state);
            if (data.state === 'playing') setInstruction('Type a letter to guess.');
        }
    };

    const handleLeave = async () => {
        await fetch('/api/leave-lobby', {
            method: 'POST',
            credentials: 'include',
        });
    };

    const handleRestart = () => {
        setShowRestart(false);
        setInstruction('Enter a word for your opponent to guess:');
        setRevealedWord('')
        setLobbyState('waiting');
        if (ws.current) { ws.current.send(JSON.stringify({ type: 'restart', payload: 'r' })); }
    };

    const handleSubmitWord = () => {
        if (!currentWord) return alert('Enter a word first.');
        if (ws.current) { ws.current.send(JSON.stringify({ type: 'submit', payload: currentWord })); }
        setInstruction('Waiting for the other player to submit their word...');
        setCurrentWord('');
        setLobbyState('ready');
    };

    const handleKeyDown = (e: KeyboardEvent) => {
        const key: string = e.key.toLowerCase();
        // if (lobbyState === 'waiting' || lobbyState === 'ready') {
        if (lobbyState === 'waiting') {
            if (/^[a-z]$/.test(key)) setCurrentWord((w: string) => w + key);
            else if (e.key === 'Backspace') setCurrentWord((w: string) => w.slice(0, -1));
        } else if (lobbyState === 'playing') {
            if (ws.current) { ws.current.send(JSON.stringify({ type: 'guess', payload: key })); }
        }
    };

    // Check if we're the host or guest
    useEffect(() => {
        if (!lobbyId || !playerId) return;

        const fetchRole = async () => {
            try {
                const res = await fetch(`/api/player-role?lobby=${lobbyId}&id=${playerId}`, {
                    credentials: 'include',
                });
                if (res.ok) {
                    const data = await res.json();
                    setIsHost(data.role === 'host');
                    isHostRef.current = data.role === 'host';
                    console.log(`isHost: ${isHost}`)
                }
            } catch (error) {
                console.error("Failed to fetch role:", error);
            }
        };

        fetchRole();
    }, [playerId, lobbyId, isHost]);

    useEffect(() => {
        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
        // eslint-disable-next-line
    }, [lobbyState, currentWord]);

    return (
        <main style={{ display: 'flex', flexDirection: 'row' }}>
            {isHost ? (
                <>
                    <div style={{ flex: 1, padding: '1rem', border: '1px solid #ccc' }}>
                        <Game
                            playerName="You"
                            revealedWord={revealedWord}
                            attemptsLeft={attemptsLeft}
                            instruction={instruction}
                            isYou={true}
                        />
                        {/* {lobbyState === 'waiting' || lobbyState === 'ready' ? ( */}
                        {lobbyState === 'waiting' ? (
                            <div>
                                <h3>
                                    <span>{currentWord}</span>
                                    <span
                                        style={{
                                            display: 'inline-block',
                                            width: '1ch',
                                            animation: 'blink 1s steps(2, start) infinite',
                                            color: 'black',
                                        }}
                                    >|</span>
                                </h3>
                                <button onClick={handleSubmitWord}>Submit Word</button>
                            </div>
                        ) : null}
                        {showRestart && <button onClick={handleRestart}>Play Again</button>}
                        <br />
                        <button onClick={handleLeave}>Leave Lobby</button>
                    </div>
                    <Game
                        playerName="Opponent"
                        revealedWord={opponentRevealed}
                        attemptsLeft={opponentAttempts}
                        instruction={opponentInstruction}
                        isYou={false}
                    />
                </>
            ) : (
                <>
                    <Game
                        playerName="Opponent"
                        revealedWord={opponentRevealed}
                        attemptsLeft={opponentAttempts}
                        instruction={opponentInstruction}
                        isYou={false}
                    />
                    <div style={{ flex: 1, padding: '1rem', border: '1px solid #ccc' }}>
                        <Game
                            playerName="You"
                            revealedWord={revealedWord}
                            attemptsLeft={attemptsLeft}
                            instruction={instruction}
                            isYou={true}
                        />
                        {lobbyState === 'waiting' ? (
                            <div>
                                <h3>
                                    <span>{currentWord}</span>
                                    <span
                                        style={{
                                            display: 'inline-block',
                                            width: '1ch',
                                            animation: 'blink 1s steps(2, start) infinite',
                                            color: 'black',
                                        }}
                                    >|</span>
                                </h3>
                                <button onClick={handleSubmitWord}>Submit Word</button>
                            </div>
                        ) : null}
                        {showRestart && <button onClick={handleRestart}>Play Again</button>}
                        <br />
                        <button onClick={handleLeave}>Leave Lobby</button>
                    </div>
                </>
            )}
            <style>{`
                @keyframes blink {
                  0%, 100% { opacity: 1; }
                  50% { opacity: 0; }
                }
          `}</style>

        </main >

        // <main>
        //     <h2>{instruction}</h2>
        //     <h3 style={{ display: lobbyState === 'playing' ? '' : 'none' }}>
        //         Attempts Left: {attemptsLeft}
        //     </h3>
        //     {(lobbyState === 'waiting' || lobbyState === 'ready') && (
        //         <div>
        //             <h3>
        //                 <span>{currentWord}</span><span style={{ display: 'inline-block', width: '1ch', animation: 'blink 1s steps(2, start) infinite', color: 'black' }}>|</span>
        //             </h3>
        //             <button onClick={handleSubmitWord}>Submit Word</button>
        //         </div>
        //     )}
        //     <h2>{revealedWord}</h2>
        //     {showRestart && <button onClick={handleRestart}>Restart Game</button>}<br></br>
        //     <button onClick={handleLeave}>Leave Lobby</button>
        //     <style>{`
        //         @keyframes blink {
        //         0%, 100% { opacity: 1; }
        //         50% { opacity: 0; }
        //         }
        //     `}</style>
        // </main>

    )
}
