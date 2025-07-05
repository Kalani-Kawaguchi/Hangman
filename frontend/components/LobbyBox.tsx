const LobbyBox = ({ name, id, playerCount, maxPlayers, onClick }) => {
    const characters = name.toLowerCase().split('');
    const nameLength = name.length;
    var widthPercent = "auto";
    if (nameLength >= 9){
        widthPercent = ((1 / name.length) * 100) + "%";
    }
     

    return (
        <button
            className="lobby-box"
            onClick={() => onClick(id)}
            style={{
                width: '100%',
                height: '10vh',
                backgroundImage: `url('/images/box.gif')`,
                backgroundSize: 'cover',
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                border: 'none',
                cursor: 'pointer',
                fontSize: '24px',
                color: 'white',
                fontFamily: 'monospace',
            }}
        >
            {/* <div className="lobby-name"> */}
            {/*     {name.split('').map((char, i) => ( */}
            {/*         <img */}
            {/*             key={i} */}
            {/*             src={`/images/letters/${char.toUpperCase()}.gif`} */}
            {/*             alt={char} */}
            {/*             style={{ height: '50px', marginRight: '2px' }} */}
            {/*         /> */}
            {/*     ))} */}
            {/* </div> */}

            <div className="lobby-name" style={{display: 'flex', justifyContent: 'center', alignItems: 'center', width: '75%', height: '50%'}}>
                {characters.map((char, index) => (
                    <img
                    key={index}
                    src={`/${char}.gif`}
                    alt={char}
                    style={{width: `${widthPercent}`, height: `auto`}}
                    />
                ))}
            </div>
            <div className="player-count" style={{width: '25%', height: '50%'}}>
                {playerCount}/{maxPlayers}
            </div>
        </button>
    )
}

export default LobbyBox
