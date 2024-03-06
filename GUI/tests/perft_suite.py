import engine_interface

from time import time


#all perfts are from here: https://www.chessprogramming.org/Perft_Results


PERFT_FENS = {
    "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" : [20, 400, 8902, 197281, 4865609],
    "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1" : [48, 2039, 97862, 4085603, 193690690],
    "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1" : [14, 191, 2812, 43238, 674624],
    "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1" : [6, 264, 9467, 422333, 15833292],
    "r2q1rk1/pP1p2pp/Q4n2/bbp1p3/Np6/1B3NBn/pPPP1PPP/R3K2R b KQ - 0 1" : [6, 264, 9467, 422333, 15833292],
    "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8" : [44, 1486, 62379, 2103487, 89941194],
    "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10" : [46, 2079, 89890, 3894594, 164075551]
}


PASS_COLOUR = "\033[92m"
FAIL_COLOUR = "\033[91m"
END_COLOUR = "\033[0m"


def show_fail(fen, depth, expected, actual):
    print(f"{FAIL_COLOUR}Position {fen} failed at depth {depth}. Expected: {expected}, got {actual}.{END_COLOUR}")


def show_pass(fen):
    print(f"{PASS_COLOUR}Position {fen} passed.{END_COLOUR}")


def test_position(engine, position, expected):
    engine.set_fen(position, [])

    for i, x in enumerate(expected):
        depth = i + 1
        nodes = engine.get_perft_nodes(depth)

        if x != nodes:
            show_fail(position, depth, x, nodes)
            return False
    
    return True


def main():
    engine = engine_interface.Engine()

    start = time()
    for fen, expected in PERFT_FENS.items():
        passed = test_position(engine, fen, expected)

        if passed:
            show_pass(fen)

    elapsed = time() - start

    print(f"Time elapsed: {elapsed :.2f}s")

    engine.kill_process()  #ensure process is killed