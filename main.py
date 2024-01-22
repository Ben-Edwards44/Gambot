from sys import argv
import src.engine as engine


def run_engine():
    engine.init()
    
    while True:
        engine.main()


def perft():
    depth = int(argv[2])
    
    engine.perft(depth)


def main():
    mode = argv[1]

    if mode == "run":
        run_engine()
    elif mode == "perft":
        perft()


if __name__ == "__main__":
    print("HAVE YOU COMPILED THE UP TO DATE GO SCRIPT??")
    main()