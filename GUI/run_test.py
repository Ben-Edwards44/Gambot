from tests import perft_suite, engine_game, speed_test

from sys import argv


def perft():
    perft_suite.main()


def eng_game():
    if len(argv) < 5:
        raise Exception("Invalid args for engine game")
    
    path1 = argv[2]
    path2 = argv[3]
    num = int(argv[4])

    engine_game.main(path1, path2, num)


def speed():
    if len(argv) < 5:
        raise Exception("Invalid args for speed test")
    
    path1 = argv[2]
    path2 = argv[3]
    num = int(argv[4])

    speed_test.main(path1, path2, num)


def main():
    if len(argv) < 1:
        raise Exception("Invalid argument number for tests")
    
    test = argv[1]  #the value at index 0 will be the name of the script

    if test == "perft":
        perft()
    elif test == "engine_game":
        eng_game()
    elif test == "speed_test":
        speed()
    else:
        raise Exception("Invalid test type")
    

if __name__ == "__main__":
    main()