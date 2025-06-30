type GameProps = {
    playerName: string;
    revealedWord: string;
    attemptsLeft: number;
    instruction: string;
    isYou: boolean;
};

export default function Game({
    playerName,
    revealedWord,
    attemptsLeft,
    instruction,
    isYou,
}: GameProps) {
    return (
        <div style={{ flex: 1, padding: '1rem', border: '1px solid #ccc' }}>
            <h2>{isYou ? 'You' : playerName}</h2>
            <p>{instruction}</p>
            <h3>Word: {revealedWord}</h3>
            <h3>Attempts Left: {attemptsLeft}</h3>
        </div>
    );
}
