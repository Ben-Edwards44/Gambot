import engine_interface

from random import randint


FEN_FILEPATH = "../data/equal_fens.txt"

MULT = {
    "ns" : 1_000_000,
    "µs" : 1_000,
    "ms" : 1,
    "s" : 1 / 1_000
}

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


def get_time(output):
    #gets the time elapsed (in ms) from an output like: Depth: 4, Searched: 1116, tt Lookups: 65, Score: 0, Elapsed: 2.226884ms

    args = output.split(", ")
    
    for i in args:
        name, value = i.split(": ")

        if name == "Elapsed":
            time = value
            break

    num = ""
    unit = ""

    for i in time:
        if i in "0123456789.":
            num += i
        else:
            unit += i

    num_ms = float(num) * MULT[unit]

    return num_ms


def engine_move(engine, fen):
    engine.set_fen(fen, [])

    engine.send_cmd(f"go movetime {MOVE_TIME}")

    output = ""
    times = []
    while len(output) < 8 or output[:8] != "bestmove":
        output = engine.read_line()

        if output[:8] != "bestmove":
            #we don't want to try and get the time for an output like "bestmove e2e4"
            time = get_time(output)
            times.append(time)

    #NOTE: only the second to last depth is reliable
    depth_searched = len(times) - 1
    final_time = times[-1]

    return depth_searched, final_time


def get_best(engine1, engine2, fen):
    #return whichever engine predicted a move the fastest

    d1, t1 = engine_move(engine1, fen)
    d2, t2 = engine_move(engine2, fen)

    #NOTE: if one engine sreaches to a greater depth, that is always better

    if d1 > d2:
        return engine1
    elif d2 > d1:
        return engine2
    else:
        if t1 < t2:
            return engine1
        else:
            return engine2


def main(path1, path2, num):
    engine1 = engine_interface.Engine(path1, True)
    engine2 = engine_interface.Engine(path2, True)

    fens = choose_fens(num)

    win1 = 0
    win2 = 0
    for i, x in enumerate(fens):
        fastest = get_best(engine1, engine2, x)

        if fastest == engine1:
            win1 += 1
        else:
            win2 += 1

        print(f"Position: {i + 1}\n{path1} speed wins: {win1}\n{path2} speed wins: {win2}")

    print(f"Final result\n{path1} speed wins: {win1}\n{path2} speed wins: {win2}")

    engine1.kill_process()
    engine2.kill_process()