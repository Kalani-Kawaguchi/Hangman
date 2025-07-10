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
        'Picking a word.': '/PickingWord.gif',
        'Ready. Waiting for you...': '/Ready.gif',
        'You win!': '/Winner.gif',
        'Opponent won!': '/Winner.gif',
        'Game Over! The word was:': '/GameOver.gif',
        'Type a letter to guess.': '/GuessLetter.gif',
        'Enter a word for your opponent to guess:': '/EnterWord.gif',
        'Waiting for the other player to submit their word...': '/Waiting.gif',
        'Wants to play again.': '/WantsToPlayAgain.gif',
        '': '/Guessing.gif',
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
    const nameChar = playerName.toLowerCase().split('');
    const revealedChar = revealedWord.toLowerCase().split('');

    const getImageName = (char) => {
        if (char === "_") return "underscore";
        if (char === " ") return "space";
        return char;
    };

    return (
        <div className="items-center text-center" style={{ flex: 1, padding: '1rem' }}>
            <h2 style={{display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
                {nameChar.map((char, index) => (
                    <img
                    key={index}
                    src={`/${char}.gif`}
                    alt={char}
                    />
                ))}
            </h2>
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
            <h3 style={{display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
                {revealedChar.map((char, index) => (
                    <img
                    key={index}
                    src={`/${getImageName(char)}.gif`}
                    alt={char}
                    />
                ))}
            </h3>
        </div>
    );
}
