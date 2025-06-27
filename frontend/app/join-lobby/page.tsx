'use client';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link'

type Lobby = {
    id: string;
    name: string;
};

export default function JoinLobby() {
    const [lobbies, setLobbies] = useState<Lobby[]>([]);
    const [name, setName] = useState('');
    const router = useRouter();

    useEffect(() => {
        const fetchLobbies = async () => {
            const res = await fetch('/api/list-lobbies', {
                method: 'GET',
                credentials: 'include',
            });
            if (res.ok) setLobbies(await res.json());
        };
        fetchLobbies();
        const interval = setInterval(fetchLobbies, 1000);
        return () => clearInterval(interval);
    }, []);

    interface JoinLobbyRequest {
        lobby_id: string;
        player_name: string;
    }

    const joinLobby = async (lobbyId: string): Promise<void> => {
        if (!name.trim()) {
            alert('Please enter your name before joining a lobby!');
            return;
        }
        const body: JoinLobbyRequest = { lobby_id: lobbyId, player_name: name };
        const res: Response = await fetch('/api/join-lobby', {
            method: 'POST',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body),
        });
        if (res.ok) {
            router.push(`/lobby?lobby=${lobbyId}`);
        } else {
            alert('Failed to join lobby');
        }
    };

    return (
        <main>
            <ul>
                {Array.isArray(lobbies) && lobbies.length > 0 ? (
                    lobbies.map(lobby => (
                        <li key={lobby.id}>
                            <button type="button" onClick={() => joinLobby(lobby.id)}>
                                Lobby: {lobby.name}
                            </button>
                        </li>
                    ))
                ) : (
                    <li>No lobbies available.</li>
                )}
            </ul>
            <input
                value={name}
                onChange={e => setName(e.target.value)}
                placeholder="Enter your name"
            /><br />
            <Link href="/">Go Back</Link>
        </main>
    );
}
