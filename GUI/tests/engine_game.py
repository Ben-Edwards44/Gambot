import draw
import board
import engine_interface

from random import randint


WHITE_PIECES = ["P", "N", "B", "R", "K", "Q"]
BLACK_PIECES = [i.lower() for i in WHITE_PIECES]

FEN_FILEPATH = "../data/equal_fens.txt"

SHOW_GRAPHICS = False
PRINT_LOSSES = False

MOVE_TIME = 500


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
        return "draw" if black_moves == 0 else "black"
    elif black_moves == 0:
        return "white"  #white moves are not 0
    else:
        return "no_win"
    

def play_game(fen, white, black):
    #play a game between two engines, return the winner or "draw"

    white.new_game()
    black.new_game()

    white_to_move = fen.split(" ")[1] == "w"

    board_obj = board.Board(fen, None, None, None)

    win = "no_win"
    move_list = []
    seen_boards = {}

    while win == "no_win":
        if white_to_move:
            engine = white
        else:
            engine = black

        engine.set_fen(fen, move_list)

        if len(move_list) > 1:
            win = check_win(white, black)
            
            if win != "no_win":
                break

        move = engine.get_move(movetime=MOVE_TIME)
        move_list.append(move)

        white_to_move = not white_to_move

        board_obj.make_move(move)
        t_board = tuple(tuple(i) for i in board_obj.board_list)

        if t_board in seen_boards:
            num = seen_boards[t_board]

            if num >= 2:
                win = "draw"  #draw by repetition
                break
            else:
                seen_boards[t_board] += 1
        else:
            seen_boards[t_board] = 1

        if SHOW_GRAPHICS:
            draw.draw_board(board_obj)

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

                    if PRINT_LOSSES:
                        print(x)
            elif winner == "black":
                if black == engine1:
                    win1 += 1
                else:
                    win2 += 1

                    if PRINT_LOSSES:
                        print(x)
            else:
                raise Exception(f"Winner ({winner}) was not draw, white or black")

        print(f"Played: {i + 1}\n{path1} wins: {win1}\n{path2} wins: {win2}\nDraws: {draws}\n\n")

    print(f"Final result:\n{path1} wins: {win1}\n{path2} wins: {win2}\nDraws: {draws}\n\n")

    engine1.kill_process()
    engine2.kill_process()