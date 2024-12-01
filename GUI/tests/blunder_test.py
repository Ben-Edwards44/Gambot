import draw
import board
import engine_interface

import chess
import chess.engine
from random import randint
from math import exp


FEN_FILEPATH = "../data/equal_fens.txt"
STOCKFISH_PATH = ""

SHOW_GRAPHICS = False

MOVE_TIME = 500
EVAL_DEPTH = 15

BLUNDER_THRESHOLD = 50  #percent


class AnalysisEngine:
    def __init__(self, path, start_fen):
        self.board = chess.Board(start_fen)
        self.engine = chess.engine.SimpleEngine.popen_uci(path)
        
        self.prev_win_frac = None
        self.prev_fens = []

    def update_board(self, move):
        current_fen = self.board.fen()
        self.prev_fens.append(current_fen)

        move_obj = chess.Move.from_uci(move)

        self.board.push(move_obj)

    def get_eval(self):
        eval = self.engine.analyse(self.board, chess.engine.Limit(depth=EVAL_DEPTH))

        score_obj = eval["score"]

        if self.board.turn == chess.WHITE:
            cp_score = score_obj.white()
        else:
            cp_score = score_obj.black()

        act_score = cp_score.score(mate_score=100000)

        return act_score

    def check_blunder(self):
        eval = self.get_eval()

        win_frac = 50 + 50 * (2 / (1 + exp(-0.00368208 * eval)) - 1)  #https://lichess.org/page/accuracy
        
        prv = self.prev_win_frac
        self.prev_win_frac = win_frac

        if prv == None:
            return False, 0
        
        accuracy = 103.1668 * exp(-0.04354 * (prv - win_frac)) - 3.1669  #https://lichess.org/page/accuracy

        return accuracy < BLUNDER_THRESHOLD, accuracy


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

    analysis = AnalysisEngine(STOCKFISH_PATH, fen)

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
        analysis.update_board(move)

        is_blunder, diff = analysis.check_blunder()

        if is_blunder:
            print(f"Blunder from {analysis.prev_fens[-2]} to {analysis.prev_fens[-1]}. Accuracy: {diff}")

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


def main(num):
    fens = choose_fens(num)

    engine1 = engine_interface.Engine()
    engine2 = engine_interface.Engine()

    if SHOW_GRAPHICS:
        draw.init()

    for i, x in enumerate(fens):
        play_game(x, engine1, engine2)

        print(f"Played: {i + 1}")

    engine1.kill_process()
    engine2.kill_process()