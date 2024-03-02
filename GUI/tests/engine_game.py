import draw
import run_gui
import engine_interface

from random import randint


FEN_FILEPATH = "../data/equal_fens.txt"

WHITE_PIECES = ["P", "N", "B", "R", "K", "Q"]
BLACK_PIECES = [i.lower() for i in WHITE_PIECES]

SHOW_GRAPHICS = True

MOVE_TIME = 500


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


def choose_fens(num):
    with open(FEN_FILEPATH, "r") as file:
        fens = file.read().split("\n")

    chosen = []
    for _ in range(num):
        inx = randint(0, len(fens) - 1)
        fen = fens.pop(inx)

        chosen.append(fen)

    return chosen


def check_win(white, black):
    #check for a win - assumes position has been set for both engines

    white_moves = white.get_perft_nodes(1)
    black_moves = black.get_perft_nodes(1)

    if white_moves == 0:
        return "draw" if black_moves == 0 else "white"
    elif black_moves == 0:
        return "black"  #white moves are not 0
    else:
        return "no_win"
    

def play_game(fen, white, black):
    #play a game between two engines, return the winner or "draw"

    white_to_move = fen.split(" ")[1] == "w"

    win = "no_win"
    move_list = []
    board = fen_to_board(fen)

    while win == "no_win":
        if white_to_move:
            engine = white
        else:
            engine = black

        engine.set_fen(fen, move_list)

        move = engine.get_move(movetime=MOVE_TIME)
        move_list.append(move)

        white_to_move = not white_to_move

        if SHOW_GRAPHICS:
            board = run_gui.make_move(move, board)
            draw.draw_board(board)

        if len(move_list) > 1:
            #TODO: draws by repetition
            win = check_win(white, black)

    return win


def main(path1, path2, num):
    fens = choose_fens(num)

    engine1 = engine_interface.Engine(path1)
    engine2 = engine_interface.Engine(path2)

    if SHOW_GRAPHICS:
        draw.init()

    win1 = 0
    draws = 0
    win2 = 0
    for i, x in enumerate(fens):
        for j in range(2):
            if j == 0:
                white = engine1
                black = engine2
            else:
                white = engine2
                black = engine1

            winner = play_game(x, white, black)

            if winner == "draw":
                draws += 1
            elif winner == "white":
                if white == engine1:
                    win1 += 1
                else:
                    win2 += 1
            else:
                if black == engine1:
                    win1 += 1
                else:
                    win2 += 1

        print(f"Played: {i + 1}\n{path1} wins: {win1}\n{path2} wins: {win2}\nDraws: {draws}\n\n")

    print(f"Final result:\n{path1} wins: {win1}\n{path2} wins: {win2}\nDraws: {draws}\n\n")