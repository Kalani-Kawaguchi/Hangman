'use client';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link'
import Image from 'next/image'
import LobbyBox from '../../components/LobbyBox';

export const metadata = {
  title: "Hangman - Join Lobby",
  description: "Join a Lobby to Play Hangman"
}

type Lobby = {
    id: string;
    name: string;
    playerCount: number;
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
        const interval = setInterval(fetchLobbies, 3000);
        return () => clearInterval(interval);
    }, []);

    interface JoinLobbyRequest {
        lobby_id: string;
        player_name: string;
    }

    interface CreateResponse {
        playerID: string;
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
            const resp: CreateResponse = await res.json();
            router.push(`/lobby?lobby=${lobbyId}&playerID=${resp.playerID}`);
        } else {
            alert('Failed to join lobby');
        }
    };

    return (
        <>
            <main>
                <div style={{ height: '25vh', width: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                    <Image src="/hangman.gif" alt="Hangman" width={0} height={0} style={{ height: 'auto', width: '75vh' }} />
                </div>
                <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', height: '10vh' }}>
                    <Image src="/lobbies.gif" alt="Lobbies" width={0} height={0} style={{ height: '100%', width: 'auto' }} />
                </div>
                <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', height: '40vh' }}>
                    <div style={{ height: '30vh', width: '100%'}}>
                        <ul>
                            {Array.isArray(lobbies) && lobbies.length > 0 ? (
                                lobbies.map(lobby => (
                                    <li key={lobby.id} style={{ marginBottom: '20px' }}>
                                        <LobbyBox
                                            name={lobby.name}
                                            id={lobby.id}
                                            playerCount={lobby.playerCount+""}
                                            maxPlayers={"2"}
                                            onClick={joinLobby}
                                        />
                                    </li>
                                ))
                            ) : (
                                <li><Image src="/noLobby.gif" alt="No Lobbies" width={0} height={0} style={{ height: '100%', width: 'auto' }} /></li>
                            )}
                        </ul>
                    </div>
                    <br />
                    <div style={{ display: 'flex', alignItems: 'center', height: '10vh' }}>
                        <label htmlFor="playerName" style={{ height: '100%' }}><Image src="/playerName.gif" alt="Lobby Name" width={0} height={0} style={{ height: '100%', width: 'auto' }} /></label>
                        <input
                            id="playerName"
                            className="border-b-2 border-white"
                            value={name}
                            onChange={e => setName(e.target.value)}
                            style={{width: '50%', marginRight: '10px'}}
                        />
                    </div>
                    <br />
                    <Link href="/" style={{ height: "10vh" }}><Image src="/goBack.gif" alt="Go Back" width={0} height={0} style={{ height: '100%', width: 'auto' }} /></Link>
                </div>
            </main>
        </>
    );
}
