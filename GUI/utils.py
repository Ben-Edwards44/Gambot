FILES = ["a", "b", "c", "d", "e", "f", "g", "h"]


def str_to_square(square):
    x = 8 - int(square[1])
    y = FILES.index(square[0])

    return x, y


def square_to_str(x, y):
    rank = 8 - x
    file = FILES[y]

    return f"{file}{rank}"


def str_to_move(move):
    start_x, start_y = str_to_square(move[:2])
    end_x, end_y = str_to_square(move[2:])

    return start_x, start_y, end_x, end_y


def move_to_str(start_x, start_y, end_x, end_y):
    start = square_to_str(start_x, start_y)
    end = square_to_str(end_x, end_y)

    return f"{start}{end}"