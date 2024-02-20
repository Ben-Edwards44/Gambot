import src.engine as engine
import src.test_engines as test_engines

from sys import argv


def run_engine():
    engine.init()
    
    while True:
        engine.main()


def perft():
    depth = int(argv[2])

    if len(argv) == 4:
        test = argv[3] == "test"
    else:
        test = False
    
    test_engines.perft(depth, test)


def engine_game():
    engine1 = argv[2]
    engine2 = argv[3]
    num_games = int(argv[4])

    test_engines.engine_game(engine1, engine2, num_games)


def speed_test():
    engine1 = argv[2]
    engine2 = argv[3]
    num_games = int(argv[4])

    test_engines.speed_test(engine1, engine2, num_games)


def main():
    mode = argv[1]

    if mode == "run":
        run_engine()
    elif mode == "perft":
        perft()
    elif mode == "engine_game":
        engine_game()
    elif mode == "speed_test":
        speed_test()


if __name__ == "__main__":
    print("HAVE YOU COMPILED THE UP TO DATE GO SCRIPT??")
    main()