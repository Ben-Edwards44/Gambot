from sys import argv


WHITE_PIECES = ["P", "N", "B", "R", "K", "Q"]
BLACK_PIECES = [i.lower() for i in WHITE_PIECES]


def convert_rank(fen):
    inx = 0
    rank = [0 for _ in range(8)]

    for i in fen:
        if i in WHITE_PIECES:
            rank[inx] = WHITE_PIECES.index(i) + 1
            inx += 1
        elif i in BLACK_PIECES:
            rank[inx] = BLACK_PIECES.index(i) + 7
            inx += 1
        else:
            inx += int(i)

    return rank


def convert_fen(fen):
    ranks = fen.split("/")

    board = []
    for i in ranks:
        rank = convert_rank(i)
        board.append(rank)

    return board


def main():
    try:
        fen = argv[1]
    except IndexError:
        raise Exception("Invalid args. Use python fen_to_board.py [FEN STRING]")
    
    fen = fen.split(" ")
    board = convert_fen(fen[0])

    print(board)


if __name__ == "__main__":
    main()