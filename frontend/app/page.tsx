import Link from 'next/link'
import Image from 'next/image'

export default function Home() {
  return (
    <main style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', minHeight: '100vh', justifyContent: 'flex-start' }}>
      <div style={{ minHeight: '25vh', width: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
        <Image src="/hangman.gif" alt="Hangman" width={0} height={0} style={{ height: '25vh', width: 'auto' }} />
      </div>
      <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '20px', justifyContent: 'center', minHeight: '50vh' }}>
        <Link href="/create-lobby"><Image src="/createLobby.gif" alt='Create Lobby!' width={0} height={0} style={{ height: '10vh', width: 'auto' }}/></Link>
        <Link href="/join-lobby"><Image src="/joinLobby.gif" alt='Join Lobby!' width={0} height={0} style={{ height: '10vh', width: 'auto' }}/></Link>
      </div>
    </main>
  );
}
