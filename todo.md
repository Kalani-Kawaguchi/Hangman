# TODO:

* If player 2 leaves after game is over, and another player joins before
player 1 clicks play again, player 2 can't submit a word.

* If player2 leaves in the middle of the game, we should not allow another player
to join until the game has been completed.
    - Right now another player can join and the lobby gets stuck
    - If player1 finishes, they don't get the option to play again since
    player2 never finishes

* If we refresh the page it looks like the lobby states on the front end don't persist.
    - ex: If a player submits a word, then refreshes the page, on refresh their
    screen will allow them to submit another word. It should refresh and stay 
    - ex: If player1 refreshes after submitting word, the screen will also revert
    to p2's game being hidden, as if a second player never joined.
    - We probably need to save these states on the server side, and have something to
    check the state on the server so that if we do refresh we get put back into 
    the correct state.

