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

function Search(field: string) {
    const searchParams = useSearchParams()

    return searchParams.get(field);
}

export default function Lobby() {
    const [currentWord, setCurrentWord] = useState('');
    const [lobbyState, setLobbyState] = useState('waiting');
    const isHostRef = useRef(false); // move isHost into a useRef to always use the updated value
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
    const [instruction, setInstruction] = useState('');
    const [guessedLetters, setGuessedLetters] = useState('');
    // Variables for player2 game
    const [opponentRevealed, setOpponentRevealed] = useState('');
    const [opponentAttempts, setOpponentAttempts] = useState("6");
    const [opponentInstruction, setOpponentInstruction] = useState('');
    const [opponentGuessedLetters, setOpponentGuessedLetters] = useState('');

    const ws = useRef<WebSocket | null>(null);
    const router = useRouter();
    const lobbyId = Search('lobby');
    const playerId = Search('playerID')

    // const [hostState, setHostState] = useState('');
    // const [oppState, setOppState] = useState('');
    const opponentExistsRef = useRef(false);
    const [opponentExists, setOpponentExists] = useState(false);
    const isMobile = useIsMobile();
    let newInstruction;
    let newOppInstruction;

    useEffect(() => {
        if (!lobbyId || !playerId) return;

        // Only create the websocket if it doesn't already exist
        if (ws.current) return;

        const socket = new WebSocket(`wss://hangman-qrdh.onrender.com/ws?lobby=${lobbyId}&id=${playerId}`);
        // const socket = new WebSocket(`ws://localhost:8080/ws?lobby=${lobbyId}&id=${playerId}`);
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
                newInstruction = 'Type a letter to guess.';
                setInstruction(newInstruction);
                setRevealedWord(msg.revealed);
                setAttemptsLeft("6");
                newOppInstruction = ''
                setOpponentInstruction(newOppInstruction);
                if (isHostRef.current) {
                    if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'One', instruction: newInstruction } })); }
                    if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'OneOpp', instruction: newOppInstruction } })); }
                } else {
                    if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'Two', instruction: newInstruction } })); }
                    if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'TwoOpp', instruction: newOppInstruction } })); }
                }
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
                    setGuessedLetters(msg.guessed_letters);
                    setOpponentGuessedLetters(msg.opponent_guessed_letters);
                }

            } else if (msg.type === 'win') {
                console.log(`Player ${msg.player} won. You are the host: ${isHostRef.current}.`)
                if (msg.player == "1" && isHostRef.current) {
                    newInstruction = 'You win!';
                    setInstruction(newInstruction);
                    if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'One', instruction: newInstruction } })); }
                    setRevealedWord(msg.word);
                } else if (msg.player == "2" && !isHostRef.current) {
                    newInstruction = 'You win!';
                    setInstruction(newInstruction);
                    if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'Two', instruction: newInstruction } })); }
                    setRevealedWord(msg.word);
                } else if (msg.player == "1" && !isHostRef.current) {
                    newOppInstruction = 'Opponent won!';
                    setOpponentInstruction(newOppInstruction);
                    if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'TwoOpp', instruction: newOppInstruction } })); }
                    setOpponentRevealed(msg.word.split('').join(' '));
                }
                else if (msg.player == "2" && isHostRef.current) {
                    newOppInstruction = 'Opponent won!';
                    setOpponentInstruction(newOppInstruction);
                    if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'OneOpp', instruction: newOppInstruction } })); }
                    setOpponentRevealed(msg.word.split('').join(' '));
                }
            } else if (msg.type === 'lost') {
                console.log(`Player ${msg.player} lost. You are the host: ${isHostRef.current}.`)
                if (msg.player == "1" && isHostRef.current) {
                    newInstruction = 'Game Over! The word was:';
                    setInstruction(newInstruction);
                    if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'One', instruction: newInstruction } })); }
                    setRevealedWord(msg.word);
                } else if (msg.player == "2" && !isHostRef.current) {
                    newInstruction = 'Game Over! The word was:';
                    setInstruction(newInstruction);
                    if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'Two', instruction: newInstruction } })); }
                    setRevealedWord(msg.word);
                } else if (msg.player == "1" && !isHostRef.current) {
                    newOppInstruction = 'Game Over! The word was:';
                    setOpponentInstruction(newOppInstruction);
                    if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'TwoOpp', instruction: newOppInstruction } })); }
                    setOpponentRevealed(msg.word.split('').join(' '));
                }
                else if (msg.player == "2" && isHostRef.current) {
                    newOppInstruction = 'Game Over! The word was:';
                    setOpponentInstruction(newOppInstruction);
                    if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'OneOpp', instruction: newOppInstruction } })); }
                    setOpponentRevealed(msg.word.split('').join(' '));
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

                    newOppInstruction = 'Ready. Waiting for you...';
                    setOpponentInstruction(newOppInstruction);
                    if (isHostRef.current) {
                        if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'OneOpp', instruction: newOppInstruction } })); }
                    } else {
                        if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'TwoOpp', instruction: newOppInstruction } })); }
                    }
                }

            } else if (msg.type === 'restart') {
                console.log(`Player ${msg.player} wants to play again.`)
                if (msg.player == "1" && !isHostRef.current) {
                    newOppInstruction = 'Wants to play again.';
                    setOpponentInstruction(newOppInstruction);
                    if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'TwoOpp', instruction: newOppInstruction } })); }
                    setOpponentRevealed("");
                    setP1Restarted(true);
                    isP1Restarted.current = true;
                    if (isP1Restarted.current) { console.log("Set P1Restarted to true"); }
                } else if (msg.player == "2" && isHostRef.current) {
                    newOppInstruction = 'Wants to play again.';
                    setOpponentInstruction(newOppInstruction);
                    if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'OneOpp', instruction: newOppInstruction } })); }
                    setOpponentRevealed("");
                    setP2Restarted(true);
                    isP2Restarted.current = true;
                    if (isP2Restarted.current) { console.log("Set P2Restarted to true"); }
                }
                if (isP1Restarted.current && isP2Restarted.current) {
                    console.log("Both Players restarted");
                    newOppInstruction = 'Picking a word.';
                    setOpponentInstruction(newOppInstruction);
                    if (isHostRef.current) {
                        if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'OneOpp', instruction: newOppInstruction } })); }
                    } else {
                        if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'TwoOpp', instruction: newOppInstruction } })); }
                    }
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
            }
        };

        return () => { };

    }, [lobbyId, playerId]);

    const fetchLobbyState = async () => {
        console.log("Fetching lobby state...");
        const res = await fetch(`/api/lobby-state?lobby=${lobbyId}`, {
            method: 'GET',
            credentials: 'include',
        });
        if (res.ok) {
            const data = await res.json();
            console.log(data);
            setLobbyState(data.state);
            if (data.state === 'playing') {
                setInstruction('Type a letter to guess.');
            }
            if (isHostRef.current) {
                if (data.player2Exists === true) {
                    setOpponentExists(true);
                    opponentExistsRef.current = true;
                    setOpponentName(data.player2Name);
                }
                setInstruction(data.player1Instruction);
                setOpponentInstruction(data.player1OppInstruction);
                setGuessedLetters(data.player1GuessedLetters);
                setOpponentGuessedLetters(data.player2GuessedLetters);
                if (data.player1Restarted === true) {
                    setP1Restarted(true);
                    console.log("Setting showRestart to true");
                } else {
                    setP1Restarted(false);
                    console.log("Setting showRestart to true");
                }
                if (data.player2Restarted === true) {
                    setP2Restarted(true);
                } else {
                    setP2Restarted(false);
                }
                if (data.player1Ready && data.player2Ready) {
                    setRevealedWord(data.player1RevealedWord);
                    setOpponentRevealed(data.player2RevealedWord);
                }
                if (data.state === 'ended') {
                    setRevealedWord('');
                    setOpponentRevealed('');
                }
                setAttemptsLeft(data.player1AttemptsLeft);
                if (data.player1Ready === false) {
                    setAttemptsLeft("6")
                }
                setOpponentAttempts(data.player2AttemptsLeft);
                if (data.player2Ready === false) {
                    setOpponentAttempts("6")
                }
            }
            if (!isHostRef.current) {
                if (data.Player1Exists === true) {
                    setOpponentExists(true)
                    opponentExistsRef.current = true;
                    setOpponentName(data.player1Name);
                }
                setInstruction(data.player2Instruction);
                setOpponentInstruction(data.player2OppInstruction);
                setGuessedLetters(data.player2GuessedLetters);
                setOpponentGuessedLetters(data.player1GuessedLetters);
                if (data.player1Restarted === true) {
                    setP1Restarted(true);
                } else {
                    setP1Restarted(false);
                }
                if (data.player2Restarted === true) {
                    setP2Restarted(true);
                    console.log("Setting showRestart to true");
                } else {
                    setP2Restarted(false);
                    console.log("Setting showRestart to true");
                }
                if (data.player1Ready && data.player2Ready) {
                    setRevealedWord(data.player2RevealedWord);
                    setOpponentRevealed(data.player1RevealedWord);
                }
                if (data.state === 'ended') {
                    setRevealedWord('');
                    setOpponentRevealed('');
                }
                setAttemptsLeft(data.player2AttemptsLeft);
                if (data.player2Ready === false) {
                    setAttemptsLeft("6")
                }
                setOpponentAttempts(data.player1AttemptsLeft);
                if (data.player1Ready === false) {
                    setOpponentAttempts("6")
                }
            }
        }
    };

    const handleLeave = async () => {
        await fetch('/api/leave-lobby', {
            method: 'POST',
            credentials: 'include',
        });
    };

    const handleRestart = () => {
        console.log("restarted");
        setRevealedWord('')
        setLobbyState('waiting');
        newInstruction = 'Enter a word for your opponent to guess:';
        setInstruction(newInstruction);
        if (isHostRef.current) {
            if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'One', instruction: newInstruction } })); }
            setP1Restarted(true);
            isP1Restarted.current = true;
        } else {
            if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'Two', instruction: newInstruction } })); }
            setP2Restarted(true);
            isP2Restarted.current = true;
        }
        if (ws.current) { ws.current.send(JSON.stringify({ type: 'restart', payload: playerId })); }
    };

    const handleSubmitWord = () => {
        if (!currentWord) return alert('Enter a word first.');
        if (ws.current) { ws.current.send(JSON.stringify({ type: 'submit', payload: currentWord })); }
        newInstruction = 'waiting for opponent word';
        setInstruction(newInstruction);
        if (isHostRef.current) {
            if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'One', instruction: newInstruction } })); }
        } else {
            if (ws.current) { ws.current.send(JSON.stringify({ type: 'instruction', payload: { player: 'Two', instruction: newInstruction } })); }
        }
        console.log("submit word");
        // send a msg to backend telling other client to update opponents instruction to "They have submit their word"
        // if (ws.current) { ws.current.send(JSON.stringify({ type: 'submit', payload: "" })); }
        setCurrentWord('');
        setLobbyState('ready');
    };

    const handleKeyDown = (e: KeyboardEvent) => {
        const key: string = e.key.toLowerCase();
        if (lobbyState === 'waiting' || instruction === 'Enter a word for your opponent to guess:') {
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
        if (!isMobile) {
            window.addEventListener('keydown', handleKeyDown);
            return () => window.removeEventListener('keydown', handleKeyDown);
        }
        // eslint-disable-next-line
    }, [lobbyState, currentWord]);



    return (
        <>
            <title>Hangman - Lobby</title>
            <meta name="description" content="Play Hangman" />
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
                            <div style={{ flex: 1, padding: '1rem', borderRight: '1px solid #ccc' }}>
                                <Game
                                    playerName={playerName}
                                    revealedWord={(lobbyState === 'waiting' || (lobbyState === 'ended' && p1Restarted)) ? currentWord : revealedWord}
                                    attemptsLeft={attemptsLeft}
                                    guessedLetters={guessedLetters}
                                    instruction={instruction}
                                    isMobile={isMobile}
                                    guessing={lobbyState === 'playing' || lobbyState === 'ended'}
                                />
                                {isMobile && (
                                    <input
                                        id="mobileKeyboardInput"
                                        type="text"
                                        inputMode="text"
                                        autoFocus
                                        onBlur={(e) => e.target.focus()} // re-focus if it blurs
                                        onChange={() => { }} // prevents React warning
                                        onKeyDown={(e) => handleKeyDown(e.nativeEvent)}
                                        style={{
                                            position: 'absolute',
                                            bottom: 0,
                                            left: 0,
                                            width: '1px',
                                            height: '1px',
                                            opacity: 0,
                                            zIndex: -1,
                                            pointerEvents: 'none',
                                        }}
                                    />
                                )}
                                {instruction === 'Enter a word for your opponent to guess:' ? (
                                    <div className="flex justify-center-safe w-full">
                                        <button className="flex justify-center" onClick={handleSubmitWord}>
                                            <img className="w-[40%]" src="/submitWord.gif" />
                                        </button>
                                    </div>
                                ) : null}
                                {(instruction == 'You win!' || instruction == 'Game Over! The word was:') ? (
                                    <div className="flex justify-center-safe w-full">
                                        <button className="flex justify-center" onClick={handleRestart}>
                                            <Image className="w-[40%]" src="/PlayAgain.gif" alt="Play again button" width={0} height={0} />
                                        </button>
                                    </div>

                                ) : (
                                    <div></div>
                                )}
                                <br />
                            </div>

                            {/* Only show opponent game if NOT on mobile */}
                            {!isMobile && (
                                <div style={{ flex: 1, padding: '1rem', borderLeft: '1px solid #ccc' }}>
                                    {opponentExists ? (
                                        <Game
                                            playerName={opponentName}
                                            revealedWord={isP2Restarted.current ? "" : opponentRevealed}
                                            attemptsLeft={opponentAttempts}
                                            guessedLetters={opponentGuessedLetters}
                                            instruction={opponentInstruction}
                                            isMobile={isMobile}
                                            guessing={lobbyState === 'playing' || lobbyState === 'ended'}
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
                                <div style={{ flex: 1, padding: '1rem', borderRight: '1px solid #ccc' }}>
                                    <Game
                                        playerName={opponentName}
                                        revealedWord={isP1Restarted.current ? "" : opponentRevealed}
                                        attemptsLeft={opponentAttempts}
                                        guessedLetters={opponentGuessedLetters}
                                        instruction={opponentInstruction}
                                        isMobile={isMobile}
                                        guessing={lobbyState === 'playing' || lobbyState === 'ended'}
                                    />
                                </div>
                            )}

                            {/* Always show your game */}
                            <div style={{ flex: 1, padding: '1rem', borderLeft: '1px solid #ccc' }}>
                                <Game
                                    playerName={playerName}
                                    revealedWord={(lobbyState === 'waiting' || (lobbyState === 'ended' && p2Restarted)) ? currentWord : revealedWord}
                                    attemptsLeft={attemptsLeft}
                                    guessedLetters={guessedLetters}
                                    instruction={instruction}
                                    isMobile={isMobile}
                                    guessing={lobbyState === 'playing' || lobbyState === 'ended'}
                                />
                                {isMobile && (
                                    <input
                                        id="mobileKeyboardInput"
                                        type="text"
                                        inputMode="text"
                                        autoFocus
                                        onBlur={(e) => e.target.focus()} // re-focus if it blurs
                                        onChange={() => { }} // prevents React warning
                                        onKeyDown={(e) => handleKeyDown(e.nativeEvent)}
                                        style={{
                                            position: 'absolute',
                                            bottom: 0,
                                            left: 0,
                                            width: '1px',
                                            height: '1px',
                                            opacity: 0,
                                            zIndex: -1,
                                            pointerEvents: 'none',
                                        }}
                                    />
                                )}
                                {instruction === 'Enter a word for your opponent to guess:' ? (
                                    < div className="flex justify-center w-full">
                                        <button className="flex justify-center" onClick={handleSubmitWord}>
                                            <img className="w-[40%]" src="/submitWord.gif" />
                                        </button>
                                    </div>
                                ) : null}
                                {(instruction == 'You win!' || instruction == 'Game Over! The word was:') ? (
                                    <div className="flex justify-center-safe">
                                        <button className="flex justify-center" onClick={handleRestart}>
                                            <Image className="w-[40%]" src="/PlayAgain.gif" alt="Play again button" width={0} height={0} />
                                        </button>
                                    </div>
                                ) : (
                                    <div></div>
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
        </>
    )
}
