import Link from 'next/link'

export default function Home() {
  return (
    <main>
      <Link href="/create-lobby">Create a Lobby!</Link><br></br>
      <Link href="/join-lobby">Join a Lobby!</Link>
    </main>
  );
}
