'use client';
import { useEffect, useRef, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import Image from 'next/image'
import Game from '../../components/Game';

function useIsMobile(breakpoint = 768) {
    const [isMobile, setIsMobile] = useState(false);

    useEffect(() => {
        const update = () => setIsMobile(window.innerWidth < breakpoint);
        update(); // initial check
        window.addEventListener('resize', update);
        return () => window.removeEventListener('resize', update);
    }, [breakpoint]);

    return isMobile;
}

export default function Lobby() {
    const [currentWord, setCurrentWord] = useState('');
    const [lobbyState, setLobbyState] = useState('waiting');
    const isHostRef = useRef(false); // mose isHost into a useRef to always use the updated value
    const [isHost, setIsHost] = useState(false);
    const isP1Restarted = useRef(false);
    const [p1Restarted, setP1Restarted] = useState(false);
    const isP2Restarted = useRef(false);
    const [p2Restarted, setP2Restarted] = useState(false);
    const [playerName, setPlayerName] = useState('');
    const [opponentName, setOpponentName] = useState('');
    // Variables for player1 game
    const [revealedWord, setRevealedWord] = useState('');
    const [attemptsLeft, setAttemptsLeft] = useState("6");
    const [instruction, setInstruction] = useState('Enter a word for your opponent to guess:');
    // Variables for player2 game
    const [opponentRevealed, setOpponentRevealed] = useState('');
    const [opponentAttempts, setOpponentAttempts] = useState("6");
    const [opponentInstruction, setOpponentInstruction] = useState('Picking a word.');
    const [showRestart, setShowRestart] = useState(false);

    const ws = useRef<WebSocket | null>(null);
    const router = useRouter();
    const params = useSearchParams();
    const lobbyId = params.get('lobby');
    const playerId = params.get('playerID')
    const opponentExistsRef = useRef(false);
    const [opponentExists, setOpponentExists] = useState(false);
    const isMobile = useIsMobile();

    // check to see when instruction is updated
    useEffect(() => {
        console.log("Instruction updated:", instruction);
    }, [instruction]);

    useEffect(() => {
        if (!lobbyId || !playerId) return;

        // Only create the websocket if it doesn't already exist
        if (ws.current) return;

        const socket = new WebSocket(`ws://localhost:8080/ws?lobby=${lobbyId}&id=${playerId}`);
        ws.current = socket;

        socket.onopen = () => {
            fetchLobbyState();
        };
        socket.onclose = () => {
            // No need to call close again, just cleanup reference
            ws.current = null;
        };
        socket.onmessage = (event: MessageEvent) => {
            const msg = JSON.parse(event.data);
            if (msg.type === 'start_game') {
                setLobbyState('playing');
                setInstruction('Type a letter to guess.');
                setRevealedWord(msg.revealed);
                setAttemptsLeft("6");
                setOpponentInstruction('');
                setOpponentRevealed(msg.opponent_revealed);
                setOpponentAttempts("6");
                setP1Restarted(false);
                isP1Restarted.current = false;
                setP2Restarted(false);
                isP2Restarted.current = false;

            } else if (msg.type === 'update') {
                if (msg.revealed) {
                    setRevealedWord(msg.revealed);
                    setOpponentRevealed(msg.opponent_revealed);
                    setAttemptsLeft(String(msg.attempts));
                    setOpponentAttempts(String(msg.opponent_attempts));
                }

            } else if (msg.type === 'win') {
                console.log(`Player ${msg.player} won. You are the host: ${isHostRef.current}.`)
                if ((msg.player == "1" && isHostRef.current) || (msg.player == "2" && !isHostRef.current)) {
                    setInstruction('You win!');
                    setRevealedWord(msg.word);
                } else if ((msg.player == "1" && !isHostRef.current) || (msg.player == "2" && isHostRef.current)) {
                    setOpponentInstruction("Opponent won!");
                    setOpponentRevealed(msg.word.split('').join(' '));
                }

            } else if (msg.type === 'lost') {
                console.log(`Player ${msg.player} lost. You are the host: ${isHostRef.current}.`)
                if ((msg.player == "1" && isHostRef.current) || (msg.player == "2" && !isHostRef.current)) {
                    setInstruction('Game Over! The word was:');
                    setRevealedWord(msg.word);
                } else if ((msg.player == "1" && !isHostRef.current) || (msg.player == "2" && isHostRef.current)) {
                    setOpponentInstruction("Game Over! The word was:");
                    setOpponentRevealed(msg.word);
                }

            } else if (msg.type === 'join') {
                console.log('A player joined the lobby');
                if (isHostRef.current) {
                    setOpponentExists(true)
                    opponentExistsRef.current = true;
                    setOpponentName(msg.message)
                }
                console.log(`opp exists: ${opponentExistsRef.current}`);

            } else if (msg.type === 'submit') {
                if ((isHostRef.current && msg.player === '2') || (!isHostRef.current && msg.player === '1')) {
                    setOpponentInstruction("Ready. Waiting for you...");
                }

            } else if (msg.type === 'restart') {
                console.log(`Player ${msg.player} wants to play again.`)
                if (msg.player == "1" && !isHostRef.current) {
                    setOpponentInstruction('Wants to play again.');
                    setOpponentRevealed("");
                    setP1Restarted(true);
                    isP1Restarted.current = true;
                    if (isP1Restarted.current) { console.log("Set P1 to true"); }
                } else if (msg.player == "2" && isHostRef.current) {
                    setOpponentInstruction('Wants to play again.');
                    setOpponentRevealed("");
                    setP2Restarted(true);
                    isP2Restarted.current = true;
                    if (isP2Restarted.current) { console.log("Set P2 to true"); }
                }
                if (isP1Restarted.current && isP2Restarted.current) {
                    console.log("Both restarted");
                    setOpponentInstruction("Picking a word.");
                }

            } else if (msg.type === 'close') {
                console.log(`Player ${msg.player} left lobby`);
                if (msg.player === '2' && isHostRef.current) {
                    setOpponentExists(false);
                    opponentExistsRef.current = false;
                    console.log(`opp exists: ${opponentExistsRef.current}`);
                }
                if (msg.player === '1' || msg.player === '2' && !isHostRef.current) {
                    if (ws.current) ws.current.close();
                    router.push('/');
                }

            } else if (msg.type === 'end') {
                setLobbyState('ended');
                setShowRestart(true);
            }
        };

        return () => { };
        // eslint-disable-next-line
    }, [lobbyId, playerId]);

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
        console.log("restarted");
        setInstruction('Enter a word for your opponent to guess:');
        setRevealedWord('')
        setLobbyState('waiting');
        if (isHostRef.current) {
            setP1Restarted(true);
            isP1Restarted.current = true;
        } else {
            setP2Restarted(true);
            isP2Restarted.current = true;
        }
        if (ws.current) { ws.current.send(JSON.stringify({ type: 'restart', payload: playerId })); }
    };

    const handleSubmitWord = () => {
        if (!currentWord) return alert('Enter a word first.');
        if (ws.current) { ws.current.send(JSON.stringify({ type: 'submit', payload: currentWord })); }
        setInstruction('waiting for opponent word');
        console.log("submit word");
        // send a msg to backend telling other client to update opponents instruction to "They have submit their word"
        // if (ws.current) { ws.current.send(JSON.stringify({ type: 'submit', payload: "" })); }
        setCurrentWord('');
        setLobbyState('ready');
    };

    const handleKeyDown = (e: KeyboardEvent) => {
        const key: string = e.key.toLowerCase();
        if (lobbyState === 'waiting') {
            if (/^[a-z]$/.test(key)) setCurrentWord((w: string) => w + key);
            else if (e.key === 'Backspace') { setCurrentWord((w: string) => w.slice(0, -1)); e.preventDefault(); }
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
                    setPlayerName(data.name)
                    setIsHost(data.role === 'host');
                    isHostRef.current = data.role === 'host';
                    if (!isHostRef.current) {
                        setOpponentExists(true);
                        opponentExistsRef.current = true; // If you're player 2, your opponent should always exist.
                        setOpponentName(data.opponent);
                        console.log(`isHost: ${isHost}`)
                    }
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
        < main >
            {/* Banner */}
            <div style={{ height: '25vh', width: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                <Image src="/hangman.gif" alt="Hangman" width={0} height={0} style={{ height: 'auto', width: '75vh' }} />
            </div>

            {/* Game Section */}
            <div style={{ display: 'flex', flexDirection: 'row' }}>
                {isHost ? (
                    <>
                        {/* Always show your game */}
                        <div style={{ flex: 1, padding: '1rem', border: '1px solid #ccc' }}>
                            <Game
                                playerName={playerName}
                                revealedWord={lobbyState === 'waiting' ? currentWord : revealedWord}
                                attemptsLeft={attemptsLeft}
                                instruction={instruction}
                                isYou={true}
                            />
                            {lobbyState === 'waiting' ? (
                                <div className="flex justify-center-safe w-full">
                                    <button className="flex justify-center" onClick={handleSubmitWord}>
                                        <img className="w-[40%]" src="/submitWord.gif" />
                                    </button>
                                </div>
                            ) : null}
                            {showRestart && (
                                <div className="flex justify-center-safe w-full">
                                    <button className="flex justify-center" onClick={handleRestart}>
                                        <Image className="w-[40%]" src="/PlayAgain.gif" alt="Play again button" width={0} height={0} />
                                    </button>
                                </div>
                            )}
                            <br />
                        </div>

                        {/* Only show opponent game if NOT on mobile */}
                        {!isMobile && (
                            <div style={{ flex: 1, padding: '1rem', border: '1px solid #ccc' }}>
                                {opponentExists ? (
                                    <Game
                                        playerName={opponentName}
                                        revealedWord={isP2Restarted.current ? "" : opponentRevealed}
                                        attemptsLeft={opponentAttempts}
                                        instruction={opponentInstruction}
                                        isYou={false}
                                    />

                                ) : (
                                    <div className="flex justify-center items-center h-full w-full">
                                        <Image src="/WaitingForOpponent.gif" alt="Waiting for an opponent" width={500} height={500} />
                                    </div>
                                )}
                            </div>
                        )}
                    </>
                ) : (
                    <>

                        {/* Only show opponent game if NOT on mobile */}
                        {!isMobile && (
                            <div style={{ flex: 1, padding: '1rem', border: '1px solid #ccc' }}>
                                <Game
                                    playerName={opponentName}
                                    revealedWord={isP1Restarted.current ? "" : opponentRevealed}
                                    attemptsLeft={opponentAttempts}
                                    instruction={opponentInstruction}
                                    isYou={false}
                                />
                            </div>
                        )}

                        {/* Always show your game */}
                        <div style={{ flex: 1, padding: '1rem', border: '1px solid #ccc' }}>
                            <Game
                                playerName={playerName}
                                revealedWord={lobbyState === 'waiting' ? currentWord : revealedWord}
                                attemptsLeft={attemptsLeft}
                                instruction={instruction}
                                isYou={true}
                            />
                            {lobbyState === 'waiting' ? (
                                <div className="flex justify-center w-full">
                                    <button className="flex justify-center" onClick={handleSubmitWord}>
                                        <img className="w-[40%]" src="/submitWord.gif" />
                                    </button>
                                </div>
                            ) : null}
                            {showRestart && (
                                <div className="flex justify-center-safe">
                                    <button className="flex justify-center" onClick={handleRestart}>
                                        <Image className="w-[40%]" src="/PlayAgain.gif" alt="Play again button" width={0} height={0} />
                                    </button>
                                </div>
                            )}
                            <br />
                        </div>
                    </>
                )}
                <style>{`
                @keyframes blink {
                  0%, 100% { opacity: 1; }
                  50% { opacity: 0; }
                }
          `}</style>
            </div>

            {/*Leave Button */}
            <div className="mt-5 w-[20%]" >
                <button onClick={handleLeave}><img src="/leaveLobby.gif" /></button>
            </div>
        </main >
    )
}
