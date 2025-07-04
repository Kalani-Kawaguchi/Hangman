const LobbyBox = ({ name, id, playerCount, maxPlayers, onClick }) => {
    return (
        <button
            className="lobby-box"
            onClick={() => onClick(id)}
            style={{
                width: '320px',
                height: '100px',
                backgroundImage: `url('/images/box.gif')`,
                backgroundSize: 'cover',
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                padding: '0 20px',
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

            <div className="lobby-name">{name}</div>
            <div className="player-count">
                {playerCount}/{maxPlayers}
            </div>
        </button>
    )
}

export default LobbyBox
