import Image from 'next/image'

type GameProps = {
    playerName: string;
    revealedWord: string;
    attemptsLeft: string;
    guessedLetters: string;
    instruction: string;
    isMobile: boolean;
    guessing: boolean;
};

const getInstructionGif = (instruction: string): string | null => {
    // console.log(`Instruction: ${instruction}`);
    const map: Record<string, string> = {
        'Picking a word.': '/PickingWord.gif',
        'waiting for opponent word': '/waiting.gif',
        'Ready. Waiting for you...': '/Ready.gif',
        'You win!': '/Winner.gif',
        'Opponent won!': '/Winner.gif',
        'Game Over! The word was:': '/GameOver.gif',
        'Type a letter to guess.': '/GuessLetter.gif',
        'Enter a word for your opponent to guess:': '/EnterWord.gif',
        'Wants to play again.': '/WantsToPlayAgain.gif',
        '': '/Guessing.gif',
    };
    return map[instruction] ?? null;
};

export default function Game({
    playerName,
    revealedWord,
    attemptsLeft,
    guessedLetters,
    instruction,
    isMobile,
    guessing,
}: GameProps) {
    const gif = getInstructionGif(instruction);
    const nameChar = playerName.toLowerCase().split('');
    const revealedChar = revealedWord.toLowerCase().split('');
    const lettersToDisplay = guessedLetters ?? '';

    const getImageName = (char: string) => {
        if (char === "_") return "Underscore";
        if (char === " ") return "space";
        return char;
    };

    return (
        <div className="items-center text-center" style={{ flex: 1, padding: '1rem' }}>
            <h2 style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
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

            {(guessing && lettersToDisplay) ? (
                <div id="guessedLetters" className="flex flex-wrap justify-center">
                    <div className="border border-white rounded px-4 py-2 flex flex-wrap justify-center">
                        {lettersToDisplay.split('').map((char: string, index: number) => (
                            <span key={index} className="border-white mx-2">
                                {char}
                            </span>
                        ))}
                    </div>
                </div>
            ) : null}

            {gif ?
                <Image
                    src={gif}
                    width={250}
                    height={100}
                    alt="Instruction GIF"
                    className="mx-auto"
                />
                :
                <p>{instruction}</p>
            }

            <div
                className={`w-full flex justify-center overflow-x-auto ${(isMobile && !guessing && instruction !== 'waiting for opponent word') ? 'border-b-2 border-white' : ''
                    }`}
            >
                <h3 className="flex flex-wrap justify-center">
                    {revealedChar.map((char, index) => (
                        <img
                            key={index}
                            src={`/${getImageName(char)}.gif`}
                            alt={char}
                            className="h-10 sm:h-12 md:h-16"
                        />
                    ))}
                </h3>
            </div>
        </div >
    );
}
