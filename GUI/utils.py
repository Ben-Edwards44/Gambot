FILES = ["a", "b", "c", "d", "e", "f", "g", "h"]

WHITE_PIECES = ["P", "N", "B", "R", "K", "Q"]
BLACK_PIECES = [i.lower() for i in WHITE_PIECES]


def str_to_square(square):
    x = 8 - int(square[1])
    y = FILES.index(square[0])

    return x, y


def square_to_str(x, y):
    rank = 8 - x
    file = FILES[y]

    return f"{file}{rank}"


def str_to_move(move):
    start_x, start_y = str_to_square(move[:2])
    end_x, end_y = str_to_square(move[2:])

    return start_x, start_y, end_x, end_y


def move_to_str(start_x, start_y, end_x, end_y):
    start = square_to_str(start_x, start_y)
    end = square_to_str(end_x, end_y)

    return f"{start}{end}"


def make_move(move, board):
    #NOTE: need to do castling / ep / promotions etc.
    start_x, start_y, end_x, end_y = str_to_move(move)

    piece_val = board[start_x][start_y]

    board[start_x][start_y] = 0
    board[end_x][end_y] = piece_val

    return board


def fen_to_board(fen):
    b_fen = fen.split(" ")[0]
    ranks = b_fen.split("/")

    board = []
    for i in ranks:
        inx = 0
        rank = [0 for _ in range(8)]

        for x in i:
            if x in WHITE_PIECES:
                rank[inx] = WHITE_PIECES.index(x) + 1
                inx += 1
            elif x in BLACK_PIECES:
                rank[inx] = BLACK_PIECES.index(x) + 7
                inx += 1
            else:
                inx += int(x)

        board.append(rank)

    return board