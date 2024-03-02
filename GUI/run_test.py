from tests import perft_suite

from sys import argv


def main():
    if len(argv) < 1:
        raise Exception("Invalid argument number for tests")
    
    test = argv[1]  #the value at index 0 will be the name of the script

    if test == "perft":
        perft_suite.main()
    else:
        raise Exception("Invalid test type")
    

if __name__ == "__main__":
    main()