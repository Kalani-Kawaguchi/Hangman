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
        <main>
            <div style={{ minHeight: '25vh', width: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                    <Image src="/hangman.gif" alt="Hangman" width={0} height={0} style={{ height: '25vh', width: 'auto' }} />
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '20px', justifyContent: 'center', minHeight: '50vh' }}>
                <form 
                    onSubmit={handleSubmit}
                    style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '10px' }}
                >
                    <div style={{ display: 'flex', alignItems: 'center' }}>
                        <label htmlFor="lobbyName"><Image src="/lobbyName.gif" alt="Lobby Name" width={0} height={0} style={{ height: 'auto', width: '10vw'}}/></label>
                        <input
                            id="lobbyName"
                            value={lobbyName}
                            onChange={e => setLobbyName(e.target.value)}
                        />
                    </div>
                    <div style={{ display: 'flex', alignItems: 'center' }}>
                        <label htmlFor="playerName"><Image src="/playerName.gif" alt="Lobby Name" width={0} height={0} style={{ height: 'auto', width: '10vw'}}/></label>
                        <input
                            id="playerName"
                            value={playerName}
                            onChange={e => setPlayerName(e.target.value)}
                        />
                    </div>
                    <button type="submit">
                        <Image src="/createLobby.gif" alt="Create Lobby" width={0} height={0} style={{ height: 'auto', width: '10vw'}}/>
                    </button>
                </form>
                <Link href="/"><Image src="/goBack.gif" alt="Go Back" width={0} height={0} style={{ height: 'auto', width: '10vw'}}/></Link>
            </div>
        </main>
    );
}
