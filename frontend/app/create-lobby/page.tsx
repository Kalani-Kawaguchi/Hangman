'use client';
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link'

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
            <form onSubmit={handleSubmit}>
                <label>Lobby Name: </label>
                <input value={lobbyName} onChange={e => setLobbyName(e.target.value)} /><br />
                <label>Player Name: </label>
                <input value={playerName} onChange={e => setPlayerName(e.target.value)} /><br />
                <input type="submit" value="Create Lobby" />
            </form>
            <Link href="/">Go Back</Link>
        </main>
    );
}
