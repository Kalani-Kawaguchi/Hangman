import Link from 'next/link'
import Image from 'next/image'
import Head from 'next/head';
import { useEffect } from "react";

export default function Home() {
  useEffect(() => {
    document.title = "Hangman"
  }, []);
  return (
    <>
      <Head>
        <title>{document.title}</title>
        <meta name="description" content="Multiplayer Hangman built with Go + React" />
      </Head>
      <main style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', minHeight: '100vh', justifyContent: 'flex-start' }}>
        <div style={{ height: '25vh', width: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
          <Image src="/hangman.gif" alt="Hangman" width={0} height={0} style={{ height: 'auto', width: '75vh' }} />
        </div>
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '20px', justifyContent: 'center', minHeight: '50vh' }}>
          <Link href="/create-lobby"><Image src="/createLobby.gif" alt='Create Lobby!' width={0} height={0} style={{ height: '10vh', width: 'auto' }}/></Link>
          <Link href="/join-lobby"><Image src="/joinLobby.gif" alt='Join Lobby!' width={0} height={0} style={{ height: '10vh', width: 'auto' }}/></Link>
        </div>
      </main>
    </>
  );
}
