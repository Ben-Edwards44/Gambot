import engine_interface


TESTS = {
    "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" : "g1h3"
}


PASS_COLOUR = "\033[92m"
FAIL_COLOUR = "\033[91m"
END_COLOUR = "\033[0m"

MOVE_TIME = 500


def show_fail(fen, expected, actual):
    print(f"{FAIL_COLOUR}Position {fen} failed. Expected move: {expected}, actual move: {actual}.{END_COLOUR}")


def show_pass(fen):
    print(f"{PASS_COLOUR}Position {fen} passed.{END_COLOUR}")


def test_pos(engine, fen, expected_move):
    engine.new_game()
    engine.set_fen(fen, [])

    move = engine.get_move(movetime=MOVE_TIME)

    if move == expected_move:
        show_pass(fen)
    else:
        show_fail(fen, expected_move, move)


def main():
    engine = engine_interface.Engine()

    for fen, move in TESTS.items():
        test_pos(engine, fen, move)

    engine.kill_process()