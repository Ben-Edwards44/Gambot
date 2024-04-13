# Gambot :chess_pawn:

![Gambot logo](assets/gambot_logo.png)

Gambot is a work in progress UCI chess engine written in Go.
The repository also contains a GUI written in Python using the Pygame library.

## Features :wrench:

- Working move generator (it passes all perfts on the [chess programming wiki](https://www.chessprogramming.org/Perft_Results))
- Minimax search with alpha-beta pruning
- Quiescence search
- PVS search
- Move ordering with MVV/LVA, hash moves and killer moves
- Zobrist hashing for board representation
- Transposition table
- Repetition table
- Middle game and endgame piece square tables
- Late move reductions
- Supported UCI commands:
    - `uci`
    - `isready`
    - `ucinewgame`
    - `position <fen | startpos> <moves>`
    - `go <wtime> <btime> <winc> <binc> <movetime>`
- Non-UCI commands:
    - `go perft <depth>` (perform perft)
    - `eval` (display the static evaluation of the current position)

## Acknowledgements :link:
- Without [this video](https://www.youtube.com/watch?v=U4ogK0MIzqk&t=1191s) and [this video](https://www.youtube.com/watch?v=_vqlIPDR2TU&t=886s) from Sebastian Lague, I never would have considered making a chess engine.
- The [source code](https://github.com/SebLague/Chess-Coding-Adventure) from the Sebastian Lague videos was a very useful reference.
- The [Tofiks](https://github.com/likeawizard/tofiks) and [Blunder](https://github.com/algerbrex/blunder) chess engines were also very useful.
- The [Chess Programming Wiki](https://www.chessprogramming.org/Main_Page) was an invaluable resource for anything and everything related to chess engine programming.
- The chess pieces used in the GUI are from [here](https://commons.wikimedia.org/wiki/Category:PNG_chess_pieces/Standard_transparent).