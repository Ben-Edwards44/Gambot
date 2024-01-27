def get_data(filename):
    with open(filename, "r") as file:
        data = file.read()

    return data.split("\n")


def build_dict(data):
    dict = {}

    for i in data:
        move, num = i.split(": ")
        dict[move] = int(num)

    return dict


def find_wrong(my_dict, stockfish_dict):
    for k, v in my_dict.items():
        if k not in stockfish_dict.keys() or stockfish_dict[k] != v:
            print(k)

    #check for missed start moves
    for i in stockfish_dict.keys():
        if i not in my_dict.keys():
            print(i)


def main():
    my_data = get_data("my_moves.txt")
    stock_data = get_data("stockfish_moves.txt")

    my_dict = build_dict(my_data)
    stock_dict = build_dict(stock_data)

    find_wrong(my_dict, stock_dict)


if __name__ == "__main__":
    main()