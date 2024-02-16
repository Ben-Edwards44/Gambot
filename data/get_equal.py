import chess
import chess.engine
from random import shuffle


EQUAL_CUTOFF = 30  #in centipawns
NUM_GAMES_SEARCH = 10000
ENGINE_PATH = "C:\\Users\\Ben Edwards\\Documents\\stockfish-windows-x86-64-avx2\\stockfish\\stockfish-windows-x86-64-avx2.exe"


def get_positions():
    with open("fen_data.txt", "r") as file:
        data = file.read()

    return data.split("\n")


def write_file(string):
    with open("equal_fens.txt", "w") as file:
        file.write(string)


def fast_remove(nums, value):
    #assumes nums is ordered

    up = len(nums) - 1
    low = 0
    while True:
        mid = (up + low) // 2

        if nums[mid] == value:
            nums.remove(value)
            return
        elif nums[mid] < value:
            up = mid - 1
        else:
            low = mid + 1

        if low > up:
            raise Exception(f"{value} not in nums list")


def choose_random(fens, num):
    shuffle(fens)

    return fens[:num]


def get_eval(fen, engine, time_limit):
    board = chess.Board(fen)
    eval = engine.analyse(board, chess.engine.Limit(time_limit))

    score_obj = eval["score"]
    cp_score = score_obj.white()
    act_score = cp_score.score(mate_score=100000)

    return act_score


def get_equal_positions(fens, engine, time_limit, num_games):
    fens = choose_random(fens, num_games)

    equal_pos = []
    for i, x in enumerate(fens):
        eval = get_eval(x, engine, time_limit)

        if abs(eval) <= EQUAL_CUTOFF:
            equal_pos.append(x)

        done_frac = i / len(fens)
        num_hash = int(75 * done_frac)

        print(f"Finding equal positions: |{'#' * num_hash}{'.' * (75 - num_hash)}| {done_frac * 100 :.2f}%", end="\r")

    print("\nDone!")
    
    return equal_pos


def main():
    engine = chess.engine.SimpleEngine.popen_uci(ENGINE_PATH)

    positions = get_positions()
    eq_pos = get_equal_positions(positions, engine, 0.1, NUM_GAMES_SEARCH)

    write_file("\n".join(eq_pos))


main()