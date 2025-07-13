'use client';
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link'
import Image from 'next/image'

export default function CreateLobby() {
    const [lobbyName, setLobbyName] = useState('');
    const [playerName, setPlayerName] = useState('');
    const router = useRouter();

    interface CreateLobbyRequest {
        lobby_name: string;
        host_name: string;
    }

    interface CreateLobbyResponse {
        id: string;
        playerID: string;
    }

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        if (!playerName.trim()) {
            alert('Please enter your name before creating a lobby!');
            return;
        }
        e.preventDefault();
        const body: CreateLobbyRequest = { lobby_name: lobbyName, host_name: playerName };
        const res: Response = await fetch('/api/create-lobby', {
            method: 'POST',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body),
        });
        if (res.ok) {
            const lobby: CreateLobbyResponse = await res.json();
            router.push(`/lobby?lobby=${lobby.id}&playerID=${lobby.playerID}`);
        } else {
            alert('Failed to create lobby');
        }
    };

    return (
        <>
            <main>
                <div style={{ height: '25vh', width: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                        <Image src="/hangman.gif" alt="Hangman" width={0} height={0} style={{ height: 'auto', width: '75vh' }} />
                </div>
                <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', height: '40vh'}}>
                    <form 
                        onSubmit={handleSubmit}
                        style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', height: '30vh'}}
                    >
                        <div style={{ display: 'flex', alignItems: 'center', height: '10vh'}}>
                            <label htmlFor="lobbyName" style={{height: '100%'}}><Image src="/lobbyName.gif" alt="Lobby Name" width={0} height={0} style={{ height: '100%', width: 'auto' }}/></label>
                            <input
                                id="lobbyName"
                                className="border-b-2 border-white"
                                value={lobbyName}
                                onChange={e => setLobbyName(e.target.value)}
                                style={{width: '50%', marginRight: '10px'}}
                            />
                        </div>
                        <div style={{ display: 'flex', alignItems: 'center', height: '10vh'}}>
                            <label htmlFor="playerName" style={{height: '100%'}}><Image src="/playerName.gif" alt="Lobby Name" width={0} height={0} style={{ height: '100%', width: 'auto'}}/></label>
                            <input
                                id="playerName"
                                className="border-b-2 border-white"
                                value={playerName}
                                onChange={e => setPlayerName(e.target.value)}
                                style={{width: '50%', marginRight: '10px'}}
                            />
                        </div>
                        <button type="submit" style={{height: '10vh'}}>
                            <Image src="/createLobby.gif" alt="Create Lobby" width={0} height={0} style={{ height: '90%', width: 'auto'}}/>
                        </button>
                    </form>
                    <Link href="/" style={{height: "10vh"}}><Image src="/goBack.gif" alt="Go Back" width={0} height={0} style={{ height: '100%', width: 'auto'}}/></Link>
                </div>
            </main>
        </>
    );
}
