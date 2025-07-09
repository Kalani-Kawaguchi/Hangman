import Image from 'next/image'

type GameProps = {
    playerName: string;
    revealedWord: string;
    attemptsLeft: string;
    instruction: string;
    isYou: boolean;
};

const getInstructionGif = (instruction: string): string | null => {
    const map: Record<string, string> = {
        'You win!': '/YouWin.gif',
        'You lost! The word was:': '/YouLose.gif',
        'Game Over! The word was:': '/GameOver.gif',
        // 'Waiting for the other player to submit their word...': '/waiting.gif',
        // 'Enter a word for your opponent to guess:': '/typeword.gif',
    };
    return map[instruction] ?? null;
};

export default function Game({
    playerName,
    revealedWord,
    attemptsLeft,
    instruction,
    isYou,
}: GameProps) {
    const gif = getInstructionGif(instruction);
    return (
        <div className="items-center text-center" style={{ flex: 1, padding: '1rem' }}>
            <h2 >{playerName}</h2>
            <Image
                src={`/Platform${attemptsLeft}.gif`}
                width={500}
                height={500}
                alt="Hangman picture"
            />
            <p>{gif !== null ? "" : instruction}</p>
            {gif && (
                <Image
                    src={gif}
                    width={250}
                    height={100}
                    alt="Instruction GIF"
                    className="mx-auto"
                />
            )}
            <h3>{revealedWord}</h3>
        </div>
    );
}
